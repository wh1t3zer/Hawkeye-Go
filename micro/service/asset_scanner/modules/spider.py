import re
import json
import memcache
import requests
from urllib.parse import urlparse

domainCN_pattern = re.compile(                                       # 中国开放注册的域名
    r'https?://[^\s]*.cn[a-zA-Z0-9/?=#&]*|'
    r'https?://[^\s]*.com[a-zA-Z0-9/?=#&]*|'
    r'https?://[^\s]*.top[a-zA-Z0-9/?=#&]*|'
    r'https?://[^\s]*.vip[a-zA-Z0-9/?=#&]*|'
    r'https?://[^\s]*.top[a-zA-Z0-9/?=#&]*|'
    r'https?://[^\s]*.xyz[a-zA-Z0-9/?=#&]*|'
    r'https?://[^\s]*.ltd[a-zA-Z0-9/?=#&]*|'
    r'https?://[^\s]*.art[a-zA-Z0-9/?=#&]*|'
    r'https?://[^\s]*.edu[a-zA-Z0-9/?=#&]*|'
    r'https?://[^\s]*.wang[a-zA-Z0-9/?=#&]*|'
    r'https?://[^\s]*.beer[a-zA-Z0-9/?=#&]*|'
    r'https?://[^\s]*.cloud[a-zA-Z0-9/?=#&]*|'
    r'https?://[^\s]*.store[a-zA-Z0-9/?=#&]*|'
    r'https?://[^\s]*.online[a-zA-Z0-9/?=#&]*'
)

# 本地资源|本域资源
# 1、 =https://antiy.cn/static/images/user.admin_1002.svg
# 2、 =//antiy.cn/static/images/user.admin_1002.svg
# src="./js-plugin/jquery-ui/jquery-ui-1.8.23.custom.min.js"
resource_local_pattern = re.compile(
    r'="\.?/?[a-zA-Z/-]*/?[a-zA-Z0-9._-]*\.js"|'                    # 编程语言
    r'="\.?/?[a-zA-Z/-]*/?[a-zA-Z0-9._-]*\.py"|'
    r'="\.?/?[a-zA-Z/-]*/?[a-zA-Z0-9._-]*\.sh"|'
    r'="\.?/?[a-zA-Z/-]*/?[a-zA-Z0-9._-]*\.go"|'
    r'="\.?/?[a-zA-Z/-]*/?[a-zA-Z0-9._-]*\.exe"|'
    r'="\.?/?[a-zA-Z/-]*/?[a-zA-Z0-9._-]*\.asp"|'
    r'="\.?/?[a-zA-Z/-]*/?[a-zA-Z0-9._-]*\.php"|'
    r'="\.?/?[a-zA-Z/-]*/?[a-zA-Z0-9._-]*\.java"|'
    r'="\.?/?[a-zA-Z/-]*/?[a-zA-Z0-9._-]*\.gz"|'                    # 包
    r'="\.?/?[a-zA-Z/-]*/?[a-zA-Z0-9._-]*\.tar"|'
    r'="\.?/?[a-zA-Z/-]*/?[a-zA-Z0-9._-]*\.jpg"|'                   # 前端静态资源
    r'="\.?/?[a-zA-Z/-]*/?[a-zA-Z0-9._-]*\.svg"|'
    r'="\.?/?[a-zA-Z/-]*/?[a-zA-Z0-9._-]*\.ico"|'
    r'="\.?/?[a-zA-Z/-]*/?[a-zA-Z0-9._-]*\.png"|'
    r'="\.?/?[a-zA-Z/-]*/?[a-zA-Z0-9._-]*\.css"|'
    r'="\.?/?[a-zA-Z/-]*/?[a-zA-Z0-9._-]*\.s?html?"|'
    r'="\.?/?[a-zA-Z/-]*/?[a-zA-Z0-9._-]*\.txt"|'                   # 配置文件
    r'="\.?/?[a-zA-Z/-]*/?[a-zA-Z0-9._-]*\.xml"|'
    r'="\.?/?[a-zA-Z/-]*/?[a-zA-Z0-9._-]*\.conf"|'
    r'="\.?/?[a-zA-Z/-]*/?[a-zA-Z0-9._-]*\.json"|'
    r'="\.?/?[a-zA-Z/-]*/?[a-zA-Z0-9._-]*\.yaml"|'
    r'="\.?/?[a-zA-Z/-]*/?[a-zA-Z0-9._-]*\.toml"|'
    r'="\.?/?[a-zA-Z/-]*/?[a-zA-Z0-9._-]*\.docx?"|'                 # 文档
    r'="\.?/?[a-zA-Z/-]*/?[a-zA-Z0-9._-]*\.pptx?"',
    re.IGNORECASE
)


