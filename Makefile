#
# Copyright 2022 Michael Graff.
#
# Licensed under the Apache License, Version 2.0 (the "License")
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#   http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#

TARGETS=test local

#
# Build targets.  Adding to these will cause magic to occur.
#

# These are targets for "make local"
BINARIES = testtool

#
# Below here lies magic...
#

all_deps := $(shell find * -name '*.go' | grep -v _test)

now := $(shell date -u +%Y%m%dT%H%M%S)

#
# Default target.
#

.PHONY: all
all: ${TARGETS}

#
# Build locally, mostly for development speed.
#

.PHONY: local
local: $(addprefix bin/,$(BINARIES))

bin/%:: ${all_deps}
	@[ -d bin ] || mkdir bin
	go build -ldflags="-s -w" -o $@ app/$(@F)/*.go

#
# Test targets
#

.PHONY: test
test:
	go test -race ./...

#
# Clean the world.
#

.PHONY: clean
clean:
	rm -f bin/*

.PHONY: really-clean
really-clean: clean
