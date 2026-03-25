build:
	go build -o bin/envy .

release:
	@if [ -z "$(VERSION)" ]; then echo "Usage: make release VERSION=v0.1.0"; exit 1; fi
	git tag $(VERSION)
	git push origin $(VERSION)
