# Certificate Expiration Exporter

## Try it out locally

First update your local gcloud credentials running and following instructions

```
$ cloud auth login
```

Then choose the project where you have stackdriver ativated

```
$ gcloud config set project PROJET_NAME
```

Activate the Monitoring API (You can also check https://cloud.google.com/monitoring/api/enable-api)

```
$ gcloud services list # just to check what we have
$ gcloud services enable monitoring
```

Then you can run

```
$ EXPORTER=stackdriver,prometheus go run ./main.go
```

To view it on the console, click on the hamburguer button, then on Monitoring under the Operations section.

Under Metrics Explorer start typing `certificate_expiration_rem...`. It should autocomplete and find the stat that you are sending to stackdriver.

# ⚠️ Please bear in mind

This is a WIP for wg-k8s-infra - Monitoring certs expiration
