# Commands
NPM=npm
NPMRUN=$(NPM) run
NPMINSTALL=$(NPM) install
NPMRUNBUILD=$(NPMRUN) build
NPMRUNTEST=$(NPMRUN) test
NPMRUNCLEAN=$(NPMRUN) clean
NPMRUNDEV=$(NPMRUN) dev

all: node_modules dist
node_modules: # gets all dependencies with npm
	$(NPMINSTALL)
dist: # build for production
	$(NPMRUNBUILD)
test: 
	$(NPMRUNTEST)
clean: 
	$(NPMRUNCLEAN)
dev:
	$(NPMRUNDEV)
