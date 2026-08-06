#!/bin/sh
set -e

# GitHub Docker Action 将 inputs 作为环境变量 INPUT_<NAME> 传入
# 本脚本将它们映射为 globalizer 参数
#
# 运行模式:
#   1. 带参数运行时 (docker run image version / serve / translate ...)
#      → 直接透传给 globalizer
#   2. 无参数运行时 (GitHub Action 模式)
#      → 读取 INPUT_* 环境变量执行 translate

# 二进制位于 /usr/local/bin/globalizer (Dockerfile 安装)
GLOBALIZER="${GLOBALIZER_BIN:-globalizer}"

# 模式 1: 有参数 → 直接透传
if [ $# -gt 0 ]; then
    exec "$GLOBALIZER" "$@"
fi

# 模式 2: GitHub Action 模式 (无参数, 读取 INPUT_*)
TARGET="${INPUT_TARGET:-README.md}"
LANGUAGES="${INPUT_LANGUAGES:-zh-CN}"
OUTPUT_DIR="${INPUT_OUTPUT_DIR:-docs}"
MODEL="${INPUT_MODEL:-gpt-4o}"
BASE_URL="${INPUT_BASE_URL:-https://api.openai.com/v1}"
API_KEY="${INPUT_API_KEY:-$OPENAI_API_KEY}"
MOCK="${INPUT_MOCK:-false}"

ARGS="translate $TARGET --lang $LANGUAGES --output $OUTPUT_DIR --model $MODEL --base-url $BASE_URL"

if [ "$MOCK" = "true" ]; then
    echo "::group::Translating $TARGET to $LANGUAGES (mock mode)"
    # shellcheck disable=SC2086
    "$GLOBALIZER" $ARGS --mock
    echo "::endgroup::"
    exit 0
fi

if [ -z "$API_KEY" ]; then
    echo "::error::OPENAI_API_KEY or input 'api-key' is required"
    exit 1
fi

echo "::group::Translating $TARGET to $LANGUAGES"
# shellcheck disable=SC2086
"$GLOBALIZER" $ARGS --api-key "$API_KEY"
echo "::endgroup::"
