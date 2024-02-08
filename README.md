# GinRCON
A lightweight, simple Go webserver build on Gin to provide a REST API frontend to sending RCON commands to a game server.

## API Call Examples
### Send command to default RCON server
Assuming you have a default RCON server specified (see environment variables)
```bash
curl --location 'http://localhost:8080/command' \
--header 'Content-Type: application/json' \
--data '{
    "command": "Save"
}'
```
### POST command to specific RCON server
```bash
curl --location 'http://localhost:8080/command' \
--header 'Content-Type: application/json' \
--data '{
    "server": "gameserver:25575",
    "password": "1234",
    "command": "Save"
}'
```
### Health check the status of the webserver
```bash
curl --location 'http://localhost:8080/status'
```
## Docker
### Image
A docker image that runs this application is available under the repo packages, or at `ghcr.io/holysoles/ginrcon`.
### Environment Variables
- `PORT`: Optional, override the port (default 8080) of the webserver within the container
- `TRUSTED_PROXIES`: Optional, set specified trusted proxy addresses
- `RCON_SERVER`: Optional, configure a default RCON server's hostname
- `RCON_PORT`: Optional, configure a default RCON server's port
- `RCON_ADMIN_PASSWORD`: Optional, configure the password to a default RCON server
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