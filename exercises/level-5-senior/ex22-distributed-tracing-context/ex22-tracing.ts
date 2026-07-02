/**
 * Ex22: Distributed Tracing + Context — TypeScript Version
 *
 * 🧠 So sánh key:
 * - TypeScript: Trong Node.js, do đặc trưng bất đồng bộ dựa trên callback/promise, ta khó truyền
 *               đối tượng context thủ công qua từng hàm.
 *               Giải pháp của Node.js là sử dụng `AsyncLocalStorage` (kế thừa từ thư viện `cls-hooked` truyền thống)
 *               để lưu trữ và truy vết Trace ID toàn cục ngầm trong luồng thực thi bất đồng bộ hiện tại.
 * - Go:         Tẩy chay hoàn toàn cơ chế "luồng ngầm/phép thuật" kiểu AsyncLocalStorage.
 *               Mọi hàm trong Go bắt buộc phải nhận tham số `ctx context.Context` đầu tiên một cách tường minh.
 *
 * 💡 Sự khác biệt lớn nhất:
 * 1. Go Context tường minh giúp mã nguồn dễ đọc, dễ hiểu dòng chảy dữ liệu, biên dịch tối ưu cực tốt.
 * 2. Node.js AsyncLocalStorage mang lại sự tiện lợi, giúp lập trình viên không cần sửa signature của hàng ngàn hàm có sẵn.
 */

import { AsyncLocalStorage } from 'async_hooks';

// Khởi tạo AsyncLocalStorage đóng vai trò lưu trữ ngầm Trace ID xuyên suốt chuỗi gọi bất đồng bộ
const asyncLocalStorage = new AsyncLocalStorage<Map<string, string>>();

// 🧠 CƠ CHẾ CỦA ASYNCLOCALSTORAGE DƯỚI NẮP CAPO (Node.js AsyncLocalStorage):
// - Node.js sử dụng cơ chế `async_hooks` của V8 Engine để theo dõi vòng đời của tất cả các tài nguyên bất đồng bộ (Promises, Timeouts, I/O).
// - Khi một Promise mới được sinh ra, Node.js sẽ liên kết (hook) và chuyển giao (propagate) vùng nhớ cục bộ của luồng cha sang cho luồng con ngầm định.
// - Điều này cho phép ta lưu `trace_id` tại middleware đầu tiên, và ở bất cứ hàm sâu thẳm nào trong project, ta chỉ cần gọi
//   `asyncLocalStorage.getStore()?.get('trace_id')` là lấy được giá trị, tương tự như biến cục bộ ThreadLocal của Java.
// - ⚠️ Điểm yếu: Việc kích hoạt `async_hooks` làm giảm hiệu năng thực thi của Event Loop (khoảng 5-10% tùy phiên bản Node.js)
//   do phải duy trì bản đồ liên kết vòng đời bất đồng bộ khổng lồ trong bộ nhớ.
function runWithTrace(traceId: string, fn: () => void) {
   const store = new Map<string, string>();
   store.set('trace_id', traceId);
   asyncLocalStorage.run(store, fn);
}

function getTraceID(): string {
   const store = asyncLocalStorage.getStore();
   return store?.get('trace_id') || '';
}

// Giả lập cuộc gọi mạng sang Service B
function callServiceB() {
   const traceId = getTraceID();
   console.log(`[Service A] Calling Service B. Injecting TraceID: ${traceId}`);

   // Chuẩn bị Header theo W3C Trace Context Standard
   const headers = {
      traceparent: `00-${traceId}-0000000000000000-01`,
   };

   // --- Giả lập truyền tải qua mạng ---
   simulateReceiveAtServiceB(headers);
}

function simulateReceiveAtServiceB(headers: Record<string, string>) {
   const traceParent = headers['traceparent'];
   let receivedTraceId = '';
   if (traceParent && traceParent.startsWith('00-')) {
      receivedTraceId = traceParent.split('-')[1];
   }

   // Khởi động vùng nhớ AsyncLocalStorage cho Service B nhận dữ liệu
   runWithTrace(receivedTraceId, () => {
      console.log(`\n[Service B] Received request at Service B.`);
      console.log(`[Service B] Extracted TraceID: ${getTraceID()}`);
      console.log(
         `[Service B] Writing logs... TraceID is automatically appended to log outputs!`,
      );
   });
}

function main() {
   const traceId = '4bf92f3577b34da6a3ce929d0e0e4736';

   runWithTrace(traceId, () => {
      console.log(
         `[Service A] Request started. TraceID is active in context: ${getTraceID()}`,
      );
      callServiceB();
   });
}

main();
