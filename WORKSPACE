workspace(name = "squzy")

load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")
load("@bazel_tools//tools/build_defs/repo:git.bzl", "git_repository")

http_archive(
    name = "com_google_protobuf",
    strip_prefix = "protobuf-3.19.3",
    urls = [
        "https://mirror.bazel.build/github.com/protocolbuffers/protobuf/archive/v3.19.3.tar.gz",
        "https://github.com/protocolbuffers/protobuf/archive/v3.19.3.tar.gz",
    ],
)

load("@com_google_protobuf//:protobuf_deps.bzl", "protobuf_deps")

protobuf_deps()

http_archive(
    name = "io_bazel_rules_go",
    sha256 = "2b1641428dff9018f9e85c0384f03ec6c10660d935b750e3fa1492a281a53b0f",
    urls = [
        "https://mirror.bazel.build/github.com/bazelbuild/rules_go/releases/download/v0.29.0/rules_go-v0.29.0.zip",
        "https://github.com/bazelbuild/rules_go/releases/download/v0.29.0/rules_go-v0.29.0.zip",
    ],
)

load("@io_bazel_rules_go//go:deps.bzl", "go_register_toolchains", "go_rules_dependencies")

go_rules_dependencies()

go_register_toolchains(version = "1.17.1")

http_archive(
    name = "bazel_gazelle",
    sha256 = "de69a09dc70417580aabf20a28619bb3ef60d038470c7cf8442fafcf627c21cb",
    urls = [
        "https://mirror.bazel.build/github.com/bazelbuild/bazel-gazelle/releases/download/v0.24.0/bazel-gazelle-v0.24.0.tar.gz",
        "https://github.com/bazelbuild/bazel-gazelle/releases/download/v0.24.0/bazel-gazelle-v0.24.0.tar.gz",
    ],
)

load("@bazel_gazelle//:deps.bzl", "gazelle_dependencies", "go_repository")

gazelle_dependencies()

load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")

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

load("@bazel_gazelle//:deps.bzl", "gazelle_dependencies")

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
    sum = "h1:dUUwHk2QECo/6vqA44rthZ8ie2QXMNeKRTHCNY2nXvo=",
    version = "v3.0.0-20200313102051-9f266ea9e77c",
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
    sum = "h1:6MnRN8NT7+YBpUIWxHtefFZOKTAPgGjpQSxqLNn0+qY=",
    version = "v2.1.1",
)

go_repository(
    name = "com_github_cncf_xds_go",
    importpath = "github.com/cncf/xds/go",
    sum = "h1:zH8ljVhhq7yC0MIeUL/IviMtY8hx2mK8cN9wEYb8ggw=",
    version = "v0.0.0-20211011173535-cb28da3451f1",
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
    sum = "h1:gmcG1KaJ57LophUzW0Hy8NmPhnMZb4M0+kPpLofRdBo=",
    version = "v1.16.0",
)

go_repository(
    name = "com_github_rogpeppe_fastuuid",
    importpath = "github.com/rogpeppe/fastuuid",
    sum = "h1:Ppwyp6VYCF1nvBTXL3trRso7mXMlRrw9ooo375wvi2s=",
    version = "v1.2.0",
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
    sum = "h1:ruQGxdhGHe7FWOJPT0mKs5+pD2Xs1Bm/kdGlHO04FmM=",
    version = "v1.2.1",
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
    sum = "h1:v+OssWQX+hTHEmOBgwxdZxK4zHq3yOs8F9J7mk0PY8E=",
    version = "v0.0.0-20201126162022-7de9c90e9dd1",
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
    sum = "h1:2vfRuCMp5sSVIDSqO8oNnWJq7mPa6KVP3iPIwFBuy8A=",
    version = "v0.1.1",
)

go_repository(
    name = "com_github_stretchr_testify",
    build_file_proto_mode = "disable_global",
    importpath = "github.com/stretchr/testify",
    sum = "h1:nwc3DEeHmmLAfoZucVR881uASk0Mfjw8xYJ99tb5CcY=",
    version = "v1.7.0",
)

