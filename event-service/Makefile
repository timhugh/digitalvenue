OUT_DIR := build
TARGET := event-service

$(OUT_DIR)/$(TARGET).zip: $(OUT_DIR)/$(TARGET)/bootstrap
	zip -rj $@ $(OUT_DIR)/$(TARGET)/bootstrap

$(OUT_DIR)/$(TARGET)/bootstrap: .
	wire
	go build -o $@ .

.PHONY: clean
clean:
	rm -f $(OUT_DIR)/$(TARGET).zip
	rm -f $(OUT_DIR)/$(TARGET)/bootstrap
	rmdir -p $(OUT_DIR)/$(TARGET)
