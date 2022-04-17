#!/usr/bin/env bash

export PLUGIN_MODE=prod
export PLUGIN_NO_CACHE=true
export PLUGIN_BUILD_ARGS=hello=world,world=hello
go run cmd/main.go