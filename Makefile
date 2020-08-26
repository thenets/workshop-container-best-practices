SHELL:=/bin/bash

BUILD_DIR=""

build-all:
	for dir in $$(ls ./steps) ; do \
		echo "# Etapa $$dir"; \
		make -s build BUILD_DIR=$$dir; \
		echo ; \
	done

build:
	cd steps/$(BUILD_DIR) && \
		docker build -t $(BUILD_DIR) .
	ls steps/$(BUILD_DIR)