go_repository(
    name = "in_gopkg_check_v1",
    build_file_proto_mode = "disable_global",
    importpath = "gopkg.in/check.v1",
    sum = "h1:qIbj1fsPNlZgppZ+VLlY7N33q108Sa+fhmuc+sWQYwY=",
    version = "v1.0.0-20180628173108-788fd7840127",
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
    sum = "h1:/hemPrYIhOhy8zYrNj+069zDB68us2sMGsfkFJO0iZs=",
    version = "v0.0.0-20190523083050-ea95bdfd59fc",
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
    sum = "h1:glEXhBS5PSLLv4IXzLA5yPRVX4bilULVyxxbrfOtDAk=",
    version = "v0.2.1",
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
    sum = "h1:fP+fF0up6oPY49OrjPrhIJ8yQfdIM85NXMLkMg1EXVs=",
    version = "v0.9.10-0.20210907150352-cf90f659a021",
)

go_repository(
    name = "com_github_envoyproxy_protoc_gen_validate",
    build_file_proto_mode = "disable_global",
    importpath = "github.com/envoyproxy/protoc-gen-validate",
    sum = "h1:EQciDnbrYxy13PgWoY8AqoxGiPrpgBZ1R8UNe3ddc+A=",
    version = "v0.1.0",
)

go_repository(
    name = "com_github_golang_glog",
    build_file_proto_mode = "disable_global",
    importpath = "github.com/golang/glog",
    sum = "h1:VKtxabqXZkF25pY9ekfRL6a582T4P37/31XEstQ5p58=",
    version = "v0.0.0-20160126235308-23def4e6c14b",
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
    sum = "h1:ROPKBNFfQgOUMifHyP+KYbvpjbdoFNs+aK7DXlji0Tw=",
    version = "v1.5.2",
)

go_repository(
    name = "com_github_google_go_cmp",
    build_file_proto_mode = "disable_global",
    importpath = "github.com/google/go-cmp",
    sum = "h1:BKbKCqvP6I+rmFHt06ZmyQtvB8xAkWdhFyr0ZUNZcxQ=",
    version = "v0.5.6",
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
    sum = "h1:eOI3/cP2VTU6uZLDYAoic+eyzzB9YyGmJ7eIjl8rOPg=",
    version = "v0.34.0",
)

go_repository(
    name = "org_golang_google_appengine",
    build_file_proto_mode = "disable_global",
    importpath = "google.golang.org/appengine",
    sum = "h1:lMO5rYAqUxkmaj76jAkRUvt5JZgFymx/+Q5Mzfivuhc=",
    version = "v1.6.6",
)

go_repository(
    name = "org_golang_google_genproto",
    build_file_proto_mode = "disable_global",
    importpath = "google.golang.org/genproto",
    sum = "h1:aCsSLXylHWFno0r4S3joLpiaWayvqd2Mn4iSvx4WZZc=",
    version = "v0.0.0-20220114231437-d2e6a121cae0",
)

go_repository(
    name = "org_golang_google_grpc",
    build_file_proto_mode = "disable",
    importpath = "google.golang.org/grpc",
    sum = "h1:Eeu7bZtDZ2DpRCsLhUlcrLnvYaMK1Gz86a+hMVvELmM=",
    version = "v1.43.0",
)

go_repository(
    name = "org_golang_x_crypto",
    build_file_proto_mode = "disable_global",
    importpath = "golang.org/x/crypto",
    sum = "h1:psW17arqaxU48Z5kZ0CQnkZWQJsqcURM6tKiBApRjXI=",
    version = "v0.0.0-20200622213623-75b288015ac9",
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
    sum = "h1:XQyxROzUlZH+WIQwySDgnISgOivlhjIEwaQaJEJrrN0=",
    version = "v0.0.0-20190313153728-d0100b6bd8b3",
)

