# Mata

监控服务器自动切换解析工具

使用 tcp 协议监控指定服务器，当服务器状态发生改变借助 CloudFlare API 进行改变解析，并支持 Telegram 通知

## 示例场景

A 服务器访问速度快但是不带防御，B 服务器带防御但是访问速度慢（或是不带防御，当回退到 B 服务器时自动开启 CloudFlare CDN）

使用 Mata 对 A 服务器进行监控，当无法连通 A 服务器时，将自动把解析切换为 B 服务器，当 A 服务器恢复时，也会自动切换回 A 服务器

## 使用示例

见`mata.sample.json`

## 使用教程

### 准备工作

**获取 Cloudflare Zone API Token**

https://dash.cloudflare.com/profile/api-tokens

路径：Create Token -> Edit zone DNS

**获取域名 Zone ID**

打开域名控制台，右侧底部如图所示

<img width="452" alt="image" src="https://github.com/csznet/mata/assets/127601663/24b3ea58-afe0-40a5-9e15-9240c5ebd1fb">

### 参数说明

**ApiKey**

CloudFlare Zone API Token

**ZoneID**

CloudFlare ZoneID

**BotToken**

Telegram bot api token  
当服务器状态改变时发送通知，不启用保持为空即可

**ChatID**

Telegram ChatID

当服务器状态改变时发送通知，不启用保持为空即可

**Corn**

监控间隔，单位 秒

**Contcp**

通过外部 API 检测，如 `http://1.1.1.1/status/`
项目见 https://github.com/csznet/contcp/

**Mata**

Main 为当服务器正常时的解析，Then 为当服务器离线时的解析

`proxied`为是否启用 CloudFlare CDN

**Target**

需要监控的服务器，采用 TCP 监控，需带上端口号 （不带端口号则直接 PING）

**Main**

当 Target 在线时的解析

**Then**

当 Target 挂掉时切换的解析

### 服务器运行

下载系统对应的编译包，将`mata.sample.json`改名为`mata.json`并运行
