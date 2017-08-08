build:
	go build -v
dev: build
	./bsupgrade
clean:
	rm bsupgrade
	rm bsupgrade.db
