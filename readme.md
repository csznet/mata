# Mata

监控服务器自动切换解析工具  

使用tcp协议监控指定服务器，当服务器状态发生改变借助CloudFlare API进行改变解析，并支持Telegram通知  

## 示例场景  

A服务器访问速度快但是不带防御，B服务器带防御但是访问速度慢  

使用Mata对A服务器进行监控，当无法连通A服务器时，将自动把解析切换为B服务器，当A服务器恢复时，也会自动切换回A服务器  


## 使用教程  

### 准备工作  

**获取Cloudflare Zone API Token**  

https://dash.cloudflare.com/profile/api-tokens  

路径：Create Token -> Edit zone DNS

**获取域名Zone ID**  

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

**Target**  

需要监控的服务器，采用TCP监控，需带上端口号  

**Mata**  

Main为当服务器正常时的解析，Then为当服务器离线时的解析  

`proxied`为是否启用CloudFlare CDN

### 服务器运行  

下载系统对应的编译包，将`mata.sample.json`改名为`mata.json`并运行  


### 无服务器运行  

1. fork本项目，将`mata.sample.json`改名为`mata.json`  

2. 将fork的项目改为private  

3. 将项目`.github/workflows/mata.yml`中的注释去掉

项目将会使用Github Action每五分钟进行一次检测    
