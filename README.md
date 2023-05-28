# minecraft-hibernate

Autostart a minecraft server and stop when noone is playing


Start using `mchibernate /path/to/config.json`.

**!!! Important: Do not use the same port for the server and in the config file !!!**


An example configuration file can be found at [example/config.json](example/config.json).

A configuration file may have the following properties:

|Property|default||
|-|-|-|
| startOnPing | false | Starts the server whenever the server is pinged (any client asks for the status in the list of multiplayer games) |
| memoryMin | 2G | The minimum memory for the minecraft server ("-Xms") |
| memoryMax | 4G | The maximum memory for the minecraft server ("-Xmx") |
| port | 25566 | The port to listen on (important: this ca not be sa same port as the one configured for the real minevraft server in the server.properties file) |
| idleTimeout | 180 | How long to wait for players to join the server again before shutting it down |
| checkInterval | 30 | How often to check the number of active players |
| workingDir | "." | The server directory |
| jarPath | "minecraft-server.jar" | The server java file |
| favicon |  | base64 version of a 64x64 png to show while the real server is not running, prefixed with "data:image/png;base64," |
