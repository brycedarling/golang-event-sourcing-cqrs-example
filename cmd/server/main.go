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
var webapiFlag = flag.Bool("webapi", true, "start web api")
var grpcapiFlag = flag.Bool("grpcapi", true, "start grpc api")

func main() {
	flag.Parse()

	env, err := config.InitializeEnv()
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to initialize env: %s\n", err)
		os.Exit(1)
	}

	conf, configCleanup, err := config.InitializeConfig(env)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to initialize config: %s\n", err)
		os.Exit(1)
	}
	defer configCleanup()

	a := application.Aggregators{}
	if *identityAggregatorFlag {
		a = append(a, identity.InitializeAggregator(conf))
	}
	if *viewingAggregatorFlag {
		a = append(a, viewing.InitializeAggregator(conf))
	}
	a.Start()
	defer a.Stop()

	c := application.Components{}
	if *identityComponentFlag {
		c = append(c, identity.InitializeComponent(conf))
	}
	c.Start()
	defer c.Stop()

	var apiShutdown func()

	if *webapiFlag {
		go func() {
			var api web.API
			api, apiShutdown, err = web.InitializeAPI(conf)
			if err != nil {
				fmt.Fprintf(os.Stderr, "failed to initialize web api: %s\n", err)
				os.Exit(1)
			}
			api.Listen()
		}()
	}

	if *grpcapiFlag {
		go func() {
			server, err := rpc.InitializeServer(conf)
			if err != nil {
				fmt.Fprintf(os.Stderr, "failed to initialize grpc api: %s\n", err)
				os.Exit(1)
			}
			server.Listen()
		}()
	}

	// Wait for ctrl-c to exit
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt)

	// Block until a signal is received
	<-ch

	if apiShutdown != nil {
		apiShutdown()
	}
}
