/**
 * Ex19: WebSocket Room + Hub Pattern — TypeScript Version
 * 
 * 🧠 So sánh key:
 * - Node.js: Dùng thư viện `ws` để tạo WS Server. Quản lý trạng thái kết nối và phân phòng bằng các cấu trúc 
 *            Map động trong JavaScript. Không lo sợ concurrent map write panic nhờ bản chất single-threaded.
 * - Go:      Dùng `gorilla/websocket` upgrade connection. Quản lý thread-safe bằng `sync.RWMutex` cho 
 *            các thao tác thêm/xoá connection và phát sóng (broadcast) đồng thời xuyên suốt các goroutine.
 * 
 * 💡 Sự khác biệt lớn nhất:
 * 1. Go bắt buộc dùng Mutex để khóa bảo vệ dữ liệu map chứa connection, nếu không sẽ panic lập tức.
 * 2. Cả hai ngôn ngữ đều bắt buộc phải có cơ chế dọn dẹp kết nối chết ngầm (zombie connections) bằng cách 
 *    sử dụng cơ chế Heartbeat Ping/Pong có định kỳ để tránh rò rỉ RAM (goroutine/connection leak).
 */

import WebSocket from 'ws';
import * as http from 'http';

interface ExtendedWebSocket extends WebSocket {
  isAlive?: boolean;
}

const server = http.createServer();
const wss = new WebSocket.Server({ noServer: true });

const connections = new Map<string, ExtendedWebSocket[]>();
const rooms = new Map<string, Set<ExtendedWebSocket>>();

function addConnection(address: string, ws: ExtendedWebSocket): void {
  if (!connections.has(address)) {
    connections.set(address, []);
  }
  connections.get(address)!.push(ws);
}

function removeConnection(address: string, ws: ExtendedWebSocket): void {
  const userConns = connections.get(address);
  if (userConns) {
    const index = userConns.indexOf(ws);
    if (index !== -1) {
      userConns.splice(index, 1);
    }
    if (userConns.length === 0) {
      connections.delete(address);
    }
  }

  for (const [roomName, wsSet] of rooms.entries()) {
    if (wsSet.has(ws)) {
      wsSet.delete(ws);
      if (wsSet.size === 0) {
        rooms.delete(roomName);
      }
    }
  }
}

function addToRoom(room: string, ws: ExtendedWebSocket): void {
  if (!rooms.has(room)) {
    rooms.set(room, new Set());
  }
  rooms.get(room)!.add(ws);
}

function broadcastToRoom(room: string, msg: any): void {
  const wsSet = rooms.get(room);
  if (wsSet) {
    const payload = JSON.stringify(msg);
    wsSet.forEach((ws) => {
      if (ws.readyState === WebSocket.OPEN) {
        ws.send(payload);
      }
    });
  }
}

server.on('upgrade', (request, socket, head) => {
  const url = new URL(request.url!, `http://${request.headers.host}`);
  const token = url.searchParams.get('token');

  if (!token) {
    socket.write('HTTP/1.1 401 Unauthorized\r\n\r\n');
    socket.destroy();
    return;
  }

  wss.handleUpgrade(request, socket, head, (ws) => {
    wss.emit('connection', ws, request, token);
  });
});

wss.on('connection', (ws: ExtendedWebSocket, request, token: string) => {
  const address = token; 
  ws.isAlive = true;

  addConnection(address, ws);

  ws.on('pong', () => {
    ws.isAlive = true;
  });

  ws.on('message', (message: string) => {
    try {
      const data = JSON.parse(message);
      if (data.action === 'join') {
        addToRoom(data.room, ws);
      } else if (data.action === 'broadcast') {
        broadcastToRoom(data.room, { sender: address, content: data.content });
      }
    } catch (err) {
      console.error('Invalid message format');
    }
  });

  ws.on('close', () => {
    removeConnection(address, ws);
  });
});

const interval = setInterval(() => {
  wss.clients.forEach((ws: ExtendedWebSocket) => {
    if (ws.isAlive === false) {
      ws.terminate();
      return;
    }
    ws.isAlive = false;
    ws.ping();
  });
}, 30000);

wss.on('close', () => {
  clearInterval(interval);
});

server.listen(8080, () => {
  console.log('WS Room Hub Server listening on port 8080');
});
