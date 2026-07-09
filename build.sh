#!/bin/bash
set -e

# ============================================================
# 镜像构建脚本
# 镜像名: jiangming/panel
# 用法: ./build.sh [tag]
#   - 不传参数: 默认 tag 为 latest
#   - 传参数:   使用指定 tag，如 ./build.sh v1.0.0
# ============================================================

IMAGE_NAME="jiangming/panel"
TAG="${1:-latest}"
FULL_TAG="${IMAGE_NAME}:${TAG}"

# 切换到脚本所在目录（项目根目录）
cd "$(dirname "$0")"

echo "=========================================="
echo "  开始构建镜像: ${FULL_TAG}"
echo "=========================================="

docker build \
    --platform linux/amd64 \
    -t "${FULL_TAG}" \
    -f Dockerfile \
    .

if [ $? -eq 0 ]; then
    echo ""
    echo "=========================================="
    echo "  构建成功: ${FULL_TAG}"
    echo "=========================================="

    # 同时打上 latest 标签（如果传入的不是 latest）
    if [ "${TAG}" != "latest" ]; then
        docker tag "${FULL_TAG}" "${IMAGE_NAME}:latest"
        echo "  已同步打上 latest 标签"
    fi

    echo ""
    echo "  运行容器示例:"
    echo "    docker run -d -p 5678:5678 -v \$(pwd)/data:/app/data ${FULL_TAG}"
else
    echo ""
    echo "=========================================="
    echo "  构建失败，请检查上方错误信息"
    echo "=========================================="
    exit 1
fi