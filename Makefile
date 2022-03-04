#
# Variables
#

# Specified from caller
OS_ARCH = local
EXT =
VERSION = 1.0.0

# Fixed
RM = rm -f
GOCMD = go
LDFLAGS = -X=main.Version=$(VERSION)


#
# Rules
#

.PHONY: clean

all: clean peco-actions

peco-actions: peco-actions.go
	$(GOCMD) build --ldflags "$(LDFLAGS)" $<


clean:
	$(RM) peco-actions

