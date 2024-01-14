workspace(name = "squzy")

load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")
load("@bazel_tools//tools/build_defs/repo:git.bzl", "git_repository")

http_archive(
    name = "bazel_gazelle",
    sha256 = "d3fa66a39028e97d76f9e2db8f1b0c11c099e8e01bf363a923074784e451f809",
    urls = [
        "https://mirror.bazel.build/github.com/bazelbuild/bazel-gazelle/releases/download/v0.33.0/bazel-gazelle-v0.33.0.tar.gz",
        "https://github.com/bazelbuild/bazel-gazelle/releases/download/v0.33.0/bazel-gazelle-v0.33.0.tar.gz",
    ],
)

http_archive(
    name = "io_bazel_rules_go",
    sha256 = "278b7ff5a826f3dc10f04feaf0b70d48b68748ccd512d7f98bf442077f043fe3",
    urls = [
        "https://mirror.bazel.build/github.com/bazelbuild/rules_go/releases/download/v0.41.0/rules_go-v0.41.0.zip",
        "https://github.com/bazelbuild/rules_go/releases/download/v0.41.0/rules_go-v0.41.0.zip",
    ],
)

load("@io_bazel_rules_go//go:deps.bzl", "go_register_toolchains", "go_rules_dependencies")
load("@bazel_gazelle//:deps.bzl", "gazelle_dependencies", "go_repository")

go_repository(
    name = "com_github_chzyer_logex",
    importpath = "github.com/chzyer/logex",
    sum = "h1:Swpa1K6QvQznwJRcfTfQJmTE72DqScAa40E+fbHEXEE=",
    version = "v1.1.10",
)

go_repository(
    name = "com_github_chzyer_readline",
    importpath = "github.com/chzyer/readline",
    sum = "h1:fY5BOSpyZCqRo5OhCuC+XN+r/bBCmeuuJtjz+bCNIf8=",
    version = "v0.0.0-20180603132655-2972be24d48e",
)

go_repository(
    name = "com_github_chzyer_test",
    importpath = "github.com/chzyer/test",
    sum = "h1:q763qf9huN11kDQavWsoZXJNW3xEE4JJyHa5Q25/sd8=",
    version = "v0.0.0-20180213035817-a1ea475d72b1",
)

go_repository(
    name = "com_github_go_redis_redismock_v8",
    importpath = "github.com/go-redis/redismock/v8",
    sum = "h1:RJFIiua58hrBrSpXhnGX3on79AU3S271H4ZhRI1wyVo=",
    version = "v8.11.5",
)

go_repository(
    name = "com_github_google_pprof",
    importpath = "github.com/google/pprof",
    sum = "h1:yAJXTCF9TqKcTiHJAE8dj7HMvPfh66eeA2JYW7eFpSE=",
    version = "v0.0.0-20210407192527-94a9f03dee38",
)

go_repository(
    name = "com_github_ianlancetaylor_demangle",
    importpath = "github.com/ianlancetaylor/demangle",
    sum = "h1:mV02weKRL81bEnm8A0HT1/CAelMQDBuQIfLw8n+d6xI=",
    version = "v0.0.0-20200824232613-28f6c0f3b639",
)

go_repository(
    name = "com_github_bitly_go_hostpool",
    importpath = "github.com/bitly/go-hostpool",
    sum = "h1:mXoPYz/Ul5HYEDvkta6I8/rnYM5gSdSV2tJ6XbZuEtY=",
    version = "v0.0.0-20171023180738-a3a6125de932",
)

go_repository(
    name = "com_github_bmizerany_assert",
    importpath = "github.com/bmizerany/assert",
    sum = "h1:DDGfHa7BWjL4YnC6+E63dPcxHo2sUxDIu8g3QgEJdRY=",
    version = "v0.0.0-20160611221934-b7ed37b82869",
)

go_repository(
    name = "com_github_gocql_gocql",
    importpath = "github.com/gocql/gocql",
    sum = "h1:IdFdOTbnpbd0pDhl4REKQDM+Q0SzKXQ1Yh+YZZ8T/qU=",
    version = "v1.6.0",
)

go_repository(
    name = "com_github_hailocab_go_hostpool",
    importpath = "github.com/hailocab/go-hostpool",
    sum = "h1:5upAirOpQc1Q53c0bnx2ufif5kANL7bfZWcc6VJWJd8=",
    version = "v0.0.0-20160125115350-e80d13ce29ed",
)

go_repository(
    name = "in_gopkg_inf_v0",
    importpath = "gopkg.in/inf.v0",
    sum = "h1:73M5CoZyi3ZLMOyDlQh031Cx6N9NDJ2Vvfl76EDAgDc=",
    version = "v0.9.1",
)

go_rules_dependencies()

go_register_toolchains(version = "1.20.5")

http_archive(
    name = "googleapis",
    sha256 = "9d1a930e767c93c825398b8f8692eca3fe353b9aaadedfbcf1fca2282c85df88",
    strip_prefix = "googleapis-64926d52febbf298cb82a8f472ade4a3969ba922",
    urls = [
        "https://github.com/googleapis/googleapis/archive/64926d52febbf298cb82a8f472ade4a3969ba922.zip",
    ],
)

load("@googleapis//:repository_rules.bzl", "switched_rules_by_language")

switched_rules_by_language(
    name = "com_google_googleapis_imports",
)

http_archive(
    name = "com_google_protobuf",
    sha256 = "39b52572da90ad54c883a828cb2ca68e5ac918aa75d36c3e55c9c76b94f0a4f7",
    strip_prefix = "protobuf-24.2",
    urls = ["https://github.com/protocolbuffers/protobuf/archive/v24.2.tar.gz"],
)

load("@com_google_protobuf//:protobuf_deps.bzl", "protobuf_deps")

protobuf_deps()

go_repository(
    name = "com_google_cloud_go_compute",
    importpath = "cloud.google.com/go/compute",
    sum = "h1:JNBsyXVoOoNJtTQcnEY5uYpZIbeCTYIeDe0Xh1bySMk=",
    version = "v1.21.0",
)

go_repository(
    name = "com_google_cloud_go_compute_metadata",
    importpath = "cloud.google.com/go/compute/metadata",
    sum = "h1:mg4jlk7mCAj6xXp9UJ4fjI9VUI5rubuGBW5aJ7UnBMY=",
    version = "v0.2.3",
)

go_repository(
    name = "org_golang_google_genproto_googleapis_api",
    importpath = "google.golang.org/genproto/googleapis/api",
    sum = "h1:FmF5cCW94Ij59cfpoLiwTgodWmm60eEV0CjlsVg2fuw=",
    version = "v0.0.0-20230711160842-782d3b101e98",
)

go_repository(
    name = "org_golang_google_genproto_googleapis_rpc",
    importpath = "google.golang.org/genproto/googleapis/rpc",
    sum = "h1:bVf09lpb+OJbByTj913DRJioFFAjf/ZGxEz7MajTp2U=",
    version = "v0.0.0-20230711160842-782d3b101e98",
)

gazelle_dependencies()

http_archive(
    name = "io_bazel_rules_docker",
    sha256 = "59536e6ae64359b716ba9c46c39183403b01eabfbd57578e84398b4829ca499a",
    strip_prefix = "rules_docker-0.22.0",
    urls = ["https://github.com/bazelbuild/rules_docker/releases/download/v0.22.0/rules_docker-v0.22.0.tar.gz"],
)

load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")
load(
    "@io_bazel_rules_docker//repositories:repositories.bzl",
    container_repositories = "repositories",
)

container_repositories()

load(
    "@io_bazel_rules_docker//go:image.bzl",
    _go_image_repos = "repositories",
)

_go_image_repos()

go_repository(
    name = "com_github_lufia_plan9stats",
    build_file_proto_mode = "disable_global",
    importpath = "github.com/lufia/plan9stats",
    sum = "h1:6E+4a0GO5zZEnZ81pIr0yLvtUWk2if982qA3F3QD6H4=",
    version = "v0.0.0-20211012122336-39d0f177ccd0",
)

go_repository(
    name = "com_github_power_devops_perfstat",
    build_file_proto_mode = "disable_global",
    importpath = "github.com/power-devops/perfstat",
    sum = "h1:ncq/mPwQF4JjgDlrVEn3C11VoGHZN7m8qihwgMEtzYw=",
    version = "v0.0.0-20210106213030-5aafc221ea8c",
)

go_repository(
    name = "com_github_shirou_gopsutil_v3",
    build_file_proto_mode = "disable_global",
    importpath = "github.com/shirou/gopsutil/v3",
    sum = "h1:VoGxEW2hpmz0Vt3wUvHIl9fquzYLNpVpgNNB7pGJimA=",
    version = "v3.21.12",
)

go_repository(
    name = "com_github_tklauser_go_sysconf",
    build_file_proto_mode = "disable_global",
    importpath = "github.com/tklauser/go-sysconf",
    sum = "h1:JeUVdAOWhhxVcU6Eqr/ATFHgXk/mmiItdKeJPev3vTo=",
    version = "v0.3.9",
)

go_repository(
    name = "com_github_tklauser_numcpus",
    build_file_proto_mode = "disable_global",
    importpath = "github.com/tklauser/numcpus",
    sum = "h1:ILuRUQBtssgnxw0XXIjKUC56fgnOrFoQQ/4+DeU2biQ=",
    version = "v0.3.0",
)

go_repository(
    name = "com_github_yusufpapurcu_wmi",
    build_file_proto_mode = "disable_global",
    importpath = "github.com/yusufpapurcu/wmi",
    sum = "h1:KBNDSne4vP5mbSWnJbO+51IMOXJB67QiYCSBrubbPRg=",
    version = "v1.2.2",
)

go_repository(
    name = "in_gopkg_yaml_v3",
    build_file_proto_mode = "disable_global",
    importpath = "gopkg.in/yaml.v3",
    sum = "h1:fxVm/GzAzEWqLHuvctI91KS9hhNmmWOoWu0XTYJS7CA=",
    version = "v3.0.1",
)

load("//:repositories.bzl", "go_repositories")

go_repository(
    name = "com_github_antihax_optional",
    importpath = "github.com/antihax/optional",
    sum = "h1:xK2lYat7ZLaVVcIuj82J8kIro4V6kDe0AUDFboUCwcg=",
    version = "v1.0.0",
)

go_repository(
    name = "com_github_cespare_xxhash_v2",
    importpath = "github.com/cespare/xxhash/v2",
    sum = "h1:DC2CZ1Ep5Y4k3ZQ899DldepgrayRUGE6BBZ/cd9Cj44=",
    version = "v2.2.0",
)

go_repository(
    name = "com_github_cncf_xds_go",
    importpath = "github.com/cncf/xds/go",
    sum = "h1:/inchEIKaYC1Akx+H+gqO04wryn5h75LSazbRlnya1k=",
    version = "v0.0.0-20230607035331-e9ce68804cb4",
)

go_repository(
    name = "com_github_ghodss_yaml",
    importpath = "github.com/ghodss/yaml",
    sum = "h1:wQHKEahhL6wmXdzwWG11gIVCkOv05bNOh+Rxn0yngAk=",
    version = "v1.0.0",
)

go_repository(
    name = "com_github_grpc_ecosystem_grpc_gateway",
    importpath = "github.com/grpc-ecosystem/grpc-gateway",
    sum = "h1:bM6ZAFZmc/wPFaRDi0d5L7hGEZEx/2u+Tmr2evNHDiI=",
    version = "v1.9.0",
)

go_repository(
    name = "com_github_rogpeppe_fastuuid",
    importpath = "github.com/rogpeppe/fastuuid",
    sum = "h1:gu+uRPtBe88sKxUCEXRoeCvVG90TJmwhiqRpvdhQFng=",
    version = "v0.0.0-20150106093220-6724a57986af",
)

go_repository(
    name = "io_opentelemetry_go_proto_otlp",
    importpath = "go.opentelemetry.io/proto/otlp",
    sum = "h1:rwOQPCuKAKmwGKq2aVNnYIibI6wnV7EvzgfTCzcdGg8=",
    version = "v0.7.0",
)

