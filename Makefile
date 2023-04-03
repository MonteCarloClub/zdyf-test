BIN_DIR:=./bin
SRC_DIR=.

BINS:=$(BIN_DIR)/zdyf-test
MAIN_SRCS:=$(SRC_DIR)/init.go $(SRC_DIR)/main.go

.PHONY:all clean zdyf-test

all:zdyf-test

clean:
	rm -rf $(BIN_DIR)

zdyf-test:
	mkdir -p $(BIN_DIR)
	go build -o $(BINS) $(MAIN_SRCS)
