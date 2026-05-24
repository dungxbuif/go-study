/**
 * Ex05: Goroutine + WaitGroup — TypeScript Version
 * 
 * 🧠 So sánh key:
 * - Node.js: Chạy phi tập trung/bất đồng bộ trên Single-threaded Event Loop. 
 *            Để thực thi song song các tác vụ I/O, ta dùng `Promise.all()`.
 *            Không cần lo lắng Race Condition khi thay đổi biến dùng chung (shared array/object) 
 *            vì JavaScript không có đa luồng CPU chạy song song thực sự trong cùng một Node process.
 * - Go:      Chạy đa luồng thực sự (multi-threaded) nhờ cơ chế Go Scheduler quản lý hàng vạn Goroutine siêu nhẹ.
 *            Cần `sync.WaitGroup` để đồng bộ hoá (đợi toàn bộ Goroutine hoàn thành).
 *            Bắt buộc sử dụng `sync.Mutex` để bảo vệ dữ liệu dùng chung (shared slice/map) khỏi Race Condition.
 */

import axios from 'axios';

const urls: string[] = [
  'https://httpbin.org/status/200',
  'https://httpbin.org/status/201',
  'https://httpbin.org/status/404',
  'https://httpbin.org/status/500',
  'https://httpbin.org/delay/1',
  'https://httpbin.org/delay/2',
];

interface URLResult {
  url: string;
  status: number;
  duration: string;
  error: string | null;
}

async function checkUrl(url: string): Promise<URLResult> {
  const start = Date.now();
  try {
    const response = await axios.get(url, { timeout: 5000 });
    const duration = Date.now() - start;
    return { url, status: response.status, duration: `${duration}ms`, error: null };
  } catch (error: any) {
    const duration = Date.now() - start;
    return { 
      url, 
      status: error.response ? error.response.status : 0, 
      duration: `${duration}ms`, 
      error: error.message 
    };
  }
}

async function main(): Promise<void> {
  console.log('=== Checking URLs in Parallel (TypeScript) ===');
  const start = Date.now();

  const promises = urls.map(url => checkUrl(url));
  const results = await Promise.all(promises);

  console.table(results);
  console.log(`\nTotal running time: ${Date.now() - start}ms`);
}

main();
