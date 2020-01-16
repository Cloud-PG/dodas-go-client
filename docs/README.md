
# DODAS client

[![Build Status](https://travis-ci.org/Cloud-PG/dodas-go-client.svg?branch=master)](https://travis-ci.org/Cloud-PG/dodas-go-client)

[Reference Manual](https://cloud-pg.github.io/dodas-go-client/dodas)

## Quick start

Download the binary from the latest release on [github](https://github.com/Cloud-PG/dodas-go-client/releases). For instance:

```bash
wget https://github.com/Cloud-PG/dodas-go-client/releases/download/v0.3.3/dodas.zip
unzip dodas.zip
cp dodas /usr/local/bin
```

> **CLI autocomplete**
>
> Autocompletion for bash and zsh is supported
>
> - **bash** add the following line to ~/.bashrc: `. <(dodas autocomplete)`
> - **zsh** add the following line to ~/.zshrc: `source <(dodas zsh-autocomplete)`

You can find now a template for creating your client configuration file in [config/client_config.yaml](https://raw.githubusercontent.com/Cloud-PG/dodas-go-client/master/config/client_config.yaml). Note that by default the client will look for `$HOME/.dodas.yaml`.

Now you are ready to go. For instance you can validate a tosca template like this:

```bash
dodas validate --template tests/tosca/valid_template.yml
```

> **Tip:** you can find supported tosca templates for applications and k8s deployments on the [DODAS template repo](https://github.com/DODAS-TS/dodas-templates)

or you can create a cluster through the InfrastructureManager configured in your configuration file:

```bash
dodas create --config my_client_conf.yaml my_template.yaml
```

To list the Infrastructure ID of all your deployments:

```bash
dodas list infIDs
```

You can eventually login into a vm in created cluster that has a public IP address with:

```bash
dodas login <infID> <vmID>
# e.g. dodas login cb585e5c-33b6-11ea-8776-0242ac150003 0
```

## Using docker image

In alternative you can create a docker image with the compiled client inside with:

```bash
make docker-img-build
```

and then you can bind the configuration files and run the previous commands as:

```bash

# list the Infrastructure ID of all your deployments
docker run -v $HOME/.dodas.yaml:/app/.dodas.yaml --rm dodas list infIDs
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
