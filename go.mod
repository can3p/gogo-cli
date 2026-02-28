module github.com/can3p/gogo-cli

go 1.26.0

require (
	github.com/antonlindstrom/pgstore v0.0.0-20220421113606-e3a6e3fed12a
	github.com/can3p/go-scarf v0.0.0-20240501104927-634d80d41254
	github.com/can3p/gogo v0.0.0-20250919103346-af8499bb4701
	github.com/can3p/kleiner v0.0.8
	github.com/can3p/pcom v0.0.0-20260218233203-00092a5182e9
	github.com/gin-contrib/sessions v1.0.4
	github.com/gin-gonic/gin v1.10.1
	github.com/google/uuid v1.6.0
	github.com/jmoiron/sqlx v1.4.0
	github.com/pelletier/go-toml/v2 v2.2.4
	github.com/pkg/errors v0.9.1
	github.com/spf13/cobra v1.8.0
	github.com/stretchr/testify v1.11.1
	github.com/volatiletech/null/v8 v8.1.2
	github.com/volatiletech/sqlboiler/v4 v4.19.1
	github.com/yuin/goldmark v1.7.16
)

// TODO: Remove this replace directive after publishing gogo with testcontainers
replace github.com/can3p/gogo => ../gogo

require (
	aead.dev/minisign v0.2.0 // indirect
	dario.cat/mergo v1.0.0 // indirect
	github.com/Azure/go-ansiterm v0.0.0-20230124172434-306776ec8161 // indirect
	github.com/JohannesKaufmann/dom v0.2.0 // indirect
	github.com/JohannesKaufmann/html-to-markdown/v2 v2.5.0 // indirect
	github.com/Microsoft/go-winio v0.6.2 // indirect
	github.com/Nvveen/Gotty v0.0.0-20120604004816-cd527374f1e5 // indirect
	github.com/PuerkitoBio/goquery v1.9.2 // indirect
	github.com/alecthomas/chroma v0.10.0 // indirect
	github.com/andybalholm/cascadia v1.3.3 // indirect
	github.com/aymerick/douceur v0.2.0 // indirect
	github.com/badoux/checkmail v1.2.4 // indirect
	github.com/bytedance/gopkg v0.1.3 // indirect
	github.com/bytedance/sonic v1.14.1 // indirect
	github.com/bytedance/sonic/loader v0.3.0 // indirect
	github.com/can3p/anti-disposable-email v0.0.0-20230623054934-598d3044afb0 // indirect
	github.com/cenkalti/backoff/v4 v4.3.0 // indirect
	github.com/cloudwego/base64x v0.1.6 // indirect
	github.com/containerd/continuity v0.4.5 // indirect
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/davidbyttow/govips/v2 v2.16.0 // indirect
	github.com/dlclark/regexp2 v1.4.0 // indirect
	github.com/docker/cli v27.4.1+incompatible // indirect
	github.com/docker/docker v28.0.0+incompatible // indirect
	github.com/docker/go-connections v0.5.0 // indirect
	github.com/docker/go-units v0.5.0 // indirect
	github.com/ericlagergren/decimal v0.0.0-20190420051523-6335edbaa640 // indirect
	github.com/fatih/color v1.16.0 // indirect
	github.com/friendsofgo/errors v0.9.2 // indirect
	github.com/gabriel-vasile/mimetype v1.4.10 // indirect
	github.com/gin-contrib/sse v1.1.0 // indirect
	github.com/go-gorp/gorp/v3 v3.1.0 // indirect
	github.com/go-playground/locales v0.14.1 // indirect
	github.com/go-playground/universal-translator v0.18.1 // indirect
	github.com/go-playground/validator/v10 v10.27.0 // indirect
	github.com/go-resty/resty/v2 v2.11.0 // indirect
	github.com/go-viper/mapstructure/v2 v2.4.0 // indirect
	github.com/goccy/go-json v0.10.5 // indirect
	github.com/gofrs/uuid v4.2.0+incompatible // indirect
	github.com/gogo/protobuf v1.3.2 // indirect
	github.com/google/go-github/v57 v57.0.0 // indirect
	github.com/google/go-querystring v1.1.0 // indirect
	github.com/google/shlex v0.0.0-20191202100458-e7afc7fbc510 // indirect
	github.com/gorilla/context v1.1.2 // indirect
	github.com/gorilla/css v1.0.1 // indirect
	github.com/gorilla/securecookie v1.1.2 // indirect
	github.com/gorilla/sessions v1.4.0 // indirect
	github.com/hashicorp/golang-lru/v2 v2.0.7 // indirect
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/klauspost/cpuid/v2 v2.3.0 // indirect
	github.com/leodido/go-urn v1.4.0 // indirect
	github.com/lib/pq v1.11.2 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/microcosm-cc/bluemonday v1.0.27 // indirect
	github.com/minio/selfupdate v0.6.0 // indirect
	github.com/mmcdole/gofeed v1.3.0 // indirect
	github.com/mmcdole/goxpp v1.1.1-0.20240225020742-a0c311522b23 // indirect
	github.com/moby/docker-image-spec v1.3.1 // indirect
	github.com/moby/sys/user v0.3.0 // indirect
	github.com/moby/term v0.5.0 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/opencontainers/go-digest v1.0.0 // indirect
	github.com/opencontainers/image-spec v1.1.0 // indirect
	github.com/opencontainers/runc v1.2.8 // indirect
	github.com/ory/dockertest/v3 v3.12.0 // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	github.com/rubenv/sql-migrate v1.8.1 // indirect
	github.com/samber/lo v1.52.0 // indirect
	github.com/samber/mo v1.16.0 // indirect
	github.com/sirupsen/logrus v1.9.3 // indirect
	github.com/spf13/cast v1.5.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/teekennedy/goldmark-markdown v0.5.1 // indirect
	github.com/twitchyliquid64/golang-asm v0.15.1 // indirect
	github.com/ugorji/go/codec v1.3.0 // indirect
	github.com/volatiletech/inflect v0.0.1 // indirect
	github.com/volatiletech/randomize v0.0.1 // indirect
	github.com/volatiletech/strmangle v0.0.8 // indirect
	github.com/xeipuuv/gojsonpointer v0.0.0-20190905194746-02993c407bfb // indirect
	github.com/xeipuuv/gojsonreference v0.0.0-20180127040603-bd5ef7bd5415 // indirect
	github.com/xeipuuv/gojsonschema v1.2.0 // indirect
	github.com/yuin/goldmark-highlighting v0.0.0-20220208100518-594be1970594 // indirect
	golang.org/x/arch v0.21.0 // indirect
	golang.org/x/crypto v0.45.0 // indirect
	golang.org/x/image v0.23.0 // indirect
	golang.org/x/net v0.47.0 // indirect
	golang.org/x/sync v0.19.0 // indirect
	golang.org/x/sys v0.38.0 // indirect
	golang.org/x/text v0.31.0 // indirect
	golang.org/x/xerrors v0.0.0-20240903120638-7835f813f4da // indirect
	google.golang.org/protobuf v1.36.9 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
