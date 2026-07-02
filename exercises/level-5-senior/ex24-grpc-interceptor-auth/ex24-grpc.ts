/**
 * Ex24: gRPC Interceptor + Auth — TypeScript Version
 * 
 * 🧠 So sánh key:
 * - TypeScript: Trong Node.js, ta sử dụng thư viện `@grpc/grpc-js` để tạo gRPC Server. 
 *               Đăng ký các Interceptors bất đồng bộ để kiểm duyệt, giải mã gRPC metadata 
 *               được truyền từ client.
 * - Go:         Sử dụng Unary Interceptor tĩnh kết hợp ép kiểu từ `context.Context` cực kỳ mạnh mẽ.
 */

import * as grpc from '@grpc/grpc-js';

// 🧠 CƠ CHẾ GIAO TIẾP NHỊ PHÂN TRONG NODE.JS (gRPC JS Core under the hood):
// - Trong Node.js, thư viện cũ `@grpc/grpc-node` sử dụng add-on viết bằng C++ (binding).
// - Thư viện mới `@grpc/grpc-js` được viết hoàn toàn bằng TypeScript thuần (Pure JS), chạy trực tiếp trên
//   gói `http2` tích hợp sẵn của Node.js để quản lý luồng dữ liệu nhị phân (HTTP/2 streams).
// - Nhờ HTTP/2 Multiplexing, hàng ngàn cuộc gọi RPC đồng thời có thể được truyền tải trên DUY NHẤT một kết nối TCP vật lý,
//   giúp giảm tải băng thông và loại bỏ chi phí bắt tay TCP Handshake liên tục của REST API.
function authInterceptor(options: any, nextCall: any) {
  return new grpc.InterceptingCall(nextCall(options), {
    start: function (metadata: grpc.Metadata, listener: any, next: any) {
      console.log('[Interceptor] Intercepting gRPC outgoing/incoming call...');
      
      // Trích xuất metadata từ gRPC Metadata class
      const tokens = metadata.get('authorization');
      
      if (tokens.length === 0 || tokens[0] !== 'Bearer valid-token-123') {
        // Hủy bỏ RPC call với mã lỗi thích hợp
        const err = {
          code: grpc.status.UNAUTHENTICATED,
          details: 'Invalid or missing authentication token',
        };
        // Trong Node.js, ta kích hoạt hàm hủy bỏ truyền dữ liệu của HTTP/2 stream
        listener.status({ code: grpc.status.UNAUTHENTICATED, details: err.details });
        return;
      }
      
      next(metadata, listener);
    },
  });
}

function main() {
  console.log('gRPC Interceptor mock loaded successfully!');
  console.log('Mô hình Interceptor trong Node.js sử dụng kiến trúc Factory Pattern (`InterceptingCall`),');
  console.log('khá trừu tượng và phức tạp hơn so với Unary Interceptor mang tính hàm số tường minh của Golang.');
}

main();
