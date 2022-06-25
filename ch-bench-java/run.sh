#!/bin/bash

#
# Java client:
#   - ./run.sh long
#   - ./run.sh record
#   - ./run.sh skip
#   - OPTS="-Xms24m -Xmx24m" ./run.sh long
#
# JDBC driver: ./run.sh
#

set -e

: ${SQL:="SELECT number FROM system.numbers_mt LIMIT 500000000"}
: ${URL:='http://localhost?!compress&format=RowBinaryWithNamesAndTypes'}
: ${VER:="0.3.2-patch10"}

: ${OPTS:=""}
#OPTS="-XX:+UnlockDiagnosticVMOptions -XX:+DebugNonSafepoints -XX:+FlightRecorder -XX:StartFlightRecording=disk=false,filename=recorded.jfr,settings=profile"
PKG="clickhouse-jdbc-$VER-all.jar"


if [ ! -f "$PKG" ]; then
    echo "Downloading $PKG..."
    curl -sOL "https://github.com/ClickHouse/clickhouse-jdbc/releases/download/v$VER/$PKG"
else
    echo "Found $PKG"
fi

if [ ! -f "Main.class" ] || [ "Main.java" -nt "Main.class" ]; then
    rm -fv Main.class
    echo "Compiling..."
    javac -cp "$PKG" Main.java
else
    echo "Found compiled class"
fi

echo "Running..."
\time -v java $OPTS -Durl="$URL" -Dsql="$SQL" -cp ".:$PKG" Main $@
