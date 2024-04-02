TARGET := square-events
OUT := build/$(TARGET).zip

$(OUT): build/bootstrap
	mkdir -p $(shell dirname $@)
	zip -rj $@ $^

build/bootstrap:
	wire
	GOOS=linux GOARCH=arm64 go build -o $@ .

.PHONY: clean
clean:
	rm -f $(TARGET)
	rm -rf build
