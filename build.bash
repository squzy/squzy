#!/usr/bin/env bash

rm -rf bin

package_path=$1
package_name=$2
version_inject=$3
if [[ -z "$package_path" ]]; then
  echo "usage: $0 <package-name>"
  exit 1
fi

platforms=("windows/amd64" "windows/386" "darwin/amd64" "linux/arm64" "linux/amd64")

for platform in "${platforms[@]}"
do
    platform_split=(${platform//\// })
    GOOS=${platform_split[0]}
    GOARCH=${platform_split[1]}
    CGO=0
    output_name=bin/$package_name'-'$GOOS'-'$GOARCH
    if [ $GOOS = "windows" ]; then
        output_name+='.exe'
    fi
    if [ $GOOS = "darwin" ] && [ $package_path = "agent_client" ]; then
      echo "one"
      CGO=1
    fi
    env CGO_ENABLED=$CGO GOOS=$GOOS GOARCH=$GOARCH go build -o $output_name -ldflags "-s -w -X squzy/apps/$package_path/version.Version=$version_inject"  apps/$package_path/main.go
    if [ $? -ne 0 ]; then
        echo 'An error has occurred! Aborting the script execution...'
        exit 1
    fi
done