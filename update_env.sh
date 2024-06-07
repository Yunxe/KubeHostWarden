#!/bin/bash

# 本地运行脚本
# 获取 en0 接口的 IPv4 地址
IP_ADDR=$(ifconfig en0 | grep "inet " | awk '{print $2}')

# 检查是否成功获取到 IP 地址
if [ -z "$IP_ADDR" ]; then
  echo "Failed to get IP address from en0 interface"
  exit 1
fi

# 更新 .env 文件中的 SSH_HOST, INFLUXDB_URL 和 MYSQL_ADDRESS 变量
ENV_FILE="./backend/.env"
if [ -f "$ENV_FILE" ]; then
  sed -i '' "s/^SSH_HOST=.*/SSH_HOST=\"$IP_ADDR\"/" "$ENV_FILE"
  sed -i '' "s|^INFLUXDB_URL=.*|INFLUXDB_URL=\"http://$IP_ADDR:8086\"|" "$ENV_FILE"
  sed -i '' "s/^MYSQL_ADDRESS=.*/MYSQL_ADDRESS=\"$IP_ADDR\"/" "$ENV_FILE"
  echo "Updated SSH_HOST, INFLUXDB_URL, and MYSQL_ADDRESS in $ENV_FILE to $IP_ADDR"
else
  echo "SSH_HOST=\"$IP_ADDR\"" > "$ENV_FILE"
  echo "INFLUXDB_URL=\"http://$IP_ADDR:8086\"" >> "$ENV_FILE"
  echo "MYSQL_ADDRESS=\"$IP_ADDR\"" >> "$ENV_FILE"
  echo "Created $ENV_FILE with SSH_HOST, INFLUXDB_URL, and MYSQL_ADDRESS set to $IP_ADDR"
fi
