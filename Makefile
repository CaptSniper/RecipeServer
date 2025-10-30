# Go settings
GO := go
BIN_DIR := bin
EXE :=

ifeq ($(OS),Windows_NT)
	EXE := .exe
endif

# Executable name
APP_NAME := recipe_tool

# Default target: build for current OS
all: build

# Create bin folder if it doesn't exist
$(BIN_DIR):
ifeq ($(OS),Windows_NT)
	if not exist $(BIN_DIR) mkdir $(BIN_DIR)
else
	mkdir -p $(BIN_DIR)
endif

# Build for current OS
build: $(BIN_DIR)
	$(GO) build -o $(BIN_DIR)/$(APP_NAME)$(EXE) .

# Build Linux 64-bit binary (from any OS)
# linux64: $(BIN_DIR)
# ifeq ($(OS),Windows_NT)
# 	@echo "Cross-compiling Linux binaries from Windows requires WSL or proper environment variables."
# else
# 	GOOS=linux GOARCH=amd64 $(GO) build -o $(BIN_DIR)/$(APP_NAME)_linux main/main.go
# endif

# Clean bin folder
clean:
ifeq ($(OS),Windows_NT)
	cmd /C "if exist $(BIN_DIR) rmdir /s /q $(BIN_DIR)"
else
	rm -rf $(BIN_DIR)
endif

.PHONY: all build linux64 clean
