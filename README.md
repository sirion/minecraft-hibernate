# minecraft-hibernate

Autostart a minecraft server and stop when noone is playing

Start using `mchibernate /path/to/config.json`.

**!!! Important: Do not use the same port for the server and in the config file if you are using proxy mode !!!**

## Modes

There are two basic modes: Proxy and Replacement

Replacement is preferable in most situations, as there is no overhead from having to copy data between two ports. The main disadvantage is that there will be a few seconds between starting/stopping the server and listening on the port, where all connections will be rejected and the server seems to be offline from the clients perspective.

### Replacement

In replacement mode (configuration `proxy: false`), minecraft-hibernate listens on the same port as the minecraft server and starts it automatically.

### Proxy

In Proxy mode (configuration `proxy: true`), the minecraft-hibernate listens on its own port and forwards all connections to the real minecraft-server.

## Configuration

An example configuration file can be found at [example/config.json](example/config.json).

A configuration file may have the following properties:

|Property|default||
|-|-|-|
| proxy | false | Whether to listen on a different port than the real minecraft server and act as a connection proxy |
| startOnPing | false | Starts the server whenever the server is pinged (when any client asks for the status in the list of multiplayer games) |
| memoryMin | 2G | The minimum memory for the minecraft server ("-Xms") |
| memoryMax | 4G | The maximum memory for the minecraft server ("-Xmx") |
| port | 25566 | ***in Proxy mode only*** The port to listen on (important: this cannot be sa same port as the one configured for the real minecraft server in the server.properties file) |
| idleTimeout | 180 | How long to wait for players to join the server again before shutting it down |
| checkInterval | 30 | How often to check the number of active players |
| workingDir | "." | The server directory |
| jarPath | "minecraft-server.jar" | The server java file |
| favicon |  | base64 version of a 64x64 png to show while the real server is not running, prefixed with "data:image/png;base64," |
