#!/usr/bin/env -S nu --stdin

# Build script using Nushell
# https://www.nushell.sh/

const PROJECT_DIR = path self .
const BINARY_NAME = ($PROJECT_DIR | path basename) ++ ".exe"

let projectEnv = {
	# Set C compiler
	CC: 'C:\cygwin64\bin\x86_64-w64-mingw32-gcc.exe',
	# Ensure CGO is enabled, just in case
	CGO_ENABLED: 1,
}

# Builds the project
def build [
	--production # adds the "production" build tag
	--headless # adds the "headless" build tag
]: nothing -> nothing {
	let $tags = [
		...(if $production { ["production"] } else { [] }),
		...(if $headless { ["headless"] } else { [] }),
	] | str join ","
	cd $PROJECT_DIR
	with-env $projectEnv {
		go build --tags $tags -o $BINARY_NAME
	}
}

# Builds the project and and runs the compiled executable
def build-run [
	--production # adds the "production" build tag
	--headless # adds the "headless" build tag
	...args # arguments to pass to the compiled executable
]: nothing -> nothing {
	print "Building..."
	build --production=$production --headless=$headless
	print "Running..."
	run ...$args
}

# Creates the .syso file for the application icon using the favicon.ico file
def build-icon []: nothing -> nothing {
	cd $PROJECT_DIR
	with-env $projectEnv {
		let $arch = $env | get --ignore-errors GOARCH
		let $flags = if $arch == null { [] } else { ["-arch", $arch] }
		glob rsrc_*.syso | each { rm $in }
		go run github.com/akavel/rsrc ...$flags -ico web/favicon.ico
	}
}

# Run the last built executable
def run [...args]: nothing -> nothing {
	cd $PROJECT_DIR
	with-env $projectEnv {
		^$BINARY_NAME ...$args
	}
}

# Runs tests for the project
def test [
	--verbose (-v) # enables verbose output
	--production # adds the "production" build tag
	--headless # adds the "headless" build tag
	...args # arguments to pass to the test binary
]: nothing -> nothing {
	let $tags = [
		...(if $production { ["production"] } else { [] }),
		...(if $headless { ["headless"] } else { [] }),
	] | str join ","
	cd $PROJECT_DIR
	let $flags = [
		...(if $verbose { ["-v"] } else { [] }),
	]
	with-env $projectEnv {
		go test ...$flags --tags $tags --args ...$args
	}
}

# Lint runs the Go linter
def lint []: nothing -> nothing {
	cd $PROJECT_DIR
	golangci-lint run --timeout 5m
}


# Cleans the project
def clean []: nothing -> nothing {
	cd $PROJECT_DIR
	rm $BINARY_NAME
}
