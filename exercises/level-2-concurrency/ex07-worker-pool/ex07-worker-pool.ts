/**
 * Ex07: Worker Pool — TypeScript Version
 * 
 * 🧠 So sánh key:
 * - Node.js: Thường sử dụng các thư viện như `p-limit` hoặc viết giải pháp Promise queue tuỳ chỉnh
 *            để khống chế số lượng Promise/I/O task được xử lý đồng thời. 
 *            Không cần khởi tạo Thread thực tế trừ phi xử lý tác vụ CPU-intensive (dùng `worker_threads`).
 * - Go:      Worker Pool là kiến trúc kinh điển trong Go. Ta spawn cứng N Goroutines (Workers).
 *            Các Goroutine này cùng lắng nghe trên một `jobs` channel chung, tự động chia sẻ công việc 
 *            và đẩy kết quả về `results` channel. Đây là giải pháp native cực kỳ hiệu năng cao và ít tốn RAM.
 */

interface ResizeResult {
  imageId: number;
  delay: number;
}

function resizeImageTask(imageId: number): Promise<ResizeResult> {
  return new Promise((resolve) => {
    const delay = Math.floor(Math.random() * 400) + 100;
    setTimeout(() => {
      resolve({ imageId, delay });
    }, delay);
  });
}

async function runWorkerPool(tasks: number[], concurrencyLimit: number): Promise<ResizeResult[]> {
  const results: Promise<ResizeResult>[] = [];
  const executing: Promise<ResizeResult>[] = [];

  console.log(`Starting Worker Pool with concurrency limit: ${concurrencyLimit}`);

  for (const task of tasks) {
    const p = resizeImageTask(task).then((res) => {
      executing.splice(executing.indexOf(p), 1);
      console.log(`[Worker Pool] Processed image_${res.imageId} in ${res.delay}ms`);
      return res;
    });

    results.push(p);
    executing.push(p);

    if (executing.length >= concurrencyLimit) {
      await Promise.race(executing);
    }
  }

  return Promise.all(results);
}

async function main(): Promise<void> {
  const images = Array.from({ length: 20 }, (_, i) => i + 1);
  const limit = 5;

  const start = Date.now();
  
  const results = await runWorkerPool(images, limit);

  console.log(`\nAll ${results.length} images processed in ${Date.now() - start}ms`);
}

main();
