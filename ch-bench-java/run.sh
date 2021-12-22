#!/bin/bash

set -e

: ${DB_HOST:="localhost"}
: ${DB_PORT:="8123"}
: ${DRIVER_VERSION:="0.3.2-test3"}

JDBC_DRIVER="clickhouse-jdbc-$DRIVER_VERSION-http.jar"

if [ ! -f "$JDBC_DRIVER" ]; then
    echo "Downloading $JDBC_DRIVER..."
    curl -sOL "https://github.com/ClickHouse/clickhouse-jdbc/releases/download/v$DRIVER_VERSION/$JDBC_DRIVER"
else
    echo "Found $JDBC_DRIVER"
fi

if [ ! -f "Main.class" ]; then
    echo "Compiling..."
    javac -cp "$JDBC_DRIVER" Main.java
else
    echo "Found compiled class"
fi

time java -DdbHost="$DB_HOST" -DdbPort="$DB_PORT" -cp ".:$JDBC_DRIVER" Main