go_repository(
    name = "com_github_yuin_goldmark",
    importpath = "github.com/yuin/goldmark",
    sum = "h1:fVcFKWvrslecOb/tg+Cc05dkeYx540o0FuFt3nUVDoE=",
    version = "v1.4.13",
)

go_repository(
    name = "com_github_cespare_xxhash",
    importpath = "github.com/cespare/xxhash",
    sum = "h1:a6HrQnmkObjyL+Gs60czilIUGqrzKutQD6XZog3p+ko=",
    version = "v1.1.0",
)

go_repository(
    name = "com_github_oneofone_xxhash",
    importpath = "github.com/OneOfOne/xxhash",
    sum = "h1:KMrpdQIwFcEqXDklaen+P1axHaj9BSKzvpUUfnHldSE=",
    version = "v1.2.2",
)

go_repository(
    name = "com_github_spaolacci_murmur3",
    importpath = "github.com/spaolacci/murmur3",
    sum = "h1:qLC7fQah7D6K1B0ujays3HV9gkFtllcxhzImRR7ArPQ=",
    version = "v0.0.0-20180118202830-f09979ecbc72",
)

go_repository(
    name = "org_golang_x_term",
    importpath = "golang.org/x/term",
    sum = "h1:3R7pNqamzBraeqj/Tj8qt1aQ2HpmlC+Cx/qL/7hn4/c=",
    version = "v0.10.0",
)

go_repository(
    name = "com_github_azure_go_ansiterm",
    importpath = "github.com/Azure/go-ansiterm",
    sum = "h1:UQHMgLO+TxOElx5B5HZ4hJQsoJ/PvUvKRhJHDQXO8P8=",
    version = "v0.0.0-20210617225240-d185dfc1b5a1",
)

go_repository(
    name = "com_github_bkaradzic_go_lz4",
    importpath = "github.com/bkaradzic/go-lz4",
    sum = "h1:RXc4wYsyz985CkXXeX04y4VnZFGG8Rd43pRaHsOXAKk=",
    version = "v1.0.0",
)

go_repository(
    name = "com_github_cenkalti_backoff",
    importpath = "github.com/cenkalti/backoff",
    sum = "h1:tNowT99t7UNflLxfYYSlKYsBpXdEet03Pg2g16Swow4=",
    version = "v2.2.1+incompatible",
)

go_repository(
    name = "com_github_clickhouse_clickhouse_go",
    importpath = "github.com/ClickHouse/clickhouse-go",
    sum = "h1:FfhyEnv6/BaWldyjgT2k4gDDmeNwJ9C4NbY/MXxJlXk=",
    version = "v1.4.5",
)

go_repository(
    name = "com_github_cloudflare_golz4",
    importpath = "github.com/cloudflare/golz4",
    sum = "h1:F1EaeKL/ta07PY/k9Os/UFtwERei2/XzGemhpGnBKNg=",
    version = "v0.0.0-20150217214814-ef862a3cdc58",
)

go_repository(
    name = "com_github_containerd_containerd",
    importpath = "github.com/containerd/containerd",
    sum = "h1:pASeJT3R3YyVn+94qEPk0SnU1OQ20Jd/T+SPKy9xehY=",
    version = "v1.4.1",
)

go_repository(
    name = "com_github_containerd_continuity",
    importpath = "github.com/containerd/continuity",
    sum = "h1:QSqfxcn8c+12slxwu00AtzXrsami0MJb/MQs9lOLHLA=",
    version = "v0.2.2",
)

go_repository(
    name = "com_github_docker_distribution",
    importpath = "github.com/docker/distribution",
    sum = "h1:dvc1KSkIYTVjZgHf/CTC2diTYC8PzhaA5sFISRfNVrE=",
    version = "v2.7.1-0.20190205005809-0d3efadf0154+incompatible",
)

go_repository(
    name = "com_github_docker_docker",
    importpath = "github.com/docker/docker",
    sum = "h1:CEeNmFM0QZIsJCZKMkZx0ZcahTiewkrgiwfYD+dfl1U=",
    version = "v20.10.12+incompatible",
)

go_repository(
    name = "com_github_docker_go_connections",
    importpath = "github.com/docker/go-connections",
    sum = "h1:El9xVISelRB7BuFusrZozjnkIM5YnzCViNKohAFqRJQ=",
    version = "v0.4.0",
)

go_repository(
    name = "com_github_docker_go_units",
    importpath = "github.com/docker/go-units",
    sum = "h1:3uh0PgVws3nIA0Q+MwDC8yjEPf9zjRfZZWXZYDct3Tw=",
    version = "v0.4.0",
)

go_repository(
    name = "com_github_fsnotify_fsnotify",
    importpath = "github.com/fsnotify/fsnotify",
    sum = "h1:hsms1Qyu0jgnwNXIxa+/V/PDsU6CfLf6CNO8H7IWoS4=",
    version = "v1.4.9",
)

go_repository(
    name = "com_github_go_redis_redis",
    importpath = "github.com/go-redis/redis",
    sum = "h1:K0pv1D7EQUjfyoMql+r/jZqCLizCGKFlFgcHWWmHQjg=",
    version = "v6.15.9+incompatible",
)

go_repository(
    name = "com_github_gorilla_context",
    importpath = "github.com/gorilla/context",
    sum = "h1:AWwleXJkX/nhcU9bZSnZoi3h/qGYqQAGhq6zZe/aQW8=",
    version = "v1.1.1",
)

go_repository(
    name = "com_github_gorilla_mux",
    importpath = "github.com/gorilla/mux",
    sum = "h1:Pgr17XVTNXAk3q/r4CpKzC5xBM/qW1uVLV+IhRZpIIk=",
    version = "v1.6.2",
)

go_repository(
    name = "com_github_hpcloud_tail",
    importpath = "github.com/hpcloud/tail",
    sum = "h1:nfCOvKYfkgYP8hkirhJocXT2+zOD8yUNjXaWfTlyFKI=",
    version = "v1.0.0",
)

go_repository(
    name = "com_github_jmoiron_sqlx",
    importpath = "github.com/jmoiron/sqlx",
    sum = "h1:41Ip0zITnmWNR/vHV+S4m+VoUivnWY5E4OJfLZjCJMA=",
    version = "v1.2.0",
)

go_repository(
    name = "com_github_microsoft_go_winio",
    importpath = "github.com/Microsoft/go-winio",
    sum = "h1:aPJp2QD7OOrhO5tQXqQoGSJc+DjDtWTGLOmNyAm6FgY=",
    version = "v0.5.1",
)

go_repository(
    name = "com_github_microsoft_hcsshim",
    importpath = "github.com/Microsoft/hcsshim",
    sum = "h1:ZfF0+zZeYdzMIVMZHKtDKJvLHj76XCuVae/jNkjj0IA=",
    version = "v0.8.6",
)

go_repository(
    name = "com_github_morikuni_aec",
    importpath = "github.com/morikuni/aec",
    sum = "h1:nXxl5PrvVm2L/wCy8dQu6DMTwH4oIuGN8GJDAlqDdVE=",
    version = "v0.0.0-20170113033406-39771216ff4c",
)

go_repository(
    name = "com_github_onsi_ginkgo",
    importpath = "github.com/onsi/ginkgo",
    sum = "h1:8xi0RTUf59SOSfEtZMvwTvXYMzG4gV23XVHOZiXNtnE=",
    version = "v1.16.5",
)

go_repository(
    name = "com_github_onsi_gomega",
    importpath = "github.com/onsi/gomega",
    sum = "h1:Vw7br2PCDYijJHSfBOWhov+8cAnUf8MfMaIOV323l6Y=",
    version = "v1.25.0",
)

go_repository(
    name = "com_github_opencontainers_go_digest",
    importpath = "github.com/opencontainers/go-digest",
    sum = "h1:apOUWs51W5PlhuyGyz9FCeeBIOUDA/6nW8Oi/yOhh5U=",
    version = "v1.0.0",
)

go_repository(
    name = "com_github_opencontainers_image_spec",
    importpath = "github.com/opencontainers/image-spec",
    sum = "h1:9yCKha/T5XdGtO0q9Q9a6T5NUCsTn/DrBg0D7ufOcFM=",
    version = "v1.0.2",
)

go_repository(
    name = "com_github_opencontainers_runc",
    importpath = "github.com/opencontainers/runc",
    sum = "h1:O9+X96OcDjkmmZyfaG996kV7yq8HsoU2h1XRRQcefG8=",
    version = "v1.1.0",
)

go_repository(
    name = "com_github_pierrec_lz4",
    importpath = "github.com/pierrec/lz4",
    sum = "h1:2xWsjqPFWcplujydGg4WmhC/6fZqK42wMM8aXeqhl0I=",
    version = "v2.0.5+incompatible",
)

go_repository(
    name = "com_github_testcontainers_testcontainers_go",
    importpath = "github.com/testcontainers/testcontainers-go",
    sum = "h1:ZyftCfROjGrKlxk3MOUn2DAzWrUtzY/mj17iAkdUIvI=",
    version = "v0.9.0",
)

go_repository(
    name = "in_gopkg_fsnotify_v1",
    importpath = "gopkg.in/fsnotify.v1",
    sum = "h1:xOHLXZwVvI9hhs+cLKq5+I5onOuwQLhQwiu63xxlHs4=",
    version = "v1.4.7",
)

go_repository(
    name = "in_gopkg_tomb_v1",
    importpath = "gopkg.in/tomb.v1",
    sum = "h1:uRGJdciOHaEIrze2W8Q3AKkepLTh2hOroT7a+7czfdQ=",
    version = "v1.0.0-20141024135613-dd632973f1e7",
)

go_repository(
    name = "org_golang_x_time",
    importpath = "golang.org/x/time",
    sum = "h1:SvFZT6jyqRaOeXpc5h/JSfZenJ2O330aBsf7JfSUXmQ=",
    version = "v0.0.0-20190308202827-9d24e82272b4",
)

go_repository(
    name = "tools_gotest",
    importpath = "gotest.tools",
    sum = "h1:VsBPFP1AI068pPrMxtb/S8Zkgf9xEmTLJjfM+P5UIEo=",
    version = "v2.2.0+incompatible",
)

go_repository(
    name = "com_github_alecthomas_template",
    importpath = "github.com/alecthomas/template",
    sum = "h1:cAKDfWh5VpdgMhJosfJnn5/FoN2SRZ4p7fJNX58YPaU=",
    version = "v0.0.0-20160405071501-a0175ee3bccc",
)

go_repository(
    name = "com_github_alecthomas_units",
    importpath = "github.com/alecthomas/units",
    sum = "h1:qet1QNfXsQxTZqLG4oE62mJzwPIB8+Tee4RNCL9ulrY=",
    version = "v0.0.0-20151022065526-2efee857e7cf",
)

go_repository(
    name = "com_github_armon_consul_api",
    importpath = "github.com/armon/consul-api",
    sum = "h1:G1bPvciwNyF7IUmKXNt9Ak3m6u9DE1rF+RmtIkBpVdA=",
    version = "v0.0.0-20180202201655-eb2c6b5be1b6",
)

go_repository(
    name = "com_github_beorn7_perks",
    importpath = "github.com/beorn7/perks",
    sum = "h1:HWo1m869IqiPhD389kmkxeTalrjNbbJTC8LXupb+sl0=",
    version = "v1.0.0",
)

go_repository(
    name = "com_github_bits_and_blooms_bitset",
    importpath = "github.com/bits-and-blooms/bitset",
    sum = "h1:Kn4yilvwNtMACtf1eYDlG8H77R07mZSPbMjLyS07ChA=",
    version = "v1.2.0",
)

go_repository(
    name = "com_github_cenkalti_backoff_v4",
    importpath = "github.com/cenkalti/backoff/v4",
    sum = "h1:6Yo7N8UP2K6LWZnW94DLVSSrbobcWdVzAYOisuDPIFo=",
    version = "v4.1.2",
)

