// Package main provides the entry point for the application
package main

import (
	"github.com/vvbbnn00/goflet/base"
	"github.com/vvbbnn00/goflet/config"
	"github.com/vvbbnn00/goflet/route"
	"github.com/vvbbnn00/goflet/task"
	"github.com/vvbbnn00/goflet/util/log"
)

func main() {
	base.PrintBanner()

	gofletCfg := config.GofletCfg

	httpConfig := gofletCfg.HTTPConfig
	router := route.RegisterRoutes()
	endpoint := gofletCfg.GetEndpoint()

	// Start the HTTP and HTTPS servers
	if httpConfig.HTTPSConfig.Enabled {
		go func() {
			err := router.RunTLS(endpoint, httpConfig.HTTPSConfig.Cert, httpConfig.HTTPSConfig.Key)
			if err != nil {
				panic(err)
			}
			log.Infof("HTTPS server started on %s", endpoint)
		}()
	} else {
		go func() {
			err := router.Run(endpoint)
			if err != nil {
				panic(err)
			}
			log.Infof("HTTP server started on %s", endpoint)
		}()
	}

	task.RunScheduledTask()

	// Wait for keyboard interrupt to stop the servers
	select {}
}
