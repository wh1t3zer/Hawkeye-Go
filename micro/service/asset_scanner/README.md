# 资产渗透

# Script 离线本地脚本扫描
desc: 扫描结果会保存至本地, 不会上传给控制端, 
running:
```java
Python3 main.py script \
    --domain           \    # 域名('example.com'), 域名/IP选其一
    --ip               \    # IP('192.168.4.10'), 域名/IP选其一
    --net              \    # 网段('192.168.4.0/24'), todo
    --scan_port        \    # 扫描端口('22,25,110-9000'), 默认热门端口
    --domain_dict      \    # 域名爆破字典('www, ftp, smtp, cdn, account'), 默认热门子域
    --webscan          \    # web扫描, 默认False
    --verify           \    # 漏洞验证, 默认False
    --exploit          \    # 漏洞利用, 默认False
    --honeypot         \    # 蜜罐识别, 默认False
    --output           \    # 结果输出路径, 默认./output
```


# Agent 分布式扫描节点(HA)
desc: 控制端发布探测任务, 选择节点进行扫描, 结果上传至控制端并在Web生成探测报告
running:
```java
Python3 main.py agent \
    --registry_host   \     # 服务注册中心主机('172.16.4.10')
    --registry_port   \     # 服务注册中心端口('8500')
    --server_name     \     # 服务名('asset_scanner')
    --server_addr     \     # 本机出口IP地址
    --server_port     \     # 默认随机端口
```
Example: python3 main.py agent --registry_host 172.31.50.249 --registry_port 8500 --server_name asset_scanner --server_addr 172.31.50.249