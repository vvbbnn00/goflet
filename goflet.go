package main

import (
	"goflet/config"
	"goflet/route"
	"log"
)

func main() {
	config.InitConfig()

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
			log.Printf("HTTPS server started on %s", endpoint)
		}()
	} else {
		go func() {
			err := router.Run(endpoint)
			if err != nil {
				panic(err)
			}
			log.Printf("HTTP server started on %s", endpoint)
		}()
	}

	// Wait for keyboard interrupt to stop the servers
	select {}
}
