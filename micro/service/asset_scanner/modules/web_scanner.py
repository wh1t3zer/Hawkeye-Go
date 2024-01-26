'''爬虫:
1、爬取网页信息分析网站结构
- index.html
- login.jsp
- static/main.css
- static/main.ico
- static/main.js
- backend/article/detail/e1022.html
2、爬取子域名
3、(todo) 加入路径字典爆破 https://www.cnblogs.com/kaiho/p/8574143.html
4、(todo) 队列及爬取深度控制
example use cnblog
'''
import re
import requests
requests.packages.urllib3.disable_warnings()

# 本域资源
local_resource = re.compile(
    r'"[a-zA-Z0-9/-]+\.py"|'        # 编程语言
    r'"[a-zA-Z0-9/-]+\.js"|'
    r'"[a-zA-Z0-9/-]+\.jsp"|'
    r'"[a-zA-Z0-9/-]+\.php"|'
    r'"[a-zA-Z0-9/-]+\.html"|'
    r'"[a-zA-Z0-9/-]+\.ico"|'       # 前端资源
    r'"[a-zA-Z0-9/-]+\.svg"|'
    r'"[a-zA-Z0-9/-]+\.css"|'
    r'"[a-zA-Z0-9/-]+\.png"|'
    r'"[a-zA-Z0-9/-]+\.jpg"|'
    r'"[a-zA-Z0-9/-]+\.yml"|'
    r'"[a-zA-Z0-9/-]+\.json"|'
    r'"[a-zA-Z0-9/-]+\.gz"|'        # 包
    r'"[a-zA-Z0-9/-]+\.tar"|'
    r'"[a-zA-Z0-9/-]+\.rar"|'
    r'"[a-zA-Z0-9/-]+\.zip"|'
    r'"[a-zA-Z0-9/-]+\.xml"|'       # office
    r'"[a-zA-Z0-9/-]+\.doc"|'
    r'"[a-zA-Z0-9/-]+\.ppt"|'
    r'"[a-zA-Z0-9/-]+\.docx"|'
    r'"[a-zA-Z0-9/-]+\.pptx"'
)


class WebScanner:
    def __init__(self, host: str, port: int):  # 爬取深度、请求超时、重试次数、正则
        """
        host可以是域名或者是IP,
        端口是80、443的直接舍去
        """
        self.target = host if port in [80, 443] else '{}:{}'.format(host, port)
        self.headers = {'User-Agent': 'Mozilla/5.0 (Windows NT 5.1; rv:5.0) Gecko/20100101 Firefox/5.0'}
        self.start_url = self.get_start_url()          # 确保携带协议的target, 默认http
        # self.login_pattern = re.compile('https?://.*{}.*signin|https?://.*{}.*login'.format(self.target, self.target))                          # 抓取登录页面链接
        self.login_pattern = re.compile('https?://[a-zA-Z0-9\.]*[a-zA-Z0-9\.\-/_]*signin|https?://[a-zA-Z0-9\.]*[a-zA-Z0-9\.\-/_]*login')
        self.subdoamin_pattern = re.compile('([a-zA-Z0-9\.]+)\.{}'.format(host))                                                           # 抓取子域页面链路
        self.route_pattern = re.compile('(https?://[\w\.]*?%s:?\d{0,5}/[a-zA-Z0-9/]+)[#\?"\s]' % host)                                 # 抓取逻辑页面链接
        self.resource_pattern = re.compile('%s/([a-zA-Z0-9/]+\w*\.\w+)[\?"#]' % host)
        self.title_parttern = re.compile('<title>(.*?)</title>')
        self.upload_pattern = re.compile('https?://[a-zA-Z0-9\.]*[a-zA-Z0-9\.\-/_]*upload[a-zA-Z0-9\.\-_]*')
        self.result = {}

    def get_start_url(self):
        url = 'https://{}'.format(self.target)
        try:
            requests.get(url, headers=self.headers, timeout=5)
            return url
        except requests.exceptions.SSLError:
            return 'http://{}'.format(self.target)
        except Exception as e:
            print(e)
            return False

    # 获取title和内容类型
    def get_base(self, headers, text):
        server = ''
        content_type = ''
        if 'server' in headers.keys():
            server = headers['server']
        if 'content-type' in headers.keys():
            content_type = headers['content-type']
        self.result['start_url'] = self.start_url
        self.result['title'] = ''.join(self.title_parttern.findall(text))
        self.result['server'] = server
        self.result['content_type'] = content_type

    def _login_find(self, text) -> list:
        return list(set(self.login_pattern.findall(text)))

    def _sub_domain(self, text) -> list:
        return list(set(self.subdoamin_pattern.findall(text)))

    def _route_find(self, text) -> list:
        return list(set(self.route_pattern.findall(text)))

    def _upload_find(self, text) -> list:
        return list(set(self.upload_pattern.findall(text)))

    def _resource_find(self, text) -> list:
        array = []
        for rr in self.resource_pattern.findall(text):
            array.append(re.sub('^/', '', rr))
        for rr in local_resource.findall(text):
            array.append(re.sub('^/', '', rr.replace('"', '')))
        result = set(array)
        return list(result)

    def parse_html(self, url):
        try:
            resp = requests.get(url, headers=self.headers, verify=False, timeout=10)
            text = resp.content.decode('utf-8')
            text = text.replace('\\u002F', '/')
            # with open('12.html', 'w') as f:
            #     f.write(text)
            self.get_base(resp.headers, text)
            self.result['login_list'] = self._login_find(text)
            self.result['upload_list'] = self._upload_find(text)
            self.result['sub_domain'] = self._sub_domain(text)
            self.result['route_list'] = self._route_find(text)
            self.result['resource_list'] = self._resource_find(text)
        except Exception as e:
            print(e)
            return False
        return True

    def Run(self):
        if not self.start_url:
            return False
        # 分析url
        resp = self.parse_html(self.start_url)
        if not resp:
            return False
        return self.result


if __name__ == "__main__":
    hs = WebScanner('csdn.net', 80)
    a = hs.Run()
    import json
    print(json.dumps(a))


'''
{
    "start_url":"https://csdn.net",
    "title":"CSDN - 专业开发者社区",
    "server":"openresty",
    "content_type":"text/html; charset=utf-8",
    "login_list":[
        "https://g.csdnimg.cn/login-box/1.1.4/login",
        "https://passport.csdn.net/account/login"
    ],
    "upload_list":[],
    "sub_domain":[
        "api",
        "alice.blog",
        "ai",
        "intel",
        "cms"
    ],
    "route_list":[
        "https://blog.csdn.net/HaaSTech/article/details/112592687",
        "https://javaedge.blog.csdn.net/article/details/110790165",
        "https://www.csdn.net/nav/arch",
        "https://bss.csdn.net/m/topic/unity3d",
        "https://www.csdn.net/nav/ops"
    ],
    "resource_list":[
        "1.png",
        "so/search/s.do"
    ]
}
'''
