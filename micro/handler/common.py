import os

BASE_DIR = os.path.dirname(os.path.dirname(os.path.dirname(os.path.abspath(__file__))))

MICRO_DIR = BASE_DIR + '/micro'

SERVICES_DIR = MICRO_DIR + '/service'

POCSUITE_SRV_DIR = SERVICES_DIR + '/asset_scanner'

POC_PLUGINS_DIR = POCSUITE_SRV_DIR + '/poc_plugins'

LOCAL_RESULR_DIR = POCSUITE_SRV_DIR + '/output'

DOMAIN_BRUTE_DICT = POCSUITE_SRV_DIR + '/data.json'

TRAP_PLUGINS_DIR = POCSUITE_SRV_DIR + '/trap_plugins'
