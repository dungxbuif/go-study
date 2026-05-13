import http from 'node:http';

function handler(request, response) {
  response.writeHead(200, { 'Content-Type': 'text/plain' });
  response.write('hello world');
  response.end();
}

const server = http.createServer(handler);
server.listen(8080, () => {
  console.log('HTTP server listening on port 8080');
});
