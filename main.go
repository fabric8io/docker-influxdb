package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"text/template"

	"github.com/GoogleCloudPlatform/kubernetes/pkg/api"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/client"
	"github.com/GoogleCloudPlatform/kubernetes/pkg/labels"
)

type Context struct {
}

func (c *Context) Env() map[string]string {
	env := make(map[string]string)
	for _, i := range os.Environ() {
		sep := strings.Index(i, "=")
		env[i[0:sep]] = i[sep+1:]
	}
	return env
}

func main() {
	if len(os.Getenv("KUBERNETES_SELECTOR")) > 0 && len(os.Getenv("KUBERNETES_MASTER")) > 0 {
		kubeMaster := os.Getenv("KUBERNETES_MASTER")
		if !(strings.HasPrefix(kubeMaster, "http://") || strings.HasPrefix(kubeMaster, "https://")) {
			kubeMaster = "http://" + kubeMaster
		}
		insecure, _ := strconv.ParseBool(os.Getenv("KUBERNETES_INSECURE"))
		kubeClient := client.NewOrDie(&client.Config{
			Host:     os.ExpandEnv(kubeMaster),
			Insecure: len(os.Getenv("KUBERNETES_INSECURE")) > 0 && insecure,
		})

		selector, err := labels.ParseSelector(os.Getenv("KUBERNETES_SELECTOR"))
		if err != nil {
			log.Println(err)
		} else {
			namespace := os.Getenv("KUBERNETES_NAMESPACE")
			if len(namespace) == 0 {
				namespace = api.NamespaceDefault
			}

			podList, err := kubeClient.Pods(namespace).List(selector)
			if err != nil {
				log.Println(err)
			} else {
				var seeds string
				currentHostname, err := os.Hostname()
				if err != nil {
					log.Println(err)
				}
				for _, pod := range podList.Items {
					if pod.Status.Phase == api.PodRunning && len(pod.Status.PodIP) > 0 && pod.Name != currentHostname {
						if len(seeds) > 0 {
							seeds += ","
						}
						seeds += fmt.Sprintf("http://%v:%v", pod.Status.PodIP, os.Getenv("INFLUXDB_BROKER_PORT"))
					}
				}
				os.Setenv("INFLUXDB_SEEDS", seeds)
			}
		}
	}

	if len(os.Getenv("INFLUXDB_SEEDS")) == 0 {
		os.Setenv("INFLUXDB_SEEDS", "")
	}

	addrs, err := net.InterfaceAddrs()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	for _, address := range addrs {
		// check the address type and if it is not a loopback the display it
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				os.Setenv("IP_ADDRESS", ipnet.IP.String())
				break
			}
		}
	}

	t, _ := template.ParseFiles("/opt/influxdb/influxdb.conf.tmpl")

	file, err := os.Create("/opt/influxdb/influxdb.conf")
	if err != nil {
		log.Fatal(err)
	}
	t.Execute(file, &Context{})
	file.Close()

	if err := os.MkdirAll(os.Getenv("INFLUXDB_DATA_DIR"), 0755); err != nil {
		if !os.IsExist(err) {
			log.Fatal(err)
		}
	}

	cmd := exec.Command("/bin/sh", "-c", os.ExpandEnv("chown -R ${INFLUXDB_USER}:${INFLUXDB_USER} ${INFLUXDB_DATA_DIR}"))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()

	cmdStr := "exec sudo -u ${INFLUXDB_USER} -H sh -c \"cd /opt/influxdb; exec ./influxd -config ${CONFIG_FILE}\""

	cmdStr = os.ExpandEnv(cmdStr)

	cmd = exec.Command("/bin/sh", "-c", cmdStr)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
}
