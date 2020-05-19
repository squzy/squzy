clean: .clean

build: .build

build_agent_server: .build_agent_server

build_agent: .build_agent

build_bin_squzy: .build_bin_squzy

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

.build_squzy:
	bazel build //apps/squzy_monitoring:squzy_monitoring_src

.test:
	bazel test --cache_test_results=no --define version="local" //...

.build_agent:
	./build.bash agent_client squzy_agent_$(version) $(version)

.build_bin_squzy:
	./build.bash squzy_monitoring squzy_monitoring_$(version) $(version)

.build_agent_server:
	./build.bash squzy_agent_server squzy_agent_server_$(version) $(version)

.test_debug:
	bazel test --define version="local" //...:all --sandbox_debug

.dep:
	bazel run //:gazelle -- update-repos -from_file=go.mod

.test_cover:
	# bazel coverage --test_arg="-test.coverprofile=c.out" //apps/...
	go test ./... -coverprofile=c.out