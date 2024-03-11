<p align="center">
    <a href="https://github.com/vvbbnn00/goflet">
        <img src="winres/icon.png" width="200" height="200" alt="goflet">
    </a>
</p>

<div align="center">
    <h1>Goflet</h1>
    <hr>
    <p>
        Goflet is a lightweight file upload and download service written in Go.
    </p>
</div>

<p align="center">
    <a href="https://github.com/vvbbnn00/goflet/blob/master/LICENSE">
        <img src="https://img.shields.io/github/license/vvbbnn00/goflet" alt="license">
    </a>
    <a href="https://app.codacy.com?utm_source=gh&utm_medium=referral&utm_content=&utm_campaign=Badge_grade">
        <img src="https://app.codacy.com/project/badge/Grade/f8c9b87f2f7a4af2baf12ddd8d09c475" alt="codacy-quality"/>
    </a>
    <a href="https://goreportcard.com/report/github.com/vvbbnn00/goflet">
        <img src="https://goreportcard.com/badge/github.com/vvbbnn00/goflet" alt="go-report"/>
    </a>
</p>

<div align="center">
    <a href="README.md">English</a>
    <span> | </span>
    <a href="README_zh.md">简体中文</a>
</div>

## 🤔 What is Goflet?

Goflet implements a lightweight file upload and download service. It is a web service based on the Go language and uses
the Gin framework. Although Goflet has simple functions, it supports features such as breakpoint resumption during
upload and download, and multi-threaded downloading. Goflet's file download interface follows various HTTP standards, so
if your service needs are relatively simple, Goflet can fully serve as your CDN service. In addition to file upload and
download, Goflet also provides simple image processing, OnlyOffice synchronization, and other functions.

## 🚀 Features

- **Lightweight**: Only one binary file is needed to run Goflet.
- **Easy to use**: Goflet's API interface follows HTTP standards, so you can easily use it.
- **Breakpoint resumption**: Supports breakpoint resumption during upload and download, improving the stability and
  efficiency of file transmission.
- **Multi-threaded download**: Supports multi-threaded downloading to speed up the file download process.
- **Image processing**: Provides simple image processing functions, such as image compression, scaling, etc.
- **OnlyOffice synchronization**: Supports synchronization editing of OnlyOffice documents.
- **JWT authentication**: Supports JWT authentication to ensure the security of your files.
- **Docker support**: Supports Docker deployment for convenience and speed.
- **Cross-platform**: Goflet supports operating systems such as Windows, Linux, macOS, etc.

## 📦 Installation and Deployment

### Docker / Docker Compose

Some example Docker Compose files are already prepared in the `docker` directory. You can directly use them to deploy
Goflet.

```bash
git clone https://github.com/vvbbnn00/goflet.git
cd goflet
docker-compose up -d
```

### Binary File

