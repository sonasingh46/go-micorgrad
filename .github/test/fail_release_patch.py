# mitmproxy addon: answer the release update call with a 502, pass everything else through.
import time
from mitmproxy import http

def request(flow):
    if flow.request.method == "PATCH" and "/releases/" in flow.request.path:
        mode = open("/tmp/mode").read().strip()
        body = b"" if mode == "empty" else b'{"message":"Server Error"}'
        flow.response = http.Response.make(502, body, {"Content-Type": "application/json"})
        print(f"{time.strftime('%H:%M:%S')} injected 502 ({mode} body) for {flow.request.method} {flow.request.path}", flush=True)