class WebSpider:
    def __init__(self):
        self.login_page = set()                                  # 集合就已经自动去重
        self.resource_path = set()
        self.subdomain = set()
        self.urls = set()
        self.headers = {'user-agent': "Mozilla/5.0 (Windows NT 6.1; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/67.0.3396.99 Safari/537.36"}

    # URL发现
    def url_find(self, text, ip=None, domain=None):              # 无需提供端口
        if domain:  # 提取本域url
            pattern = re.compile('https?://[^\s]*{}:?[a-zA-Z0-9/\?=#&_\.]*'.format(domain))
            result = set(re.findall(pattern, text))
            self.urls = self.urls | result
        elif ip:
            pattern = re.compile('https?://[^\s]*{}:?[a-zA-Z0-9/\?=#&_\.]*'.format(ip))
            result = set(re.findall(pattern, text))
            self.urls = self.urls | result

    # 子域发现 域名专供
    def subdomain_find(self, domain, text):
        '''
        1、从url分析取的
        2、从信息取的 cdn.www.cnblogs.cn
        '''
        if not domain:
            return []
        pattern = re.compile('[a-zA-Z0-9\.]+{}'.format(domain))
        result = set(re.findall(pattern, text))
        self.subdomain = self.subdomain | result

    def resource_find(self, urls=None, text=None):
        # 1、本域|本IP资源发现, 从url中进行parse
        if urls:
            for url in urls:
                attr = urlparse(url)
                path = attr.path
                # login页面发现
                if len(path.split('login')) > 1 or len(path.split('sign')) > 1:
                    print(attr)
                    self.login_page.add("{}://{}{}".format(attr.scheme, attr.netloc, path.strip()))
                # 资源发现
                array = re.findall('.*\.[a-zA-Z]*', path)
                if array:
                    self.resource_path.add(self.strip_suffix(array[0]))

        # 2、本地资源发现 '<script type="text/javascript" src="./ss/jquery-ui-1.8.23.custom.min.pptx"></script>'
        elif text:
            array = re.findall(resource_local_pattern, text)
            array = map(self.strip_suffix, array)
            self.resource_path = self.resource_path | set(array)

    def strip_suffix(self, item):
        aa = re.sub('^="\.?/?|"$|^\./|^/', '', item)
        return aa

    def test(self, url, domain=None, ip=None):
        resp = requests.get(url, headers=self.headers, verify=False, timeout=5)
        text = resp.text
        if domain:
            self.url_find(text, domain=domain)   # url发现
            self.subdomain_find(domain, text)    # 子域匹配
            self.resource_find(urls=self.urls)   # 用url来匹配
            self.resource_find(text=text)        # 用网页text来匹配
        elif ip:
            self.url_find(text, ip=ip)
            self.resource_find(urls=self.urls)
            self.resource_find(text=text)
        print('==' * 40)
        print(json.dumps(list(self.login_page)))
        print('==' * 40)
        print(json.dumps(list(self.urls)))
        print('==' * 40)
        print(json.dumps(list(self.resource_path)))
        print('==' * 40)
        print(json.dumps(list(self.subdomain)))


class MemcacheHoneypot:
    def __init__(self):
        self.result = {'name': '', 'desc': []}

    def get_stats(self, stats):
        for item in stats:
            if type(item) != dict:
                continue
            if item['version'] == '1.4.25':
                self.result['name'] = "Memcache honeypot of dionaea"
                self.result["desc"].append("Non randomized features=1.4.25")
            if item['libevent'] == '2.0.22-stable':
                self.result['name'] = "Memcache honeypot of dionaea"
                self.result["desc"].append("Non randomized features=1.4.25")
            # if item['rusage_system'] = "0.233":
            if item['rusage_system'] == "0.110544":
                self.result['name'] = "Memcache honeypot of dionaea"
                self.result["desc"].append("Non randomized features=1.4.25")

    def run(self):
        mc = memcache.Client(['127.0.0.1:11211'], debug=True)
        stats = mc.get_stats()

        # 黑洞, 只进不出
        if not stats:
            self.result["name"] = "low interaction honeypot of memcache",
            self.result["desc"] = ["nothing data, it looks like a low interaction honeypot"]

        # 校验参数特征
        self.get_stats()

        return self.result


if __name__ == "__main__":
    ws = WebSpider()
    # ws.test("https://antiy.cn/", "antiy.cn")
    # ws.test("https://cnblogs.com/", "cnblogs.com")
    ws.test("http://127.0.0.1:8081/", ip="127.0.0.1")
