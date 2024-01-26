import sys

sys.path.append("/Users/wh1t3zer/Public/Hawkeye")
import json
import argparse
from argparse import RawTextHelpFormatter
from micro.service.asset_scanner.handler.agent import PocScanServer
from micro.service.asset_scanner.handler.script import AssetScanner
from micro.handler.common import DOMAIN_BRUTE_DICT


# 分布式扫描节点
def agent(args):
    x = PocScanServer(args.registry_host, args.registry_port, args.server_name, args.server_addr, args.server_port)
    x.ListenAndServer()


# 离线脚本扫描
def script(args):
    if not args.domain and not args.ip:
        print('python3 main.py: error: missing arguments: --domain or --ip')
        return
    '''
    domain='', ip='', scan_port='22,25,53,80,3306,5900,6379,7001,8080,8081,9000,27017', domain_dict=['www', 'mail', 'smtp'],
        webscan=False, verify=False, exploit=False, honeypot=False
    '''
    domain_dict = None
    scan_ports = None
    conf_path = args.conf
    if not conf_path:
        conf_path = DOMAIN_BRUTE_DICT
    try:
        with open(conf_path, "r") as f:
            conf_data = json.loads(f.read())
            domain_dict = conf_data['domain_dict']
            scan_ports = ','.join(conf_data['ports'])
    except Exception as e:
        print("Failed loading conf_data, Info: ", e)
        return
    print("scan_ports", scan_ports, domain_dict, args.webscan)
    if args.domain:
        as_scanner = AssetScanner(
            domain=args.domain,
            scan_port=scan_ports,
            domain_dict=domain_dict,
            webscan=args.webscan,
            verify=args.verify,
            exploit=args.exploit,
            honeypot=args.honeypot
        )
    elif args.ip:
        as_scanner = AssetScanner(
            ip=args.ip,
            scan_port=scan_ports,
            webscan=args.webscan,
            verify=args.verify,
            exploit=args.exploit,
            honeypot=args.honeypot
        )
    else:
        print("[*] Unknow Error")
        return
    resp = as_scanner.Run()
    if not resp:
        print("[*] Response data nothing")
        return
    filename = as_scanner.save_result(resp, dirs=args.output)
    if filename:
        print("[*] Save result data at ", filename)
        # data = json.dumps(resp)
        # print(data)


def main():
    print("\033[33m\t\t---------- Asset Scanner ----------\033[0m")
    description = "\033[32mProvide Service of Airth by GRPC.\033[0m"
    example = "\n\nexample:\n\t[-] 向注册中心(127.0.0.1:8500)注册主机(192.168.1.11:50051)服务:\n\t" \
              "python3 agent.py agent --registry_host 127.0.0.1 --registry_port 8500 --server_name asset_scanner --server_addr 127.0.0.1\n\t" \
              "\n\nexample:\n\t[-] 离线扫描脚本:\n\t" \
              "python3 agent.py  --registry_host 127.0.0.1 --registry_port 8500 --server_name asset_scanner --server_addr 127.0.0.1\n\t"
    description += example
    parser = argparse.ArgumentParser(description=description, prog='python3 main.py',
                                     formatter_class=RawTextHelpFormatter)
    # parser = argparse.ArgumentParser(prog='python3 main.py')
    subparsers = parser.add_subparsers(help='sub-command help')
    # 添加子命令 agent
    parser_a = subparsers.add_parser('agent', help='agent help')
    parser_a.add_argument('--registry_host', required=True,
                          help='\t输入注册中心IP(必填)')  # 添加--verbose标签，标签别名可以为-v，这里action的意思是当读取的参数中出现--verbose/-v的时候
    parser_a.add_argument('--registry_port', required=True, help='\t输入注册中心端口(必填)')
    parser_a.add_argument('--server_name', required=True, help='\t输入服务名(必填)')
    parser_a.add_argument('--server_addr', required=True, help='\t输入绑定本机出口IP, 非回环IP或0.0.0.0(必填)')
    parser_a.add_argument('--server_port', type=int, required=False, help='\t输入服务端口(选填), 默认随机端口')
    # 设置默认函数
    parser_a.set_defaults(func=agent)
    # 添加子命令 script
    parser_s = subparsers.add_parser('script', help='script help')
    parser_s.add_argument('--domain', required=False,
                          help='\t目标域名: example.com (选填), domain与ip必选其一')  # 添加--verbose标签，标签别名可以为-v，这里action的意思是当读取的参数中出现--verbose/-v的时候
    parser_s.add_argument('--ip', required=False, help='\t目标IP 8.8.8.8 (选填), ip与domain必选其一')
    parser_s.add_argument('--net', required=False, help='\t目标网段 192.168.4.0/24 (选填), ip与domain必选其一')
    parser_s.add_argument('--webscan', type=bool, required=False, help='\t是否开启web扫描(选填), 默认False')
    parser_s.add_argument('--verify', type=bool, required=False, help='\t是否开启漏洞验证(选填), 默认False')
    parser_s.add_argument('--exploit', type=bool, required=False, help='\t是否开启漏洞利用(选填), 默认False')
    parser_s.add_argument('--honeypot', type=bool, required=False, help='\t输入蜜罐识别(选填), 默认False')
    parser_s.add_argument('--conf', required=False, help='\t指定数据源(选填), 默认./data.json')
    parser_s.add_argument('--output', required=False, help='\t存储扫描结果的本地目录(选填), 默认./output/')
    # 设置默认函数
    parser_s.set_defaults(func=script)
    args = parser.parse_args()
    # 执行函数功能
    args.func(args)


if __name__ == "__main__":
    main()
