'''
识别蜜罐服务
link = "https://mp.weixin.qq.com/s?__biz=Mzk0NzE4MDE2NA==&mid=2247483908&idx=1&sn=e6a319e22c3cd54650bdbba511e58a43" \
    "&chksm=c37b85eff40c0cf91b6ceb75b07b80374380c010abb921db7685da59156dfc7a287f0c475cfe&mpshare=1&scene=23&srcid=0108kdZnMqp2EF3XSqvnNzxQ" \
    "&sharer_sharetime=1610091882557&sharer_shareid=22d54e05201de510604d4da97dc84192#rd"
'''
import os
import sys
sys.path.append('/root/go/src/github.com/wh1t3zer/Hawkeye')
from micro.handler import common
from pocsuite3.api import init_pocsuite, start_pocsuite, get_results


class TrapScanner:
    def __init__(self, target: list, trap_id: str = None, plugin_text: str = None, plugin_list: list = []):
        self.target = target
        self.trap_id = trap_id
        self.plugin_list = self._set_attr(plugin_text, plugin_list)

    def _set_attr(self, plugin_text, plugin_list):
        if plugin_text:
            filepath = "{}/{}.py".format(common.TRAP_PLUGINS_DIR, self.trap_id)
            with open(filepath, "w") as f:
                f.write(plugin_text)
            return set([filepath])
        else:
            return self._get_plugins(plugin_list)

    def _get_plugins(self, array) -> set:
        result = []
        # 若没有填写插件, 加载本地插件
        if not array:
            for root, dirs, files in os.walk(common.TRAP_PLUGINS_DIR):
                print("[*] 加载本地插件...")
                if not files:
                    print("[*] 没有检测到Poc插件！！！")
                    return
                return self._get_plugins(files)
        for poc in array:
            result.append("{}/{}".format(common.TRAP_PLUGINS_DIR, poc))
        return set(result)

    def run(self, expolit=False) -> list:
        result = []
        config = {
            'url': set(self.target),     # 不能传输协议
            'poc': set(self.plugin_list),
        }
        init_pocsuite(config)
        start_pocsuite()
        for item in get_results():
            if not item.result:
                continue
            result.append(item.result)
        return result


if __name__ == "__main__":
    ts = TrapScanner(['172.31.50.249:11211'], plugin_list=None)
    aa = ts.run()
    print(aa)
