# Cleaning up MiniO Objects as Kubernetes CronJob

Cleaning up MiniO Objects older than 5 days only except `loki_cluster_seed.json`.  


## 1. Set your MiniO endpoint
Replace variable `endpoint` in the `main.go` file with your MiniO endpoint so the Pod can reach your MiniO.

## 2. Build image
Run following command to build container image.
```sh
docker build -t yourContainerRegistry/imageName .
```

## 3. Create Kubernetes secret

```sh
kubectl create secret generic minio-credential --from-literal=ACCESS_KEY_ID=<MiniO-id> --from-literal=SECRET_ACCESS_KEY=<MiniO-password>
```

This command will populate Kubernetes secret with your MiniO credential.

## 4. Deploy CronJob
First, you have to replace
`spec.jobTemplate.spec.template.spec.containers[0].image` field with your container image that you build at #1.  

Then you apply CronJob manifest.

```sh
kubectl apply -f kubernetes/go-minio-cleanup-cronjob
```

Now It will automatically start to cleaning at 00:00 (It is depends on your Kubernetes cluster TimeZone).
