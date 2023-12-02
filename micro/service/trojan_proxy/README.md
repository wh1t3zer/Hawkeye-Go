# 资产的木马服务

# Runing  
"""shell
./trojan1.0 \
--server_metadata id=5fb8e3dc11109ffb5e8cdc3e \
--server_metadata mq=127.0.0.1:6379 \
--server_name 5fb8e3dc11109ffb5e8cdc3e \
--registry_address 172.31.50.249:8500
"""

"""shell
./trojan1.0 --server_metadata id=5fb8e3dc11109ffb5e8cdc3e --server_metadata mq=127.0.0.1:6379 --server_name 5fb8e3dc11109ffb5e8cdc3e --registry_address 172.31.50.249:8500
"""

字段解释: 
registry_address(必填): 注册中心的地址, 直连通信
server_name(选填):      服务名(资产ID), 用于控制中心的发现与会话
mq(选填)：              双向消息推送队列, 备用通信
id(选填):               服务ID(资产ID), 通知对端当前木马服务所属资产
