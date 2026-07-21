"""OCI image + Docker Hub push macro"""

load("@io_bazel_rules_go//go:def.bzl", "go_cross_binary")
load("@rules_oci//oci:defs.bzl", "oci_image", "oci_push")
load("@rules_pkg//pkg:tar.bzl", "pkg_tar")

# Tags for both the manual and CI push targets are stamped from the `version`
# make variable (`--define version=...`), defaulting to `dev` via .bazelrc.
_PUSH_NAMES = ["squzy_push_hub", "squzy_push_hub_ci"]

def squzy_image(binary, repository):
    """Builds a linux/amd64 OCI image for `binary` and Docker Hub push targets.
    Args:
        binary: the `go_binary` label to containerize, e.g. ":squzy_monitoring".
        repository: Docker Hub repository, e.g. "squzy/squzy_monitoring".
    """
    name = binary.split(":")[-1]
    linux_bin = name + "_linux"

    # Force linux/amd64 so the image is correct regardless of host platform.
    go_cross_binary(
        name = linux_bin,
        target = binary,
        platform = "@io_bazel_rules_go//go/toolchain:linux_amd64",
    )

    pkg_tar(
        name = name + "_layer",
        srcs = [":" + linux_bin],
    )

    oci_image(
        name = name + "_image",
        base = "@distroless_static_linux_amd64",
        tars = [":" + name + "_layer"],
        entrypoint = ["/" + linux_bin],
    )

    # `$(version)` expands from --define version=... (default "dev" in .bazelrc).
    native.genrule(
        name = name + "_version_tag",
        outs = [name + "_version_tag.txt"],
        cmd = "echo \"$(version)\" > $@",
    )

    for push_name in _PUSH_NAMES:
        oci_push(
            name = push_name,
            image = ":" + name + "_image",
            repository = "index.docker.io/" + repository,
            remote_tags = ":" + name + "_version_tag",
        )
