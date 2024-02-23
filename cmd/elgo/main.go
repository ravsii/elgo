package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/integrii/flaggy"
	"github.com/ravsii/elgo"
	"github.com/ravsii/elgo/grpc"
	"github.com/ravsii/elgo/socket"
)

const version = "v0.0.1"

// Common flags for both client and server
// It's not possible to use both sides in one instance,
// so it's safe to reuse them.
var (
	addr       string = ""
	port       uint16 = 8080
	network    string = "tcp"
	serverType string = "grpc"
)

// Server flags.
var (
	playerRetry       time.Duration = time.Second
	globalRetry       time.Duration = time.Second
	increaseBordersBy float64       = 100
)

// Client flags (TODO).
var ()

func main() {
	flaggy.SetName("elgo")
	flaggy.SetDescription("The matckmaker")
	flaggy.DefaultParser.ShowHelpOnUnexpected = true
	flaggy.DefaultParser.AdditionalHelpPrepend = "https://github.com/ravsii/elgo"

	serverCommand := flaggy.NewSubcommand("server")
	flaggy.AttachSubcommand(serverCommand, 1)
	serverCommand.Description = "Server-side related commands"

	serverCommand.String(&network, "n", "network",
		"type of network to listen (tcp|udp)")
	serverCommand.String(&addr, "a", "addr",
		"host/address to accept connections at")
	serverCommand.UInt16(&port, "p", "port",
		"port to accept connections at (1 to 65535)")
	serverCommand.String(&serverType, "t", "type",
		"type of server (grpc|socket)")
	serverCommand.Float64(&increaseBordersBy, "i", "increase-borders-by",
		"amount of ELO points to increase player's search range, if no match was found.")
	serverCommand.Duration(&globalRetry, "gr", "global-retry",
		"global retry interval, duration reference: https://pkg.go.dev/time#ParseDuration")
	serverCommand.Duration(&playerRetry, "pr", "player-retry",
		"player retry interval, duration reference: https://pkg.go.dev/time#ParseDuration")

	clientCommand := flaggy.NewSubcommand("client")
	clientCommand.Description = "NOT IMPLEMENTED (only usable as a lib now)"
	flaggy.AttachSubcommand(clientCommand, 1)

	flaggy.SetVersion(version)
	flaggy.Parse()

	switch {
	case serverCommand.Used:
		serverUsed()
	default:
		exitErr("specify side: server or client")
	}
}

func serverUsed() {
	if port == 0 {
		exitErr("port should be in range from 1 to 65535")
	}

	if network != "tcp" && network != "udp" {
		exitErr("invalid listener type: %s. (tcp|udp)", serverType)
	}

	if serverType != "grpc" && serverType != "socket" {
		exitErr("invalid server type: %s. (grpc|socket)", serverType)
	}

	pool := elgo.NewPool(
		elgo.WithPlayerRetryInterval(playerRetry),
		elgo.WithGlobalRetryInterval(globalRetry),
		elgo.WithIncreasePlayerBorderBy(increaseBordersBy),
	)

	go pool.Run()

	switch serverType {
	case "grpc":
		startGrpcServer(pool)
	case "socket":
		startSockerServer(pool)
	default:
		log.Fatalf("%s not supported", serverType)
	}
}

func startGrpcServer(pool *elgo.Pool) {
	srvAddr := fmt.Sprintf("%s:%d", addr, port)
	grpcPool := grpc.NewPoolServer(pool)
	server, err := grpc.NewListener(network, srvAddr, grpcPool)
	if err != nil {
		exitErr(err.Error())
	}

	// Add graceful shutdown listener
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		server.Close()

		log.Println("\nShutting down the server...")

		if players := pool.Close(); len(players) > 0 {
			plrs := make([]string, 0, len(players))
			for p := range players {
				plrs = append(plrs, p)
			}

			log.Println("Players were dropped from the queue:", strings.Join(plrs, ", "))
		}

		os.Exit(0)
	}()

	log.Println("Starting the server...")
	if err := server.Listen(); err != nil {
		log.Fatalf("listen: %s", err)
	}
}

func startSockerServer(pool *elgo.Pool) {
	server := socket.NewServer(pool)

	// Add graceful shutdown listener
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		err := server.Close()
		if err != nil {
			log.Fatalln("socket server close:", err)
		}

		log.Println("\nShutting down the server...")

		if players := pool.Close(); len(players) > 0 {
			plrs := make([]string, 0, len(players))
			for p := range players {
				plrs = append(plrs, p)
			}

			log.Println("Players were dropped from the queue:", strings.Join(plrs, ", "))
		}

		os.Exit(0)
	}()

	log.Println("Starting the server...")
	srvAddr := fmt.Sprintf("%s:%d", addr, port)
	if err := server.Listen(network, srvAddr); err != nil {
		log.Fatalf("listen: %s", err)
	}
}

// exitErr shows a given errors and exits with the code 2.
func exitErr(format string, args ...any) {
	flaggy.ShowHelpAndExit(fmt.Sprintf("err: "+format, args...))
	os.Exit(2)
}
