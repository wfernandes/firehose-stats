#!/bin/bash

set -e

(cf uninstall-plugin "FirehoseStats" || true) && go build && cf install-plugin firehose-stats