go_repository(
    name = "com_github_checkpoint_restore_go_criu_v5",
    importpath = "github.com/checkpoint-restore/go-criu/v5",
    sum = "h1:wpFFOoomK3389ue2lAb0Boag6XPht5QYpipxmSNL4d8=",
    version = "v5.3.0",
)

go_repository(
    name = "com_github_cilium_ebpf",
    importpath = "github.com/cilium/ebpf",
    sum = "h1:1k/q3ATgxSXRdrmPfH8d7YK0GfqVsEKZAX9dQZvs56k=",
    version = "v0.7.0",
)

go_repository(
    name = "com_github_containerd_console",
    importpath = "github.com/containerd/console",
    sum = "h1:lIr7SlA5PxZyMV30bDW0MGbiOPXwc63yRuCP0ARubLw=",
    version = "v1.0.3",
)

go_repository(
    name = "com_github_coreos_bbolt",
    importpath = "github.com/coreos/bbolt",
    sum = "h1:wZwiHHUieZCquLkDL0B8UhzreNWsPHooDAG3q34zk0s=",
    version = "v1.3.2",
)

go_repository(
    name = "com_github_coreos_etcd",
    importpath = "github.com/coreos/etcd",
    sum = "h1:jFneRYjIvLMLhDLCzuTuU4rSJUjRplcJQ7pD7MnhC04=",
    version = "v3.3.10+incompatible",
)

go_repository(
    name = "com_github_coreos_go_semver",
    importpath = "github.com/coreos/go-semver",
    sum = "h1:3Jm3tLmsgAYcjC+4Up7hJrFBPr+n7rAqYeSw/SZazuY=",
    version = "v0.2.0",
)

go_repository(
    name = "com_github_coreos_go_systemd",
    importpath = "github.com/coreos/go-systemd",
    sum = "h1:Wf6HqHfScWJN9/ZjdUKyjop4mf3Qdd+1TvvltAvM3m8=",
    version = "v0.0.0-20190321100706-95778dfbb74e",
)

go_repository(
    name = "com_github_coreos_go_systemd_v22",
    importpath = "github.com/coreos/go-systemd/v22",
    sum = "h1:D9/bQk5vlXQFZ6Kwuu6zaiXJ9oTPe68++AzAJc1DzSI=",
    version = "v22.3.2",
)

go_repository(
    name = "com_github_coreos_pkg",
    importpath = "github.com/coreos/pkg",
    sum = "h1:lBNOc5arjvs8E5mO2tbpBpLoyyu8B6e44T7hJy6potg=",
    version = "v0.0.0-20180928190104-399ea9e2e55f",
)

go_repository(
    name = "com_github_cpuguy83_go_md2man_v2",
    importpath = "github.com/cpuguy83/go-md2man/v2",
    sum = "h1:EoUDS0afbrsXAZ9YQ9jdu/mZ2sXgT1/2yyNng4PGlyM=",
    version = "v2.0.0",
)

go_repository(
    name = "com_github_creack_pty",
    importpath = "github.com/creack/pty",
    sum = "h1:07n33Z8lZxZ2qwegKbObQohDhXDQxiMMz1NOUGYlesw=",
    version = "v1.1.11",
)

go_repository(
    name = "com_github_cyphar_filepath_securejoin",
    importpath = "github.com/cyphar/filepath-securejoin",
    sum = "h1:YX6ebbZCZP7VkM3scTTokDgBL2TY741X51MTk3ycuNI=",
    version = "v0.2.3",
)

go_repository(
    name = "com_github_dgrijalva_jwt_go",
    importpath = "github.com/dgrijalva/jwt-go",
    sum = "h1:7qlOGliEKZXTDg6OTjfoBKDXWrumCAMpl/TFQ4/5kLM=",
    version = "v3.2.0+incompatible",
)

go_repository(
    name = "com_github_dgryski_go_sip13",
    importpath = "github.com/dgryski/go-sip13",
    sum = "h1:RMLoZVzv4GliuWafOuPuQDKSm1SJph7uCRnnS61JAn4=",
    version = "v0.0.0-20181026042036-e10d5fee7954",
)

go_repository(
    name = "com_github_docker_cli",
    importpath = "github.com/docker/cli",
    sum = "h1:lZlz0uzG+GH+c0plStMUdF/qk3ppmgnswpR5EbqzVGA=",
    version = "v20.10.12+incompatible",
)

go_repository(
    name = "com_github_dustin_go_humanize",
    importpath = "github.com/dustin/go-humanize",
    sum = "h1:VSnTsYCnlFHaM2/igO1h6X3HA71jcobQuxemgkq4zYo=",
    version = "v1.0.0",
)

go_repository(
    name = "com_github_frankban_quicktest",
    importpath = "github.com/frankban/quicktest",
    sum = "h1:8sXhOn0uLys67V8EsXLc6eszDs8VXWxL3iRvebPhedY=",
    version = "v1.11.3",
)

go_repository(
    name = "com_github_godbus_dbus_v5",
    importpath = "github.com/godbus/dbus/v5",
    sum = "h1:mkgN1ofwASrYnJ5W6U/BxG15eXXXjirgZc7CLqkcaro=",
    version = "v5.0.6",
)

go_repository(
    name = "com_github_golang_groupcache",
    importpath = "github.com/golang/groupcache",
    sum = "h1:veQD95Isof8w9/WXiA+pa3tz3fJXkt5B7QaRBrM62gk=",
    version = "v0.0.0-20190129154638-5b532d6fd5ef",
)

go_repository(
    name = "com_github_google_btree",
    importpath = "github.com/google/btree",
    sum = "h1:0udJVsspx3VBr5FwtLhQQtuAsVc79tTq0ocGIPAU6qo=",
    version = "v1.0.0",
)

go_repository(
    name = "com_github_google_shlex",
    importpath = "github.com/google/shlex",
    sum = "h1:El6M4kTTCOh6aBiKaUGG7oYTSPP8MxqL4YI3kZKwcP4=",
    version = "v0.0.0-20191202100458-e7afc7fbc510",
)

go_repository(
    name = "com_github_gotestyourself_gotestyourself",
    importpath = "github.com/gotestyourself/gotestyourself",
    sum = "h1:AQwinXlbQR2HvPjQZOmDhRqsv5mZf+Jb1RnSLxcqZcI=",
    version = "v2.2.0+incompatible",
)

go_repository(
    name = "com_github_grpc_ecosystem_go_grpc_prometheus",
    importpath = "github.com/grpc-ecosystem/go-grpc-prometheus",
    sum = "h1:Ovs26xHkKqVztRpIrF/92BcuyuQ/YW4NSIpoGtfXNho=",
    version = "v1.2.0",
)

go_repository(
    name = "com_github_hashicorp_hcl",
    importpath = "github.com/hashicorp/hcl",
    sum = "h1:0Anlzjpi4vEasTeNFn2mLJgTSwt0+6sfsiTG8qcWGx4=",
    version = "v1.0.0",
)

go_repository(
    name = "com_github_imdario_mergo",
    importpath = "github.com/imdario/mergo",
    sum = "h1:b6R2BslTbIEToALKP7LxUvijTsNI9TAe80pLWN2g/HU=",
    version = "v0.3.12",
)

go_repository(
    name = "com_github_jonboulle_clockwork",
    importpath = "github.com/jonboulle/clockwork",
    sum = "h1:VKV+ZcuP6l3yW9doeqz6ziZGgcynBVQO+obU0+0hcPo=",
    version = "v0.1.0",
)

go_repository(
    name = "com_github_julienschmidt_httprouter",
    importpath = "github.com/julienschmidt/httprouter",
    sum = "h1:TDTW5Yz1mjftljbcKqRcrYhd4XeOoI98t+9HbQbYf7g=",
    version = "v1.2.0",
)

go_repository(
    name = "com_github_magiconair_properties",
    importpath = "github.com/magiconair/properties",
    sum = "h1:LLgXmsheXeRoUOBOjtwPQCWIYqM/LU1ayDtDePerRcY=",
    version = "v1.8.0",
)

go_repository(
    name = "com_github_matttproud_golang_protobuf_extensions",
    importpath = "github.com/matttproud/golang_protobuf_extensions",
    sum = "h1:4hp9jkHxhMHkqkrB3Ix0jegS5sx/RkqARlsWZ6pIwiU=",
    version = "v1.0.1",
)

go_repository(
    name = "com_github_mitchellh_go_homedir",
    importpath = "github.com/mitchellh/go-homedir",
    sum = "h1:lukF9ziXFxDFPkA1vsr5zpc1XuPDn/wFntq5mG+4E0Y=",
    version = "v1.1.0",
)

go_repository(
    name = "com_github_mitchellh_mapstructure",
    importpath = "github.com/mitchellh/mapstructure",
    sum = "h1:OVowDSCllw/YjdLkam3/sm7wEtOy59d8ndGgCcyj8cs=",
    version = "v1.4.3",
)

go_repository(
    name = "com_github_moby_sys_mountinfo",
    importpath = "github.com/moby/sys/mountinfo",
    sum = "h1:2Ks8/r6lopsxWi9m58nlwjaeSzUX9iiL1vj5qB/9ObI=",
    version = "v0.5.0",
)

go_repository(
    name = "com_github_moby_term",
    importpath = "github.com/moby/term",
    sum = "h1:dcztxKSvZ4Id8iPpHERQBbIJfabdt4wUm5qy3wOL2Zc=",
    version = "v0.0.0-20210619224110-3f7ff695adc6",
)

go_repository(
    name = "com_github_mrunalp_fileutils",
    importpath = "github.com/mrunalp/fileutils",
    sum = "h1:NKzVxiH7eSk+OQ4M+ZYW1K6h27RUV3MI6NUTsHhU6Z4=",
    version = "v0.5.0",
)

go_repository(
    name = "com_github_mwitkow_go_conntrack",
    importpath = "github.com/mwitkow/go-conntrack",
    sum = "h1:F9x/1yl3T2AeKLr2AMdilSD8+f9bvMnNN8VS5iDtovc=",
    version = "v0.0.0-20161129095857-cc309e4a2223",
)

go_repository(
    name = "com_github_nvveen_gotty",
    importpath = "github.com/Nvveen/Gotty",
    sum = "h1:TngWCqHvy9oXAN6lEVMRuU21PR1EtLVZJmdB18Gu3Rw=",
    version = "v0.0.0-20120604004816-cd527374f1e5",
)

go_repository(
    name = "com_github_oklog_ulid",
    importpath = "github.com/oklog/ulid",
    sum = "h1:EGfNDEx6MqHz8B3uNV6QAib1UR2Lm97sHi3ocA6ESJ4=",
    version = "v1.3.1",
)

go_repository(
    name = "com_github_opencontainers_runtime_spec",
    importpath = "github.com/opencontainers/runtime-spec",
    sum = "h1:3snG66yBm59tKhhSPQrQ/0bCrv1LQbKt40LnUPiUxdc=",
    version = "v1.0.3-0.20210326190908-1c3f411f0417",
)

go_repository(
    name = "com_github_opencontainers_selinux",
    importpath = "github.com/opencontainers/selinux",
    sum = "h1:rAiKF8hTcgLI3w0DHm6i0ylVVcOrlgR1kK99DRLDhyU=",
    version = "v1.10.0",
)

go_repository(
    name = "com_github_ory_dockertest",
    importpath = "github.com/ory/dockertest",
    sum = "h1:iLLK6SQwIhcbrG783Dghaaa3WPzGc+4Emza6EbVUUGA=",
    version = "v3.3.5+incompatible",
)

go_repository(
    name = "com_github_ory_dockertest_v3",
    importpath = "github.com/ory/dockertest/v3",
    sum = "h1:vU/8d1We4qIad2YM0kOwRVtnyue7ExvacPiw1yDm17g=",
    version = "v3.8.1",
)

go_repository(
    name = "com_github_prometheus_client_golang",
    importpath = "github.com/prometheus/client_golang",
    sum = "h1:9iH4JKXLzFbOAdtqv/a+j8aewx2Y8lAjAydhbaScPF8=",
    version = "v0.9.3",
)

