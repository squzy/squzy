clean: .clean

build: .build

build_agent: .build_agent

run_agent: .run_agent

push: .push

push_hub: .push_hub

build_squzy: .build_squzy

run_squzy: .run_squzy

test: .test

test_d: .test_debug

test_cover: .test_cover

dep: .dep

lint: .lint

default: build

.lint:
	golangci-lint run

.run_squzy:
	bazel run //apps/squzy:squzy_app

.build_squzy:
	bazel build //apps/squzy:squzy

.push:
	bazel run //apps/squzy:squzy_push

.push_hub:
	bazel run //apps/squzy:squzy_push_hub

.build:
	bazel build //apps/...

.test:
	bazel test --define tag="" //apps/...

.build_agent:
	./build.bash apps/agent/main.go squzy_agent_$(version)
	# bazel build //apps/agent:agent

.run_agent:
	bazel run //apps/agent:squzy_agent_app

.test_debug:
	bazel test --define tag="" //apps/...:all --sandbox_debug

.dep:
	bazel run //:gazelle -- update-repos -from_file=go.mod

.test_cover:
	# bazel coverage --test_arg="-test.coverprofile=c.out" //apps/...
	go test ./... -coverprofile=c.out