go_repository(
    name = "org_golang_x_net",
    build_file_proto_mode = "disable_global",
    importpath = "golang.org/x/net",
    sum = "h1:4nGaVu0QrbjT/AK2PRLuQfQuh6DJve+pELhqTdAj3x0=",
    version = "v0.0.0-20210405180319-a5a99cb37ef4",
)

go_repository(
    name = "org_golang_x_oauth2",
    build_file_proto_mode = "disable_global",
    importpath = "golang.org/x/oauth2",
    sum = "h1:TzXSXBo42m9gQenoE3b9BGiEpg5IG2JkU5FkPIawgtw=",
    version = "v0.0.0-20200107190931-bf48bf16ab8d",
)

go_repository(
    name = "org_golang_x_sync",
    build_file_proto_mode = "disable_global",
    importpath = "golang.org/x/sync",
    sum = "h1:SQFwaSi55rU7vdNs9Yr0Z324VNlrF+0wMqRXT4St8ck=",
    version = "v0.0.0-20201020160332-67f06af15bc9",
)

go_repository(
    name = "org_golang_x_sys",
    build_file_proto_mode = "disable_global",
    importpath = "golang.org/x/sys",
    sum = "h1:taxlMj0D/1sOAuv/CbSD+MMDof2vbyPTqz5FNYKpXt8=",
    version = "v0.0.0-20211013075003-97ac67df715c",
)

go_repository(
    name = "org_golang_x_text",
    build_file_proto_mode = "disable_global",
    importpath = "golang.org/x/text",
    sum = "h1:i6eZZ+zk0SOf0xgBpEpPD18qWcJda6q1sxt3S0kzyUQ=",
    version = "v0.3.5",
)

go_repository(
    name = "org_golang_x_tools",
    build_file_proto_mode = "disable_global",
    importpath = "golang.org/x/tools",
    sum = "h1:CB3a9Nez8M13wwlr/E2YtwoU+qYHKfC+JrDa45RXXoQ=",
    version = "v0.0.0-20210106214847-113979e3529a",
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
    sum = "h1:uZLBy9jA8q0KALC3FGSzNecm71B3KQXeZQHGzfwPVNY=",
    version = "v1.13.0",
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
    sum = "h1:U+CaK85mrNNb4k8BNOfgJtJ/gr6kswUCFj6miSzVC6M=",
    version = "v1.9.5",
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
    sum = "h1:9VEQWz6LLMUsUl6PueE49ir4Ka6CzLymOAZDxpFsTDc=",
    version = "v1.6.0",
)

go_repository(
    name = "com_github_tidwall_match",
    build_file_proto_mode = "disable_global",
    importpath = "github.com/tidwall/match",
    sum = "h1:PnKP62LPNxHKTwvHHZZzdOAOCtsJTjo6dZLCwpKm5xc=",
    version = "v1.0.1",
)

go_repository(
    name = "com_github_tidwall_pretty",
    build_file_proto_mode = "disable_global",
    importpath = "github.com/tidwall/pretty",
    sum = "h1:HsD+QiTn7sK6flMKIvNmpqz1qrpP3Ps6jOKIKMooyg4=",
    version = "v1.0.0",
)

go_repository(
    name = "com_github_squzy_squzy_proto",
    importpath = "github.com/squzy/squzy_proto",
    sum = "h1:S73+nz2/3EbJLIJ3S4mUgylzvl9TFCddfEDEWF2an/Y=",
    version = "v0.0.0-20200501175525-0ca4b022b9f0",
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
    sum = "h1:Qgr9rKW7uDUkrbSmQeiDsGa8SjGyCOGtuasMWwvp2P4=",
    version = "v0.0.1",
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
    sum = "h1:L/CwN0zerZDmRFUapSPitk6f+Q3+0za1rQkzVuMiMFI=",
    version = "v0.1.0",
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
    sum = "h1:iURUrRGxPUNPdy5/HRSm+Yj6okJ6UtLINN0Q9M4+h3I=",
    version = "v0.8.1",
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
    sum = "h1:SPIRibHv4MatM3XXNO2BJeFLZwZ2LvZgfQ5+UNI2im4=",
    version = "v1.4.2",
)

