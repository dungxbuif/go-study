import net from 'node:net';

function handler(socket) {
  socket.write('Received: ');
  socket.pipe(socket);
}

const server = net.createServer(handler);
server.listen(3000, () => {
  console.log('TCP server listening on port 3000');
  // Auto-close for verification purposes if needed, 
  // but usually we just want to see it starts.
});
