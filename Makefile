.PHONY: remote build

all: build

build:
	./build.sh $(filter-out $@,$(MAKECMDGOALS))

remote:
	./remote.sh $(filter-out $@,$(MAKECMDGOALS))

%:
	@:
