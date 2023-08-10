#!/usr/bin/env sh

if command -v podman; then
	COMMAND="podman"
elif command -v docker; then
	COMMAND="docker"
else
	printf >&2 "Did not find podman or docker command \n"
	exit 1
fi

$COMMAND build -t quickbeam-build .
$COMMAND run --rm -v "./:/src" quickbeam-build
