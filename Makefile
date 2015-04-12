build/start-influxdb: *.go
	godep go build -o build/start-influxdb

image:
	docker build -t influxdb .

release:
	docker tag influxdb fabric8/influxdb
	docker push fabric8/influxdb

.PHONY: clean
clean:
	rm -rf build
