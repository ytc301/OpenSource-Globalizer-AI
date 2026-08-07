#!/bin/sh
set -e

# GitHub Docker Action 将 inputs 作为环境变量 INPUT_<NAME> 传入
# (GitHub 自动为每个 input 设置 INPUT_<NAME>=value, 含默认值)
#
# 运行模式:
#   1. GitHub Action 模式: 检测到 INPUT_TARGET 环境变量 → 读取 INPUT_* 执行 translate
#   2. CLI 透传模式: 其余情况 → 将参数直接透传给 globalizer
#      (docker run image version / serve / translate ...)

# 二进制位于 /usr/local/bin/globalizer (Dockerfile 安装)
GLOBALIZER="${GLOBALIZER_BIN:-globalizer}"

# 模式 1: GitHub Action 模式 (存在 INPUT_TARGET 环境变量)
if [ -n "$INPUT_TARGET" ]; then
    TARGET="$INPUT_TARGET"
    LANGUAGES="${INPUT_LANGUAGES:-zh-CN}"
    # 注意: GitHub Docker Action 的 input 名保留连字符 (INPUT_API-KEY),
    # 而 sh 变量名不能含连字符, 需从环境变量中查找
    get_env() {
        env | grep "^$1=" | head -1 | cut -d= -f2-
    }
    OUTPUT_DIR="$(get_env INPUT_OUTPUT-DIR)"
    [ -n "$OUTPUT_DIR" ] || OUTPUT_DIR="${INPUT_OUTPUT_DIR:-docs}"
    MODEL="${INPUT_MODEL:-gpt-4o}"
    BASE_URL="$(get_env INPUT_BASE-URL)"
    [ -n "$BASE_URL" ] || BASE_URL="${INPUT_BASE_URL:-${OPENAI_BASE_URL:-https://api.openai.com/v1}}"
    API_KEY="$(get_env INPUT_API-KEY)"
    [ -n "$API_KEY" ] || API_KEY="${INPUT_API_KEY:-$OPENAI_API_KEY}"
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
    exit 0
fi

# 模式 2: CLI 透传 (docker run image <args>)
# 注意: 镜像 CMD 为 ["serve"], docker run 无命令时会自动追加 → 默认启动 HTTP 服务
exec "$GLOBALIZER" "$@"