go_repository(
    name = "com_github_prometheus_common",
    importpath = "github.com/prometheus/common",
    sum = "h1:7etb9YClo3a6HjLzfl6rIQaU+FDfi0VSX39io3aQ+DM=",
    version = "v0.4.0",
)

go_repository(
    name = "com_github_prometheus_procfs",
    importpath = "github.com/prometheus/procfs",
    sum = "h1:sofwID9zm4tzrgykg80hfFph1mryUeLRsUfoocVVmRY=",
    version = "v0.0.0-20190507164030-5867b95ac084",
)

go_repository(
    name = "com_github_prometheus_tsdb",
    importpath = "github.com/prometheus/tsdb",
    sum = "h1:YZcsG11NqnK4czYLrWd9mpEuAJIHVQLwdrleYfszMAA=",
    version = "v0.7.1",
)

go_repository(
    name = "com_github_russross_blackfriday_v2",
    importpath = "github.com/russross/blackfriday/v2",
    sum = "h1:lPqVAte+HuHNfhJ/0LC98ESWRz8afy9tM/0RK8m9o+Q=",
    version = "v2.0.1",
)

go_repository(
    name = "com_github_seccomp_libseccomp_golang",
    importpath = "github.com/seccomp/libseccomp-golang",
    sum = "h1:58EBmR2dMNL2n/FnbQewK3D14nXr0V9CObDSvMJLq+Y=",
    version = "v0.9.2-0.20210429002308-3879420cc921",
)

go_repository(
    name = "com_github_shurcool_sanitized_anchor_name",
    importpath = "github.com/shurcooL/sanitized_anchor_name",
    sum = "h1:PdmoCO6wvbs+7yrJyMORt4/BmY5IYyJwS/kOiWx8mHo=",
    version = "v1.0.0",
)

go_repository(
    name = "com_github_soheilhy_cmux",
    importpath = "github.com/soheilhy/cmux",
    sum = "h1:0HKaf1o97UwFjHH9o5XsHUOF+tqmdA7KEzXLpiyaw0E=",
    version = "v0.1.4",
)

go_repository(
    name = "com_github_spf13_afero",
    importpath = "github.com/spf13/afero",
    sum = "h1:m8/z1t7/fwjysjQRYbP0RD+bUIF/8tJwPdEZsI83ACI=",
    version = "v1.1.2",
)

go_repository(
    name = "com_github_spf13_cast",
    importpath = "github.com/spf13/cast",
    sum = "h1:oget//CVOEoFewqQxwr0Ej5yjygnqGkvggSE/gB35Q8=",
    version = "v1.3.0",
)

go_repository(
    name = "com_github_spf13_jwalterweatherman",
    importpath = "github.com/spf13/jwalterweatherman",
    sum = "h1:XHEdyB+EcvlqZamSM4ZOMGlc93t6AcsBEu9Gc1vn7yk=",
    version = "v1.0.0",
)

go_repository(
    name = "com_github_spf13_viper",
    importpath = "github.com/spf13/viper",
    sum = "h1:yXHLWeravcrgGyFSyCgdYpXQ9dR9c/WED3pg1RhxqEU=",
    version = "v1.4.0",
)

go_repository(
    name = "com_github_syndtr_gocapability",
    importpath = "github.com/syndtr/gocapability",
    sum = "h1:kdXcSzyDtseVEc4yCz2qF8ZrQvIDBJLl4S1c3GCXmoI=",
    version = "v0.0.0-20200815063812-42c35b437635",
)

go_repository(
    name = "com_github_tmc_grpc_websocket_proxy",
    importpath = "github.com/tmc/grpc-websocket-proxy",
    sum = "h1:LnC5Kc/wtumK+WB441p7ynQJzVuNRJiqddSIE3IlSEQ=",
    version = "v0.0.0-20190109142713-0ad062ec5ee5",
)

go_repository(
    name = "com_github_tv42_httpunix",
    importpath = "github.com/tv42/httpunix",
    sum = "h1:u6SKchux2yDvFQnDHS3lPnIRmfVJ5Sxy3ao2SIdysLQ=",
    version = "v0.0.0-20191220191345-2ba4b9c3382c",
)

go_repository(
    name = "com_github_urfave_cli",
    importpath = "github.com/urfave/cli",
    sum = "h1:+mkCCcOFKPnCmVYVcURKps1Xe+3zP90gSYGNfRkjoIY=",
    version = "v1.22.1",
)

go_repository(
    name = "com_github_vishvananda_netlink",
    importpath = "github.com/vishvananda/netlink",
    sum = "h1:1iyaYNBLmP6L0220aDnYQpo1QEV4t4hJ+xEEhhJH8j0=",
    version = "v1.1.0",
)

go_repository(
    name = "com_github_vishvananda_netns",
    importpath = "github.com/vishvananda/netns",
    sum = "h1:OviZH7qLw/7ZovXvuNyL3XQl8UFofeikI1NW1Gypu7k=",
    version = "v0.0.0-20191106174202-0a2b9b5464df",
)

go_repository(
    name = "com_github_xeipuuv_gojsonpointer",
    importpath = "github.com/xeipuuv/gojsonpointer",
    sum = "h1:zGWFAtiMcyryUHoUjUJX0/lt1H2+i2Ka2n+D3DImSNo=",
    version = "v0.0.0-20190905194746-02993c407bfb",
)

go_repository(
    name = "com_github_xeipuuv_gojsonreference",
    importpath = "github.com/xeipuuv/gojsonreference",
    sum = "h1:EzJWgHovont7NscjpAxXsDA8S8BMYve8Y5+7cuRE7R0=",
    version = "v0.0.0-20180127040603-bd5ef7bd5415",
)

go_repository(
    name = "com_github_xeipuuv_gojsonschema",
    importpath = "github.com/xeipuuv/gojsonschema",
    sum = "h1:LhYJRs+L4fBtjZUfuSZIKGeVu0QRy8e5Xi7D17UxZ74=",
    version = "v1.2.0",
)

go_repository(
    name = "com_github_xiang90_probing",
    importpath = "github.com/xiang90/probing",
    sum = "h1:eY9dn8+vbi4tKz5Qo6v2eYzo7kUS51QINcR5jNpbZS8=",
    version = "v0.0.0-20190116061207-43a291ad63a2",
)

go_repository(
    name = "com_github_xordataexchange_crypt",
    importpath = "github.com/xordataexchange/crypt",
    sum = "h1:ESFSdwYZvkeru3RtdrYueztKhOBCSAAzS4Gf+k0tEow=",
    version = "v0.0.3-0.20170626215501-b2862e3d0a77",
)

go_repository(
    name = "in_gopkg_alecthomas_kingpin_v2",
    importpath = "gopkg.in/alecthomas/kingpin.v2",
    sum = "h1:jMFz6MfLP0/4fUyZle81rXUoxOBFi19VUFKVDOQfozc=",
    version = "v2.2.6",
)

go_repository(
    name = "in_gopkg_resty_v1",
    importpath = "gopkg.in/resty.v1",
    sum = "h1:CuXP0Pjfw9rOuY6EP+UvtNvt5DSqHpIxILZKT/quCZI=",
    version = "v1.12.0",
)

go_repository(
    name = "io_etcd_go_bbolt",
    importpath = "go.etcd.io/bbolt",
    sum = "h1:Z/90sZLPOeCy2PwprqkFa25PdkusRzaj9P8zm/KNyvk=",
    version = "v1.3.2",
)

go_repository(
    name = "org_bazil_fuse",
    importpath = "bazil.org/fuse",
    sum = "h1:SRsZGA7aFnCZETmov57jwPrWuTmaZK6+4R4v5FUe1/c=",
    version = "v0.0.0-20200407214033-5883e5a4b512",
)

go_repository(
    name = "tools_gotest_v3",
    importpath = "gotest.tools/v3",
    sum = "h1:4AuOwCGf4lLR9u3YOe2awrHygurzhO/HeQ6laiA6Sx0=",
    version = "v3.0.3",
)

go_repository(
    name = "com_github_xdg_go_pbkdf2",
    importpath = "github.com/xdg-go/pbkdf2",
    sum = "h1:Su7DPu48wXMwC3bs7MCNG+z4FhcyEuz5dlvchbq0B0c=",
    version = "v1.0.0",
)

go_repository(
    name = "com_github_xdg_go_scram",
    importpath = "github.com/xdg-go/scram",
    sum = "h1:akYIkZ28e6A96dkWNJQu3nmCzH3YfwMPQExUYDaRv7w=",
    version = "v1.0.2",
)

go_repository(
    name = "com_github_xdg_go_stringprep",
    importpath = "github.com/xdg-go/stringprep",
    sum = "h1:6iq84/ryjjeRmMJwxutI51F2GIPlP5BfTvXHeYjyhBc=",
    version = "v1.0.2",
)

go_repository(
    name = "com_github_youmark_pkcs8",
    importpath = "github.com/youmark/pkcs8",
    sum = "h1:splanxYIlg+5LfHAM6xpdFEAYOk8iySO56hMFq6uLyA=",
    version = "v0.0.0-20181117223130-1be2e3e5546d",
)

go_repository(
    name = "com_github_cockroachdb_apd",
    importpath = "github.com/cockroachdb/apd",
    sum = "h1:3LFP3629v+1aKXU5Q37mxmRxX/pIu1nijXydLShEq5I=",
    version = "v1.1.0",
)

go_repository(
    name = "com_github_go_kit_log",
    importpath = "github.com/go-kit/log",
    sum = "h1:DGJh0Sm43HbOeYDNnVZFl8BvcYVvjD5bqYJvp0REbwQ=",
    version = "v0.1.0",
)

go_repository(
    name = "com_github_gofrs_uuid",
    importpath = "github.com/gofrs/uuid",
    sum = "h1:1SD/1F5pU8p29ybwgQSwpQk+mwdRrXCYuPhW6m+TnJw=",
    version = "v4.0.0+incompatible",
)

go_repository(
    name = "com_github_jackc_chunkreader",
    importpath = "github.com/jackc/chunkreader",
    sum = "h1:4s39bBR8ByfqH+DKm8rQA3E1LHZWB9XWcrz8fqaZbe0=",
    version = "v1.0.0",
)

go_repository(
    name = "com_github_jackc_chunkreader_v2",
    importpath = "github.com/jackc/chunkreader/v2",
    sum = "h1:i+RDz65UE+mmpjTfyz0MoVTnzeYxroil2G82ki7MGG8=",
    version = "v2.0.1",
)

go_repository(
    name = "com_github_jackc_pgconn",
    importpath = "github.com/jackc/pgconn",
    sum = "h1:DzdIHIjG1AxGwoEEqS+mGsURyjt4enSmqzACXvVzOT8=",
    version = "v1.10.1",
)

go_repository(
    name = "com_github_jackc_pgio",
    importpath = "github.com/jackc/pgio",
    sum = "h1:g12B9UwVnzGhueNavwioyEEpAmqMe1E/BN9ES+8ovkE=",
    version = "v1.0.0",
)

go_repository(
    name = "com_github_jackc_pgmock",
    importpath = "github.com/jackc/pgmock",
    sum = "h1:DadwsjnMwFjfWc9y5Wi/+Zz7xoE5ALHsRQlOctkOiHc=",
    version = "v0.0.0-20210724152146-4ad1a8207f65",
)

go_repository(
    name = "com_github_jackc_pgpassfile",
    importpath = "github.com/jackc/pgpassfile",
    sum = "h1:/6Hmqy13Ss2zCq62VdNG8tM1wchn8zjSGOBJ6icpsIM=",
    version = "v1.0.0",
)

go_repository(
    name = "com_github_jackc_pgproto3",
    importpath = "github.com/jackc/pgproto3",
    sum = "h1:FYYE4yRw+AgI8wXIinMlNjBbp/UitDJwfj5LqqewP1A=",
    version = "v1.1.0",
)

go_repository(
    name = "com_github_jackc_pgproto3_v2",
    importpath = "github.com/jackc/pgproto3/v2",
    sum = "h1:r7JypeP2D3onoQTCxWdTpCtJ4D+qpKr0TxvoyMhZ5ns=",
    version = "v2.2.0",
)

