# Cloud Events Viewer

Websocket service to view your cloud events.

## Build

```
ko publish -B github.com/itsmurugappan/cloud-events-viewer/cmd/cev --platform=all
```

## Deploy

```
kn service create cev --image=docker.io/murugappans/cev:latest
```