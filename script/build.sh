#!/bin/sh

DIRECTORY="$(dirname "$0")"
ROOT_DIRECTORY="$DIRECTORY/.."

swiftc \
    -target wasm32-unknown-wasi \
    "$ROOT_DIRECTORY/checker/checker.swift" -o "$ROOT_DIRECTORY/checker/checker.wasm" \
    -Xlinker --export=is_valid \
    -Xlinker --export=is_can_be_modified \
    -Xlinker --export=allocate \
    -Xlinker --export=deallocate \
    -Xclang-linker -mexec-model=reactor
