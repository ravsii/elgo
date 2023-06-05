package main

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/integrii/flaggy"
	"github.com/ravsii/elgo"
	"github.com/ravsii/elgo/grpc"
)

const version = "v0.0.1"

// Common flags for both client and server
// It's not possible to use both sides in one instance,
// so it's safe to reuse them.
var (
	addr   string = ""
	port   uint16 = 8080
	listen string = "tcp"
	typ    string = "grpc"
)

// Server flags
var (
	playerRetry       time.Duration = time.Second
	globalRetry       time.Duration = time.Second
	increaseBordersBy float64       = 100
)

// Client flags (TODO):
var ()

func main() {
	flaggy.SetName("elgo")
	flaggy.SetDescription("The matckmaker")
	flaggy.DefaultParser.ShowHelpOnUnexpected = true
	flaggy.DefaultParser.AdditionalHelpPrepend = "https://github.com/ravsii/elgo"

	serverCommand := flaggy.NewSubcommand("server")
	flaggy.AttachSubcommand(serverCommand, 1)
	serverCommand.Description = "Server-side related commands"

	serverCommand.String(&listen, "l", "listen", "type of network to listen (tcp|udp)")
	serverCommand.String(&addr, "a", "addr", "host/address to accept connections at")
	serverCommand.UInt16(&port, "p", "port", "port to accept connections at (1 to 65535)")
	serverCommand.String(&typ, "t", "type", "type of server (grpc|socket)")
	serverCommand.Float64(&increaseBordersBy, "i", "increase-borders", "amount of ELO points to increase player's search range, if no match was found.")
	serverCommand.Duration(&globalRetry, "gr", "global-retry", "global retry interval, duration reference: https://pkg.go.dev/time#ParseDuration")
	serverCommand.Duration(&playerRetry, "pr", "player-retry", "player retry interval, duration reference: https://pkg.go.dev/time#ParseDuration")

	clientCommand := flaggy.NewSubcommand("client")
	clientCommand.Description = "NOT IMPLEMENTED (only usable as a lib now)"
	flaggy.AttachSubcommand(clientCommand, 1)

	flaggy.SetVersion(version)
	flaggy.Parse()

	switch {
	case serverCommand.Used:
		handleStartServer()
	default:
		exitErr("specify side: server or client")
	}
}

func handleStartServer() {
	if port == 0 {
		exitErr("port should be in range from 1 to 65535")
	}

	if typ != "grpc" && typ != "socket" {
		exitErr("invalid server type: %s. (grpc|socket)", typ)
	}

	if listen != "tcp" && listen != "udp" {
		exitErr("invalid listener type: %s. (tcp|udp)", typ)
	}

	poolServer := grpc.NewPoolServer(
		elgo.WithPlayerRetryInterval(playerRetry),
		elgo.WithGlobalRetryInterval(globalRetry),
		elgo.WithIncreasePlayerBorderBy(increaseBordersBy),
	)

	addr := fmt.Sprintf("%s:%d", addr, port)
	server, err := grpc.NewListener(listen, addr, poolServer)
	if err != nil {
		exitErr(err.Error())
	}

	// Add graceful shutdown listener
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, os.Kill, syscall.SIGTERM, syscall.SIGABRT)
	go func() {
		<-c
		players, err := server.Close()
		if err != nil {
			exitErr(err.Error())
		}

		fmt.Println("\nShutting down the server...")

		if len(players) > 0 {
			plrs := make([]string, 0, len(players))
			for p := range players {
				plrs = append(plrs, p)
			}

			fmt.Println("Players were dropped from the queue:", strings.Join(plrs, ", "))
		}

		os.Exit(0)
	}()

	fmt.Println("Starting the server...")
	exitErr(server.Listen().Error())
}

// exitErr shows a given errors and exits with the code 2.
func exitErr(format string, args ...any) {
	flaggy.ShowHelpAndExit(fmt.Sprintf("err: "+format, args...))
	os.Exit(2)
}
