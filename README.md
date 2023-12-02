# Hawkeye 攻击面暴露系统

# 一、程序功能
### 1.1 资产收集
|||
|--|--|
|设备类型	|识别Dell,、QEMU、VMware等|
|操心系统	|识别Linux、Windows等|
|IP信息	|运营商、GPS、区域位置信息等|
|域名信息|	子域及关联域名、域名注册商、注册日期、DNS解析地址、域名解析器、域名状态等|
|资产信息|	服务、端口、协议及其状态|
### 1.2 Web扫描
|||
|--|--|
|登录页面	|支持爬取、爆破的方式收录登录页面(login、signin)
|上传页面	|支持爬取、爆破的方式收录上传页面(upload)
|子域名	|支持爬取、爆破的方式收录子域名
|本域链接	|支持爬取、爆破的方式收录本域链接
|资源链接	|支持爬取、爆破的方式收录资源链接(js,css,jpg...)
|Web服务器	|如: nginx、openresty、apache、jsp
|主域标题	|如: 502 Bad Gateway
|网页内容类型|	如: text/html; charset=utf-8
### 1.3 漏洞扫描
|||
|--|--|
|漏洞验证	|支持基于漏洞原理检测并验证漏洞
|漏洞利用	|支持对已验证漏洞进行攻击利用
|服务代理|	支持对目标服务进行代理。下次仅需访问代理即可，无需直接访问目标
|后渗透	|对漏洞进行植入木马，建立c/s通信
|Session	|支持对取得控制权的主机进行本地攻击
|能力扩展|	支持用户自定义Poc插件进行能力扩展(Web提交即可,无需重启服务)
### 1.4 权限扫描
|||
|--|--|
|资产服务|	支持ssh、ftp、vnc、mysql、redis、mongodb等服务的权限扫描爆破
|用户字典	|支持用户自定义权限用户字典
|密码字典|	支持用户自定义权限密码字典
### 1.5 蜜罐识别
|||
|--|--|
|开源蜜罐|	支持Dionaea、conpot、Amun、Nependthes、Cowrie等蜜罐识别
|蜜罐服务|	支持SSH、FTP、S7、IMAP、Telnet等十余种蜜罐服务识别
|能力扩展|	支持用户自定义脚本插件进行能力扩展(Web提交即可,无需重启服务)
### 1.6 系统管理
|||
|--|--|
|攻击面展示|	支持任务扫描实时展示、支持资产任务聚合展示、支持首页大盘实时数据展示|
|任务配置|	支持配置限定任务执行时间, 并立即生效
||支持配置任务并发数配置, 并立即生效
||支持配置任务失败重试次数, 并立即生效
||支持配置自定义探测端口列表, 并立即生效
||支持配置自定义域名字典, 并立即生效
||支持爬虫模式选择(网页渲染型、多线程爬虫), 并立即生效
||支持配置Web扫描的线程数、超时时间、失败重试次数、超时重试次数, 并立即生效
||支持cookie爬虫、事件触发型爬虫, 并立即生效
||支持配置后台登录用户密码, 并立即生效
||支持配置Web扫描的黑白名单(如后缀、路径), 并立即生效
||支持自定义选择Poc插件列表, 并立即生效
||支持自定义选择蜜罐识别插件列表, 并立即生效
||支持配置权限爆破的线程数, 并立即生效
||支持配置权限爆破的超时时间、失败重试、超时重试次数, 并立即生效
||支持配置权限爆破的服务字典、用户字典、密码字典等, 并立即生效
|系统配置|	支持目标配置、网络参数配置、攻击模式配置等
||	支持系统数据备份和恢复
||	支持定时备份与手动备份
|漏洞管理	|支持认证凭据配置
||	支持历史漏洞检索及后渗透、会话操作
||	支持查询在线木马(肉鸡)服务及通讯
||	支持建立服务代理
||	支持Poc插件的管理(CURD)
||	支持蜜罐识别插件的管理(CURD)

# 二、主要功能演示(GIF)
详情请看说明书.docx

# 三、程序部署
### 3.1 Redis服务
[root@localhost Hawkeye]# docker run -itd -p 6379:6379 --name redis redis
### 3.2 Mysql服务
启动mysql服务  
[root@localhost Hawkeye]# docker run -itd -p 3306:3306 --mame mysql mysql:5.7
使用备份库进行恢复
[root@localhost Hawkeye]# docker cp database.sql mysql:/root  
[root@localhost Hawkeye]# docker exec -it mysql /bin/bash  
root@63xeqdar98d7w7ed ~# mysql -uroot -p  
mysql> create database Hawkeye;  
mysql> use Hawkeye;  
mysql> source /root/database.sql;  
### 3.2 注册中心
[root@localhost Hawkeye]# docker run -itd -p 8300:8300 -p 8301:8301 -p 8301:8301/udp -p 8302:8302 -p8400:8400 -p 8500:8500 -p 53:53/udp --name consul consul

### 3.3 主节点部署
确认如下配置文件参数
- 基础配置文件 conf/dev/base.toml
- 服务配置文件 conf/dev/micro.toml
- Mysql配置文件 conf/dev/mysql_map.toml
- Redis配置文件 conf/dev/redis_map.toml
- 扫描器配置文件 micro/service/asset_scanner/settings.py
源码安装部署
- [root@localhost Hawkeye]# go mod tidy  
[root@localhost Hawkeye]# go build .  
[root@localhost Hawkeye]# ./Hawkeye  

### 3.4 从节点部署
Python库安装：pip3 install -r requirement.txt  
工具安装: yum install hydra nmap -y  
配置修改: settings.py  
```shell
Python3 main.py agent \
    --registry_host   \     # 服务注册中心主机('172.16.4.10')
    --registry_port   \     # 服务注册中心端口('8500')
    --server_name     \     # 服务名('asset_scanner')
    --server_addr     \     # 本机出口IP地址
    --server_port     \     # 默认随机端口
Example: python3 main.py agent --registry_host 172.31.x.x --registry_port 8500 --server_name asset_scanner --server_addr 172.31.x.x
```
### 3.5 Web前端部署
Docker部署(主要修改vue.config.js第43行proxy-target)  
[root@localhost ~]# docker run -itd -v /home/views/Hawkeye-view/vue.config.js:/opt/vue.config.js -p 9527:9527 --name Hawkeye-view jstang/Hawkeye-view:2.0  
源码部署  
git clone https://github.com/wh1t3zer/Hawkeye-view  
npm run dev  

### 3.6 API接口文档
访问地址: http://127.0.0.1:8700/swagger/index.html  
