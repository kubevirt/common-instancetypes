SHELL=/bin/bash

all: lint generate

lint:
	yamllint .

generate:
	./scripts/generate.sh

clean: 
	rm -f *.yaml

.PHONY: all lint generate
