/**
 * Ex23: Object Pooling — TypeScript Version
 * 
 * 🧠 So sánh key:
 * - TypeScript: Trong Node.js/V8 Engine, ta tự tạo ra một Class ObjectPool thủ công. 
 *               V8 có Garbage Collector cực kỳ tối ưu, nhưng nếu ứng dụng phân bổ hàng triệu objects 
 *               mỗi giây (như game server, bộ đệm buffer lớn), V8 sẽ bị dính các đợt "Stop-The-World" 
 *               do GC dọn dẹp bộ nhớ (GC Minor/Major Sweeps).
 * - Go:         Hỗ trợ sẵn cấu trúc tối tân `sync.Pool` ở tầng thư viện chuẩn của ngôn ngữ.
 */

class UserRequest {
  public id: number;
  public data: Buffer;
  public processed: boolean;

  constructor() {
    this.id = 0;
    this.data = Buffer.alloc(1024); // Khởi tạo buffer 1KB
    this.processed = false;
  }

  // reset: Xóa sạch dữ liệu ô nhớ để tránh rò rỉ dữ liệu cũ sang request tiếp theo (Security Vulnerability)
  public reset(): void {
    this.id = 0;
    this.data.fill(0);
    this.processed = false;
  }
}

// ObjectPool: Triển khai bộ tái sử dụng đối tượng trong JavaScript.
//
// 🧠 QUẢN LÝ GC TRONG V8 ENGINE (Node.js GC vs Go GC Object Pooling):
// - V8 chia bộ nhớ Heap thành 2 thế hệ chính: New Space (Chứa các objects trẻ) và Old Space (Chứa các objects thọ).
// - **New Space** được dọn dẹp rất nhanh bằng thuật toán Scavenge (Copying GC). 
//   Tuy nhiên, nếu ta liên tục tạo và bỏ rơi các buffer dung lượng lớn, bộ nhớ New Space sẽ nhanh chóng đầy,
//   đẩy các đối tượng này lên **Old Space** trước thời hạn.
// - Dọn dẹp ở Old Space sử dụng thuật toán Mark-Sweep-Compact cực kỳ tốn CPU và gây trễ (Latency Spikes).
// - Bằng việc tự viết `ObjectPool` để lưu trữ và luân chuyển các đối tượng `UserRequest`, ta giúp V8 duy trì lượng
//   cấp phát bộ nhớ gần như bằng 0 (Zero-Allocation), tránh hoàn toàn các đợt GC Sweeps ở Old Space, bảo đảm latency ổn định mượt mà.
class ObjectPool {
  private pool: UserRequest[];

  constructor() {
    this.pool = [];
  }

  public get(): UserRequest {
    if (this.pool.length > 0) {
      return this.pool.pop()!;
    }
    console.log('[Pool] Allocating new UserRequest on V8 Heap.');
    return new UserRequest();
  }

  public put(obj: UserRequest): void {
    obj.reset(); // Bắt buộc phải dọn dẹp
    this.pool.push(obj);
  }
}

function main() {
  const pool = new ObjectPool();
  const payload = Buffer.from('API-Payload-Data');

  console.log('--- 1. First iteration (Pool is empty, allocating) ---');
  const activeRequests: UserRequest[] = [];
  for (let i = 1; i <= 3; i++) {
    const req = pool.get();
    req.id = i;
    payload.copy(req.data);
    req.processed = true;
    activeRequests.push(req);
  }

  // Trả lại pool sau khi xử lý xong
  activeRequests.forEach(req => pool.put(req));

  console.log('\n--- 2. Second iteration (Reusing recycled objects) ---');
  for (let i = 4; i <= 6; i++) {
    const req = pool.get();
    req.id = i;
    payload.copy(req.data);
    req.processed = true;
    // Xử lý nốt ...
    pool.put(req);
  }
}

main();
