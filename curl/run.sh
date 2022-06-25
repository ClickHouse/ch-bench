#!/bin/bash

set -e

: ${FMT:="RowBinaryWithNamesAndTypes"}
: ${SQL:="SELECT number FROM system.numbers_mt LIMIT 500000000"}
: ${URL:="http://localhost:8123"}

\time -v curl -s -H "X-ClickHouse-Format: $FMT" --get --data-urlencode "query=$SQL" "$URL" > /dev/null