You can download the binary file of Goflet from the [Releases](https://github.com/vvbbnn00/goflet/releases) page and
then run it on your server.

### Source Code Compilation

You can also install Goflet through source code compilation. First, you need to install the Go language runtime
environment, and then execute the following commands:

```bash
git clone https://github.com/vvbbnn00/goflet.git
cd goflet
go build -o goflet
```

## 📄 Configuration File

> **Warning**
>
> Please be sure to modify the JWT signature key in goflet.json before deploying Goflet.

The configuration file of Goflet is goflet.json. When Goflet is run for the first time, it will be automatically
generated. You can configure various parameters of Goflet in goflet.json, such as the listening port, file storage path,
etc. After you have finished modifying the configuration, you need to restart Goflet for the configuration to take
effect.

Below is a detailed explanation of the goflet.json configuration file:

```json5
{
  // Whether to enable debug mode
  "debug": false,
  // Log configuration
  "logConfig": {
    // Whether to enable logging
    "enabled": true,
    // Log level (debug, info, warn, error, fatal)
    "level": "info"
  },
  // Whether to enable Swagger documentation
  "swaggerEnabled": true,
  // HTTP service configuration
  "httpConfig": {
    // Listening address
    "host": "0.0.0.0",
    // Listening port
    "port": 8080,
    // Whether to enable CORS
    "cors": {
      "enabled": true,
      // Allowed origins
      "origins": [
        "*"
      ],
      // Allowed methods
      "methods": [
        "HEAD",
        "GET",
        "POST",
        "PUT",
        "DELETE",
        "OPTIONS"
      ],
      // Allowed headers
      "headers": [
        "Content-Type",
        "Authorization"
      ]
    },
    // Browser cache configuration
    "clientCache": {
      // Whether to enable browser caching
      "enabled": true,
      // Cache time
      "maxAge": 3600
    },
    // HTTPS configuration
    "httpsConfig": {
      "enabled": false,
      // Certificate path
      "cert": "cert/cert.pem",
      // Private key path
      "key": "cert/key.pem"
    }
  },
  // File storage configuration
  "fileConfig": {
    // File storage path
    "baseFileStoragePath": "data",
    // Whether to allow folder creation
    "allowFolderCreation": true,
    // Whether to allow file upload
    "uploadPath": "upload",
    // Upload file size limit
    "uploadLimit": 1073741824,
    // Upload timeout
    "uploadTimeout": 7200
  },
  // Cache configuration
  "cacheConfig": {
    // Cache type (Cache, ...)
    "cacheType": "Cache",
    // In-memory cache configuration
    "Cache": {
      // Maximum number of cache entries
      "maxEntries": 1000,
      // Default TTL
      "defaultTTL": 60
    }
  },
  // Image processing configuration
  "imageConfig": {
    // Default format (jpeg, png, gif)
    "defaultFormat": "jpeg",
    // Allowed formats (png, jpeg, gif)
    "allowedFormats": [
      "png",
      "jpeg",
      "gif"
    ],
    // Strict mode, if enabled, only allowed sizes are permitted
    "strictMode": true,
    // Allowed sizes
    "allowedSizes": [
      16,
      32,
      64,
      128,
      256,
      512,
      1024
    ],
    // Maximum file size
    "maxFileSize": 20971520,
    // Maximum width
    "maxWidth": 4096,
    // Maximum height
    "maxHeight": 4096
  },
  // JWT configuration
  "jwtConfig": {
    // Whether to enable JWT (strongly recommended if you are deploying on the public network)
    "enabled": true,
    // JWT algorithm, supports HS256, HS384, HS512, RS256, RS384, RS512, ES256, ES384, ES512
    "algorithm": "HS256",
    // JWT security configuration
    "Security": {
      // JWT signing key, please be sure to modify it when deploying
      "signingKey": "goflet",
      // PKCS8 format public key
      "publicKey": "",
      // PKCS8 format private key (since this project does not actively distribute JWT at present, it can be left empty)
      "privateKey": ""
    },
    // Trusted issuers, if empty, will not be verified
    "trustedIssuers": null
  },
  // Automatic task configuration (if 0, then not enabled, in seconds)
  "cronConfig": {
    // Delete empty folders
    "deleteEmptyFolder": 3600,
    // Clean outdated upload files
    "cleanOutdatedFile": 3600
  }
}

```

## 📝 API Documentation

The integrated Swagger documentation can be accessed at `http://<host>:<port>/swagger/index.html`. You can use it to
easily understand and use Goflet's API.

If you want to disable the Swagger documentation, you can set `swaggerEnabled` to `false` in `goflet.json`.

### Authentication Method

Goflet's authentication method is JWT. You can configure the JWT-related parameters in `goflet.json`. If you are
deploying
on the public network, it is strongly recommended that you enable JWT.

### JWT Format Explanation

Here is an example of a JWT. You can use [JWT.io](https://jwt.io) to parse it. It uses the HS256 algorithm, and the
signing key is `goflet`

```text
eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJnb2ZsZXQiLCJpYXQiOjE3MTAwMDgxNTk
sImV4cCI6MTcxMDA5NDU1OSwibmJmIjoxNzEwMDA4MTU5LCJwZXJtaXNzaW9ucyI6W3sicGF0aCI6Ii9
hcGkvaW1hZ2UvdGVzdC93YWxscGFwZXIucG5nIiwibWV0aG9kcyI6WyJHRVQiXSwicXVlcnkiOnsiaCI
6IjY0IiwidyI6IjY0IiwiZiI6ImpwZWcifX0seyJwYXRoIjoiL2ZpbGUvaW1hZ2VzLyoiLCJtZXRob2R
zIjpbIkdFVCJdfV19.DYSxViy895kjycdt3k0MRNV9UCPho4QsAlRbZ4D9Kik
```

Payload:

```json5
{
  "iss": "goflet",
  "iat": 1710008159,
  "exp": 1710094559,
  "nbf": 1710008159,
  // Permission list, here you can configure the permissions of this JWT, multiple permissions can be configured
  "permissions": [
    {
      // Allowed path, excluding domain and port, supports wildcard (paths after the wildcard will not be verified)
      "path": "/api/image/test/wallpaper.png",
      // Allowed methods (GET, POST, PUT, DELETE, HEAD)
      "methods": [
        "GET"
      ],
      // Allowed parameters, if empty, will not be verified, the parameters configured here need to be exactly matched in the request, other parameters will not be verified, wildcards are not supported
      "query": {
        "h": "64",
        "w": "64",
        "f": "jpeg"
      }
    },
    // query is not mandatory, you can also not set parameters, this will allow all requests
    {
      "path": "/file/images/*",
      // However, you need to set methods, otherwise, it will not be verified
      "methods": [
        "GET"
      ]
    }
  ]
}

```

## 📜 License

Goflet is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

