dep:
	yarn install

run:
	yarn build
	sudo go run *.go	

dev:
	yarn dev