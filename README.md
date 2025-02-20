🏗️ WIP 🚧 This repo is an open source rewrite (using [templ](https://templ.guide/) and [Tailwind CSS](https://tailwindcss.com/)) of the closed prototype ([hugo](https://gohugo.io/) plugin) hosted at [shopd.link](https://shopd.link/) 👷


---
# shopd

Portable e-commerce
- Mobile first
- Fast, Hypermedia API and SSR
- Cross compile to run on your host OS
- Easy to move, uses an SQLite DB file
- Stores data, user content, and config in the same directory
- Self-contained web service
- Does not require other software to be installed


## Quick Start

Install dev dependencies
```bash
go mod vendor
pnpm install
```

Config
```bash
export APP_DIR=$(pwd)
mage EnvGen dev example.com
source .env.dev-example-com.sh
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


