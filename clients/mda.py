import requests
import simplejson
import os
from time import sleep
from datetime import datetime, timedelta
from dateutil.tz import tzlocal
from schedule import Scheduler
from cachecontrol import CacheControl
from cachecontrol.caches.file_cache import FileCache

cache_l = os.path.join(os.path.dirname(__file__),'.web_cache')
sess = CacheControl(requests.Session(),
                    cache=FileCache(cache_l))


data = {
    "name": "test_job",
    "command": "echo 'HELLO'".split(' '),
    "epsilon": "PT5S",
    "pipeoutput":True
}

dt = datetime.isoformat(datetime.now(tzlocal()) + timedelta(0, 10))
data["schedule"] = "%s/%s/%s" % ("R2", dt, "PT10S")

if __name__ == "__main__":
    sess = requests.session()
    cached_sess = CacheControl(sess)

    response = cached_sess.get('http://darklordwill:4008/mda')
    r = response.json()
    print(r)
    s = Scheduler()
    s.add(data)