# mitmproxy addon: answer the first N release update calls with an empty 502.
import time
from mitmproxy import http

def request(flow):
    if flow.request.method == "PATCH" and "/releases/" in flow.request.path:
        left = int(open("/tmp/fail_left").read())
        if left > 0:
            open("/tmp/fail_left", "w").write(str(left - 1))
            flow.response = http.Response.make(502, b"", {"Content-Type": "application/json"})
            print(f"{time.strftime('%H:%M:%S')} injected empty 502 for PATCH {flow.request.path}", flush=True)
