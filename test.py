import requests
import json
import os

#import pdb; pdb.set_trace()

CACHE_PORT = 4000
CACHE_SERVER_URL = "http://localhost:{}/".format(CACHE_PORT)

key = "ram"

resp = requests.get(CACHE_SERVER_URL+key)

if resp.content:
	resp = json.loads(resp.content)
	print(resp)
else:
	print("Key doesn't exists.")

	


