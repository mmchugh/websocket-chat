import asyncio
import sys
import websockets


class ChatClient(object):
    def send_message(self):
        message = sys.stdin.readline().rstrip()  # strip off the newline

        asyncio.async(self.websocket.send(message))

    @asyncio.coroutine
    def listener(self):
        self.websocket = yield from websockets.connect("ws://localhost:5678")
        print("Type + enter to send message; ctrl-c to quit")

        while True:
            message = yield from self.websocket.recv()

            if message is None:
                print("Lost connection to server, quitting.")
                break

            print("< {}".format(message))


client = ChatClient()
asyncio.get_event_loop().add_reader(sys.stdin, client.send_message)
asyncio.get_event_loop().run_until_complete(client.listener())
