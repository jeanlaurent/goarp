dep:
	yarn install

css:
	cp node_modules/bootstrap/dist/css/bootstrap.min.css app/style/bootstrap/css/

build: css
	yarn build
	go build *.go

runhttp: css
	yarn dev

run:
	sudo go run *.go	