#!/bin/bash

set -e

echo "🚫 正在卸载 Go..."

# 1. 删除 Go 安装目录
if [ -d "/usr/local/go" ]; then
    sudo rm -rf /usr/local/go
    echo "✅ 已删除 /usr/local/go"
else
    echo "⚠️ 目录 /usr/local/go 不存在，跳过删除。"
fi

# 2. 清理环境变量（从 .bashrc 或 .zshrc 中）
PROFILE_FILE="$HOME/.bashrc"
if [ -n "$ZSH_VERSION" ] || [ -f "$HOME/.zshrc" ]; then
    PROFILE_FILE="$HOME/.zshrc"
fi

echo "🧹 正在清除环境变量配置..."

sed -i '/export PATH=\$PATH:\/usr\/local\/go\/bin/d' "$PROFILE_FILE"
sed -i '/export GOPATH=\$HOME\/go/d' "$PROFILE_FILE"
sed -i '/export PATH=\$PATH:\$GOPATH\/bin/d' "$PROFILE_FILE"

echo "🔄 已清除环境变量中的 Go 配置。"

# 3. 可选：删除 GOPATH 下的所有内容（慎用）
read -p "是否删除 GOPATH ($HOME/go)？[y/N]: " del_gopath
if [[ "$del_gopath" == "y" || "$del_gopath" == "Y" ]]; then
    rm -rf "$HOME/go"
    echo "✅ 已删除 $HOME/go"
fi

echo "🎉 Go 已彻底卸载完成。"

