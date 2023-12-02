
CREATE TABLE `Hawkeye_admin` (
  `id` bigint(20) NOT NULL COMMENT '自增id',
  `user_name` varchar(255) NOT NULL DEFAULT '' COMMENT '用户名',
  `salt` varchar(50) NOT NULL DEFAULT '' COMMENT '盐',
  `password` varchar(255) NOT NULL DEFAULT '' COMMENT '密码',
  `create_at` datetime NOT NULL DEFAULT '1971-01-01 00:00:00' COMMENT '新增时间',
  `update_at` datetime NOT NULL DEFAULT '1971-01-01 00:00:00' COMMENT '更新时间',
  `is_delete` tinyint(4) NOT NULL DEFAULT '0' COMMENT '是否删除'
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='管理员表';

INSERT INTO `Hawkeye_admin` (`id`, `user_name`, `salt`, `password`, `create_at`, `update_at`, `is_delete`) VALUES
(1, 'admin', 'admin', '2823d896e9822c0833d41d4904f0c00756d718570fce49b9a379a62c804689d3', '2020-04-10 16:42:05', '2020-04-21 06:35:08', 0);

CREATE TABLE `poc_plugin` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '自增id',
  `vul_id` varchar(255) NOT NULL DEFAULT '' COMMENT 'ssvid',
  `vul_name` varchar(255) NOT NULL DEFAULT '' COMMENT '漏洞名',
  `vul_type` varchar(255) NOT NULL DEFAULT '' COMMENT '漏洞类型',
  `vul_date` datetime NOT NULL DEFAULT '1971-01-01 00:00:00' COMMENT '漏洞发布日期',
  `version` varchar(255) NOT NULL DEFAULT '1.0' COMMENT '插件版本',
  `author` varchar(255) NOT NULL DEFAULT '' COMMENT '编写者',
  `app_powerLink` varchar(255) NOT NULL DEFAULT '' COMMENT '产商链接',
  `app_name` varchar(255) NOT NULL DEFAULT '' COMMENT '应用名',
  `app_version` varchar(255) NOT NULL DEFAULT '' COMMENT '应用版本',
  `desc` varchar(255) NOT NULL DEFAULT '' COMMENT '漏洞描述',
  `cnnvd` varchar(255) NOT NULL DEFAULT '' COMMENT 'cnnvd',
  `cve_id` varchar(255) NOT NULL DEFAULT '' COMMENT 'cve_id',
  `rank` tinyint(4) NOT NULL DEFAULT '5' COMMENT '危险等级',
  `default_ports` varchar(255) NOT NULL DEFAULT '' COMMENT '默认端口',
  `default_service` varchar(255) NOT NULL DEFAULT '' COMMENT '默认服务',
  `content` text COMMENT '脚本内容',
  `create_at` datetime NOT NULL DEFAULT '1971-01-01 00:00:00' COMMENT '新增时间',
  `update_at` datetime NOT NULL DEFAULT '1971-01-01 00:00:00' COMMENT '更新时间',
  `is_delete` tinyint(4) NOT NULL DEFAULT '0' COMMENT '是否删除',
  primary key(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='Poc插件表';

CREATE TABLE `trap_plugin` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '自增id',
  `trap_id` varchar(255) NOT NULL DEFAULT '0000' COMMENT '蜜罐ID',
  `name` varchar(255) NOT NULL DEFAULT '' COMMENT '插件名',
  `author` varchar(255) NOT NULL DEFAULT '' COMMENT '编写者',
  `protocol` varchar(255) NOT NULL DEFAULT 'TCP' COMMENT '协议',
  `app_name` varchar(255) NOT NULL DEFAULT '' COMMENT '应用名',
  `honeypot` varchar(255) NOT NULL DEFAULT '' COMMENT '蜜罐名',
  `desc` varchar(1024) NOT NULL DEFAULT '' COMMENT '描述',
  `content` text COMMENT '脚本内容',
  `create_at` datetime NOT NULL DEFAULT '1971-01-01 00:00:00' COMMENT '新增时间',
  `update_at` datetime NOT NULL DEFAULT '1971-01-01 00:00:00' COMMENT '更新时间',
  `is_delete` tinyint(4) NOT NULL DEFAULT '0' COMMENT '是否删除',
  primary key(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='trap插件表';

CREATE TABLE `poc_task` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '自增id',
  `asset_id` bigint(20) NOT NULL DEFAULT '0' COMMENT '资产id',
  `portinfo_id` bigint(20) NOT NULL DEFAULT '0' COMMENT '端口id',
  `status` varchar(255) NOT NULL DEFAULT 'New' COMMENT '任务状态',
  `target_list` varchar(255) NOT NULL DEFAULT '' COMMENT '目标列表',
  `task_name` varchar(255) NOT NULL DEFAULT '' COMMENT '任务名',
  `recursion` tinyint(4) NOT NULL DEFAULT '0' COMMENT '周期',
  `create_at` datetime NOT NULL DEFAULT '1971-01-01 00:00:00' COMMENT '新增时间',
  `update_at` datetime NOT NULL DEFAULT '1971-01-01 00:00:00' COMMENT '更新时间',
  `plugin_list` varchar(255) NOT NULL DEFAULT '' COMMENT '插件列表',
  `is_delete` tinyint(4) NOT NULL DEFAULT '0' COMMENT '是否删除',
  primary key(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='Poc任务表';

CREATE TABLE `Hawkeye_task` (
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '自增id',
  `rule_id` bigint(20) NOT NULL DEFAULT '0' COMMENT '规则id',
  `name` varchar(255) NOT NULL DEFAULT '' COMMENT '任务名',
  `target_list` varchar(255) NOT NULL DEFAULT '' COMMENT '目标列表',
  `web_scan` tinyint(4) NOT NULL DEFAULT '0' COMMENT 'web扫描',
  `poc_scan` tinyint(4) NOT NULL DEFAULT '0' COMMENT '漏洞渗透',
  `auth_scan` tinyint(4) NOT NULL DEFAULT '0' COMMENT '权限爆破',
  `trap_scan` tinyint(4) NOT NULL DEFAULT '0' COMMENT '蜜罐识别',
  `recursion` tinyint(4) NOT NULL DEFAULT '0' COMMENT '周期',
  `progress` varchar(255) NOT NULL DEFAULT '' COMMENT '进程',
  `percent` tinyint(4) NOT NULL DEFAULT '0' COMMENT '百分比',
  `status` varchar(255) NOT NULL DEFAULT 'New' COMMENT '任务状态(New,Stop,Failed,Successfully)',
  `create_at` datetime NOT NULL DEFAULT '1971-01-01 00:00:00' COMMENT '新增时间',
  `update_at` datetime NOT NULL DEFAULT '1971-01-01 00:00:00' COMMENT '更新时间',
  `is_delete` tinyint(4) NOT NULL DEFAULT '0' COMMENT '是否删除',
  primary key(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='任务表';

CREATE TABLE `Hawkeye_asset` ( /**/
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '自增id',
  `task_id` bigint(20) NOT NULL DEFAULT '0' COMMENT '任务id',
  `ip` varchar(255) NOT NULL DEFAULT '' COMMENT 'ip地址',
  `gps` varchar(255) NOT NULL DEFAULT '' COMMENT 'GPS',
  `area` varchar(255) NOT NULL DEFAULT '' COMMENT '区域',
  `isp` varchar(255) NOT NULL DEFAULT '' COMMENT '运营商',
  `os` varchar(255) NOT NULL DEFAULT '' COMMENT '操作系统',
  `vendor` varchar(255) NOT NULL DEFAULT '' COMMENT '设备',
  `create_at` datetime NOT NULL DEFAULT '1971-01-01 00:00:00' COMMENT '新增时间',
  `is_delete` tinyint(4) NOT NULL DEFAULT '0' COMMENT '是否删除',
  primary key(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='资产表';

CREATE TABLE `Hawkeye_domain` ( /**/
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '自增id',
  `asset_id` bigint(20) NOT NULL DEFAULT '0' COMMENT '资产id',
  `domain` varchar(255) NOT NULL DEFAULT '' COMMENT '域名',
  `subdomain_list` varchar(255) NOT NULL DEFAULT '' COMMENT '子域列表',
  `registrar` varchar(255) NOT NULL DEFAULT '' COMMENT '注册商',
  `register_date` varchar(255) NOT NULL DEFAULT '' COMMENT '注册日期',
  `name_server` varchar(255) NOT NULL DEFAULT '' COMMENT 'DNS解析地址',
  `domain_server` varchar(255) NOT NULL DEFAULT '' COMMENT '域名解析器',
  `status` varchar(255) NOT NULL DEFAULT '' COMMENT '状态',
  primary key(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='域名表';

CREATE TABLE `Hawkeye_portinfo` ( /**/
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '自增id',
  `asset_id` bigint(20) NOT NULL DEFAULT '0' COMMENT '资产id',
  `port` varchar(255) NOT NULL DEFAULT '' COMMENT '端口',
  `name` varchar(255) NOT NULL DEFAULT '' COMMENT '服务(ssh,mysql,http)',
  `state` varchar(255) NOT NULL DEFAULT '' COMMENT '状态',
  `product` varchar(255) NOT NULL DEFAULT '' COMMENT '应用(OpenSSH,nginx)',
  `version` varchar(255) NOT NULL DEFAULT '' COMMENT '版本(7.4,syn-ack)',
  `extrainfo` varchar(255) NOT NULL DEFAULT '' COMMENT 'protocol 2.0,Servlet 2.5',
  `conf` varchar(255) NOT NULL DEFAULT '' COMMENT '10',
  `cpe` varchar(255) NOT NULL DEFAULT '' COMMENT '指纹cpe:/a:openbsd:openssh:7.4,10',
  primary key(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='服务表';

CREATE TABLE `Hawkeye_webinfo` ( /**/
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '自增id',
  `port_id` bigint(20) NOT NULL DEFAULT '0' COMMENT '端口id',
  `start_url` varchar(255) NOT NULL DEFAULT '' COMMENT '起始URL',
  `title` varchar(255) NOT NULL DEFAULT '' COMMENT '站点标题',
  `server` varchar(255) NOT NULL DEFAULT '' COMMENT 'Web服务器',
  `content_type` varchar(255) NOT NULL DEFAULT '' COMMENT '内容类型',
  `login_list` varchar(1024) NOT NULL DEFAULT '' COMMENT '登录页列表',
  `upload_list` varchar(1024) NOT NULL DEFAULT '' COMMENT '上传页面列表',
  `sub_domain` varchar(1024) NOT NULL DEFAULT '' COMMENT '子域名列表',
  `route_list` text COMMENT 'URL列表',
  `resource_list` text COMMENT '资源列表',
  primary key(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='Web信息表';

CREATE TABLE `Hawkeye_vulinfo` ( /**/
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '自增id',
  `asset_id` bigint(20) NOT NULL DEFAULT '0' COMMENT '资产id',
  `port_id` bigint(20) NOT NULL DEFAULT '0' COMMENT '端口id',
  `plugin_id` bigint(20) NOT NULL DEFAULT '0' COMMENT '插件ID',
  `verify_url` varchar(1024) NOT NULL DEFAULT '' COMMENT '漏洞验证URL',
  `verify_payload` text COMMENT '漏洞验证Payload',
  `verify_result` varchar(1024) NOT NULL DEFAULT '' COMMENT '漏洞验证Result',
  `exploit_url` varchar(1024) NOT NULL DEFAULT '' COMMENT '漏洞利用URL',
  `exploit_payload` text COMMENT '漏洞利用Payload',
  `exploit_result` varchar(1024) NOT NULL DEFAULT '' COMMENT '漏洞利用Result',
  `webshell_url` varchar(1024) NOT NULL DEFAULT '' COMMENT 'Webshell URL',
  `webshell_payload` text COMMENT 'Webshell Payload',
  `webshell_result` varchar(1024) NOT NULL DEFAULT '' COMMENT 'Webshell Result',
  `trojan_url` varchar(1024) NOT NULL DEFAULT '' COMMENT 'Trojan URL',
  `trojan_payload` text COMMENT 'Trojan Payload',
  `trojan_result` varchar(1024) NOT NULL DEFAULT '' COMMENT 'Trojan Result',
  `create_at` datetime NOT NULL DEFAULT '1971-01-01 00:00:00' COMMENT '新增时间',
  `is_delete` tinyint(4) NOT NULL DEFAULT '0' COMMENT '是否删除',
  primary key(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='漏洞表';

CREATE TABLE `Hawkeye_trap` ( /**/
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '自增id',
  `asset_id` bigint(20) NOT NULL DEFAULT '0' COMMENT '资产id',
  `port_id` bigint(20) NOT NULL DEFAULT '0' COMMENT '端口id',
  `plugin_id` bigint(20) NOT NULL DEFAULT '0' COMMENT '插件id',
  `verify` varchar(1024) NOT NULL DEFAULT '' COMMENT '验证项',
  `trap_id` varchar(255) NOT NULL DEFAULT '0000' COMMENT '蜜罐ID',
  `name` varchar(255) NOT NULL DEFAULT '' COMMENT '插件名',
  `protocol` varchar(255) NOT NULL DEFAULT 'TCP' COMMENT '协议',
  `app_name` varchar(255) NOT NULL DEFAULT '' COMMENT '应用名',
  `honeypot` varchar(255) NOT NULL DEFAULT '' COMMENT '蜜罐名',
  `desc` varchar(1024) NOT NULL DEFAULT '' COMMENT '描述',
  `create_at` datetime NOT NULL DEFAULT '1971-01-01 00:00:00' COMMENT '新增时间',
  primary key(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='蜜罐识别表';

CREATE TABLE `Hawkeye_auth` ( /**/
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '自增id',
  `asset_id` bigint(20) NOT NULL DEFAULT '0' COMMENT '资产id',
  `port_id` bigint(20) NOT NULL DEFAULT '0' COMMENT '端口id',
  `target` varchar(255) NOT NULL DEFAULT '' COMMENT '目标',
  `service` varchar(255) NOT NULL DEFAULT '' COMMENT '服务',
  `username` varchar(255) NOT NULL DEFAULT '' COMMENT '用户名',
  `password` varchar(255) NOT NULL DEFAULT '' COMMENT '密码',
  `command` varchar(255) NOT NULL DEFAULT '' COMMENT '验证命令',
  primary key(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='服务权限表';

CREATE TABLE `task_rule` ( /**/
  `id` bigint(20) NOT NULL AUTO_INCREMENT COMMENT '自增id',
  `trojan_cmd` varchar(255) NOT NULL DEFAULT '' COMMENT '木马注入payload的命令',
  `port_list` varchar(255) NOT NULL DEFAULT '' COMMENT '端口列表',
  `domain_dict` varchar(255) NOT NULL DEFAULT '' COMMENT '域名字典',
  `user_dict` varchar(255) NOT NULL DEFAULT '' COMMENT '用户字典',
  `passwd_dict` varchar(255) NOT NULL DEFAULT '' COMMENT '密码字典',
  `create_at` datetime NOT NULL DEFAULT '1971-01-01 00:00:00' COMMENT '新增时间',
  primary key(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COMMENT='服务权限表';