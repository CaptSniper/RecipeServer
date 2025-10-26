# Go settings
GO := go
BIN_DIR := bin

# Binary base names
READER_BIN := read_recipe
WRITER_BIN := create_recipe

# Detect OS for extension
ifeq ($(OS),Windows_NT)
    EXE := .exe
else
    EXE :=
endif

# Default target: build for current OS
all: $(BIN_DIR)/$(READER_BIN)$(EXE) $(BIN_DIR)/$(WRITER_BIN)$(EXE)

# Ensure bin directory exists
$(BIN_DIR):
	mkdir -p $(BIN_DIR)

# Build reader binary
$(BIN_DIR)/$(READER_BIN)$(EXE): mainRead/mainRead.go rfp/*.go | $(BIN_DIR)
	$(GO) build -o $@ mainRead/mainRead.go

# Build writer binary
$(BIN_DIR)/$(WRITER_BIN)$(EXE): mainWrite/mainWrite.go rfp/*.go | $(BIN_DIR)
	$(GO) build -o $@ mainWrite/mainWrite.go

# Cross-compile for Linux
linux: | $(BIN_DIR)
	GOOS=linux GOARCH=amd64 $(GO) build -o $(BIN_DIR)/$(READER_BIN)_linux mainRead/mainRead.go
	GOOS=linux GOARCH=amd64 $(GO) build -o $(BIN_DIR)/$(WRITER_BIN)_linux mainWrite/mainWrite.go

# Cross-compile for Windows
windows: | $(BIN_DIR)
	GOOS=windows GOARCH=amd64 $(GO) build -o $(BIN_DIR)/$(READER_BIN)_windows.exe mainRead/mainRead.go
	GOOS=windows GOARCH=amd64 $(GO) build -o $(BIN_DIR)/$(WRITER_BIN)_windows.exe mainWrite/mainWrite.go

# Clean build artifacts
clean:
	rm -rf $(BIN_DIR)

.PHONY: all clean linux windows
