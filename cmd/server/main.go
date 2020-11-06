package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"

	"github.com/brycedarling/go-practical-microservices/internal/application"
	"github.com/brycedarling/go-practical-microservices/internal/application/identity"
	"github.com/brycedarling/go-practical-microservices/internal/application/viewing"
	"github.com/brycedarling/go-practical-microservices/internal/infrastructure/config"
	"github.com/brycedarling/go-practical-microservices/internal/presentation/rpc"
	"github.com/brycedarling/go-practical-microservices/internal/presentation/web"
)

var identityAggregatorFlag = flag.Bool("ia", true, "start identity aggregator")
var viewingAggregatorFlag = flag.Bool("va", true, "start viewing aggregator")
var identityComponentFlag = flag.Bool("ic", true, "start identity component")
var webFlag = flag.Bool("web", true, "start web api")
var grpcFlag = flag.Bool("grpc", true, "start grpc server")

func main() {
	flag.Parse()

	env := initializeEnv()

	conf, configCleanup := initializeConfig(env)
	defer configCleanup()

	a := initializeAggregators(conf)
	a.Start()
	defer a.Stop()

	c := initializeComponents(conf)
	c.Start()
	defer c.Stop()

	webapiShutdown := initializeWebAPI(conf)

	grpcServerShutdown := initializeGRPCServer(conf)

	// Wait for ctrl-c to exit
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)

	// Block until a signal is received
	<-ch

	if webapiShutdown != nil {
		webapiShutdown()
	}
	if grpcServerShutdown != nil {
		grpcServerShutdown()
	}
}

func initializeEnv() *config.Env {
	env, err := config.InitializeEnv()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to initialize env: %s\n", err)
		os.Exit(1)
	}
	return env
}

func initializeConfig(env *config.Env) (*config.Config, func()) {
	conf, configCleanup, err := config.InitializeConfig(env)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to initialize config: %s\n", err)
		os.Exit(1)
	}
	return conf, configCleanup
}

func initializeAggregators(conf *config.Config) application.Aggregators {
	a := application.Aggregators{}
	if *identityAggregatorFlag {
		a = append(a, identity.InitializeAggregator(conf))
	}
	if *viewingAggregatorFlag {
		a = append(a, viewing.InitializeAggregator(conf))
	}
	return a
}

func initializeComponents(conf *config.Config) application.Components {
	c := application.Components{}
	if *identityComponentFlag {
		c = append(c, identity.InitializeComponent(conf))
	}
	return c
}

func initializeWebAPI(conf *config.Config) func() {
	var webapiShutdown func()

	if *webFlag {
		go func() {
			var api web.API
			var err error
			api, webapiShutdown, err = web.InitializeAPI(conf)
			if err != nil {
				fmt.Fprintf(os.Stderr, "failed to initialize web api: %s\n", err)
				os.Exit(1)
			}
			api.Listen()
		}()
	}

	return webapiShutdown
}

func initializeGRPCServer(conf *config.Config) func() {
	var grpcServerShutdown func()

	if *grpcFlag {
		go func() {
			var server *rpc.Server
			var err error
			server, grpcServerShutdown, err = rpc.InitializeServer(conf)
			if err != nil {
				fmt.Fprintf(os.Stderr, "failed to initialize grpc server: %s\n", err)
				os.Exit(1)
			}
			server.Listen()
		}()
	}

	return grpcServerShutdown
}
