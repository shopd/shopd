# shopd

Portable e-commerce https://shopd.link
- Mobile first
- Fast, Hypermedia API and SSR
- Cross compile to run on your host OS
- Easy to move, uses an SQLite DB file
- Stores user content and config in the same directory
- Self-contained web service
- Does not require other software to be installed<sup>[1]</sup>


## Quick Start

Install dev dependencies
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


## Footnotes

[1] Plugins may depend on additional software


