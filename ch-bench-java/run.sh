#!/bin/bash

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

if [ ! -f "Main.class" ]; then
    echo "Compiling..."
    javac -cp "$PKG" Main.java
else
    echo "Found compiled class"
fi

echo "Running..."
\time -v java $OPTS -Durl="$URL" -Dsql="$SQL" -cp ".:$PKG" Main $@
