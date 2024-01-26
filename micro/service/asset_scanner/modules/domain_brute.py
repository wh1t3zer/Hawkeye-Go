# 1、域名解析
# 2、爆破
# 3、whois查询域名详情
import re
import dns.resolver
import requests
import json


# 公网域名解析
def resolution(domain):
    # 最后一位必须是字母组合, 返回结果或空列表
    domain_pattern = re.compile(r'[0-9a-zA-Z]{1,}\.[a-z]{2,}$')
    if not domain_pattern.findall(domain):
        print('[*] Invalid domain ({})'.format(domain))
        return False

    result = {}
    record_a = []
    record_cname = []
    try:
        respond = dns.resolver.query(domain.strip())
        for record in respond.response.answer:
            for i in record.items:
                if i.rdtype == dns.rdatatype.from_text('A'):
                    record_a.append(str(i))
                    result[domain] = record_a
                elif i.rdtype == dns.rdatatype.from_text('CNAME'):
                    record_cname.append(str(i))
                    result[domain] = record_cname
    except Exception as e:
        print('[*] Failed reslove domain({})'.format(e))
        return False
    return result 


def whois_query(domain) -> dict:
    print("[+] Whois Query Running....")
    url = 'https://api.devopsclub.cn/api/whoisquery?domain={}&type=json'.format(domain)
    data = ""
    for i in range(10):
        try:
            resp = requests.get(url, headers={'User-Agent': 'Mozilla/5.0 (Windows NT 5.1; rv:5.0) Gecko/20100101 Firefox/5.0'})
            data = json.loads(resp.text)
            data = data['data']['data']
            break
        except Exception:
            continue
    print("[*] Whois Query Finished....")
    return data


class DomainBrute:
    def __init__(self, domain, domain_dict=['www', 'mail', 'smtp']):
        self.domain = domain
        self.subdomain_dict_2 = domain_dict
        # self.subdomain_dict_3 = ['www', 'mail', 'smtp']
        self.sub_domain = []

    # 加一层域名爆破
    def domain_handle(self):
        for sub_domain_2 in self.subdomain_dict_2:
            self.sub_domain.append(sub_domain_2.strip() + '.' + self.domain.strip())

    # 多线程解析器
    def resolver(self):
        print("[+] Domain Brute Running....")
        result = []
        self.domain_handle()
        for sub in self.sub_domain:
            resp = resolution(sub)
            if not resp:
                continue
            result.append(resp)
        print("[+] Domain Brute Finished....")
        return result  


if __name__ == "__main__":
    # 测试域名解析
    # print(resolution('47.92.255.39'))
    # print(resolution('antiy.cn'))

    # 测试域名字典爆破
    # db = DomainBrute(['mail', 'www', 'cdn', 'smtp'])
    # resp = db.resolver()
    # print(resp)
    # 测试whois查询
    # data = whois_query('')
    # if data:
    #     print(data)
    aa = ['DNS7.HICHINA.COM', 'DNS8.HICHINA.COM']
    dd = ','.join(aa)
    print(dd)
'''
{
    'creationDate': '2020-03-06T09:55:08Z',
    'domainName': '',
    'domainStatus': 'ok https://icann.org/epp#ok',
    'nameServer': ['DNS7.HICHINA.COM', 'DNS8.HICHINA.COM'],
    'registrar': 'Alibaba Cloud Computing (Beijing) Co., Ltd.',
    'registrarAbuseContactEmail': 'DomainAbuse@service.aliyun.com',
    'registrarAbuseContactPhone': '+86.95187',
    'registrarIANAID': '420',
    'registrarURL': 'http://www.net.cn',
    'registrarWHOISServer': 'grs-whois.hichina.com',
    'registryDomainID': '2500400073_DOMAIN_COM-VRSN',
    'registryExpiryDate': '2021-03-06T09:55:08Z',
    'updatedDate': '2020-03-06T10:05:12Z'
}
'''