go_repository(
    name = "com_github_jackc_pgservicefile",
    importpath = "github.com/jackc/pgservicefile",
    sum = "h1:C8S2+VttkHFdOOCXJe+YGfa4vHYwlt4Zx+IVXQ97jYg=",
    version = "v0.0.0-20200714003250-2b9c44734f2b",
)

go_repository(
    name = "com_github_jackc_pgtype",
    importpath = "github.com/jackc/pgtype",
    sum = "h1:/SH1RxEtltvJgsDqp3TbiTFApD3mey3iygpuEGeuBXk=",
    version = "v1.9.0",
)

go_repository(
    name = "com_github_jackc_pgx_v4",
    importpath = "github.com/jackc/pgx/v4",
    sum = "h1:TgdrmgnM7VY72EuSQzBbBd4JA1RLqJolrw9nQVZABVc=",
    version = "v4.14.0",
)

go_repository(
    name = "com_github_jackc_puddle",
    importpath = "github.com/jackc/puddle",
    sum = "h1:DNDKdn/pDrWvDWyT2FYvpZVE81OAhWrjCv19I9n108Q=",
    version = "v1.2.0",
)

go_repository(
    name = "com_github_masterminds_semver_v3",
    importpath = "github.com/Masterminds/semver/v3",
    sum = "h1:hLg3sBzpNErnxhQtUy/mmLR2I9foDujNK030IGemrRc=",
    version = "v3.1.1",
)

go_repository(
    name = "com_github_mattn_go_colorable",
    importpath = "github.com/mattn/go-colorable",
    sum = "h1:6Su7aK7lXmJ/U79bYtBjLNaha4Fs1Rg9plHpcH+vvnE=",
    version = "v0.1.6",
)

go_repository(
    name = "com_github_rs_xid",
    importpath = "github.com/rs/xid",
    sum = "h1:mhH9Nq+C1fY2l1XIpgxIiUOfNpRBYH1kKcr+qfKgjRc=",
    version = "v1.2.1",
)

go_repository(
    name = "com_github_rs_zerolog",
    importpath = "github.com/rs/zerolog",
    sum = "h1:uPRuwkWF4J6fGsJ2R0Gn2jB1EQiav9k3S6CSdygQJXY=",
    version = "v1.15.0",
)

go_repository(
    name = "com_github_satori_go_uuid",
    importpath = "github.com/satori/go.uuid",
    sum = "h1:0uYX9dsZ2yD7q2RtLRtPSdGDWzjeM3TbMJP9utgA0ww=",
    version = "v1.2.0",
)

go_repository(
    name = "com_github_shopspring_decimal",
    importpath = "github.com/shopspring/decimal",
    sum = "h1:abSATXmQEYyShuxI4/vyW3tV1MrKAJzCZ/0zLUXYbsQ=",
    version = "v1.2.0",
)

go_repository(
    name = "com_github_zenazn_goji",
    importpath = "github.com/zenazn/goji",
    sum = "h1:RSQQAbXGArQ0dIDEq+PI6WqN6if+5KHu6x2Cx/GXLTQ=",
    version = "v0.9.0",
)

go_repository(
    name = "in_gopkg_inconshreveable_log15_v2",
    importpath = "gopkg.in/inconshreveable/log15.v2",
    sum = "h1:RlWgLqCMMIYYEVcAR5MDsuHlVkaIPDAF+5Dehzg8L5A=",
    version = "v2.0.0-20180818164646-67afb5ed74ec",
)

go_repository(
    name = "com_github_jinzhu_gorm",
    importpath = "github.com/jinzhu/gorm",
    sum = "h1:Drgk1clyWT9t9ERbzHza6Mj/8FY/CqMyVzOiHviMo6Q=",
    version = "v1.9.12",
)

go_repository(
    name = "com_github_go_logr_logr",
    importpath = "github.com/go-logr/logr",
    sum = "h1:2DntVwHkVopvECVRSlL5PSo9eG+cAkDCuckLubN+rq0=",
    version = "v1.2.3",
)

go_repository(
    name = "com_github_go_task_slim_sprig",
    importpath = "github.com/go-task/slim-sprig",
    sum = "h1:p104kn46Q8WdvHunIJ9dAyjPVtrBPhSr3KT2yUst43I=",
    version = "v0.0.0-20210107165309-348f09dbbbc0",
)

go_repository(
    name = "com_github_nxadm_tail",
    importpath = "github.com/nxadm/tail",
    sum = "h1:nPr65rt6Y5JFSKQO7qToXr7pePgD6Gwiw05lkbyAQTE=",
    version = "v1.4.8",
)

go_repository(
    name = "com_github_onsi_ginkgo_v2",
    importpath = "github.com/onsi/ginkgo/v2",
    sum = "h1:/XxtEV3I3Eif/HobnVx9YmJgk8ENdRsuUmM+fLCFNow=",
    version = "v2.7.0",
)

go_repository(
    name = "com_github_dgryski_go_rendezvous",
    importpath = "github.com/dgryski/go-rendezvous",
    sum = "h1:lO4WD4F/rVNCu3HqELle0jiPLLBs70cWOduZpkS1E78=",
    version = "v0.0.0-20200823014737-9f7001d12a5f",
)

go_repository(
    name = "com_github_go_redis_redis_v8",
    importpath = "github.com/go-redis/redis/v8",
    sum = "h1:AcZZR7igkdvfVmQTPnu9WE37LRrO/YrBH5zWyjDC0oI=",
    version = "v8.11.5",
)

# gazelle:repository_macro repositories.bzl%go_repositories
go_repositories()

gazelle_dependencies()

go_repository(
    name = "com_github_davecgh_go_spew",
    build_file_proto_mode = "disable_global",
    importpath = "github.com/davecgh/go-spew",
    sum = "h1:vj9j/u1bqnvCEfJOwUhtlOARqs3+rkHYY13jYWTU97c=",
    version = "v1.1.1",
)

go_repository(
    name = "com_github_pmezard_go_difflib",
    build_file_proto_mode = "disable_global",
    importpath = "github.com/pmezard/go-difflib",
    sum = "h1:4DBwDE0NGyQoBHbLQYPwSUPoCMWR5BEzIk/f1lZbAQM=",
    version = "v1.0.0",
)

go_repository(
    name = "com_github_stretchr_objx",
    build_file_proto_mode = "disable_global",
    importpath = "github.com/stretchr/objx",
    sum = "h1:1zr/of2m5FGMsad5YfcqgdqdWrIhu+EBEJRhR1U7z/c=",
    version = "v0.5.0",
)

go_repository(
    name = "com_github_stretchr_testify",
    build_file_proto_mode = "disable_global",
    importpath = "github.com/stretchr/testify",
    sum = "h1:w7B6lhMri9wdJUVmEZPGGhZzrYTPvgJArz7wNPgYKsk=",
    version = "v1.8.1",
)

go_repository(
    name = "in_gopkg_check_v1",
    build_file_proto_mode = "disable_global",
    importpath = "gopkg.in/check.v1",
    sum = "h1:YR8cESwS4TdDjEe65xsg0ogRM/Nc3DYOhEAlW+xobZo=",
    version = "v1.0.0-20190902080502-41f04d3bba15",
)

go_repository(
    name = "in_gopkg_yaml_v2",
    build_file_proto_mode = "disable_global",
    importpath = "gopkg.in/yaml.v2",
    sum = "h1:D8xgwECY7CYvx+Y2n4sBz93Jn9JRvxdiyyo8CTfuKaY=",
    version = "v2.4.0",
)

go_repository(
    name = "co_honnef_go_tools",
    build_file_proto_mode = "disable_global",
    importpath = "honnef.co/go/tools",
    sum = "h1:3JgtbtFHMiCmsznwGVTUWbgGov+pVqnlf1dEJTNAXeM=",
    version = "v0.0.1-2019.2.3",
)

go_repository(
    name = "com_github_burntsushi_toml",
    build_file_proto_mode = "disable_global",
    importpath = "github.com/BurntSushi/toml",
    sum = "h1:WXkYYl6Yr3qBf1K79EBnL4mak0OimBfB0XUf9Vl28OQ=",
    version = "v0.3.1",
)

go_repository(
    name = "com_github_census_instrumentation_opencensus_proto",
    build_file_proto_mode = "disable_global",
    importpath = "github.com/census-instrumentation/opencensus-proto",
    sum = "h1:iKLQ0xPNFxR/2hzXZMrBo8f1j86j5WHzznCCQxV/b8g=",
    version = "v0.4.1",
)

go_repository(
    name = "com_github_client9_misspell",
    build_file_proto_mode = "disable_global",
    importpath = "github.com/client9/misspell",
    sum = "h1:ta993UF76GwbvJcIo3Y68y/M3WxlpEHPWIGDkJYwzJI=",
    version = "v0.3.4",
)

go_repository(
    name = "com_github_envoyproxy_go_control_plane",
    build_file_proto_mode = "disable_global",
    importpath = "github.com/envoyproxy/go-control-plane",
    sum = "h1:wSUXTlLfiAQRWs2F+p+EKOY9rUyis1MyGqJ2DIk5HpM=",
    version = "v0.11.1",
)

go_repository(
    name = "com_github_envoyproxy_protoc_gen_validate",
    build_file_proto_mode = "disable_global",
    importpath = "github.com/envoyproxy/protoc-gen-validate",
    sum = "h1:QkIBuU5k+x7/QXPvPPnWXWlCdaBFApVqftFV6k087DA=",
    version = "v1.0.2",
)

go_repository(
    name = "com_github_golang_glog",
    build_file_proto_mode = "disable_global",
    importpath = "github.com/golang/glog",
    sum = "h1:/d3pCKDPWNnvIWe0vVUpNP32qc8U3PDVxySP/y360qE=",
    version = "v1.1.0",
)

go_repository(
    name = "com_github_golang_mock",
    build_file_proto_mode = "disable_global",
    importpath = "github.com/golang/mock",
    sum = "h1:G5FRp8JnTd7RQH5kemVNlMeyXQAztQ3mOWV95KxsXH8=",
    version = "v1.1.1",
)

go_repository(
    name = "com_github_golang_protobuf",
    build_file_proto_mode = "disable_global",
    importpath = "github.com/golang/protobuf",
    sum = "h1:KhyjKVUg7Usr/dYsdSqoFveMYd5ko72D+zANwlG1mmg=",
    version = "v1.5.3",
)

go_repository(
    name = "com_github_google_go_cmp",
    build_file_proto_mode = "disable_global",
    importpath = "github.com/google/go-cmp",
    sum = "h1:O2Tfq5qg4qc4AmwVlvv0oLiVAGB7enBSJ2x2DqQFi38=",
    version = "v0.5.9",
)

go_repository(
    name = "com_github_prometheus_client_model",
    build_file_proto_mode = "disable_global",
    importpath = "github.com/prometheus/client_model",
    sum = "h1:gQz4mCbXsO+nc9n1hCxHcGA3Zx3Eo+UHZoInFGUIXNM=",
    version = "v0.0.0-20190812154241-14fe0d1b01d4",
)

go_repository(
    name = "com_google_cloud_go",
    build_file_proto_mode = "disable_global",
    importpath = "cloud.google.com/go",
    sum = "h1:e0WKqKTd5BnrG8aKH3J3h+QvEIQtSUcf2n5UZ5ZgLtQ=",
    version = "v0.26.0",
)

go_repository(
    name = "org_golang_google_appengine",
    build_file_proto_mode = "disable_global",
    importpath = "google.golang.org/appengine",
    sum = "h1:FZR1q0exgwxzPzp/aF+VccGrSfxfPpkBqjIIEq3ru6c=",
    version = "v1.6.7",
)

go_repository(
    name = "org_golang_google_genproto",
    build_file_proto_mode = "disable_global",
    importpath = "google.golang.org/genproto",
    sum = "h1:Z0hjGZePRE0ZBWotvtrwxFNrNE9CUAGtplaDK5NNI/g=",
    version = "v0.0.0-20230711160842-782d3b101e98",
)

