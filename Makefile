clean: .clean

generate_build: .gen_build

build: .build

build_agent_server: .build_agent_server

build_incident: .build_incident

build_agent: .build_agent

build_bin_api: .build_bin_api

build_agent_mac: .build_agent_mac

build_bin_squzy: .build_bin_squzy

build_bin_storage: .build_bin_storage

build_application_monitoring: .build_application_monitoring

build_notification: .build_notification

run_agent: .run_agent

run_squzy: .run_squzy

test: .test

test_d: .test_debug

test_cover: .test_cover

dep: .dep

lint: .lint

default: build

.lint:
	golangci-lint run

.gen_build:
	bazel run gazelle -- fix

.build_squzy:
	bazel build //apps/squzy_monitoring:squzy_monitoring_src

.test:
	bazel test --cache_test_results=no --define version="local" //...


.build_agent:
	./build.bash agent_client squzy_agent_$(version) $(version)

.build_incident:
	./build.bash squzy_incident squzy_incident_$(version) $(version)

.build_notification:
	./build.bash squzy_notification squzy_notification_$(version) $(version)

.build_bin_squzy:
	./build.bash squzy_monitoring squzy_monitoring_$(version) $(version)

.build_bin_storage:
	./build.bash squzy_storage squzy_storage_$(version) $(version)

.build_bin_api:
	./build.bash squzy_api squzy_api_$(version) $(version)

.build_agent_server:
	./build.bash squzy_agent_server squzy_agent_server_$(version) $(version)

.build_application_monitoring:
	./build.bash squzy_application_monitoring squzy_application_monitoring_$(version) $(version)

.test_debug:
	bazel test --define version="local" //...:all --sandbox_debug

.dep:
	bazel run //:gazelle -- update-repos -from_file=go.mod

.build_agent_mac:
	env CGO_ENABLED=1 GOOS=darwin GOARCH=amd64 go build -o bin/squzy_agent_$(version)-darwin-amd64 -ldflags "-s -w -X squzy/apps/agent_client/version.Version=$(version)"  apps/agent_client/main.go

.test_cover:
	# bazel coverage --test_arg="-test.coverprofile=c.out" //apps/...
	go test ./... -coverprofile=c.out