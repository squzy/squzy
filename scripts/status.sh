#!/usr/bin/env bash

# shellcheck disable=SC2046
echo VERSION $($GITHUB_REF | cut -d / -f 3)