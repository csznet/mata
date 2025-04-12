# Mata

监控服务器自动切换解析工具

使用 tcp 协议监控指定服务器，当服务器状态发生改变借助 CloudFlare API 进行改变解析，并支持 Telegram 通知

## 示例场景

A 服务器访问速度快但是不带防御，B 服务器带防御但是访问速度慢（或是不带防御，当回退到 B 服务器时自动开启 CloudFlare CDN）

使用 Mata 对 A 服务器进行监控，当无法连通 A 服务器时，将自动把解析切换为 B 服务器，当 A 服务器恢复时，也会自动切换回 A 服务器

## 使用教程

https://n.csz.net/proj/mata/

## 生成配置

支持在线生成配置

https://mata-start.csz.net/
