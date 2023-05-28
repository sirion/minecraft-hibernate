package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"minecraftmanager/lib/mc"
	"net"
	"os"
	"os/exec"
	"time"
)

const DEFAULT_FAVICON = "data:image/png;base64,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"

var serverFlags = []string{
	// Optimizations
	// "-XX:+UseG1GC",
	// "-XX:MaxGCPauseMillis=130",
	// "-XX:+UnlockExperimentalVMOptions",
	// "-XX:+DisableExplicitGC",
	// "-XX:+AlwaysPreTouch",
	// "-XX:G1NewSizePercent=28",
	// "-XX:G1HeapRegionSize=16M",
	// "-XX:G1ReservePercent=20",
	// "-XX:G1MixedGCCountTarget=3",
	// "-XX:InitiatingHeapOccupancyPercent=10",
	// "-XX:G1MixedGCLiveThresholdPercent=90",
	// "-XX:G1RSetUpdatingPauseTimePercent=0",
	// "-XX:SurvivorRatio=32",
	// "-XX:MaxTenuringThreshold=1",
	// "-XX:G1SATBBufferEnqueueingThresholdPercent=30",
	// "-XX:G1ConcMarkStepDurationMillis=5",
	// "-XX:G1ConcRSHotCardLimit=16",
	// "-XX:G1ConcRefinementServiceIntervalMillis=150",
	// Memory allocation
}

type ProcessStatus uint8

const (
	PROC_STAT_STOPPED ProcessStatus = iota
	PROC_STAT_RUNNING
	PROC_STAT_STARTING
	PROC_STAT_STOPPING
	PROC_STAT_ERROR
)

type Server struct {
	Active            bool
	Process           *exec.Cmd
	ProcessStatus     ProcessStatus
	Commander         *Commander
	Version           string
	ProtocolVersion   int
	IdleSince         time.Time
	config            serverConfig
	activeConnections int
}

type serverConfig struct {
	StartOnPing   bool   `json:"startOnPing"`
	Port          int    `json:"port"`
	IdleTimeout   int    `json:"idleTimeout"`
	CheckInterval int    `json:"checkInterval"`
	WorkingDir    string `json:"workingDir"`
	JarPath       string `json:"jarPath"`
	Favicon       string `json:"favicon"`
	MemoryMin     string `json:"memoryMin"`
	MemoryMax     string `json:"memoryMax"`
}

func (srv *Server) Listen() {
	srv.Commander = NewCommander()

	listener, err := net.Listen("tcp", fmt.Sprintf(":%d", srv.config.Port))
	if err != nil {
		panic(err)
	}

	done := make(chan bool)

	go srv.check()

	go func() {
		for srv.Active {
			conn, err := listener.Accept()
			if err != nil {
				fmt.Fprintf(os.Stderr, "[ERROR] %s\n", err.Error())
				return
			}

			var pkg *mc.Package

			if srv.ProcessStatus != PROC_STAT_RUNNING {
				connReader := bufio.NewReader(conn)
				pkg = readPackage(connReader)
				hs, err := pkg.Handshake()
				if err != nil {
					// No Handshake
					fmt.Fprintf(os.Stderr, "Wrong handshake from client. PKG: %#v", pkg)
					conn.Close()
					continue
				}

				// hs.ServerAddress and hs.ServerPort could be used for routing in the future

				// TODO: This should be taken from the server when started and then saved for next time
				srv.ProtocolVersion = hs.ProtocolVersion
				srv.Version = "1.19.4"

				if hs.NextState == 1 {
					// Requests Status
					go srv.answerSleeping(connReader, conn)
					if !srv.config.StartOnPing {
						continue
					} else {
						pkg = nil // Don send handshake to started server
					}
				} else if hs.NextState != 2 {
					fmt.Fprintf(os.Stderr, "Wrong state from client: %d --> PKG: %#v", hs.NextState, pkg)
					conn.Close()
					continue
				}

			}

			if srv.ProcessStatus == PROC_STAT_STOPPING {
				srv.Process.Wait()
				srv.Process = nil
				srv.ProcessStatus = PROC_STAT_STOPPED
			}

			if srv.ProcessStatus == PROC_STAT_STOPPED {
				srv.start()
			}

			// if srv.ProcessStatus == PROC_STAT_STARTING {
			// 	// TODO: Wait for server to listen
			// 	time.Sleep(1 * time.Second)
			// }

			if srv.ProcessStatus == PROC_STAT_ERROR {
				fmt.Fprintf(os.Stderr, "Server process error: %#v", srv.Process.ProcessState.ExitCode())
				panic("PROC ERROR")
			}

			go func() {
				serverconn, err := net.Dial("tcp", fmt.Sprintf("localhost:%d", srv.Commander.ServerPort()))
				if err != nil {
					fmt.Fprintf(os.Stderr, "[ERROR] %s\n", err.Error())
					conn.Close()
					return
				}

				srv.WaitUntilStarted()
				if pkg != nil {
					pkg.WriteTo(serverconn)
				}

				srv.activeConnections += 1
				go io.Copy(conn, serverconn)
				io.Copy(serverconn, conn)
				srv.activeConnections -= 1
				conn.Close()

			}()
		}

		srv.Commander.Stop()
		srv.Process.Wait()
		done <- true
	}()

	//<-time.After(300 * time.Second)
	//srv.Active = false
	<-done
}

