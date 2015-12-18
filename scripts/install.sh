#!/bin/bash

set -e

(cf uninstall-plugin "FirehoseStats" || true) && go build -o firehose-stats main.go && cf install-plugin firehose-stats
