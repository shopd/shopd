# shopd

Portable e-commerce https://shopd.link
- Mobile first
- Easy to move, uses an SQLite DB file
- Fast, Hypermedia API and SSR
- Self-contained web service
- Cross compile to run on your host OS


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
# https://github.com/mozey/config#toggling-env
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
