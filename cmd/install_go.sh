#!/bin/bash

set -e

GO_VERSION=1.20
GO_TARBALL=go$GO_VERSION.linux-amd64.tar.gz
DOWNLOAD_URL=https://go.dev/dl/$GO_TARBALL

# 1. 下载 Go 安装包
echo "📦 正在下载 Go $GO_VERSION..."
wget -q --show-progress $DOWNLOAD_URL

# 2. 删除旧版本
echo "🧹 正在移除旧版本（如果存在）..."
sudo rm -rf /usr/local/go

# 3. 解压到 /usr/local
echo "📂 正在解压 Go 到 /usr/local..."
sudo tar -C /usr/local -xzf $GO_TARBALL

# 4. 设置环境变量
PROFILE_FILE="$HOME/.bashrc"
if [ -n "$ZSH_VERSION" ]; then
  PROFILE_FILE="$HOME/.zshrc"
fi

echo "🌍 正在配置环境变量..."
grep -q 'export PATH=$PATH:/usr/local/go/bin' $PROFILE_FILE || echo 'export PATH=$PATH:/usr/local/go/bin' >> $PROFILE_FILE
grep -q 'export GOPATH=$HOME/go' $PROFILE_FILE || echo 'export GOPATH=$HOME/go' >> $PROFILE_FILE
grep -q 'export PATH=$PATH:$GOPATH/bin' $PROFILE_FILE || echo 'export PATH=$PATH:$GOPATH/bin' >> $PROFILE_FILE

# 5. 生效配置
echo "🔄 使环境变量生效..."
source $PROFILE_FILE

# 6. 清理安装包
rm $GO_TARBALL

# 7. 测试
echo "✅ Go 安装完成，版本如下："
go version
