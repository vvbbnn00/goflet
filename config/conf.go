// Package config provides the configuration for the application
package config

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"github.com/vvbbnn00/goflet/util/confutil"
)

// PathOfConfig is the path of the configuration file
const PathOfConfig = "goflet.json"

//go:embed goflet.json
var defaultConfig string

// CacheType is the type of the cache
type CacheType string

const (
	// CacheTypeMemory is the memory cache type
	CacheTypeMemory CacheType = "MemoryCache"
)

// GofletConfig contains the configuration for the application
type GofletConfig struct {
	Debug          bool `json:"debug" default:"false"`         // Enable debug mode
	SwaggerEnabled bool `json:"swaggerEnabled" default:"true"` // Enable swagger doc
	LogConfig      struct {
		// Log configuration
		Enabled bool   `json:"enabled" default:"true"` // Enable log
		Level   string `json:"level" default:"info"`   // The log level
	} `json:"logConfig"`
	HTTPConfig struct {
		Host string `json:"host" default:"0.0.0.0"` // The host to bind the server
		Port int    `json:"port" default:"8080"`    // The port to bind the server
		Cors struct {
			// CORS configuration
			Enabled bool     `json:"enabled" default:"true"`                     // Enable CORS
			Origins []string `json:"origins" default:"*"`                        // The list of allowed origins
			Methods []string `json:"methods" default:"GET,POST,PUT,HEAD,DELETE"` // The list of allowed methods
			Headers []string `json:"headers"`                                    // The list of allowed headers
		} `json:"cors"`
		ClientCache struct {
			// Client cache configuration
			Enabled bool `json:"enabled" default:"true"` // Enable client cache
			MaxAge  int  `json:"maxAge" default:"3600"`  // The maximum age of the client cache
		} `json:"clientCache"`
		HTTPSConfig struct {
			// HTTPS configuration
			Enabled bool   `json:"enabled" default:"false"` // Enable HTTPS
			Cert    string `json:"cert"`                    // The certificate file
			Key     string `json:"key"`                     // The key file
		} `json:"httpsConfig"`
	} `json:"httpConfig"`
	FileConfig struct {
		// File configuration
		BaseFileStoragePath string `json:"baseFileStoragePath" default:"data"` // The base path where the files will be stored
		UploadPath          string `json:"uploadPath" default:"upload"`        // The path where the files will be temporarily stored before moving to the base path
		AllowFolderCreation bool   `json:"allowFolderCreation" default:"true"` // Allow the creation of folders, otherwise the files will be stored in the base path
		UploadLimit         int64  `json:"uploadLimit" default:"1073741824"`   // The maximum size of the file to be uploaded
		UploadTimeout       int    `json:"uploadTimeout" default:"7200"`       // The maximum time to wait for the file to be uploaded
		MaxPostSize         int64  `json:"maxPostSize" default:"20971520"`     // The maximum size of the post request
	} `json:"fileConfig"`
	CacheConfig struct {
		// Cache configuration
		CacheType CacheType `json:"cacheType" default:"1"` // The cache type to be used
		// Cache configuration for memory
		MemoryCache struct {
			MaxEntries int `json:"maxEntries" default:"100"` // The maximum number of entries to be stored in the cache
			DefaultTTL int `json:"defaultTTL" default:"60"`  // The default time to live for the cache
		}
	} `json:"cacheConfig"`
	ImageConfig struct {
		// Image configuration
		DefaultFormat  string   `json:"defaultFormat" default:"png"` // The default format for the image
		AllowedFormats []string `json:"allowedFormats"`              // The list of allowed formats for the image

		StrictMode   bool  `json:"strictMode" default:"true"` // If true, the image size will only accept the allowed sizes
		AllowedSizes []int `json:"allowedSizes"`              // The list of allowed sizes for the image, like 32, 64, 128, 256

		MaxWidth    int   `json:"maxWidth" default:"4096"`        // The maximum width of the image
		MaxHeight   int   `json:"maxHeight" default:"4096"`       // The maximum height of the image
		MaxFileSize int64 `json:"maxFileSize" default:"20971520"` // The maximum size of the image file
	} `json:"imageConfig"`
	JWTConfig struct {
		// JWT configuration
		Enabled   bool   `json:"enabled" default:"true"`    // Enable JWT
		Algorithm string `json:"algorithm" default:"HS256"` // The algorithm to be used for the JWT
		Security  struct {
			// Security configuration
			SigningKey string `json:"signingKey" default:"goflet"` // The signing key for the JWT when the algorithm is HS256/HS384/HS512
			PublicKey  string `json:"publicKey"`                   // The public key for the JWT when the algorithm is RS256/RS384/RS512
			PrivateKey string `json:"privateKey"`                  // The private key for the JWT when the algorithm is RS256/RS384/RS512
		}
		TrustedIssuers []string `json:"trustedIssuers"` // The list of trusted issuers for the JWT, if empty, it will trust any issuer
	} `json:"jwtConfig"`
	CronConfig struct {
		// Cron configuration, if the value le 0, the cron job will be disabled
		DeleteEmptyFolder int `json:"deleteEmptyFolder" default:"3600"` // The interval to delete empty folders, in seconds
		CleanOutdatedFile int `json:"cleanOutdatedFile" default:"3600"` // The interval to clean outdated files, in seconds
	} `json:"cronConfig"`
}

// GetEndpoint returns the endpoint for the HTTP/S server
func (c *GofletConfig) GetEndpoint() string {
	porti := c.HTTPConfig.Port
	port := strconv.Itoa(porti)
	return c.HTTPConfig.Host + ":" + port
}

var (
	// GofletCfg is the configuration for the application
	GofletCfg GofletConfig
)

func init() {
	InitConfig()
}

// InitConfig initializes the configuration
func InitConfig() {
	// Export default config into file if not exists
	if _, err := os.Stat(PathOfConfig); os.IsNotExist(err) {
		file, err := os.Create(PathOfConfig)
		if err != nil {
			panic(err)
		}
		defer func(file *os.File) {
			err := file.Close()
			if err != nil {
				panic(err)
			}
		}(file)

		// Write the default configuration to the file
		_, err = file.WriteString(defaultConfig)
		if err != nil {
			panic(err)
		}
	}

	// Load the configuration from the file
	err := loadConfig()
	if err != nil {
		panic(err)
	}

	// Set the default values
	confutil.SetDefaults(&GofletCfg)

	// Set the default value for the cache type
	if !GofletCfg.JWTConfig.Enabled {
		fmt.Println("[WARN] JWT is disabled, the security of the application is not guaranteed.")
	}

	// Print the configuration if the debug mode is enabled
	if GofletCfg.Debug {
		fmt.Println("[WARN] Debug mode is enabled, the application is not suitable for production.")
		fmt.Printf("%+v\n", GofletCfg)
	}
}

// loadConfig loads the configuration from the file
func loadConfig() error {
	// Load the configuration from the file
	file, err := os.Open(PathOfConfig)
	if err != nil {
		return err
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			panic(err)
		}
	}(file)

	jsonErr := json.NewDecoder(file).Decode(&GofletCfg)

	if jsonErr != nil {
		return jsonErr
	}
	return nil
}
