#
# Variables
#

# Specified from caller
OS_ARCH = local
EXT =
VERSION = 1.0.0

# Fixed
RM = rm -rf
GOCMD = go
SRCDIR = ./cmd
DESTDIR = .
LDFLAGS = -X=main.Version=$(VERSION)


#
# Rules
#

.PHONY: clean

all: clean peco-actions

peco-actions: ${SRCDIR}/peco-actions.go
	$(GOCMD) build --ldflags "$(LDFLAGS)" $<


clean:
	$(RM) $(DESTDIR)/peco-actions

