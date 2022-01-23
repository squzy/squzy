#!/bin/sh

git config core.hooksPath hooks

command_exists () {
    type "$1" &> /dev/null ;
}

if command_exists bazelisk; then
    echo "Try to install bazelisk it will help with work with bazel"
fi