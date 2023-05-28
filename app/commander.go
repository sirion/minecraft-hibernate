package main

import (
	"bytes"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type ActivePlayers struct {
	Number int
	Names  []string
}

type Commander struct {
	started          bool
	startedCallbacks []func()
	buffer           []byte
	//	messages []string
	commands       []string
	cActivePlayers chan ActivePlayers
	version        string

	cServerPort chan int
	serverPort  int
}

func NewCommander() *Commander {
	return &Commander{
		//		messages: make([]string, 0, 100),
		started:          false,
		startedCallbacks: make([]func(), 0, 1),
		commands:         make([]string, 0, 100),
		buffer:           make([]byte, 0, 10*1024),
		cActivePlayers:   make(chan ActivePlayers, 1),
		cServerPort:      make(chan int, 1),
		version:          "1.19.4",
	}
}

func (c *Commander) ServerPort() int {
	if c.serverPort == 0 {
		c.serverPort = <-c.cServerPort
	}
	return c.serverPort
}

func (c *Commander) Write(p []byte) (n int, err error) {
	c.buffer = append(c.buffer, p...)
	c.parseBuffer()
	return len(p), nil
}

func (c *Commander) parseBuffer() {
	pos := bytes.IndexByte(c.buffer, '\n')
	for pos != -1 {
		c.parseMessage(c.buffer[0:pos])
		c.buffer = c.buffer[pos+1:]
		pos = bytes.IndexByte(c.buffer, '\n')
	}
}

var regNumPlayers = regexp.MustCompile(`^\[.*\] \[.*\]: There are ([0-9]*) of a max of [0-9] players online: (.*)$`)
var regStarted = regexp.MustCompile(`^\[.*\] \[.*\]: Done \(([0-9]*\.[0-9]*)s\)!`)
var regPort = regexp.MustCompile(`^\[.*\] \[.*\]: Starting Minecraft server on .*:([0-9]*)`)
var regVersion = regexp.MustCompile(`^\[.*\] \[.*\]: Starting minecraft server version ([0-9\.]*)`)

func (c *Commander) parseMessage(msg []byte) {
	// fmt.Fprint(os.Stdout, string(msg), "\n")

	// [13:33:04] [Server thread/INFO]: There are 1 of a max of 4 players online: xxx
	matches := regNumPlayers.FindStringSubmatch(string(msg))
	if len(matches) > 0 {
		// Num players response
		num, err := strconv.ParseInt(matches[1], 10, 32)
		if err != nil {
			fmt.Fprintf(os.Stderr, "SERVERLOG: %s", msg)
			fmt.Fprintf(os.Stderr, "Error parsing number of players: %s", err.Error())
			return
		}

		np := ActivePlayers{
			Number: int(num),
			Names:  strings.Split(string(matches[2]), " "),
		}
		c.cActivePlayers <- np
		return
	}

	// [18:38:50] [Server thread/INFO]: Done (15.852s)! For help, type "help"
	if regStarted.Match(msg) {
		c.setStarted()
	}

	// [18:38:34] [Server thread/INFO]: Starting Minecraft server on *:25565
	matches = regPort.FindStringSubmatch(string(msg))
	if len(matches) > 0 {
		port, err := strconv.ParseInt(matches[1], 10, 32)
		if err != nil {
			fmt.Fprintf(os.Stderr, "SERVERLOG: %s\n", msg)
			fmt.Fprintf(os.Stderr, "Error parsing port: %s\n", err.Error())
			return
		}

		c.cServerPort <- int(port)
	}

	// [18:38:34] [Server thread/INFO]: Starting minecraft server version 1.19.4
	matches = regVersion.FindStringSubmatch(string(msg))
	if len(matches) > 0 {
		c.version = matches[1]
	}

	// Starting net.minecraft.server.Main
	// [18:38:31] [ServerMain/INFO]: Environment: authHost='https://authserver.mojang.com', accountsHost='https://api.mojang.com', sessionHost='https://sessionserver.mojang.com', servicesHost='https://api.minecraftservices.com', name='PROD'
	// [18:38:33] [ServerMain/INFO]: Loaded 7 recipes
	// [18:38:33] [ServerMain/INFO]: Loaded 1179 advancements
	// [18:38:34] [Server thread/INFO]: Starting minecraft server version 1.19.4
	// [18:38:34] [Server thread/INFO]: Loading properties
	// [18:38:34] [Server thread/INFO]: Default game type: SURVIVAL
	// [18:38:34] [Server thread/INFO]: Generating keypair
	// [18:38:34] [Server thread/INFO]: Starting Minecraft server on *:25565
	// [18:38:34] [Server thread/INFO]: Using epoll channel type
	// [18:38:34] [Server thread/INFO]: Preparing level "pakora"
	// [18:38:42] [Server thread/INFO]: Preparing start region for dimension minecraft:overworld
	// [18:38:48] [Worker-Main-1/INFO]: Preparing spawn area: 0%
	// [18:38:48] [Worker-Main-1/INFO]: Preparing spawn area: 0%
	// [18:38:48] [Worker-Main-1/INFO]: Preparing spawn area: 0%
	// [18:38:48] [Worker-Main-1/INFO]: Preparing spawn area: 0%
	// [18:38:48] [Worker-Main-1/INFO]: Preparing spawn area: 0%
	// [18:38:48] [Worker-Main-1/INFO]: Preparing spawn area: 0%
	// [18:38:48] [Worker-Main-1/INFO]: Preparing spawn area: 0%
	// [18:38:48] [Worker-Main-1/INFO]: Preparing spawn area: 0%
	// [18:38:48] [Worker-Main-1/INFO]: Preparing spawn area: 0%
	// [18:38:48] [Worker-Main-1/INFO]: Preparing spawn area: 0%
	// [18:38:48] [Worker-Main-1/INFO]: Preparing spawn area: 0%
	// [18:38:48] [Worker-Main-1/INFO]: Preparing spawn area: 0%
	// [18:38:48] [Worker-Main-1/INFO]: Preparing spawn area: 0%
	// [18:38:49] [Worker-Main-1/INFO]: Preparing spawn area: 0%
	// [18:38:49] [Worker-Main-1/INFO]: Preparing spawn area: 2%
	// [18:38:50] [Server thread/INFO]: Time elapsed: 7704 ms
	// [18:38:50] [Server thread/INFO]: Done (15.852s)! For help, type "help"

	fmt.Fprintf(os.Stdout, "SERVERLOG: %s\n", msg)
}

func (c *Commander) setStarted() {
	c.started = true
	for _, c := range c.startedCallbacks {
		c()
	}
}

func (c *Commander) onStarted(callback func()) {
	if c.started {
		callback()
	} else {
		c.startedCallbacks = append(c.startedCallbacks, callback)
	}
}

func (c *Commander) ActivePlayers() ActivePlayers {
	c.onStarted(func() {
		c.List()
	})
	return <-c.cActivePlayers
}

func (c *Commander) Read(p []byte) (n int, err error) {
	if len(c.commands) == 0 {
		time.Sleep(500 * time.Millisecond)
		return 0, nil
	}

	if len(c.commands[0]) < len(p)-1 {
		// Command fits in buffer, send full command

		for i, b := range c.commands[0] {
			p[i] = byte(b)
		}

		length := len(c.commands[0])
		p[length] = '\n'
		c.commands = c.commands[1:]
		return length + 1, nil

	} else {
		// Part of command
		for i := 0; i < len(p); i++ {
			p[i] = byte(c.commands[0][i])
		}
		c.commands[0] = c.commands[0][len(p):]
		return len(p), nil
	}
}

func (c *Commander) List() {
	c.commands = append(c.commands, "list")
}

func (c *Commander) Stop() {
	c.commands = append(c.commands, "stop")
}

/**
 COMMANDS
help

/experience (add|set|query)
/xp -> experience
/fill <from> <to> <block> [replace|keep|outline|hollow|destroy]
/fillbiome <from> <to> <biome> [replace]
/forceload (add|remove|query)
/function <name>
/gamemode <gamemode> [<target>]
/gamerule (announceAdvancements|blockExplosionDropDecay|commandBlockOutput|commandModificationBlockLimit|disableElytraMovementCheck|disableRaids|doDaylightCycle|doEntityDrops|doFireTick|doImmediateRespawn|doInsomnia|doLimitedCrafting|doMobLoot|doMobSpawning|doPatrolSpawning|doTileDrops|doTraderSpawning|doVinesSpread|doWardenSpawning|doWeatherCycle|drowningDamage|fallDamage|fireDamage|forgiveDeadPlayers|freezeDamage|globalSoundEvents|keepInventory|lavaSourceConversion|logAdminCommands|maxCommandChainLength|maxEntityCramming|mobExplosionDropDecay|mobGriefing|naturalRegeneration|playersSleepingPercentage|randomTickSpeed|reducedDebugInfo|sendCommandFeedback|showDeathMessages|snowAccumulationHeight|spawnRadius|spectatorsGenerateChunks|tntExplosionDropDecay|universalAnger|waterSourceConversion)
/give <targets> <item> [<count>]
/help [<command>]
/item (replace|modify)
/kick <targets> [<reason>]
/kill [<targets>]
/list [uuids]
/locate (structure|biome|poi)
/loot (replace|insert|give|spawn)
/msg <targets> <message>
/tell -> msg
/w -> msg
/particle <name> [<pos>]
/place (feature|jigsaw|structure|template)
/playsound <sound> (master|music|record|weather|block|hostile|neutral|player|ambient|voice)
/reload
/recipe (give|take)
/ride <target> (mount|dismount)
/say <message>
/schedule (function|clear)
/scoreboard (objectives|players)
/seed
/setblock <pos> <block> [destroy|keep|replace]
/spawnpoint [<targets>]
/setworldspawn [<pos>]
/spectate [<target>]
/spreadplayers <center> <spreadDistance> <maxRange> (<respectTeams>|under)
/stopsound <targets> [*|master|music|record|weather|block|hostile|neutral|player|ambient|voice]
/summon <entity> [<pos>]
/tag <targets> (add|remove|list)
/team (list|add|remove|empty|join|leave|modify)
/teammsg <message>
/tm -> teammsg
/teleport (<location>|<destination>|<targets>)
/tp -> teleport
/tellraw <targets> <message>
/time (set|add|query)
/title <targets> (clear|reset|title|subtitle|actionbar|times)
/trigger <objective> [add|set]
/weather (clear|rain|thunder)
/worldborder (add|set|center|damage|get|warning)
/jfr (start|stop)
/ban-ip <target> [<reason>]
/banlist [ips|players]
/ban <targets> [<reason>]
/deop <targets>
/op <targets>
/pardon <targets>
/pardon-ip <target>
/perf (start|stop)
/save-all [flush]
/save-off
/save-on
/setidletimeout <minutes>
/stop
/whitelist (on|off|list|add|remove|reload)
**/
