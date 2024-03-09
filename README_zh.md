<p align="center">
    <a href="https://github.com/vvbbnn00/goflet">
        <img src="winres/icon.png" width="200" height="200" alt="goflet">
    </a>
</p>

<div align="center">
    <h1>Goflet</h1>
    <hr>
    <p>
        Goflet 是一个基于 Go 语言的轻量级文件上传和下载服务。
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

## 🤔 Goflet是什么？

Goflet实现了一个轻量级的文件上传和下载服务，它是一个基于Go语言的Web服务，使用了Gin框架。虽然Goflet的功能很简单，但它可以在上传下载时支持断点续传、
多线程下载等功能。Goflet的文件下载接口遵循了HTTP的各个标准，因此，若您的服务需求较为简单，Goflet完全可以胜任为您的CDN服务。除了文件上传和下载，
Goflet也提供了简单的图像处理、OnlyOffice同步等功能。

## 🚀 功能特点

- **轻量级**: 只需要一个二进制文件，即可运行Goflet。
- **简单易用**: Goflet的API接口遵循了HTTP标准，因此您可以很容易地使用它。
- **断点续传**: 支持在上传和下载时断点续传，提高文件传输的稳定性和效率。
- **多线程下载**: 支持多线程下载，加速文件下载过程。
- **图像处理**: 提供简单的图像处理功能，如图片压缩、缩放等。
- **OnlyOffice同步**: 支持OnlyOffice文档的同步编辑功能。
- **JWT鉴权**: 支持JWT鉴权，保障您的文件安全。
- **Docker支持**: 支持Docker部署，方便快捷。
- **跨平台**: Goflet支持Windows、Linux、macOS等操作系统。

## 📦 安装和部署

### Docker / Docker Compose

一些示例Docker Compose文件已经在`docker`目录下准备好了，您可以直接使用它们来部署Goflet。

```bash
git clone https://github.com/vvbbnn00/goflet.git
cd goflet
docker-compose up -d
```

### 二进制文件

