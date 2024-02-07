# GinRCON
A lightweight, simple Go webserver build on Gin to provide a REST API frontend to sending RCON commands to a game server.

## Calling the API

## Docker
### Image
A docker image that runs this application is available, `ghcr.io/holysoles/ginrcon`
### Building
`docker build --tag ginrcon .`
### Testing
`docker run --rm -d -p 8080:8080 ginrcon:latest`

## Credits
Special thanks to the following projects for providing the essential libraries:
- [gorcon/rcon](https://github.com/gorcon/rcon)
- [gin-gonic/gin](https://github.com/gin-gonic/gin)