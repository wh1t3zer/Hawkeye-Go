import nmap
import json


class NmapScanner:
    def __init__(self):
        self.target = ''

    def _run(self, target, ports, args):
        self.target = target
        nm = nmap.PortScanner()
        try:
            if ports:
                return nm.scan(hosts=target, ports=ports, arguments=args)
            else:
                return nm.scan(hosts=target, arguments=args)
        except Exception as e:
            print(self.target, e)
            return False

    def parse_scanresult(self, scan_result):
        result = {}
        portinfo_list = []
        data = scan_result['scan'][self.target]

        # 1.获取设备信息
        if 'vendor' in data.keys():
            result['vendor'] = list(data['vendor'].values())[0] if len(data['vendor'].values()) else "Unknow"

        # 2.获取系统信息
        if 'osmatch' in data.keys():
            os_list = []
            for item in data['osmatch']:
                for os in item['osclass']:
                    os_list.append(os['osfamily'])
            if len(os_list) == 0:
                result['os'] = 'Unknow'
            else:
                result['os'] = max(os_list, key=os_list.count)
                result['vendor'] = result['os'] if result['vendor'] == 'Unknow' else result['vendor']

        # 3.获取服务列表
        for key, val in data['tcp'].items():
            if val['state'] == 'closed':
                continue
            portinfo_list.append({
                'port': key, 'name': val['name'], 'state': val['state'], 'product': val['product'],
                'version': val['version'], 'extrainfo': val['extrainfo'], 'conf': val['conf'], 'cpe': val['cpe']
            })
        result['portinfo_list'] = portinfo_list
        return result

    def get_detail(self, ip, ports='22-65535', args='-sV -O') -> dict:
        # 返回系统信息、端口列表等
        print("[+] Nmap get_detail Running....")
        scan_result = self._run(ip, ports, args)
        if not scan_result:
            return False
        print("[+] Nmap get_detail Finished....")
        return self.parse_scanresult(scan_result)

    def get_alive(self, net, args='-sn') -> list:
        # 返回网段存活的IP列表
        print("[+] Nmap get_alive Running....")
        scan_result = self._run(net, ports='', args=args)
        if not scan_result:
            return False
        data = scan_result['scan']
        print("[+] Nmap get_alive Finished....")
        return list(data.keys())


if __name__ == '__main__':
    scanner = NmapScanner()
    # detail
    ports_val = '22,25,53,80,111,631,3306,3389,5900,5901,7001,9200'
    data = scanner.get_detail('172.31.50.239', ports_val, '-sV -O')
    print(json.dumps(data))

    # alive
    # net = '172.31.50.0/24'
    # data = scanner.get_alive('172.31.50.0/24', '-sn')
    # print(data)

'''
args: -sU -sX -sC
22/tcp   open   ssh           OpenSSH 7.4 (protocol 2.0)
25/tcp   open   smtp          Postfix smtpd
80/tcp   open   http          nginx 1.18.0
111/tcp  open   rpcbind       2-4 (RPC #100000)
631/tcp  open   ipp           CUPS 1.6
3306/tcp open   mysql         MySQL 5.7.30

alive_list
{
    "Realtek Semiconductor": ["172.31.50.20"],
    "Evoc Intelligent Technology Co." : ["172.31.50.86"],
    "SAE IT-systems GmbH & Co. KG": ["172.31.50.122"],
    "Intel Corporate": ["172.31.50.119", "172.31.50.151"],
    "PC Partner":[
        "172.31.50.55", "172.31.50.102", "172.31.50.160", "172.31.50.222"
    ],
    "Dell": [
        "172.31.50.9", "172.31.50.28", "172.31.50.41"
        "172.31.50.108", "172.31.50.173", "172.31.50.224"
    ],
    "Dell Pcba Test": ["172.31.50.185", "172.31.50.186"],
    "VMware": [
        "172.31.50.4", "172.31.50.6", "172.31.50.46", "172.31.50.66", "172.31.50.72", "172.31.50.90",
        "172.31.50.92", "172.31.50.109", "172.31.50.145", "172.31.50.152", "172.31.50.182",
        "172.31.50.190","172.31.50.192", "172.31.50.225", "172.31.50.236", "172.31.50.239"
    ]
    "QEMU Virtual NIC": [
        "172.31.50.115", "172.31.50.141", "172.31.50.144", "172.31.50.146",
        "172.31.50.148", "172.31.50.178", "172.31.50.252", "172.31.50.253", "172.31.50.254",
    ]
}

'''