您可以在[Releases](https://github.com/vvbbnn00/goflet/releases)页面下载Goflet的二进制文件，然后在您的服务器上运行它。

### 源码编译

您也可以通过源码编译的方式来安装Goflet。首先，您需要安装Go语言的运行环境，然后执行以下命令：

```bash
git clone https://github.com/vvbbnn00/goflet.git
cd goflet
go build -o goflet
```

## 📄 配置文件

> **Warning**
>
> 请务必在部署Goflet前修改`goflet.json`中的JWT签名密钥。

Goflet的配置文件是`goflet.json`，在初次运行Goflet时，它会自动生成。您可以在`goflet.json`中配置Goflet的各种参数，例如监听端口、
文件存储路径等配置。在您修改完毕配置后，需要重启Goflet才能使配置生效。

下面是针对`goflet.json`配置文件的详细说明：

```json5
{
  // 是否开启调试模式
  "debug": false,
  // 日志配置
  "logConfig": {
    // 是否开启日志
    "enabled": true,
    // 日志级别(debug, info, warn, error, fatal)
    "level": "info"
  },
  // 是否开启Swagger文档，未实装
  "swaggerEnabled": true,
  // HTTP服务配置
  "httpConfig": {
    // 监听地址
    "host": "0.0.0.0",
    // 监听端口
    "port": 8080,
    // 是否开启CORS
    "cors": {
      "enabled": true,
      // 允许的来源
      "origins": [
        "*"
      ],
      // 允许的方法
      "methods": [
        "HEAD",
        "GET",
        "POST",
        "PUT",
        "DELETE",
        "OPTIONS"
      ],
      // 允许携带的头部
      "headers": [
        "Content-Type",
        "Authorization"
      ]
    },
    // 浏览器缓存配置
    "clientCache": {
      // 是否开启浏览器缓存
      "enabled": true,
      // 缓存时间
      "maxAge": 3600
    },
    // HTTPS配置
    "httpsConfig": {
      "enabled": false,
      // 证书路径
      "cert": "cert/cert.pem",
      // 私钥路径
      "key": "cert/key.pem"
    }
  },
  // 文件存储配置
  "fileConfig": {
    // 文件存储路径
    "baseFileStoragePath": "data",
    // 是否允许创建文件夹
    "allowFolderCreation": true,
    // 是否允许上传文件
    "uploadPath": "upload",
    // 上传文件大小限制
    "uploadLimit": 1073741824,
    // 上传超时时间
    "uploadTimeout": 7200
  },
  // 缓存配置
  "cacheConfig": {
    // 缓存类型（Cache, ...）
    "cacheType": "Cache",
    // 内存缓存配置
    "Cache": {
      // 最大缓存条目数
      "maxEntries": 1000,
      // 默认TTL
      "defaultTTL": 60
    }
  },
  // 图像处理配置
  "imageConfig": {
    // 默认格式(jpeg, png, gif)
    "defaultFormat": "jpeg",
    // 允许的格式(png, jpeg, gif)
    "allowedFormats": [
      "png",
      "jpeg",
      "gif"
    ],
    // 严格模式，若启用，则只允许允许的尺寸
    "strictMode": true,
    // 允许的尺寸
    "allowedSizes": [
      16,
      32,
      64,
      128,
      256,
      512,
      1024
    ],
    // 最大文件大小
    "maxFileSize": 20971520,
    // 最大宽度
    "maxWidth": 4096,
    // 最大高度
    "maxHeight": 4096
  },
  // JWT配置
  "jwtConfig": {
    // 是否开启JWT（若您部署在公网，强烈建议开启）
    "enabled": true,
    // JWT算法，支持HS256, HS384, HS512, RS256, RS384, RS512, ES256, ES384, ES512
    "algorithm": "HS256",
    // JWT安全配置
    "Security": {
      // JWT签名密钥，部署时请务必修改
      "signingKey": "goflet",
      // PKCS8格式公钥
      "publicKey": "",
      // PKCS8格式私钥（由于目前本项目不主动分发JWT，因此可以不填）
      "privateKey": ""
    },
    // 信任的发行者，若为空，则不验证
    "trustedIssuers": null
  },
  // 自动任务配置（若为0，则不启用，单位为秒）
  "cronConfig": {
    // 清理空文件夹
    "deleteEmptyFolder": 3600,
    // 清理过期的上传文件
    "cleanOutdatedFile": 3600
  }
}

```

## 📝 API文档

文档尚未完成，不过这里提供了目前版本下，Goflet的API列表。

### 接口列表

- [HEAD/GET/PUT/POST/DELETE] /file/{path}
- [GET] /api/meta/{path}
- [GET] /api/image/{path}?w={width}&h={height}&f={format}&q={quality}&a={angle}&s={scaleType:
  fit,fill,resize,fit_width,fit_height}
- [POST] /api/onlyoffice/{path}

### 鉴权方式

Goflet的鉴权方式是JWT，您可以在`goflet.json`中配置JWT的相关参数。若您部署在公网上，强烈建议您开启JWT。

#### JWT格式说明

此处提供了一个JWT的示例，您可以使用[JWT.io](https://jwt.io)来解析它，它使用了HS256算法，签名密钥为`goflet`。

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
  // 权限列表，此处可以配置这个JWT的权限，可以配置多个权限
  "permissions": [
    {
      // 允许访问的路径，不包含域名和端口，支持通配符（通配符之后的路径不会被验证）
      "path": "/api/image/test/wallpaper.png",
      // 允许的方法（GET, POST, PUT, DELETE, HEAD）
      "methods": [
        "GET"
      ],
      // 允许的参数，若为空，则不验证，配置在此处的参数需要在请求中完全匹配，其他的参数则不会被验证，不支持通配符
      "query": {
        "h": "64",
        "w": "64",
        "f": "jpeg"
      }
    },
    // query不是必要的，您也可以不设置参数，这样就会允许所有的请求
    {
      "path": "/file/images/*",
      // 但是，您需要设置methods，否则将不会被验证
      "methods": [
        "GET"
      ]
    }
  ]
}
```

## 📜 许可证

Goflet是一个开源项目，它使用了MIT许可证。您可以在[这里](LICENSE)找到许可证的详细内容。
