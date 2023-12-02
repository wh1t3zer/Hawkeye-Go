import time
from concurrent import futures
from micro.proto.grpc import aquaman_pb2  # import PocResponse, ResolvResponse, AnlsResponse, LocResponse, DetlResponse, AlvResponse
from micro.proto.grpc.aquaman_pb2_grpc import add_DomainServicer_to_server, DomainServicer, grpc
from micro.proto.grpc.aquaman_pb2_grpc import add_VulServicer_to_server, VulServicer
from micro.proto.grpc.aquaman_pb2_grpc import add_HostServicer_to_server, HostServicer
from micro.proto.grpc.aquaman_pb2_grpc import add_WebScrapServicer_to_server, WebScrapServicer
from micro.handler.grpc_registry import GRPCServiceBase
from modules.verify import Pocsuite
from modules.domain_brute import resolution, DomainBrute, whois_query
from modules.ip_discern import getipinfo
from modules.pynmap import NmapScanner
from modules.web_scanner import WebScanner
from modules.hydra_brute import HydraScanner
from modules.trap_scanner import TrapScanner


class HostUtils(HostServicer):
    def Location(self, request, context):
        resp = getipinfo(request.ip)
        for i in range(5):
            if resp:
                return aquaman_pb2.LocResponse(area=resp['area'], isp=resp['isp'], gps=resp['gps'])
            else:
                time.sleep(1)
                resp = getipinfo(request.ip)
        return aquaman_pb2.LocResponse()

    def Detail(self, request, context):
        ports_val = ",".join(request.ports)
        scanner = NmapScanner()
        resp = scanner.get_detail(request.ip, ports_val)
        if not resp:
            return aquaman_pb2.DetlResponse()

        array = []
        for item in resp['portinfo_list']:
            array.append({
                'port': str(item['port']), 'name': item['name'], 'state': item['state'], 'product': item['product'],
                'version': item['version'], 'extrainfo': item['extrainfo'], 'conf': item['conf'], 'cpe': item['cpe'],
            })
        print(resp)
        # return DetlResponse(os="zan71.com", vendor="linux", array=[])
        return aquaman_pb2.DetlResponse(os=resp["os"], vendor=resp["vendor"], array=array)

    def Alive(self, request, context):
        scanner = NmapScanner()
        resp = scanner.get_alive(request.net)
        return aquaman_pb2.AlvResponse(hosts=resp)


class DomainUtils(DomainServicer):
    def Resolv(self, request, context):
        resp = resolution(request.domain)
        if not resp:
            return aquaman_pb2.ResolvResponse()
        return aquaman_pb2.ResolvResponse(ip=resp[request.domain][0])

    def Analysis(self, request, context):
        print(request.domain, type(request.domain_dict))
        # 1、域名字典爆破
        resp = DomainBrute(request.domain, request.domain_dict).resolver()
        array = []
        for item in resp:
            array.append(list(item.keys())[0].split(".")[0])

        # 2、whois查询
        resp = whois_query(request.domain)
        if not resp:
            return aquaman_pb2.AnlsResponse()
        return aquaman_pb2.AnlsResponse(
            registrar=resp['registrar'], register_date=resp['creationDate'],
            name_server=','.join(resp['nameServer']), domain_server=resp['registrarWHOISServer'],
            status=','.join(resp['domainStatus']), subdomain_list=','.join(array)
        )


class WebUtils(WebScrapServicer):
    def Spider(self, request, context):
        ws = WebScanner(request.host, request.port)
        resp = ws.Run()
        if not resp:
            return aquaman_pb2.SpiResponse()
        return aquaman_pb2.SpiResponse(
            start_url=resp['start_url'], title=resp['title'], server=resp['server'],
            content_type=resp['content_type'], login_list=resp['login_list'],
            upload_list=resp['upload_list'], sub_domain=resp['sub_domain'],
            route_list=resp['route_list'], resource_list=resp['resource_list'],
        )


