from pypresence import Presence
import datetime as dt
import time


class RichPresenceManager:
    def __init__(self):
        self.client_id = "1144714449563955302"
        self.rpc = None
        self.running = False

    def start_presence(self):
        if not self.running:
            self.running = True
            self.rpc = Presence(self.client_id)
            try:
                self.rpc.connect()
                self.update_presence()
            except KeyboardInterrupt as error:
                print(error)
                self.rpc = None
                self.running = False
            except Exception as e:
                print(f"Error: Unable to connect to Rich Presence. {e}")
                self.rpc = None
                self.running = False

    def update_presence(self):
        if self.rpc:
            self.rpc.update(
                state="applio.org",
                details="Ultimate voice cloning tool.",
                buttons=[
                    {"label": "Home", "url": "https://applio.org"},
                    {"label": "Download", "url": "https://applio.org/download"},
                ],
                large_image="logo",
                large_text="Experimenting with applio",
                start=dt.datetime.now().timestamp(),
            )

    def stop_presence(self):
        self.running = False
        if self.rpc:
            self.rpc.close()
            self.rpc = None


RPCManager = RichPresenceManager()
