import consul
import socket
import random
from random import sample
from string import digits, ascii_lowercase


class GRPCServiceBase:
    def __init__(self, registry_addr, registry_port, server_addr, server_port):
        # 连接注册中心
        self._consul = consul.Consul(host=registry_addr, port=registry_port)
        self.server_addr = server_addr
        self.server_port = self._setattr(server_port)

    def _setattr(self, server_port):
        if not server_port:  # 没有服务提供端口
            return self._generate_port()
        else:  # 提供服务端口, 校验有无被占用
            return self._generate_port(server_port)

    def _generate_port(self, port=None):
        sock = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        try:
            while True:
                if not port:
                    port = random.randrange(50000, 59000)
                result = sock.connect_ex(('127.0.0.1', port))
                if result != 0:
                    return port
                elif port:  # 被占用并且是用户提供的
                    print("Failed bind port [::]:{}, Already in runing...".format(port))
                    return False
        except Exception as e:
            print("Failed generate server port, info: ", e)
            return False

    def RegisterService(self, name, host, port, tags=None):
        # 注册服务
        tags = tags or []
        srv_id = ''.join(sample(digits + ascii_lowercase, 10))
        self._consul.agent.service.register(name, service_id=srv_id, address=host, port=port, tags=tags, check=consul.Check().tcp(host, port, "5s", "30s", "30s"))

    def GetService(self, name):
        services = self._consul.agent.services()
        service = services.get(name)
        if not service:
            return None, None
        addr = "{0}:{1}".format(service['Address'], service['Port'])
        return service, addr
