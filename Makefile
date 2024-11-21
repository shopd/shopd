# Live reload with other tools
# - Tailwind CSS for generating a css bundle
# - esbuild for bundling JavaScript or TypeScript
# - air for re-building Go source, sending a reload event to the templ proxy
# https://templ.guide/commands-and-tools/live-reload-with-other-tools/#setting-up-the-makefile

# dev/site generates _templ.txt with watch mode
# https://github.com/a-h/templ/pull/366
# Don't use the live reload proxy, rather make make app.js poll /api 
# and reload the page on if the build timestamp changed
# https://templ.guide/commands-and-tools/live-reload
dev/site:
	templ generate -v --watch --path www/content

# dev/shopd detects go file changes to re-build and re-run the server
dev/shopd:
	air \
	--build.cmd "go build -o build/shopd ./cmd/shopd/main.go" \
	--build.bin "build/shopd run" \
	--build.delay "100" \
	--build.exclude_dir "node_modules" \
	--build.exclude_dir "vendor" \
	--build.include_ext "go" \
	--build.stop_on_error "false" \
	--misc.clean_on_exit true

# TODO Don't minify for dev
# dev/css uses tailwind generates the app.css bundle
dev/css:
	pnpx tailwindcss -i ./src/app.css -o ./build/app.css \
	--minify --watch

# dev/app uses esbuild generates the app.js bundle
dev/app:
	pnpx esbuild src/app.ts --bundle --outdir=build/ \
	--watch

# TODO Create sync dev cmd?
# dev/sync detects changes in the build folder, 
# then reloads the browser via templ proxy
dev/sync:
	air \
	--build.cmd "go run ./cmd/shopd/main.go sync dev" \
	--build.bin "true" \
	--build.delay "100" \
	--build.exclude_dir "" \
	--build.include_dir "build" \
	--build.include_ext "js,css"

# TODO Don't serve static files in dev?
# Rather make use of templ watch mode,
# this requires a Go backend service cmd with routes
# corresponding to the static site paths.
# Another cmd is then used to render static files for prod
# https://templ.guide/static-rendering/generating-static-html-files-with-templ
dev/caddy:
	caddy run

# dev runs all the dev tasks with live reload
# make -j6 dev/caddy dev/site dev/shopd dev/css dev/app dev/sync
dev: 
	make -j2 dev/caddy dev/shopd
