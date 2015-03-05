build/start-influxdb: *.go
	godep go build -o build/start-influxdb

image:
	docker build --no-cache -t start-influxdb .

release:
	docker tag start-influxdb fabric8/start-influxdb
	docker push fabric8/start-influxdb

.PHONY: clean
clean:
	rm -rf build
