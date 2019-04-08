
# DODAS client

[![Build Status](https://travis-ci.org/Cloud-PG/dodas-go-client.svg?branch=master)](https://travis-ci.org/Cloud-PG/dodas-go-client)

## Quick start

Download the binary from the latest release on [github](https://github.com/Cloud-PG/dodas-go-client/releases). For instance:

```bash
wget https://github.com/Cloud-PG/dodas-go-client/releases/download/v0.1.1/dodas.zip
unzip dodas.zip
cp dodas /usr/local/bin
```

You can find now a template for creating your client configuration file in [config/client_config.yaml](https://raw.githubusercontent.com/Cloud-PG/dodas-go-client/master/config/client_config.yaml). Note that by default the client will look for `$HOME/.dodas.yaml`.

Now you are ready to go. For instance you can validate a tosca template like this:

```bash
dodas validate --template tests/tosca/valid_template.yml
```

or you can create a cluster through the InfrastructureManager configured in your configuration file:

```bash
dodas create --config my_client_conf.yaml --template my_template.yaml
```

## Building from source

To compile on a linux machine (go version that supports `go modules` is required for building from source: e.g. >= v1.12):

```bash
make build
```

while to compile with Docker:

```bash
make docker-build
```

It's also possible to cross compile for windows and macOS with:

```bash
make windows-build
make macos-build
```

## Contributing

If you want to contribute:

1. create a branch
2. upload your changes
3. create a pull request

Thanks!
