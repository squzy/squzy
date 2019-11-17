clean: .clean

build: .build

test: .test

test_d: .test_debug

test_cover: .test_cover

dep: .dep

default: build

.test:
	bazel test //apps/...

.test_debug:
	bazel test //apps/...:all --sandbox_debug

.dep:
	bazel run //:gazelle -- update-repos -from_file=go.mod

.test_cover:
	# bazel coverage --test_arg="-test.coverprofile=c.out" //apps/...
	go test ./... -coverprofile=c.out
