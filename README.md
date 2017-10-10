# Last Mile Transfer
LMT is a proxy service that extends the [File Transfer Service](http://fts3-docs.web.cern.ch/fts3-docs/) in order to enable local data transfers on the [WLCG](http://wlcg-public.web.cern.ch/) infrastructure.

On one end, LMT connects to clients using [WebFTS](http://fts3-service.web.cern.ch/documentation/webfts) via a Websocket connection, on the other, it listens to incoming connections from [FTS](http://fts3-docs.web.cern.ch/fts3-docs/). When FTS is ready to begin the transfer, it start forwarding the data from the client accordingly. 

## Deployment
### Build from sources
In order to ensure reproducible builds, all of the packages that LMT depends on have been included in a /vendor directory.
To build the project:
```bash
go build -o lmt .
```
### Usage
```shell
Usage of ./lmt:
  -cert string
        path to the server's certificate in PEM format (default "/etc/grid-security/hostcert.pem")
  -key string
        path to the server's private key in PEM format (default "/etc/grid-security/hostkey.pem")
  -port string
        port to listen on (default "8080")
```
Example:
```bash
./lmt -port=8080
```
### Run it inside Docker
1. Build the Docker image:

    ```bash
    docker build -t lmt .
    ```
2. Run it:

    ```bash
    docker run -it lmt
    ```

## Development
A development environment has been set up using [Vagrant](https://www.vagrantup.com/docs/installation/). It contains all the dependencies and configuration files necessary to work on the development of the project.

If you have Vagrant installed, just hit `vagrant up` from within the project directory and everything is installed and configured for you to work.