go_repository(
    name = "org_golang_google_grpc",
    build_file_proto_mode = "disable",
    importpath = "google.golang.org/grpc",
    sum = "h1:32JY8YpPMSR45K+c3o6b8VL73V+rR8k+DeMIr4vRH8o=",
    version = "v1.58.0",
)

go_repository(
    name = "org_golang_x_crypto",
    build_file_proto_mode = "disable_global",
    importpath = "golang.org/x/crypto",
    sum = "h1:6Ewdq3tDic1mg5xRO4milcWCfMVQhI4NkqWWvqejpuA=",
    version = "v0.11.0",
)

go_repository(
    name = "org_golang_x_exp",
    build_file_proto_mode = "disable_global",
    importpath = "golang.org/x/exp",
    sum = "h1:c2HOrn5iMezYjSlGPncknSEr/8x5LELb/ilJbXi9DEA=",
    version = "v0.0.0-20190121172915-509febef88a4",
)

go_repository(
    name = "org_golang_x_lint",
    build_file_proto_mode = "disable_global",
    importpath = "golang.org/x/lint",
    sum = "h1:5hukYrvBGR8/eNkX5mdUezrA6JiaEZDtJb9Ei+1LlBs=",
    version = "v0.0.0-20190930215403-16217165b5de",
)

go_repository(
    name = "org_golang_x_net",
    build_file_proto_mode = "disable_global",
    importpath = "golang.org/x/net",
    sum = "h1:cfawfvKITfUsFCeJIHJrbSxpeu/E81khclypR0GVT50=",
    version = "v0.12.0",
)

go_repository(
    name = "org_golang_x_oauth2",
    build_file_proto_mode = "disable_global",
    importpath = "golang.org/x/oauth2",
    sum = "h1:zHCpF2Khkwy4mMB4bv0U37YtJdTGW8jI0glAApi0Kh8=",
    version = "v0.10.0",
)

go_repository(
    name = "org_golang_x_sync",
    build_file_proto_mode = "disable_global",
    importpath = "golang.org/x/sync",
    sum = "h1:ftCYgMx6zT/asHUrPw8BLLscYtGznsLAnjq5RH9P66E=",
    version = "v0.3.0",
)

go_repository(
    name = "org_golang_x_sys",
    build_file_proto_mode = "disable_global",
    importpath = "golang.org/x/sys",
    sum = "h1:SqMFp9UcQJZa+pmYuAKjd9xq1f0j5rLcDIk0mj4qAsA=",
    version = "v0.10.0",
)

go_repository(
    name = "org_golang_x_text",
    build_file_proto_mode = "disable_global",
    importpath = "golang.org/x/text",
    sum = "h1:LAntKIrcmeSKERyiOh0XMV39LXS8IE9UL2yP7+f5ij4=",
    version = "v0.11.0",
)

go_repository(
    name = "org_golang_x_tools",
    build_file_proto_mode = "disable_global",
    importpath = "golang.org/x/tools",
    sum = "h1:BOw41kyTf3PuCW1pVQf8+Cyg8pMlkYB1oo9iJ6D/lKM=",
    version = "v0.6.0",
)

go_repository(
    name = "com_github_google_uuid",
    build_file_proto_mode = "disable_global",
    importpath = "github.com/google/uuid",
    sum = "h1:t6JiXgmwXMjEs8VusXIJk2BXHsn+wx8BZdTaoZ5fu7I=",
    version = "v1.3.0",
)

go_repository(
    name = "com_github_squzy_squzy_generated",
    build_file_proto_mode = "disable_global",
    importpath = "github.com/squzy/squzy_generated",
    sum = "h1:iLA2vnDSwXAMorXkzl3WRvDsAGt0wPucMh0J0K8zS+c=",
    version = "v1.14.1-0.20231227142833-f4d58fa8fb6f",
)

go_repository(
    name = "com_github_go_ole_go_ole",
    build_file_proto_mode = "disable_global",
    importpath = "github.com/go-ole/go-ole",
    sum = "h1:/Fpf6oFPoeFik9ty7siob0G6Ke8QvQEuVcuChpwXzpY=",
    version = "v1.2.6",
)

go_repository(
    name = "com_github_shirou_w32",
    importpath = "github.com/shirou/w32",
    sum = "h1:udFKJ0aHUL60LboW/A+DfgoHVedieIzIXE8uylPue0U=",
    version = "v0.0.0-20160930032740-bb4de0191aa4",
)

go_repository(
    name = "com_github_stackexchange_wmi",
    importpath = "github.com/StackExchange/wmi",
    sum = "h1:G0m3OIz70MZUWq3EgK3CesDbo8upS2Vm9/P3FtgI+Jk=",
    version = "v0.0.0-20190523213315-cbe66965904d",
)

go_repository(
    name = "com_github_klauspost_compress",
    build_file_proto_mode = "disable_global",
    importpath = "github.com/klauspost/compress",
    sum = "h1:P76CopJELS0TiO2mebmnzgWaajssP/EszplttgQxcgc=",
    version = "v1.13.6",
)

go_repository(
    name = "com_github_klauspost_cpuid",
    importpath = "github.com/klauspost/cpuid",
    sum = "h1:vJi+O/nMdFt0vqm8NZBI6wzALWdA2X+egi0ogNyrC/w=",
    version = "v1.2.1",
)

go_repository(
    name = "com_github_valyala_bytebufferpool",
    importpath = "github.com/valyala/bytebufferpool",
    sum = "h1:GqA5TC/0021Y/b9FG4Oi9Mr3q7XYx6KllzawFIhcdPw=",
    version = "v1.0.0",
)

go_repository(
    name = "com_github_valyala_fasthttp",
    importpath = "github.com/valyala/fasthttp",
    sum = "h1:uWF8lgKmeaIewWVPwi4GRq2P6+R46IgYZdxWtM+GtEY=",
    version = "v1.6.0",
)

go_repository(
    name = "com_github_valyala_tcplisten",
    importpath = "github.com/valyala/tcplisten",
    sum = "h1:0R4NLDRDZX6JcmhJgXi5E4b8Wg84ihbmUKp/GvSPEzc=",
    version = "v0.0.0-20161114210144-ceec8f93295a",
)

go_repository(
    name = "com_github_tidwall_gjson",
    build_file_proto_mode = "disable_global",
    importpath = "github.com/tidwall/gjson",
    sum = "h1:3TFY9yxOQShrvmjdM76K+jc66zJeT6D3/VFFYCGQf7M=",
    version = "v1.13.0",
)

go_repository(
    name = "com_github_tidwall_match",
    build_file_proto_mode = "disable_global",
    importpath = "github.com/tidwall/match",
    sum = "h1:+Ho715JplO36QYgwN9PGYNhgZvoUSc9X2c80KVTi+GA=",
    version = "v1.1.1",
)

go_repository(
    name = "com_github_tidwall_pretty",
    build_file_proto_mode = "disable_global",
    importpath = "github.com/tidwall/pretty",
    sum = "h1:RWIZEg2iJ8/g6fDDYzMpobmaoGh5OLl4AXtGUGPcqCs=",
    version = "v1.2.0",
)

go_repository(
    name = "com_github_squzy_squzy_proto",
    importpath = "github.com/squzy/squzy_proto",
    sum = "h1:6WkBEez7NLD474ojNd3f0RTg3pPZkdxOUVy4KUHBXgw=",
    version = "v1.14.1-0.20230108131950-dfd182da38bc",
)

go_repository(
    name = "com_github_go_stack_stack",
    build_file_proto_mode = "disable_global",
    importpath = "github.com/go-stack/stack",
    sum = "h1:5SgMzNM5HxrEjV0ww2lTmX6E2Izsfxas4+YHWRs3Lsk=",
    version = "v1.8.0",
)

go_repository(
    name = "com_github_gobuffalo_attrs",
    build_file_proto_mode = "disable_global",
    importpath = "github.com/gobuffalo/attrs",
    sum = "h1:hSkbZ9XSyjyBirMeqSqUrK+9HboWrweVlzRNqoBi2d4=",
    version = "v0.0.0-20190224210810-a9411de4debd",
)

go_repository(
    name = "com_github_gobuffalo_depgen",
    build_file_proto_mode = "disable_global",
    importpath = "github.com/gobuffalo/depgen",
    sum = "h1:31atYa/UW9V5q8vMJ+W6wd64OaaTHUrCUXER358zLM4=",
    version = "v0.1.0",
)

go_repository(
    name = "com_github_gobuffalo_envy",
    build_file_proto_mode = "disable_global",
    importpath = "github.com/gobuffalo/envy",
    sum = "h1:GlXgaiBkmrYMHco6t4j7SacKO4XUjvh5pwXh0f4uxXU=",
    version = "v1.7.0",
)

go_repository(
    name = "com_github_gobuffalo_flect",
    build_file_proto_mode = "disable_global",
    importpath = "github.com/gobuffalo/flect",
    sum = "h1:3GQ53z7E3o00C/yy7Ko8VXqQXoJGLkrTQCLTF1EjoXU=",
    version = "v0.1.3",
)

go_repository(
    name = "com_github_gobuffalo_genny",
    build_file_proto_mode = "disable_global",
    importpath = "github.com/gobuffalo/genny",
    sum = "h1:iQ0D6SpNXIxu52WESsD+KoQ7af2e3nCfnSBoSF/hKe0=",
    version = "v0.1.1",
)

go_repository(
    name = "com_github_gobuffalo_gitgen",
    build_file_proto_mode = "disable_global",
    importpath = "github.com/gobuffalo/gitgen",
    sum = "h1:mSVZ4vj4khv+oThUfS+SQU3UuFIZ5Zo6UNcvK8E8Mz8=",
    version = "v0.0.0-20190315122116-cc086187d211",
)

go_repository(
    name = "com_github_gobuffalo_gogen",
    build_file_proto_mode = "disable_global",
    importpath = "github.com/gobuffalo/gogen",
    sum = "h1:dLg+zb+uOyd/mKeQUYIbwbNmfRsr9hd/WtYWepmayhI=",
    version = "v0.1.1",
)

go_repository(
    name = "com_github_gobuffalo_logger",
    build_file_proto_mode = "disable_global",
    importpath = "github.com/gobuffalo/logger",
    sum = "h1:8thhT+kUJMTMy3HlX4+y9Da+BNJck+p109tqqKp7WDs=",
    version = "v0.0.0-20190315122211-86e12af44bc2",
)

go_repository(
    name = "com_github_gobuffalo_mapi",
    build_file_proto_mode = "disable_global",
    importpath = "github.com/gobuffalo/mapi",
    sum = "h1:fq9WcL1BYrm36SzK6+aAnZ8hcp+SrmnDyAxhNx8dvJk=",
    version = "v1.0.2",
)

go_repository(
    name = "com_github_gobuffalo_packd",
    build_file_proto_mode = "disable_global",
    importpath = "github.com/gobuffalo/packd",
    sum = "h1:4sGKOD8yaYJ+dek1FDkwcxCHA40M4kfKgFHx8N2kwbU=",
    version = "v0.1.0",
)

go_repository(
    name = "com_github_gobuffalo_packr_v2",
    build_file_proto_mode = "disable_global",
    importpath = "github.com/gobuffalo/packr/v2",
    sum = "h1:Ir9W9XIm9j7bhhkKE9cokvtTl1vBm62A/fene/ZCj6A=",
    version = "v2.2.0",
)

go_repository(
    name = "com_github_gobuffalo_syncx",
    build_file_proto_mode = "disable_global",
    importpath = "github.com/gobuffalo/syncx",
    sum = "h1:tpom+2CJmpzAWj5/VEHync2rJGi+epHNIeRSWjzGA+4=",
    version = "v0.0.0-20190224160051-33c29581e754",
)

go_repository(
    name = "com_github_golang_snappy",
    build_file_proto_mode = "disable_global",
    importpath = "github.com/golang/snappy",
    sum = "h1:fHPg5GQYlCeLIPB9BZqMVR5nR9A+IM5zcgeTdjMYmLA=",
    version = "v0.0.3",
)

