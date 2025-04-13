#!/bin/bash

# 定义变量
RELEASE_URL="https://github.com/csznet/mata/releases/latest/download"
INSTALL_DIR="/opt/mata"
SERVICE_FILE="/etc/systemd/system/mata.service"

# 检测系统架构
ARCH=$(uname -m)
if [[ "$ARCH" == "x86_64" ]]; then
    FILE_NAME="mata.zip"
elif [[ "$ARCH" == "aarch64" ]]; then
    FILE_NAME="mata-arm.zip"
else
    echo "不支持的架构: $ARCH"
    exit 1
fi

# 检查是否安装 unzip
if ! command -v unzip &>/dev/null; then
    echo "未检测到 unzip，请先安装 unzip 后再运行此脚本。"
    exit 1
fi

# 如果安装目录已存在，提示用户是否替换
if [[ -d "$INSTALL_DIR" ]]; then
    read -p "检测到安装目录已存在，是否替换？(y/n): " REPLACE
    if [[ "$REPLACE" != "y" ]]; then
        echo "安装已取消。"
        exit 0
    fi

    # 如果服务已存在，先停止服务
    if systemctl list-units --type=service | grep -q "mata.service"; then
        echo "停止现有服务..."
        sudo systemctl stop mata
    fi

    # 删除旧目录
    echo "删除旧的安装目录..."
    sudo rm -rf "$INSTALL_DIR"
fi

# 下载文件
echo "正在下载 $FILE_NAME..."
curl -L -o /tmp/$FILE_NAME "$RELEASE_URL/$FILE_NAME"

# 创建安装目录
echo "创建安装目录 $INSTALL_DIR..."
sudo mkdir -p $INSTALL_DIR

# 解压文件
echo "解压到 $INSTALL_DIR..."
sudo unzip -o /tmp/$FILE_NAME -d $INSTALL_DIR

# 清理临时文件
rm /tmp/$FILE_NAME

# 配置 systemctl 服务
echo "配置 systemctl 服务..."
sudo bash -c "cat > $SERVICE_FILE" <<EOL
[Unit]
Description=Mata Service
After=network.target

[Service]
ExecStart=$INSTALL_DIR/mata
WorkingDirectory=$INSTALL_DIR
Restart=always
User=root

[Install]
WantedBy=multi-user.target
EOL

# 重新加载 systemctl 并启动服务
echo "重新加载 systemctl ..."
sudo systemctl daemon-reload
sudo systemctl enable mata

echo "安装完成！"
echo "配置文件位置: $INSTALL_DIR/mata.json"
echo "启动服务命令: systemctl start mata"
echo "查看日志命令: journalctl -u mata -f"
