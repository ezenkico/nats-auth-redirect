from http.server import BaseHTTPRequestHandler, HTTPServer
import json


HOST = "localhost"
PORT = 8000


class AuthHandler(BaseHTTPRequestHandler):
    def do_GET(self):
        auth_header = self.headers.get("Authorization", "")

        bearer_token = None
        if auth_header.startswith("Bearer "):
            bearer_token = auth_header.removeprefix("Bearer ").strip()

        print("---- Auth Request ----")
        print(f"Path: {self.path}")
        print(f"Authorization: {auth_header}")
        print(f"Bearer token: {bearer_token}")

        response = {
            "account": "APP",
            "perms": {
                "sub": [
                    {"allow": "requests"},
                ],
            }
        }

        data = json.dumps(response).encode("utf-8")

        self.send_response(200)
        self.send_header("Content-Type", "application/json")
        self.send_header("Content-Length", str(len(data)))
        self.end_headers()
        self.wfile.write(data)


if __name__ == "__main__":
    server = HTTPServer((HOST, PORT), AuthHandler)
    print(f"Listening on http://{HOST}:{PORT}")
    server.serve_forever()