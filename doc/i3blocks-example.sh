#!/bin/bash
set -euo pipefail

# i3blocks.conf
# [topic]
# command=$HOME/.config/i3/i3blocks/topic
# interval=1

TMOUT=1  # Time out if this hangs on reading for some reason

IFS=, read -r DURATION TOPIC <<< "$(topic latest -truncate 1s -duration -topic)"

if [[ "${#TOPIC}" -eq 0 ]]; then
	echo "🗒"
	exit 0
elif [ "${#TOPIC}" -gt 5 ]; then
	printf -v SHRT '%.5s…' "$TOPIC"
else
	SHRT="$TOPIC"
fi

echo "🗒${TOPIC} ${DURATION}"
echo "🗒${SHRT} ${DURATION}"
echo "#ffafaf"