go_repository(
    name = "com_github_spf13_cobra",
    build_file_proto_mode = "disable_global",
    importpath = "github.com/spf13/cobra",
    sum = "h1:ZlrZ4XsMRm04Fr5pSFxBgfND2EBVa1nLpiy1stUsX/8=",
    version = "v0.0.3",
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
    sum = "h1:IYppNjEV/C+/3VPbhHVxQ4t04eVW0cLp0/pNdW++6Ug=",
    version = "v1.3.2",
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
    sum = "h1:cxzIVoETapQEqDhQu3QfnvXAV4AlzcvUCxkVUFw3+EU=",
    version = "v1.4.0",
)

go_repository(
    name = "org_uber_go_multierr",
    build_file_proto_mode = "disable_global",
    importpath = "go.uber.org/multierr",
    sum = "h1:HoEmRHQPVSqub6w2z2d2EOVs2fjyFRGyofhKuyDq0QI=",
    version = "v1.1.0",
)

go_repository(
    name = "org_uber_go_zap",
    build_file_proto_mode = "disable_global",
    importpath = "go.uber.org/zap",
    sum = "h1:ORx85nbTijNz8ljznvCMR1ZBIPKFn3jQrag10X2AsuM=",
    version = "v1.10.0",
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
    sum = "h1:RM4zey1++hCTbCVQfnWeKs9/IEsaBLA8vTkd0WVtmH4=",
    version = "v0.3.0",
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
    sum = "h1:hzAQntlaYRkVSFEfj9OTWlVV1H155FMD8BTKktLv0QI=",
    version = "v0.0.0-20210930031921-04548b0d99d4",
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
    sum = "h1:ahKqKTFpO5KTPHxWZjEdPScmYaGtLo8Y4DMHoEsnp14=",
    version = "v1.6.3",
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
    sum = "h1:KgJ0snyC2R9VXYN2rneOtQcw5aHQB1Vv0sFl1UcHBOY=",
    version = "v10.2.0",
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
    name = "com_github_jinzhu_gorm",
    build_file_proto_mode = "disable_global",
    importpath = "github.com/jinzhu/gorm",
    sum = "h1:Drgk1clyWT9t9ERbzHza6Mj/8FY/CqMyVzOiHviMo6Q=",
    version = "v1.9.12",
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
    sum = "h1:HjfetcXq097iXP0uoPCdnM4Efp5/9MsM0/M+XOTeR3M=",
    version = "v1.0.1",
)

go_repository(
    name = "com_github_json_iterator_go",
    build_file_proto_mode = "disable_global",
    importpath = "github.com/json-iterator/go",
    sum = "h1:9yzud/Ht36ygwatGx56VwCZtlI/2AD15T1X2sjSuGns=",
    version = "v1.1.9",
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
    sum = "h1:sJZmqHoEaY7f+NPP8pgLB/WxulyR3fewgCM2qaSlBb4=",
    version = "v1.1.1",
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
    sum = "h1:ZqeYNhU3OHLH3mGKHDcjJRFFRrJa6eAM5H+CtDdOsPc=",
    version = "v0.0.0-20180228061459-e0a39a4cb421",
)

go_repository(
    name = "com_github_modern_go_reflect2",
    build_file_proto_mode = "disable_global",
    importpath = "github.com/modern-go/reflect2",
    sum = "h1:Esafd1046DLDQ0W1YjYsBW+p8U2u7vzgW2SQVmlNazg=",
    version = "v0.0.0-20180701023420-4b7aa43c6742",
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
    sum = "h1:SnqbnDw1V7RiZcXPx5MEeqPv2s79L9i7BJUlG/+RurQ=",
    version = "v1.27.1",
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
