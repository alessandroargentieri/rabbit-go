# Instructions

---

## With Minikube

- Go to [Katacoda Minikube environment](https://www.katacoda.com/scenario-examples/courses/environment-usages/minikube):

- Start minikube with: `minikube start`

- Clone the [Github project](https://github.com/alessandroargentieri/rabbit-go.git)

- enter the folder `cd rabbit-go`

### Publish docker images of the two services: publisher and consumer

You can skip this step if the images are already present to the public repository of Dockerhub according to the name specified in the _yaml_ files.

```
# docker login
docker login -u=<username-here> -p=<password-here>

# create publisher image and push it on DockerHub
cd publisher
docker build -f Dockerfile -t alessandroargentieri/rabbit-go-publisher .
docker push alessandroargentieri/rabbit-go-publisher
cd ..

# create consumer image and push it on DockerHub
cd consumer
docker build -f Dockerfile -t alessandroargentieri/rabbit-go-consumer .
docker push alessandroargentieri/rabbit-go-consumer
cd ..
```

### Create k8s secret for the rabbitmq credentials

```
# create secret
kubectl create secret generic rabbit-credentials --from-literal=password=password --from-literal=user=myuser

# get yaml of the created secret
kubectl get secret rabbit-credentials -o yaml
```

### Deploy yaml to k8s

```
kubectl apply -f k8s/yamls/rabbit.yaml
kubectl apply -f k8s/yamls/publisher.yaml
kubectl apply -f k8s/yamls/consumer.yaml
```

### Get the NodePort addresses

```
# copy the addresses for the publisher nodeport and the consumer nodeport
minikube service list
```

### Call the publisher service through its NodePort service (on port 32000)

```
curl http://<publisher-copied-ip>:32000/publisher
```
### Call the consumer service through its NodePort service (on port 32100)

```
curl http://<consumer-copied-ip>:32100/consumer
```
---

## With Go and Docker

- Go to [Katacoda Ubuntu environment](https://www.katacoda.com/courses/ubuntu/playground):

- Clone the [Github project](https://github.com/alessandroargentieri/rabbit-go.git)

- enter the folder `cd rabbit-go`

- start everything with the script `./init.sh`