go_repository(
    name = "com_github_inconshreveable_mousetrap",
    build_file_proto_mode = "disable_global",
    importpath = "github.com/inconshreveable/mousetrap",
    sum = "h1:Z8tu5sraLXCXIcARxBp/8cbvlwVa7Z1NHg9XEKhtSvM=",
    version = "v1.0.0",
)

go_repository(
    name = "com_github_joho_godotenv",
    build_file_proto_mode = "disable_global",
    importpath = "github.com/joho/godotenv",
    sum = "h1:Zjp+RcGpHhGlrMbJzXTrZZPrWj+1vfm90La1wgB6Bhc=",
    version = "v1.3.0",
)

go_repository(
    name = "com_github_karrick_godirwalk",
    build_file_proto_mode = "disable_global",
    importpath = "github.com/karrick/godirwalk",
    sum = "h1:lOpSw2vJP0y5eLBW906QwKsUK/fe/QDyoqM5rnnuPDY=",
    version = "v1.10.3",
)

go_repository(
    name = "com_github_kisielk_errcheck",
    build_file_proto_mode = "disable_global",
    importpath = "github.com/kisielk/errcheck",
    sum = "h1:e8esj/e4R+SAOwFwN+n3zr0nYeCyeweozKfO23MvHzY=",
    version = "v1.5.0",
)

go_repository(
    name = "com_github_konsorten_go_windows_terminal_sequences",
    build_file_proto_mode = "disable_global",
    importpath = "github.com/konsorten/go-windows-terminal-sequences",
    sum = "h1:DB17ag19krx9CFsz4o3enTrPXyIXCl+2iCXH/aMAp9s=",
    version = "v1.0.2",
)

go_repository(
    name = "com_github_kr_pretty",
    build_file_proto_mode = "disable_global",
    importpath = "github.com/kr/pretty",
    sum = "h1:Fmg33tUaq4/8ym9TJN1x7sLJnHVwhP33CNkpYV/7rwI=",
    version = "v0.2.1",
)

go_repository(
    name = "com_github_kr_pty",
    build_file_proto_mode = "disable_global",
    importpath = "github.com/kr/pty",
    sum = "h1:VkoXIwSboBpnk99O/KFauAEILuNHv5DVFKZMBN/gUgw=",
    version = "v1.1.1",
)

go_repository(
    name = "com_github_kr_text",
    build_file_proto_mode = "disable_global",
    importpath = "github.com/kr/text",
    sum = "h1:45sCR5RtlFHMR4UwH9sdQ5TC8v0qDQCHnXt+kaKSTVE=",
    version = "v0.1.0",
)

go_repository(
    name = "com_github_markbates_oncer",
    build_file_proto_mode = "disable_global",
    importpath = "github.com/markbates/oncer",
    sum = "h1:JgVTCPf0uBVcUSWpyXmGpgOc62nK5HWUBKAGc3Qqa5k=",
    version = "v0.0.0-20181203154359-bf2de49a0be2",
)

go_repository(
    name = "com_github_markbates_safe",
    build_file_proto_mode = "disable_global",
    importpath = "github.com/markbates/safe",
    sum = "h1:yjZkbvRM6IzKj9tlu/zMJLS0n/V351OZWRnF3QfaUxI=",
    version = "v1.0.1",
)

go_repository(
    name = "com_github_montanaflynn_stats",
    build_file_proto_mode = "disable_global",
    importpath = "github.com/montanaflynn/stats",
    sum = "h1:iruDEfMl2E6fbMZ9s0scYfZQ84/6SPL6zC8ACM2oIL0=",
    version = "v0.0.0-20171201202039-1bf9dbcd8cbe",
)

go_repository(
    name = "com_github_pelletier_go_toml",
    build_file_proto_mode = "disable_global",
    importpath = "github.com/pelletier/go-toml",
    sum = "h1:u3Z1r+oOXJIkxqw34zVhyPgjBsm6X2wn21NWs/HfSeg=",
    version = "v1.4.0",
)

go_repository(
    name = "com_github_pkg_errors",
    build_file_proto_mode = "disable_global",
    importpath = "github.com/pkg/errors",
    sum = "h1:FEBLx1zS214owpjy7qsBeixbURkuhQAwrK5UwLGTwt4=",
    version = "v0.9.1",
)

go_repository(
    name = "com_github_rogpeppe_go_internal",
    build_file_proto_mode = "disable_global",
    importpath = "github.com/rogpeppe/go-internal",
    sum = "h1:RR9dF3JtopPvtkroDZuVD7qquD0bnHlKSqaQhgwt8yk=",
    version = "v1.3.0",
)

go_repository(
    name = "com_github_sirupsen_logrus",
    build_file_proto_mode = "disable_global",
    importpath = "github.com/sirupsen/logrus",
    sum = "h1:dJKuHgqk1NNQlqoA6BTlM1Wf9DOH3NBjQyu0h9+AZZE=",
    version = "v1.8.1",
)

go_repository(
    name = "com_github_spf13_cobra",
    build_file_proto_mode = "disable_global",
    importpath = "github.com/spf13/cobra",
    sum = "h1:6m/oheQuQ13N9ks4hubMG6BnvwOeaJrqSPLahSnczz8=",
    version = "v1.0.0",
)

go_repository(
    name = "com_github_spf13_pflag",
    build_file_proto_mode = "disable_global",
    importpath = "github.com/spf13/pflag",
    sum = "h1:zPAT6CGy6wXeQ7NtTnaTerfKOsV6V6F8agHXFiazDkg=",
    version = "v1.0.3",
)

go_repository(
    name = "com_github_xdg_scram",
    build_file_proto_mode = "disable_global",
    importpath = "github.com/xdg/scram",
    sum = "h1:u40Z8hqBAAQyv+vATcGgV0YCnDjqSL7/q/JyPhhJSPk=",
    version = "v0.0.0-20180814205039-7eeb5667e42c",
)

go_repository(
    name = "com_github_xdg_stringprep",
    build_file_proto_mode = "disable_global",
    importpath = "github.com/xdg/stringprep",
    sum = "h1:n+nNi93yXLkJvKwXNP9d55HC7lGK4H/SRcwB5IaUZLo=",
    version = "v0.0.0-20180714160509-73f8eece6fdc",
)

go_repository(
    name = "in_gopkg_errgo_v2",
    build_file_proto_mode = "disable_global",
    importpath = "gopkg.in/errgo.v2",
    sum = "h1:0vLT13EuvQ0hNvakwLuFZ/jYrLp5F3kcWHXdRggjCE8=",
    version = "v2.1.0",
)

go_repository(
    name = "org_mongodb_go_mongo_driver",
    build_file_proto_mode = "disable_global",
    importpath = "go.mongodb.org/mongo-driver",
    sum = "h1:8ssUXufb90ujcIvR6MyE1SchaNj0SFxsakiZgxIyrMk=",
    version = "v1.8.2",
)

go_repository(
    name = "com_github_go_kit_kit",
    build_file_proto_mode = "disable_global",
    importpath = "github.com/go-kit/kit",
    sum = "h1:wDJmvq38kDhkVxi50ni9ykkdUr1PKgqKOoi01fa0Mdk=",
    version = "v0.9.0",
)

go_repository(
    name = "com_github_go_logfmt_logfmt",
    build_file_proto_mode = "disable_global",
    importpath = "github.com/go-logfmt/logfmt",
    sum = "h1:MP4Eh7ZCb31lleYCFuwm0oe4/YGak+5l1vA2NOE80nA=",
    version = "v0.4.0",
)

go_repository(
    name = "com_github_gogo_protobuf",
    build_file_proto_mode = "disable_global",
    importpath = "github.com/gogo/protobuf",
    sum = "h1:Ov1cvc58UF3b5XjBnZv7+opcTcQFZebYjWzi34vdm4Q=",
    version = "v1.3.2",
)

go_repository(
    name = "com_github_grpc_ecosystem_go_grpc_middleware",
    build_file_proto_mode = "disable_global",
    importpath = "github.com/grpc-ecosystem/go-grpc-middleware",
    sum = "h1:+9834+KizmvFV7pXQGSXQTsaWhq2GjuNUt0aUU0YBYw=",
    version = "v1.3.0",
)

go_repository(
    name = "com_github_kisielk_gotool",
    build_file_proto_mode = "disable_global",
    importpath = "github.com/kisielk/gotool",
    sum = "h1:AV2c/EiW3KqPNT9ZKl07ehoAGi4C5/01Cfbblndcapg=",
    version = "v1.0.0",
)

go_repository(
    name = "com_github_kr_logfmt",
    build_file_proto_mode = "disable_global",
    importpath = "github.com/kr/logfmt",
    sum = "h1:T+h1c/A9Gawja4Y9mFVWj2vyii2bbUNDw3kt9VxK2EY=",
    version = "v0.0.0-20140226030751-b84e30acd515",
)

go_repository(
    name = "com_github_opentracing_opentracing_go",
    build_file_proto_mode = "disable_global",
    importpath = "github.com/opentracing/opentracing-go",
    sum = "h1:pWlfV3Bxv7k65HYwkikxat0+s3pV4bsqf19k25Ur8rU=",
    version = "v1.1.0",
)

go_repository(
    name = "org_uber_go_atomic",
    build_file_proto_mode = "disable_global",
    importpath = "go.uber.org/atomic",
    sum = "h1:Ezj3JGmsOnG1MoRWQkPBsKLe9DwWD9QeXzTRzzldNVk=",
    version = "v1.6.0",
)

go_repository(
    name = "org_uber_go_multierr",
    build_file_proto_mode = "disable_global",
    importpath = "go.uber.org/multierr",
    sum = "h1:KCa4XfM8CWFCpxXRGok+Q0SS/0XBhMDbHHGABQLvD2A=",
    version = "v1.5.0",
)

go_repository(
    name = "org_uber_go_zap",
    build_file_proto_mode = "disable_global",
    importpath = "go.uber.org/zap",
    sum = "h1:nR6NoDBgAf67s68NhaXbsojM+2gxp3S1hWkHDl27pVU=",
    version = "v1.13.0",
)

go_repository(
    name = "com_github_google_renameio",
    importpath = "github.com/google/renameio",
    sum = "h1:GOZbcHa3HfsPKPlmyPyN2KEohoMXOhdMbHrvbpl2QaA=",
    version = "v0.1.0",
)

go_repository(
    name = "com_github_squzy_mongo_helper",
    build_file_proto_mode = "disable_global",
    importpath = "github.com/squzy/mongo_helper",
    sum = "h1:NL6hCxYRb9Y24COaTkwklWcQjZTlsPoGKw/yd60KFiM=",
    version = "v0.0.0-20200713232419-037a870c9d06",
)

go_repository(
    name = "org_golang_x_mod",
    importpath = "golang.org/x/mod",
    sum = "h1:LUYupSeNrTNCGzR/hVBk2NHZO4hXcVaW1k4Qx7rjPx8=",
    version = "v0.8.0",
)

go_repository(
    name = "org_golang_x_xerrors",
    build_file_proto_mode = "disable_global",
    importpath = "golang.org/x/xerrors",
    sum = "h1:go1bK/D/BFZV2I8cIQd1NKEZ+0owSTG1fDTci4IqFcE=",
    version = "v0.0.0-20200804184101-5ec99f83aff1",
)

go_repository(
    name = "org_uber_go_tools",
    importpath = "go.uber.org/tools",
    sum = "h1:0mgffUl7nfd+FpvXMVz4IDEaUSmT1ysygQC7qYo7sG4=",
    version = "v0.0.0-20190618225709-2cfd321de3ee",
)

go_repository(
    name = "com_github_cncf_udpa_go",
    build_file_proto_mode = "disable_global",
    importpath = "github.com/cncf/udpa/go",
    sum = "h1:QQ3GSy+MqSHxm/d8nCtnAiZdYFd45cYZPs8vOOIYKfk=",
    version = "v0.0.0-20220112060539-c52dc94e7fbe",
)

