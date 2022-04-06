#
# Variables
#

# Specified from caller
OS_ARCH = local
EXT =

# Fixed
RM = rm -f
GOCMD = go


#
# Rules
#

.PHONY: clean

all: clean peco-actions

peco-actions: peco-actions.go
	$(GOCMD) build --ldflags "$(LDFLAGS)" $<


clean:
	$(RM) peco-actions

