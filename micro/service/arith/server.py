import sys
sys.path.append("/root/go/src/github.com/wh1t3zer/Hawkeye")
import re
import time
from concurrent import futures
import argparse
from argparse import RawTextHelpFormatter
from micro.proto.grpc.arith_pb2 import ArithResponse
from micro.proto.grpc.arith_pb2_grpc import ArithServicer, add_ArithServicer_to_server, grpc
from micro.handler.grpc_registry import GRPCServiceBase


class Arither(ArithServicer):
    def XiangJia(self, request, context):
        data = re.sub(r"[\n\r\t]", ", ", "{}".format(request))
        print("Recv Data From Client, Data: ", data)
        return ArithResponse(result=request.num1 + request.num2)

    def XiangJian(self, request, context):
        print("Recv Data From Client, Data: ", request)
        return ArithResponse(result=request.num1 - request.num2)


class ArithServer(GRPCServiceBase):
    def __init__(self, registry_host, registry_port, server_addr, server_port=None):
        super(ArithServer, self).__init__(registry_host, registry_port, server_addr, server_port)

    def ListenAndServer(self):
        if not self.server_port:
            print("Faild generate service port2.")
            return
        # 1.运行守护程序
        print("listen server on [::]:{}".format(self.server_port))
        server = grpc.server(futures.ThreadPoolExecutor(max_workers=10))
        add_ArithServicer_to_server(Arither(), server)
        server.add_insecure_port("[::]:{}".format(self.server_port))
        server.start()
        # 1.注册服务
        # self.RegisterService("test_python3", "172.31.50.249", self.server_port)
        self.RegisterService("test_python3", self.server_addr, self.server_port)
        try:
            while True:
                time.sleep(100)
        except KeyboardInterrupt:
            print("stop")
            server.stop(0)


def user_args():
    print("\033[33m\t\t---------- ArithServer ----------\033[0m")
    description = "\033[32mProvide Service of Airth by GRPC.\033[0m"
    example = "\n\nexample:\n\t[-] 向注册中心(127.0.0.1:8500)注册主机(192.168.1.11:50051)服务:\n\t" \
        "python3 server.py --registry_host 127.0.0.1 --registry_port 8500 --server_addr 192.168.1.11 --server_port 50051\n\t"
    description += example
    parser = argparse.ArgumentParser(description=description, prog='python2 main.py', formatter_class=RawTextHelpFormatter)                        # description参数可以用于插入描述脚本用途的信息，可以为空
    parser.add_argument('--registry_host', required=True, help='\tinput registry host')               # 添加--verbose标签，标签别名可以为-v，这里action的意思是当读取的参数中出现--verbose/-v的时候
    parser.add_argument('--registry_port', required=True, help='\tinput registry port')
    parser.add_argument('--server_addr', required=True, help='\tinput server port')
    parser.add_argument('--server_port', type=int, required=False, help='\tinput server port')

    args = parser.parse_args(sys.argv[1:])                                             # 将变量以标签-值的字典形式存入args字典
    return{'registry_host': args.registry_host, 'registry_port': args.registry_port, 'server_addr': args.server_addr, 'server_port': args.server_port}


if __name__ == "__main__":
    args = user_args()
    x = ArithServer(args['registry_host'], args['registry_port'], args['server_addr'], args['server_port'])
    x.ListenAndServer()
