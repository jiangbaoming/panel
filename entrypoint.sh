#!/bin/sh
set -e

# 默认 PUID/PGID
PUID=${PUID:-1000}
PGID=${PGID:-1000}

# 创建用户组（如不存在）
if ! getent group "$PGID" >/dev/null 2>&1; then
    addgroup -g "$PGID" appgroup
fi

# 创建用户（如不存在）
if ! getent passwd "$PUID" >/dev/null 2>&1; then
    adduser -D -u "$PUID" -G "$(getent group "$PGID" | cut -d: -f1)" appuser
fi

# 获取用户名
USERNAME=$(getent passwd "$PUID" | cut -d: -f1)

# 修正数据目录权限
chown -R "$PUID:$PGID" /app/data

# 以指定用户身份启动应用
exec su-exec "$USERNAME" "$@"