# Kubernetes Task Runner



## Spin up local k8s cluster

```
k3d registry create registry.localhost --port 5000
k3d cluster create newcluster --registry-use k3d-registry.localhost:5000 -p "8080:80@loadbalancer" --agents 2
```

Add entry to `/etc/hosts
```
localhost k3d-registry.localhost
```

## Push image to k3d Docker registry

See https://k3d.io/usage/guides/registries/#using-a-local-registry

```
docker tag <image> k3d-registry.localhost:5000/<image>
docker push k3d-registry.localhost:5000/<image>
```

Test
```
curl -X POST -H "Content-Type: application/json" --data '{"image":"k3d-registry.localhost:5000/simple_pytest:0.1","command":["sh", "-c"],"args":["pytest ./tests/ --junitxml output.xml && cat output.xml"]}' localhost:8080/

# Successful test
curl -X POST -H "Content-Type: application/json" --data '{"image":"k3d-registry.localhost:5000/simple_pytest:0.1","command":["sh", "-c"],"args":["pytest ./tests/ -k test_sample_1 --junitxml output.xml && cat output.xml"]}' localhost:8080/

# Store Id to variable
JOB_ID=$(curl -X POST -H "Content-Type: application/json" --data '{"image":"k3d-registry.localhost:5000/simple_pytest:0.1","command":["sh", "-c"],"args":["pytest ./tests/ -k test_sample_1 --junitxml output.xml && cat output.xml"]}' localhost:8080/ | jq -r .id)
```

Get test results
```
curl -s localhost:8080/$JOB_ID | jq .

# Get test logs
curl -s localhost:8080/$JOB_ID | jq -r .logs
```