go_repository(
    name = "com_github_data_dog_go_sqlmock",
    build_file_proto_mode = "disable_global",
    importpath = "github.com/DATA-DOG/go-sqlmock",
    sum = "h1:ThlnYciV1iM/V0OSF/dtkqWb6xo5qITT1TJBG1MRDJM=",
    version = "v1.4.1",
)

go_repository(
    name = "com_github_denisenkom_go_mssqldb",
    build_file_proto_mode = "disable_global",
    importpath = "github.com/denisenkom/go-mssqldb",
    sum = "h1:83Wprp6ROGeiHFAP8WJdI2RoxALQYgdllERc3N5N2DM=",
    version = "v0.0.0-20191124224453-732737034ffd",
)

go_repository(
    name = "com_github_erikstmartin_go_testdb",
    build_file_proto_mode = "disable_global",
    importpath = "github.com/erikstmartin/go-testdb",
    sum = "h1:Yzb9+7DPaBjB8zlTR87/ElzFsnQfuHnVUVqpZZIcV5Y=",
    version = "v0.0.0-20160219214506-8d10e4a1bae5",
)

go_repository(
    name = "com_github_gin_contrib_sse",
    build_file_proto_mode = "disable_global",
    importpath = "github.com/gin-contrib/sse",
    sum = "h1:Y/yl/+YNO8GZSjAhjMsSuLt29uWRFHdHYUb5lYOV9qE=",
    version = "v0.1.0",
)

go_repository(
    name = "com_github_gin_gonic_gin",
    build_file_proto_mode = "disable_global",
    importpath = "github.com/gin-gonic/gin",
    sum = "h1:3DoBmSbJbZAWqXJC3SLjAPfutPJJRN1U5pALB7EeTTs=",
    version = "v1.7.7",
)

go_repository(
    name = "com_github_go_playground_assert_v2",
    build_file_proto_mode = "disable_global",
    importpath = "github.com/go-playground/assert/v2",
    sum = "h1:MsBgLAaY856+nPRTKrp3/OZK38U/wa0CcBYNjji3q3A=",
    version = "v2.0.1",
)

go_repository(
    name = "com_github_go_playground_locales",
    build_file_proto_mode = "disable_global",
    importpath = "github.com/go-playground/locales",
    sum = "h1:HyWk6mgj5qFqCT5fjGBuRArbVDfE4hi8+e8ceBS/t7Q=",
    version = "v0.13.0",
)

go_repository(
    name = "com_github_go_playground_universal_translator",
    build_file_proto_mode = "disable_global",
    importpath = "github.com/go-playground/universal-translator",
    sum = "h1:icxd5fm+REJzpZx7ZfpaD876Lmtgy7VtROAbHHXk8no=",
    version = "v0.17.0",
)

go_repository(
    name = "com_github_go_playground_validator_v10",
    build_file_proto_mode = "disable_global",
    importpath = "github.com/go-playground/validator/v10",
    sum = "h1:pH2c5ADXtd66mxoE0Zm9SUhxE20r7aM3F26W0hOn+GE=",
    version = "v10.4.1",
)

go_repository(
    name = "com_github_go_sql_driver_mysql",
    build_file_proto_mode = "disable_global",
    importpath = "github.com/go-sql-driver/mysql",
    sum = "h1:g24URVg0OFbNUTx9qqY1IRZ9D9z3iPyi5zKhQZpNwpA=",
    version = "v1.4.1",
)

go_repository(
    name = "com_github_golang_sql_civil",
    build_file_proto_mode = "disable_global",
    importpath = "github.com/golang-sql/civil",
    sum = "h1:lXe2qZdvpiX5WZkZR4hgp4KJVfY3nMkvmwbVkpv1rVY=",
    version = "v0.0.0-20190719163853-cb61b32ac6fe",
)

go_repository(
    name = "com_github_google_gofuzz",
    build_file_proto_mode = "disable_global",
    importpath = "github.com/google/gofuzz",
    sum = "h1:A8PeW59pxE9IoFRqBp37U+mSNaQoZ46F1f0f863XSXw=",
    version = "v1.0.0",
)

go_repository(
    name = "com_github_jinzhu_inflection",
    build_file_proto_mode = "disable_global",
    importpath = "github.com/jinzhu/inflection",
    sum = "h1:K317FqzuhWc8YvSVlFMCCUb36O/S9MCKRDI7QkRKD/E=",
    version = "v1.0.0",
)

go_repository(
    name = "com_github_jinzhu_now",
    build_file_proto_mode = "disable_global",
    importpath = "github.com/jinzhu/now",
    sum = "h1:tHnRBy1i5F2Dh8BAFxqFzxKqqvezXrL2OW1TnX+Mlas=",
    version = "v1.1.4",
)

go_repository(
    name = "com_github_json_iterator_go",
    build_file_proto_mode = "disable_global",
    importpath = "github.com/json-iterator/go",
    sum = "h1:Kz6Cvnvv2wGdaG/V8yMvfkmNiXq9Ya2KUv4rouJJr68=",
    version = "v1.1.10",
)

go_repository(
    name = "com_github_kazegusuri_grpc_panic_handler",
    importpath = "github.com/kazegusuri/grpc-panic-handler",
    sum = "h1:oW3n7kE84CWfrnc9rcK3mBy3XtSLy2VNuI4pQFD+IKc=",
    version = "v0.0.0-20160502122501-093ec776affc",
)

go_repository(
    name = "com_github_leodido_go_urn",
    build_file_proto_mode = "disable_global",
    importpath = "github.com/leodido/go-urn",
    sum = "h1:hpXL4XnriNwQ/ABnpepYM/1vCLWNDfUNts8dX3xTG6Y=",
    version = "v1.2.0",
)

go_repository(
    name = "com_github_lib_pq",
    build_file_proto_mode = "disable_global",
    importpath = "github.com/lib/pq",
    sum = "h1:AqzbZs4ZoCBp+GtejcpCpcxM3zlSMx29dXbUSeVtJb8=",
    version = "v1.10.2",
)

go_repository(
    name = "com_github_mattn_go_isatty",
    build_file_proto_mode = "disable_global",
    importpath = "github.com/mattn/go-isatty",
    sum = "h1:wuysRhFDzyxgEmMf5xjvJ2M9dZoWAXNNr5LSBS7uHXY=",
    version = "v0.0.12",
)

go_repository(
    name = "com_github_mattn_go_sqlite3",
    build_file_proto_mode = "disable_global",
    importpath = "github.com/mattn/go-sqlite3",
    sum = "h1:xQ15muvnzGBHpIpdrNi1DA5x0+TcBZzsIDwmw9uTHzw=",
    version = "v2.0.1+incompatible",
)

go_repository(
    name = "com_github_modern_go_concurrent",
    build_file_proto_mode = "disable_global",
    importpath = "github.com/modern-go/concurrent",
    sum = "h1:TRLaZ9cD/w8PVh93nsPXa1VrQ6jlwL5oN8l14QlcNfg=",
    version = "v0.0.0-20180306012644-bacd9c7ef1dd",
)

go_repository(
    name = "com_github_modern_go_reflect2",
    build_file_proto_mode = "disable_global",
    importpath = "github.com/modern-go/reflect2",
    sum = "h1:9f412s+6RmYXLWZSEzVVgPGK7C2PphHj5RJrvfx9AWI=",
    version = "v1.0.1",
)

go_repository(
    name = "com_github_ugorji_go",
    build_file_proto_mode = "disable_global",
    importpath = "github.com/ugorji/go",
    sum = "h1:/68gy2h+1mWMrwZFeD1kQialdSzAb432dtpeJ42ovdo=",
    version = "v1.1.7",
)

go_repository(
    name = "com_github_ugorji_go_codec",
    build_file_proto_mode = "disable_global",
    importpath = "github.com/ugorji/go/codec",
    sum = "h1:2SvQaVZ1ouYrrKKwoSk2pzd4A9evlKJb9oTL+OaLUSs=",
    version = "v1.1.7",
)

go_repository(
    name = "org_golang_google_protobuf",
    build_file_proto_mode = "disable_global",
    importpath = "google.golang.org/protobuf",
    sum = "h1:g0LDEJHgrBl9N9r17Ru3sqWhkIx2NB67okBHPwC7hs8=",
    version = "v1.31.0",
)

go_repository(
    name = "com_github_antonmedv_expr",
    build_file_proto_mode = "disable_global",
    importpath = "github.com/antonmedv/expr",
    sum = "h1:uVwIkIBNO2yn4vY2u2DQUqXTmv9jEEMCEcHa19G5weY=",
    version = "v1.8.8",
)

go_repository(
    name = "com_github_araddon_dateparse",
    build_file_proto_mode = "disable_global",
    importpath = "github.com/araddon/dateparse",
    sum = "h1:TEBmxO80TM04L8IuMWk77SGL1HomBmKTdzdJLLWznxI=",
    version = "v0.0.0-20200409225146-d820a6159ab1",
)

go_repository(
    name = "com_github_gdamore_encoding",
    build_file_proto_mode = "disable_global",
    importpath = "github.com/gdamore/encoding",
    sum = "h1:+7OoQ1Bc6eTm5niUzBa0Ctsh6JbMW6Ra+YNuAtDBdko=",
    version = "v1.0.0",
)

go_repository(
    name = "com_github_gdamore_tcell",
    build_file_proto_mode = "disable_global",
    importpath = "github.com/gdamore/tcell",
    sum = "h1:r35w0JBADPZCVQijYebl6YMWWtHRqVEGt7kL2eBADRM=",
    version = "v1.3.0",
)

go_repository(
    name = "com_github_lucasb_eyer_go_colorful",
    build_file_proto_mode = "disable_global",
    importpath = "github.com/lucasb-eyer/go-colorful",
    sum = "h1:QIbQXiugsb+q10B+MI+7DI1oQLdmnep86tWFlaaUAac=",
    version = "v1.0.3",
)

go_repository(
    name = "com_github_mattn_go_runewidth",
    build_file_proto_mode = "disable_global",
    importpath = "github.com/mattn/go-runewidth",
    sum = "h1:3tS41NlGYSmhhe/8fhGRzc+z3AYCw1Fe1WAyLuujKs0=",
    version = "v0.0.8",
)

go_repository(
    name = "com_github_rivo_tview",
    build_file_proto_mode = "disable_global",
    importpath = "github.com/rivo/tview",
    sum = "h1:4CFNy7/q7P06AsIONZzuWy7jcdqEmYQvOZ9FAFZdbls=",
    version = "v0.0.0-20200219210816-cd38d7432498",
)

go_repository(
    name = "com_github_rivo_uniseg",
    build_file_proto_mode = "disable_global",
    importpath = "github.com/rivo/uniseg",
    sum = "h1:+2KBaVoUmb9XzDsrx/Ct0W/EYOSFf/nWTauy++DprtY=",
    version = "v0.1.0",
)

go_repository(
    name = "com_github_sanity_io_litter",
    build_file_proto_mode = "disable_global",
    importpath = "github.com/sanity-io/litter",
    sum = "h1:DGJO0bxH/+C2EukzOSBmAlxmkhVMGqzvcx/rvySYw9M=",
    version = "v1.2.0",
)

go_repository(
    name = "com_github_go_test_deep",
    build_file_proto_mode = "disable_global",
    importpath = "github.com/go-test/deep",
    sum = "h1:u2CU3YKy9I2pmu9pX0eq50wCgjfGIt539SqR7FbHiho=",
    version = "v1.0.4",
)

go_repository(
    name = "com_github_gorilla_websocket",
    build_file_proto_mode = "disable_global",
    importpath = "github.com/gorilla/websocket",
    sum = "h1:+/TMaTYc4QFitKJxsQ7Yye35DkWvkdLcvGKqM+x0Ufc=",
    version = "v1.4.2",
)

go_repository(
    name = "com_github_slack_go_slack",
    build_file_proto_mode = "disable_global",
    importpath = "github.com/slack-go/slack",
    sum = "h1:IkDKtJ2IROJNoe3d6mW870/NRKvq2fhLB/Q5XmzWk00=",
    version = "v0.6.5",
)
