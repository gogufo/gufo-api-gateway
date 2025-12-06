#!/bin/bash
set -e

PROTO_FILE="microservice.proto"
OUT_ROOT="."

LANGUAGES=(
  go
  java
  python
  node
  ruby
  php
  cpp
  csharp
  objc
)

if [ ! -f "$PROTO_FILE" ]; then
  echo "❌ Proto file not found: $PROTO_FILE"
  exit 1
fi

echo "🚀 Generating proto SDKs for all languages..."

for LANG in "${LANGUAGES[@]}"; do
  echo "➡ Generating for: $LANG"

 # mkdir -p "$OUT_ROOT/$LANG"

 docker run -v $PWD:/defs namely/protoc-all \
   -f microservice.proto -o "$LANG/" -l "$LANG"


done

echo "✅ All SDKs generated successfully"
