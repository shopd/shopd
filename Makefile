# Live reload with other tools
# - Tailwind CSS for generating a css bundle
# - esbuild for bundling JavaScript or TypeScript
# - air for re-building Go source, sending a reload event to the templ proxy
# https://templ.guide/commands-and-tools/live-reload-with-other-tools/#setting-up-the-makefile

# dev/site uses templ to detect changes to .templ files 
# creates _templ.txt files (to reduce Go code to be re-generated),
# and sends the reload event to the browser
# Default url: http://localhost:7331
dev/site:
	templ generate -v --watch \
	--proxy="http://localhost:8443" \
	--open-browser=false

# dev/shopd detects go file changes to re-build and re-run the server
dev/shopd:
	air \
	--build.cmd "go build -o www/build/shopd ./cmd/shopd/main.go" \
	--build.bin "www/build/shopd run" \
	--build.delay "100" \
	--build.exclude_dir "node_modules" \
	--build.exclude_dir "vendor" \
	--build.include_ext "go" \
	--build.stop_on_error "false" \
	--misc.clean_on_exit true

# dev/css uses tailwind generates the app.css bundle
dev/css:
	pnpx tailwindcss -i ./src/app.css -o ./www/build/app.css \
	--minify --watch

# dev/app uses esbuild generates the app.js bundle
dev/app:
	pnpx esbuild src/app.ts --bundle --outdir=www/build/ \
	--watch

# dev/sync detects changes in the build folder, 
# then reloads the browser via templ proxy
dev/sync:
	air \
	--build.cmd "templ generate --notify-proxy" \
	--build.bin "true" \
	--build.delay "100" \
	--build.exclude_dir "" \
	--build.include_dir "www/build" \
	--build.include_ext "js,css"

dev/caddy:
	caddy run

# dev runs all the dev tasks with live reload
# make -j6 dev/caddy dev/site dev/shopd dev/css dev/app dev/sync
dev: 
	make -j2 dev/caddy dev/shopd
