
PONY: build push clean local-test
IMAGE=
VERSION=0.0.1

build:
	GOOS=linux go build -o bin/check-state main.go
	docker build -t $(IMAGE):$(VERSION) .

push:
	docker push $(IMAGE):$(VERSION)

clean:
	docker rmi $(IMAGE):$(VERSION) || true
	rm bin/check-config || true
	rm file_test_1.txt file_test_2.txt || true
	rm empty_file_test || true

local-test: clean
	echo "hello 1" > file_test_1.txt
	echo "hello 2" > file_test_2.txt
	export SIDE_CAR_CONFIG_FILES="file_test_1.txt,file_test_2.txt"; go run main.go

k8s-test: clean
	kubectl create -f k8s/ns-configMap.yaml
	kubectl create -f k8s/sidecar-example.yaml

k8s-clean:
	kubectl delete -f k8s/sidecar-example.yaml
	kubectl delete -f k8s/ns-configMap.yaml
