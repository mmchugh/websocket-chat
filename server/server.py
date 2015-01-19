import asyncio
import websockets

class BroadcastHandler(object):
    clients = []

    @asyncio.coroutine
    def listener(self, websocket, path):
        self.clients.append(websocket)
        print("{} clients are connected".format(len(self.clients)))

        while True:
            message = yield from websocket.recv()
            print("> {}".format(message))

            if message is None:  # empty string is ok, so check None explictly
                self.clients.remove(websocket)
                print("Client left. {} clients are connected.".format(
                    len(self.clients)
                ))
                break

            yield from self.broadcast(message)

    def broadcast(self, message):
        # iterate over a copy of the list so that removing closed connections
        # won't cause some clients to be skipped
        for client in self.clients[:]:
            if not client.open:  # connection closed but didn't get removed
                self.clients.remove(client)
                continue

            yield from client.send(message)


if __name__ == "__main__":
    host, port = "localhost", "5678"
    handler = BroadcastHandler()
    start_server = websockets.serve(handler.listener, host, port)
    print("Starting server at {host}:{port}".format(host=host, port=port))

    asyncio.get_event_loop().run_until_complete(start_server)
    asyncio.get_event_loop().run_forever()
