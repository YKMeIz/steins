# Steins

[![License](https://img.shields.io/github/license/YKMeIz/steins.svg?color=%232b2b2b&style=flat-square)](https://github.com/YKMeIz/Steins/blob/main/LICENSE)

Steins is a reverse proxy server that atomically redirect traffic to container instances. Currently, it only serves HTTPS request, while communication between steins and local container instances is via HTTP request. 

### Build

Steins is built as docker container image named steins:$(date), e.g. steins:20210427

To build basic steins:
```
$ scripts/build
```
Basic steins container image does not contain public key and private key for HTTPS service.

To build steins example:
```
$ scripts/build example
```
Steins example contains a pre-generated public key and private key for test purpose.

To build steins bundle:
```
$ scripts/build bundle /path/to/public/key /path/to/private/key
```
Steins bundle put given public key and private key into container image.

### Before Run

In order to allow container instances can communicate each other, it is required that these container instances are located in same virtual network. Simply create a docker network and mention it when run a container instance:
```
$ docker network create my-net
```

### Run Virtual Host Web Application (Container Instance)

Web applications are run as container instances. They are started with given labels so that steins can detect configuration for the web application.
```
$ docker run --restart always -l proxy.steins.server.name=yourDomain.com -l proxy.steins.server.proxy_port=8080 --network my-net -d -i my-app
```
End user outside access web application via address `proxy.steins.server.name`, while your web application inside container instance listens on `proxy.steins.server.proxy_port`.

> In this case, you don't need to publish specific port for container instance.

> `proxy.steins.server.name` is mandatory, while `proxy.steins.server.proxy_port` is optional. If no `proxy.steins.server.proxy_port` is defined, the default port value is `80`.

### Run

Steins should be run after starting all web application instances.
To run steins:
```
$ docker run --restart always -p 443:443 -v /var/run/docker.sock:/var/run/docker.sock --network my-net -d -i steins:20210427
```

If public key and private key are not bundled, run:
```
$ docker run --restart always -p 443:443 -v /path/to/public/key:/server.crt -v /path/to/private/key/server.key -v /var/run/docker.sock:/var/run/docker.sock --network my-net -d -i steins:20210427
```

> Please be aware that steins only detects virtual host web applications at startup. Thus, you may want to restart steins if you have changed domain or port configuration.
