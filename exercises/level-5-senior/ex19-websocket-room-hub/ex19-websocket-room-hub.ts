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

import * as http from 'http';
import WebSocket from 'ws';

interface ExtendedWebSocket extends WebSocket {
   isAlive?: boolean;
}

const server = http.createServer();
const wss = new WebSocket.Server({ noServer: true });

// 🧠 AN TOÀN TRẠNG THÁI KHÔNG CẦN MUTEX (Node.js single-threaded maps):
// - Trong Node.js, `connections` và `rooms` được định nghĩa bằng `Map`.
// - Vì Javascript chạy đơn luồng, các thao tác thêm, xóa, cập nhật các map này xảy ra tuần tự.
// - Không bao giờ có hai dòng code cùng thay đổi Map cùng lúc, do đó ta hoàn toàn không cần các cơ chế khóa phức tạp
//   như `sync.RWMutex` trong Go.
// - Tuy nhiên, nếu ta thực hiện các thao tác bất đồng bộ (như await DB) ở giữa quá trình thêm/xóa,
//   vẫn có nguy cơ xảy ra Race Condition logic (lỗi logic nghiệp vụ), nhưng lỗi hỏng bộ nhớ hệ thống (Memory Corruption) thì được bảo vệ tuyệt đối.
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

// 🧠 CƠ CHẾ NÂNG CẤP KẾT NỐI TRONG EXPRESS/NODE.JS (upgrade HTTP Event):
// - Trong Node.js, http.Server phát ra sự kiện 'upgrade' khi nhận yêu cầu nâng cấp kết nối từ Client.
// - Ta bắt sự kiện này, xác thực token từ URL, rồi truyền quyền kiểm soát socket TCP qua `wss.handleUpgrade(...)`
//   để thư viện `ws` hoàn thành handshake nâng cấp lên giao thức WebSocket.
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
            broadcastToRoom(data.room, {
               sender: address,
               content: data.content,
            });
         }
      } catch (err) {
         console.error('Invalid message format');
      }
   });

   ws.on('close', () => {
      removeConnection(address, ws);
   });
});

// 🧠 DỌN DẸP ZOMBIE SOCKETS BẰNG HEARTBEAT:
// - Zombie connections xảy ra khi client tắt nguồn đột ngột, router trung gian ngắt kết nối mà không gửi tín hiệu TCP FIN.
// - Cả Express/Node.js và Go đều giải quyết bằng Ping/Pong định kỳ.
// - Ở đây, ta dùng một timer `setInterval` chạy mỗi 30s. Ta duyệt qua toàn bộ clients:
//   - Nếu `isAlive` là false (nghĩa là trong 30s qua không trả lời Ping bằng Pong), ta chủ động ngắt kết nối bằng `ws.terminate()`.
//   - Nếu còn sống, ta đặt lại `isAlive = false` và gửi tin nhắn `ws.ping()` kiểm tra.
// - Việc dọn dẹp này cực kỳ quan trọng để giải phóng File Descriptors của hệ điều hành, tránh lỗi "Too many open files" làm sập server.
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
