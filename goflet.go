// Package main provides the entry point for the application
package main

import (
	"github.com/vvbbnn00/goflet/base"
	"github.com/vvbbnn00/goflet/config"
	"github.com/vvbbnn00/goflet/route"
	"github.com/vvbbnn00/goflet/task"
	"github.com/vvbbnn00/goflet/util/log"
)

// @title           Goflet API
// @version         unknown
// @description     Goflet is a lightweight file upload and download service written in Go.

// @contact.name   vvbbnn00
// @contact.url    https://github.com/vvbbnn00/goflet
// @contact.email  vvbbnn00@gmail.com

// @license.name  MIT
// @license.url   https://opensource.org/licenses/MIT

// @BasePath  /

// @securityDefinitions.apikey Authorization
// @in header
// @name Authorization
// @description You need to provide a valid jwt token in the header, in headers, you should provide a key-value pair like this: Authorization: Bearer xxxxxx; The token has the same effect as the token in the query string, but it is more secure than the token in the query string. Or you can just provide the token in the query string, like this: ?token=xxxxxx. More info about jwt: https://github.com/vvbbnn00/goflet?tab=readme-ov-file#authentication-method

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/
func main() {
	base.PrintBanner()

	gofletCfg := config.GofletCfg

	httpConfig := gofletCfg.HTTPConfig
	router := route.RegisterRoutes()
	endpoint := gofletCfg.GetEndpoint()

	// Start the HTTP and HTTPS servers
	if *httpConfig.HTTPSConfig.Enabled {
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
