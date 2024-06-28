# Ctrl+Revise

This README file contains all the necessary information about:

- [Project overview](#project-overview)
- [Starting Ctrl+Revise](#starting-ctrlrevise)
- [Developing Ctrl+Revise](#developing-ctrlrevise)
- [Deploying Ctrl+Revise](#deploying-ctrlrevise)

And some words [about Ctrl+Revise](#about-ctrlrevise).

## Overview

Ctrl+Revise is locally-run artificial intelligence (AI) tool designed to elevate your writing standards. This desktop application leverages AI-driven suggestions to refine and improve your written content. With real-time clipboard monitoring, it detects new text being copied and offers expert insights to enhance its clarity, coherence, and overall professionalism.


Additionally, Ctrl+Revise features customizable keyboard shortcuts for streamlined workflow, allowing you to focus on the creative process. This powerful tool is compatible with a range of platforms, including Windows, Linux, and macOS, supporting AMD, Nvidia, and Apple M1 chip architectures.

Frontend:
- GUI toolkit: `fyne`

Learn more about the Fyne Toolkit at [fyne.io](https://fyne.io/).

Tools:
- Linter: `golangci-lint`

Learn more about the GolangCI-Lint at [golangci-lint.run](https://golangci-lint.run/).

Dependencies:

- [Ollama](https://ollama.com/)

Ollama is a tool for interacting with various LLMs. [Docker](https://docker.com) is being used to run Ollama.

Ctrl+Revise will pull that latest Ollama image and manage running it.

Running Ollama natively **is** supported, but it is currently up to the user to download and start it. Managing Ollama natively is on the roadmap.

For users who would like to run Ollama natively, download the latest release from the [Ollama.com](https://ollama.com/download) website.

Ctrl+Revise will attempt to connect to Ollama on startup. If it is not running it will attempt to start Ollama using Docker, first looking to see it the image is already downloaded, and if not, it will pull the latest image and start the container. Currently, the container image it pulls down is `ollama/ollama:rocm` which provides support for AMD GPUs. 
The docker command that is run is:
```bash
docker run -d --device /dev/kfd --device /dev/dri -v ollama:/root/.ollama -p 11434:11434 --name ollama --restart=always ollama/ollama:rocm
```


## Starting

To start your project, us the `go run` command in your terminal or the make recipe `make run`

After cloning the repository, navigate to the project folder and run the following command:
```console
go run .
```

## Developing
To develop the project, you need to have the following tools installed on your machine:
- [Go](https://golang.org/dl/)
- [Docker](https://docs.docker.com/get-docker/)

#### The Stringer tool
This project uses the [stringer tool](https://pkg.go.dev/golang.org/x/tools/cmd/stringer), this will generate a `<type>_string.go` file with the `PromptMsg` type and its `String()` method. To generate the `string.go` file, use the make recipe `make stringer` or run the following command:
```bash
go install golang.org/x/tools/cmd/stringer@latest
stringer -linecomment -type=PromptMsg
```

## Deploying

All deploy settings are located in the `Dockerfile` and `docker-compose.yml` files in the project root folder.

To deploy from the git repo, follow these steps:

1. Go to your hosting/cloud provider and create a new VDS/VPS.
2. Update all OS packages on the server and install Docker, Docker Compose and Git packages.
3. Use `git clone` command to clone the repository with your project to the server and navigate to its folder.
4. Run the `docker-compose up` command to start your project on your server.> ❗️ Don't forget to generate Go files from `*.templ` templates before run the `docker-compose up` command.


## About Ctrl+Revise

The [**Ctrl+Revise**](https://ctrlplusrevise.com) is in early development and there are many features that are planned to be added. The project is open-source and you can contribute to it by submitting a pull request.
