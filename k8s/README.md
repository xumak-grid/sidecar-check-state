## Test the example

Create the ns, configMap and the pod

`kubectl create -f k8s/ns-configMap.yaml`

`kubectl create -f k8s/sidecar-example.yaml`

Open another terminal to see logs from the two containers

`kubectl -n demo logs -f -c test-container sidecar-check-state`

`kubectl -n demo logs -f -c test-sidecar sidecar-check-state`

Edit the content of the configMap using the dashboard of with kubectl apply...

Delete stuff

`kubectl delete -f k8s/sidecar-example.yaml`

`kubectl delete -f k8s/ns-configMap.yaml`
