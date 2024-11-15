# Live reload with other tools
# - Tailwind CSS for generating a css bundle
# - esbuild for bundling JavaScript or TypeScript
# - air for re-building Go source, sending a reload event to the templ proxy
# https://templ.guide/commands-and-tools/live-reload-with-other-tools/#setting-up-the-makefile

# watch templ detects changes to .templ files 
# and re-creates _templ.go files,
# then send reload event to browser.
# Default url: http://localhost:7331
watch/templ:
	templ generate -v --watch \
	--proxy="http://localhost:8080" \
	--open-browser=false

# watch shopd detects go file changes to re-build and re-run the server
watch/server:
	go run github.com/cosmtrek/air@v1.51.0 \
	--build.cmd "go build -o www/build/shopd ./cmd/shopd/..." \
	--build.bin "www/build/shopd run" \
	--build.delay "100" \
	--build.exclude_dir "node_modules" \
	--build.exclude_dir "vendor" \
	--build.include_ext "go" \
	--build.stop_on_error "false" \
	--misc.clean_on_exit true

# watch tailwind generates the app.css bundle
watch/tailwind:
	pnpx tailwindcss -i ./src/app.css -o ./www/build/app.css \
	--minify --watch

# watch esbuild generates the app.js bundle
watch/esbuild:
	pnpx esbuild src/app.ts --bundle --outdir=www/build/ \
	--watch

# watch sync detects changes in the build folder, 
# then reloads the browser via templ proxy
watch/sync:
	go run github.com/cosmtrek/air@v1.51.0 \
	--build.cmd "templ generate --notify-proxy" \
	--build.bin "true" \
	--build.delay "100" \
	--build.exclude_dir "" \
	--build.include_dir "www/build" \
	--build.include_ext "js,css"

# dev runs the server with live reload
dev: 
	make -j5 watch/templ watch/server watch/tailwind watch/esbuild watch/sync