# 实现RPC类接口服务
class VulUtils(VulServicer):
    def Hydra(self, request, context):
        print(request.service, request.args, request.target_list, request.username_list, request.password_list)
        hs = HydraScanner(request.service, request.args, request.target_list, request.username_list, request.password_list)
        array = hs.run()
        # 一个服务可能有多个权限信息
        result = []
        for item in array:
            result.append({  # target只有IP
                'target': item['target'], 'service': item['service'], 'username': item['username'],
                'password': item['password'], 'command': item['command']
            })
        return aquaman_pb2.AuthResponse(array=result)

    def Trap(self, request, context):
        ts = TrapScanner(target=request.target_list, trap_id=request.trap_id, plugin_text=request.plugin_text)
        resp = ts.run()  # [{'verify': 'Non randomized features: version=1.4.25'}]
        result = []
        for item in resp:
            result.append({'verify': item['verify']})
        return aquaman_pb2.TrapResponse(array=result)

    def Verify(self, request, context):
        result = {
            'verify_url': "",
            'verify_payload': "",
            'verify_result': "",
            'exploit_url': "",
            'exploit_payload': "",
            'exploit_result': "",
            'webshell_url': "",
            'webshell_payload': "",
            'webshell_result': "",
            'trojan_url': "",
            'trojan_payload': "",
            'trojan_result': ""
        }
        # print("Recv Data From Client, Data: ", re.sub(r"[\n\r\t]", ",", "{}".format(request)))

        # request.command
        p = Pocsuite(target=request.target, vul_id=request.vul_id, poc_content=request.poc_content, asset_id=request.asset_id)
        # p = Pocsuite(request.target, request.poc_plugins, request.asset_id)
        resp = p.Verify(request.exploit)

        if not resp:
            return aquaman_pb2.PocResponse()
        if 'VerifyInfo' in resp.keys():
            result['verify_url'] = resp['VerifyInfo']['URL']
            result['verify_payload'] = resp['VerifyInfo']['PostData']
            result['verify_result'] = resp['VerifyInfo']['Result']
        if 'ExploitInfo' in resp.keys():
            result['exploit_url'] = resp['ExploitInfo']['URL']
            result['exploit_payload'] = resp['ExploitInfo']['PostData']
            result['exploit_result'] = resp['ExploitInfo']['Result']
        if 'WebshellInfo' in resp.keys():
            result['webshell_url'] = resp['WebshellInfo']['URL']
            result['webshell_payload'] = resp['WebshellInfo']['PostData']
            result['webshell_result'] = resp['WebshellInfo']['Result']
        if 'TrojanInfo' in resp.keys():
            result['trojan_url'] = resp['TrojanInfo']['URL']
            result['trojan_payload'] = resp['TrojanInfo']['PostData']
            result['trojan_result'] = resp['TrojanInfo']['Result']
        return aquaman_pb2.PocResponse(
            verify_url=result['verify_url'],
            verify_payload=result['verify_payload'],
            verify_result=result['verify_result'],
            exploit_url=result['exploit_url'],
            exploit_payload=result['exploit_payload'],
            exploit_result=result['exploit_result'],
            webshell_url=result['webshell_url'],
            webshell_payload=result['webshell_payload'],
            webshell_result=result['webshell_result'],
            trojan_url=result['trojan_url'],
            trojan_payload=result['trojan_payload'],
            trojan_result=result['trojan_result'],
        )


# 开启HTTP端口并提供服务
class PocScanServer(GRPCServiceBase):
    def __init__(self, registry_host, registry_port, server_name, server_addr, server_port=None):
        super(PocScanServer, self).__init__(registry_host, registry_port, server_addr, server_port)
        self.server_name = server_name

    def ListenAndServer(self):
        if not self.server_port:
            print("Faild generate service port2.")
            return
        # 1.运行守护程序
        print("listen server on [::]:{}".format(self.server_port))
        server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
        add_VulServicer_to_server(VulUtils(), server)
        add_DomainServicer_to_server(DomainUtils(), server)
        add_HostServicer_to_server(HostUtils(), server)
        add_WebScrapServicer_to_server(WebUtils(), server)
        server.add_insecure_port("[::]:{}".format(self.server_port))
        server.start()
        # 2.注册服务
        # self.RegisterService("test_python3", "172.31.50.249", self.server_port)
        self.RegisterService(self.server_name, self.server_addr, self.server_port)
        try:
            while True:
                time.sleep(100)
        except KeyboardInterrupt:
            print("stop")
            server.stop(0)
