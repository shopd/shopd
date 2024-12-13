
# TODO Don't minify for dev
# dev/css uses tailwind generates the app.css bundle
dev/css:
	pnpx tailwindcss -i ./src/app.css -o ./build/app.css \
	--minify --watch

# dev/app uses esbuild generates the app.js bundle
dev/app:
	pnpx esbuild src/app.ts --bundle --outdir=build/ --watch
	pnpx esbuild /Users/mozey/pro/shopd/shopd/src/app.ts --bundle --outdir=/Users/mozey/pro/shopd/shopd/build --watch
