#!/usr/bin/env bash
set -u                    # exit on undefined variable
bash -c 'set -o pipefail' # return code of first cmd to fail in a pipeline

APP_DIR="$APP_DIR"

# Reset all generated code related to templ.
# When refactoring, the generated code might break,
# and prevent the magefiles etc from compiling
rm "$APP_DIR/go/router/init_api_templ.go"
rm "$APP_DIR/go/router/init_content_templ.go"
