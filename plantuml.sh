#!/bin/bash
SCRIPT_DIR=$(cd $(dirname $0); pwd)
java -jar "$SCRIPT_DIR/lib/plantuml-1.2024.8.jar" "$@" 