# GinRCON
A lightweight, simple Go webserver build on Gin to provide a REST API frontend to sending RCON commands to a game server.

## Calling the API

## Docker
### Image
A docker image that runs this application is available, `ghcr.io/holysoles/ginrcon`.
### Environment Variables
- `PORT`: Optional, override the port of the webserver within the container
- `TRUSTED_PROXIES`: Optional, set specified trusted proxy addresses
### Compose
An example compose file can be found in the repo [here](https://github.com/holysoles/ginrcon/blob/master/compose.yaml)
### Building
`docker build --tag ghcr.io/holysoles/ginrcon:<tag> .`
### Testing
`docker run --rm -d -p 8080:8081 -e "PORT=8081" ghcr.io/holysoles/ginrcon:<tag>`
### Upload
`docker push ghcr.io/holysoles/ginrcon:<tag>`
## Credits
Special thanks to the following projects for providing essential libraries:
- [gorcon/rcon](https://github.com/gorcon/rcon)
- [gin-gonic/gin](https://github.com/gin-gonic/gin)