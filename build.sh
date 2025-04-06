#!/bin/sh
env GOOS=js GOARCH=wasm go build -o sisyphos.wasm sisyphos.optimisticotter.me
