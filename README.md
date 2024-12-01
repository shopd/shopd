# shopd

Mostly static e-commerce https://shopd.link


## Quick Start

Install dependencies
```bash
go mod vendor
pnpm install
```

Config
```bash
APP_DIR=$(pwd) mage EnvGen dev example.com
conf dev-example-com
mage CaddyfileGenDev
```

Run dev
```bash
mage dev
```

[Preview on localhost](https://localhost:8443)

Stop
```bash
mage down dev
```

**TODO** Debug templ static gen cmd
```bash
mage DebugTemplStaticGen
# templ generate -v --watch --path /Users/mozey/pro/shopd/shopd/www --cmd "go run /Users/mozey/pro/shopd/shopd/cmd/shopd/... static gen --env dev"
```
