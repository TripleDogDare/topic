#!/bin/bash
set -euo pipefail

./build.sh
cp ./bin/topic /usr/local/bin/topic
