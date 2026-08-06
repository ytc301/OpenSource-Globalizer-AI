#!/bin/sh
set -e

# GitHub Docker Action 将 inputs 作为环境变量 INPUT_<NAME> 传入
# 本脚本将它们映射为 globalizer 参数，并处理 serve 模式

# 判断运行模式：serve 或 translate
if [ "$1" = "serve" ]; then
    exec /globalizer serve
fi

# 默认 translate 模式
TARGET="${INPUT_TARGET:-README.md}"
LANGUAGES="${INPUT_LANGUAGES:-zh-CN}"
OUTPUT_DIR="${INPUT_OUTPUT_DIR:-docs}"
MODEL="${INPUT_MODEL:-gpt-4o}"
BASE_URL="${INPUT_BASE_URL:-https://api.openai.com/v1}"
API_KEY="${INPUT_API_KEY:-$OPENAI_API_KEY}"

if [ -z "$API_KEY" ]; then
    echo "::error::OPENAI_API_KEY or input 'api-key' is required"
    exit 1
fi

echo "::group::Translating $TARGET to $LANGUAGES"
/globalizer translate "$TARGET" \
    --lang "$LANGUAGES" \
    --output "$OUTPUT_DIR" \
    --model "$MODEL" \
    --base-url "$BASE_URL" \
    --api-key "$API_KEY"
echo "::endgroup::"
