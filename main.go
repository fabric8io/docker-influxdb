package main

import (
	"fmt"
	"log"
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
				for index, pod := range podList.Items {
					if index > 0 {
						seeds += ","
					}
					seeds += fmt.Sprintf("http://%v:%v", pod.Status.PodIP, os.Getenv("INFLUXDB_BROKER_PORT"))
				}
				os.Setenv("INFLUXDB_SEEDS", seeds)
			}
		}
	}

	if len(os.Getenv("INFLUXDB_SEEDS")) == 0 {
		os.Setenv("INFLUXDB_SEEDS", "")
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
