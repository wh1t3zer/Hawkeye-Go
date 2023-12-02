# 资产扫描
# 1、主机信息
# 2、端口及服务
# 3、对[域名]的进行域名字典爆破、whois查询
import os
import re
from modules.domain_brute import resolution, whois_query, DomainBrute
from modules.pynmap import NmapScanner
from modules.verify import Pocsuite
import json
from modules.web_scanner import WebScanner
from modules.trap_scanner import TrapScanner
from modules.ip_discern import getipinfo
from random import sample
from string import digits, ascii_lowercase
from micro.handler.common import LOCAL_RESULR_DIR


class AssetScanner:
    '''
    单个主机的扫描(不能输入端口及协议)
    Example: 192.168.1.11 or zan71.com.cn
    '''
    def __init__(
        self, domain='', ip='', scan_port='22,25,53,80,3306,5900,6379,7001,8080,8081,9000,27017', domain_dict=['www', 'mail', 'smtp'],
        webscan=False, verify=False, exploit=False, honeypot=False
    ):
        self.domain = self._format_domain(domain)
        self.ip = self._format_ip(ip)
        self.host = self._format_host()         # 匹配域名 or IP
        self.scan_port = scan_port
        self.honeypot = honeypot
        self.webscan = webscan
        self.verify = verify
        self.exploit = exploit
        self.domain_dict = domain_dict

    # 判定输入的是否合法域名
    def _format_domain(self, domain):
        pattern = re.compile(
            r'^(([a-zA-Z]{1})|([a-zA-Z]{1}[a-zA-Z]{1})|'
            r'([a-zA-Z]{1}[0-9]{1})|([0-9]{1}[a-zA-Z]{1})|'
            r'([a-zA-Z0-9][-_.a-zA-Z0-9]{0,61}[a-zA-Z0-9]))\.'
            r'([a-zA-Z]{2,13}|[a-zA-Z0-9-]{2,30}.[a-zA-Z]{2,3})$'
        )
        return domain if pattern.match(domain) else False

    # 判定输入的是否合法IP
    def _format_ip(self, ip):
        pattern = re.compile('^((25[0-5]|2[0-4]\d|[01]?\d\d?)\.){3}(25[0-5]|2[0-4]\d|[01]?\d\d?)$')
        if pattern.match(ip):
            return ip
        else:
            return False

    def _format_host(self):
        if not self.domain and not self.ip:
            print("Failed Scan Asset, Info: Invalid Domain Or IP.")
            return False
        return self.domain if self.domain else self.ip

    # 主机信息扫描
    def hostinfo_scan(self):
        '''
        1、域名解析、域名查询、IP信息、子域名爆破
        2、主机扫描(主机信息、端口)
        '''
        result = NmapScanner().get_detail(self.host, ports=self.scan_port)
        result['domain'] = None
        result['whois'] = None
        result['sub_domain'] = None
        if not result:
            print("[*] Failed Execute NmapScanner.")
            return

        if self.domain:  # 域名需要额外扫描
            # 1.1 解析域名
            resp = resolution(self.domain)
            if resp:
                self.ip = resp[self.domain][0]
            # 2.对主机进行扫描
            result['domain'] = self.domain
            result['whois'] = whois_query(self.domain)
            result['sub_domain'] = DomainBrute(self.domain, domain_dict=self.domain_dict).resolver()

        # 3.IP信息查询
        result['ip_info'] = getipinfo(self.ip)
        return result

    def Run(self):
        # 0.目标是否合法 (example.com | 192.168.4.195)
        if not self.host:
            return False
        # 1.主机扫描
        result = self.hostinfo_scan()
        if not result:
            return False

        # 2.web扫描与漏洞验证[]
        result['web_info'] = []
        result['vulnerability'] = []
        result['honeypot_info'] = []
        ports = result['ports']
        for port_info in ports:
            print("[-] scan port: ", port_info['port'])
            # web扫描
            if self.webscan:
                resp = None
                if self.domain:
                    resp = WebScanner(port=port_info['port'], domain=self.domain).Run()
                else:
                    resp = WebScanner(port=port_info['port'], ip=self.ip).Run()
                if resp:
                    result['web_info'].append(resp)

            # 漏洞验证及利用
            if self.verify:
                target = "{}:{}".format(self.ip, port_info['port'])
                ps = Pocsuite(target=target)
                resp = ps.Verify(self.exploit)
                result['vulnerability'].append(resp)

            # 蜜罐检测
            if self.honeypot:
                resp = TrapScanner(self.ip, port_info['port']).Run()
                if resp:
                    result['honeypot_info'].append()

        return result

    # 保存到文件夹中
    def save_result(self, result, dirs=None):
        if not result:
            print("Failed save resut, info: null data")
            return

        if not dirs:
            dirs = LOCAL_RESULR_DIR

        if not os.path.exists(dirs):
            os.makedirs(dirs)

        filepath = "{}/{}.json".format(dirs, ''.join(sample(digits + ascii_lowercase, 10)))
        try:
            with open(filepath, 'w') as f:
                f.write(json.dumps(result))
            return filepath
        except Exception as e:
            print("Failed save result to file, info: ", e)
        return None


if __name__ == "__main__":
    ass = AssetScanner(domain='zan71.com')
    resp = ass.Run()
    resp = json.dumps(resp)
    print(resp)

'''
result = {
    "domain": "",                 # add
    "sub_domain": "",             # add
    "hostname":"",
    "host":"172.31.50.254",
    "hostname_type":"",
    "vendor":"QEMU Virtual NIC",
    "mac":"52:54:00:AD:67:27",
    "ports":[
        {
            "protocol":"tcp",
            "port":"22",
            "name":"ssh",
            "state":"open",
            "product":"OpenSSH",
            "extrainfo":"protocol 2.0",
            "reason":"syn-ack",
            "version":"7.4",
            "conf":"10",
            "cpe":"cpe:/a:openbsd:openssh:7.4"
        }
    ]
}
'''
