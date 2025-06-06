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
        systemctl stop mata
    fi

    # 删除旧目录
    echo "删除旧的安装目录..."
    rm -rf "$INSTALL_DIR"
fi

# 下载文件
echo "正在下载 $FILE_NAME..."
curl -L -o /tmp/$FILE_NAME "$RELEASE_URL/$FILE_NAME"

# 创建安装目录
echo "创建安装目录 $INSTALL_DIR..."
mkdir -p $INSTALL_DIR

# 解压文件
echo "解压到 $INSTALL_DIR..."
unzip -o /tmp/$FILE_NAME -d $INSTALL_DIR

# 清理临时文件
rm /tmp/$FILE_NAME

# 配置 systemctl 服务
echo "配置 systemctl 服务..."

SERVICE_CONTENT="[Unit]
Description=Mata Service
After=network.target

[Service]
ExecStart=$INSTALL_DIR/mata
WorkingDirectory=$INSTALL_DIR
Restart=always
User=root

[Install]
WantedBy=multi-user.target
"

if [[ -f "$SERVICE_FILE" ]]; then
    # 检查内容是否一致
    EXISTING_CONTENT=$(cat "$SERVICE_FILE")
    if [[ "$EXISTING_CONTENT" == "$SERVICE_CONTENT" ]]; then
        echo "$SERVICE_FILE 已存在且内容一致，跳过写入。"
    else
        read -p "$SERVICE_FILE 已存在，内容不同，是否覆盖？(y/n): " OVERWRITE
        if [[ "$OVERWRITE" == "y" ]]; then
            echo "$SERVICE_CONTENT" >"$SERVICE_FILE"
            echo "已覆盖 $SERVICE_FILE"
        else
            echo "未覆盖 $SERVICE_FILE"
        fi
    fi
else
    echo "$SERVICE_CONTENT" >"$SERVICE_FILE"
    echo "已创建 $SERVICE_FILE"
fi

# 重新加载 systemctl 并启动服务
echo "重新加载 systemctl ..."
systemctl daemon-reload
systemctl enable mata

echo "安装完成！"
echo "配置文件位置: $INSTALL_DIR/mata.json"
echo "启动服务命令: systemctl start mata"
echo "查看日志命令: journalctl -u mata -f"
