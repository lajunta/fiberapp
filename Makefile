assets:
	yarn run build-css
	esbuild --bundle ./src/app.js --outdir=./public/ --minify --sourcemap

dev: assets
	ConfigDir=${HOME}/.fiberapp PORT=3000 go run .

build: assets
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o bin/app
	upx -q bin/app
