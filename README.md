# Cleaning up MiniO Objects as Kubernetes CronJob

Cleaning up MiniO Objects older than desired retentionPeriod.


## 1. Set your MiniO envs
Set envs in the `.env` file    
<div>
</div>

`ENDPOINT`: Endpoint of your object Storage  
`ACCESSKEYID`: ID of your object Storage  
`SECRETACCESSKEY`: Password of your object Storage  
`BUCKETNAME`: Bucket you want clean up  
`RETENTIONPERIOD`: retention period in hour (3 => 3hour)

## 2. Build image
Run following command to build container image.
```sh
docker build -t yourContainerRegistry/imageName .
```

## 3. Create Kubernetes secret

```sh
kubectl create secret generic minio-credential --from-literal=ACCESS_KEY_ID=<MiniO-id> --from-literal=SECRET_ACCESS_KEY=<MiniO-password> --from-literal=BUCKETNAME=<Bucket name> --from-literal=RETENTIONPERIOD=<retention period>
```

This command will populate Kubernetes secret with your MiniO credential.

## 4. Deploy CronJob
First, you have to replace
`spec.jobTemplate.spec.template.spec.containers[0].image` field with your container image that you build at #1.  

Then you apply CronJob manifest.

```sh
kubectl apply -f kubernetes/go-minio-cleanup-cronjob
```

Now It will automatically start to cleaning at 00:00 (It depends on your Kubernetes cluster TimeZone).
