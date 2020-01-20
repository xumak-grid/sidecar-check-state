# Sidecar check state

Check the state of a container inside a pod to allow a restart when is needed

### Steps to follow Version 0.0.1:

1. when the service is starting loads in memory the hash number of each config Files added to be checked
2. if any of the confgFiles aren't present waits and try again (* PENDINNG)
3. the service should expose a port 9090 with the path `/check/liveness`
4. the handler `checkLiveness` checks again the sha number of the configFiles and compairs with the previous
5. if some of the files have changed or they aren't present `restartRequired=true`
6. the handler responds with a status message (status code 200 or 500)

```go
type Status struct {
    Message string
    RestartRequired bool
}
```

7. configFiles can be register using env var e.g `SIDE_CAR_CONFIG_FILES="/location/file1,/location/fil2"`
the env var registers two options: single file or many using `SIDE_CAR_CONFIG_FILES="/location/*"` or conbination of the two options
8. default port `9090`
9. register the logs

Copyright © 2016 Tikal Technologies, Inc.