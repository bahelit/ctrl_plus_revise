# Ctrl+Revise

- [Project overview](#project-overview)
- [Starting Ctrl+Revise](#starting-ctrlrevise)
- [Developing Ctrl+Revise](#developing-ctrlrevise)
- [Deploying Ctrl+Revise](#deploying-ctrlrevise)
- [About Ctrl+Revise](#about-ctrlrevise)


## Overview

Ctrl+Revise is locally-run artificial intelligence (AI) tool designed to elevate your writing standards. This desktop application works in any text program, providing real-time suggestions for improving your writing. Ctrl+Revise uses a large language model to analyze your text and suggest improvements, such as grammar corrections, sentence structure, and vocabulary enhancements.


Ctrl+Revise uses keyboard shortcuts that when pressed have the AI act on the text that is currently selected (highlighted) in the text editor and the response replaces the highlighted text and is added to the clipboard to be pasted elsewhere.

This tool is compatible Windows, Linux, and macOS, supporting AMD, Nvidia, and Apple M1 chip architectures.

None of the text you copy is stored or sent to any server. The AI model runs locally on your machine, ensuring your privacy and data security.

Frontend:
- GUI toolkit: `fyne`

Learn more about the Fyne Toolkit at [fyne.io](https://fyne.io/).

Tools:
- Linter: `golangci-lint`

Learn more about the GolangCI-Lint at [golangci-lint.run](https://golangci-lint.run/).

Dependencies:

- [Ollama](https://ollama.com/)

Ollama is a tool for interacting with various LLMs.

Ctrl+Revise will pull that latest Ollama [Docker](https://docker.com) image and manage running it if docker is running.

Docker command for AMD GPUs:
```bash
docker run -d --device /dev/kfd --device /dev/dri -v ollama:/root/.ollama -p 11434:11434 --name ollama --restart=always ollama/ollama:rocm
```

Running Ollama natively **is** supported, but it is currently up to the user to download and start it. Managing Ollama is on the roadmap.

For users who would like to run Ollama natively, download the latest release from the [Ollama.com](https://ollama.com/download) website.

Arch Linux users can install Ollama from the official repository.

CPU Only:
```bash
pacman -S ollama
```
AMD Only:
```bash
pacman -S ollama-rocm
```
Nvidia Only:
```bash
pacman -S ollama-cuda
```


## Starting

Ctrl+Revise will attempt to connect to Ollama on startup. If it is not running it will attempt to start Ollama using Docker, first looking to see it the image is already downloaded, and if not, it will pull the latest image and start the container. Currently, the container image it pulls down is `ollama/ollama:rocm` which provides support for AMD GPUs.

To start your project, use the `go run` command in your terminal or the make recipe `make run` from the project's root directory.

After cloning the repository, navigate to the project folder and run the following command:
```console
go run .
```

## Developing

To develop the project, you need to have the following tools installed on your machine:
- [Go](https://golang.org/dl/)
- [Stringer](https://pkg.go.dev/golang.org/x/tools/cmd/stringer)
- [Docker](https://docs.docker.com/get-docker/)

#### The Stringer tool
This project uses the stringer tool, this will generate a `<type>_string.go` file with the `PromptMsg` type and its `String()` method. To generate the `string.go` file, use the make recipe `make stringer` or run the following command:
```bash
go run golang.org/x/tools/cmd/stringer@latest -linecomment -type=PromptMsg
```


## Screenshots

![Main Menu](images/Screenshot_Main_Window.png)

![Shortcuts Window](images/Screenshot_Keyboard_Shortcuts.png)

## About Ctrl+Revise

The [**Ctrl+Revise**](https://ctrlplusrevise.com) is in early development and there are many features that are planned to be added. This project is open-source and you can contribute to it by submitting a pull request.