func ServerFromConfig(path string) (*Server, error) {
	content, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	conf := serverConfig{}
	err = json.Unmarshal(content, &conf)
	if err != nil {
		return nil, err
	}

	if conf.MemoryMin == "" {
		conf.MemoryMin = "2G"
	}
	if conf.MemoryMax == "" {
		conf.MemoryMax = "4G"
	}
	if conf.Port == 0 {
		conf.Port = 25566
	}
	if conf.IdleTimeout == 0 {
		conf.IdleTimeout = 180
	}
	if conf.CheckInterval == 0 {
		conf.CheckInterval = 30
	}
	if conf.JarPath == "" {
		conf.JarPath = "minecraft-server.jar"
	}

	srv := &Server{
		Active:    true,
		Process:   nil,
		IdleSince: time.Now().Add(time.Hour),
		config:    conf,
	}

	return srv, nil
}

/// check runs in the background
func (srv *Server) check() {

	go func() {
		for srv.Active {
			if srv.Process == nil {
				srv.ProcessStatus = PROC_STAT_STOPPED
			} else if srv.Process.ProcessState != nil {
				srv.ProcessStatus = PROC_STAT_STOPPED
			}

			time.Sleep(5 * time.Second)
		}
	}()

	for srv.Active {

		// Check if players are active
		if srv.activeConnections == 0 && srv.Commander.ActivePlayers().Number == 0 {
			if srv.IdleSince.Add(time.Duration(srv.config.IdleTimeout) * time.Second).Before(time.Now()) {
				// Shutdown server
				srv.Commander.Stop()
				srv.ProcessStatus = PROC_STAT_STOPPING
				srv.Process.Wait()
				srv.Process = nil
				srv.ProcessStatus = PROC_STAT_STOPPED
				srv.IdleSince = time.Now().Add(99999 * time.Hour)

			} else if srv.IdleSince.After(time.Now()) {
				srv.IdleSince = time.Now()
			}
		} else {
			srv.IdleSince = time.Now().Add(99999 * time.Hour)
		}

		time.Sleep(time.Duration(srv.config.IdleTimeout) * time.Second)
	}
}

func (srv *Server) start() {
	if srv.Process == nil {
		// Start for the first time
		srv.ProcessStatus = PROC_STAT_STOPPED
		cmdFlags := append(serverFlags,
			fmt.Sprintf("-Xms%s", srv.config.MemoryMin),
			fmt.Sprintf("-Xmx%s", srv.config.MemoryMax),
			"-jar", srv.config.JarPath,
			"nogui",
		)
		srv.Process = exec.Command("java", cmdFlags...)
		srv.Process.Dir = srv.config.WorkingDir
		srv.Process.Stdin = srv.Commander
		srv.Process.Stdout = srv.Commander
	}

	err := srv.Process.Start()
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %v", err)
		srv.ProcessStatus = PROC_STAT_ERROR
	} else {
		srv.ProcessStatus = PROC_STAT_STARTING
	}

	srv.Commander.onStarted(func() {
		srv.ProcessStatus = PROC_STAT_RUNNING
	})
}

func (srv *Server) WaitUntilStarted() {
	started := make(chan bool)
	go srv.Commander.onStarted(func() {
		started <- true
	})

	<-started
}

func (srv *Server) answerSleeping(r mc.Reader, w io.Writer) {
	// This should be either Status or Ping Package
	pkg := readPackage(r)
	if pkg.ID == 0 {
		// Status request
		srv.answerStatus(w)
		pkg = readPackage(r)
	}

	if pkg.ID == 1 {
		if DEBUG {
			fmt.Printf("[PING] %#v\n", pkg)
		}
		srv.answerPong(w, pkg)
	}

}

func (srv *Server) answerPong(w io.Writer, pkg *mc.Package) {
	_, err := pkg.WriteTo(w)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[ERROR] %s\n", err.Error())
		return
	}
	if DEBUG {
		fmt.Printf("[PONG] %#v\n", pkg)
	}
}

func (srv *Server) answerStatus(w io.Writer) {

	description := "Server is sleeping. Connect to start. Retry after 10s if connection fails."
	if srv.config.StartOnPing {
		description = "Server is starting. Please wait and refresh."
	}

	favicon := DEFAULT_FAVICON
	if srv.config.Favicon != "" {
		favicon = srv.config.Favicon
	}

	resp := mc.StatusResponse{
		Version: mc.Version{
			Name:     srv.Version,
			Protocol: srv.ProtocolVersion,
		},
		Players: mc.Players{
			Max:    0,
			Online: 0,
		},
		Description: mc.Description{
			Text: description,
		},
		Favicon:            favicon,
		EnforcesSecureChat: true,
	}

	written, err := resp.AsPackage().WriteTo(w)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[ERROR] %s\n", err.Error())
		return
	}

	if DEBUG {
		fmt.Printf("[OUT] %#v\n", written)
	}
}

// Helper

func readPackage(r mc.Reader) *mc.Package {
	pkg := mc.NewPackage(0)

	err := pkg.Read(r)
	if err != nil {
		fmt.Fprintf(os.Stderr, "[ERROR] %s\n", err.Error())
	}
	if DEBUG {
		fmt.Printf("[IN] %#v\n", pkg)
	}

	return pkg
}
