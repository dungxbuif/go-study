/**
 * Ex07: Worker Pool — TypeScript Version
 *
 * ============================================================================
 * 💡 KHÁI NIỆM WORKER POOL & SO SÁNH GIỮA NODE.JS VS GO
 * ============================================================================
 *
 * 1. Worker Pool là gì?
 *    - Là một mẫu thiết kế (Design Pattern) kiểm soát mức độ đồng thời (Concurrency Control).
 *    - Bạn chỉ sinh ra một số lượng "Worker" giới hạn cứng (ví dụ 5 workers) để xử lý một số lượng
 *      công việc rất lớn (ví dụ 10,000 tasks).
 *
 * 2. Tại sao phải dùng Worker Pool?
 *    - **Bảo vệ tài nguyên (Rate Limiting):** Nếu bạn có 10,000 ảnh cần resize, việc kích hoạt 10,000
 *      luồng xử lý đồng thời sẽ làm sập máy chủ (đứt kết nối Database, cạn kiệt CPU/RAM, tràn File Descriptors).
 *    - Giới hạn số lượng Worker giúp hệ thống hoạt động ổn định và có thể dự báo trước mức độ tiêu hao RAM/CPU.
 *
 * 3. Sự khác biệt về cách cài đặt giữa Node.js vs Go:
 *    - **Go (Goroutines + Channels):**
 *      - Go hỗ trợ đa luồng thực sự. Chúng ta chỉ cần chạy `N` Goroutines (Workers).
 *      - Tất cả các Workers này cùng trỏ chung vào một channel công việc (`jobs`). Go Runtime sẽ tự động
 *        phân phối công việc xuống cho worker nào đang rảnh (Load Balancing ngầm định cực kỳ hiệu quả).
 *    - **Node.js (Event Loop + Promises):**
 *      - Node.js mặc định chạy **đơn luồng (Single-threaded)**.
 *      - Để đạt được hiệu năng song song (Concurrent), chúng ta không tạo ra nhiều luồng CPU thực tế
 *        mà chúng ta quản lý **hàng đợi các Promises bất đồng bộ đang thực thi**.
 *      - Chúng ta cho phép chạy tối đa `limit` Promises cùng lúc. Khi đạt ngưỡng limit, chúng ta sử dụng
 *        `Promise.race()` để chờ cho đến khi **chỉ cần ít nhất 1 Promise hoàn thành** ➡️ Giải phóng 1 slot
 *        ➡️ Kéo công việc tiếp theo vào chạy. Đây là thuật toán Pool động (Dynamic Concurrency Pool).
 * ============================================================================
 */

interface ResizeResult {
   imageId: number;
   delay: number;
}

// Giả lập một tác vụ xử lý ảnh bất đồng bộ (resizing) tiêu tốn thời gian ngẫu nhiên 100-500ms
function resizeImageTask(imageId: number): Promise<ResizeResult> {
   return new Promise((resolve) => {
      const delay = Math.floor(Math.random() * 400) + 100;
      setTimeout(() => {
         resolve({ imageId, delay });
      }, delay);
   });
}

/**
 * Hàm quản lý chạy thử Worker Pool động trong TypeScript
 * @param tasks Danh sách ID các công việc cần xử lý
 * @param concurrencyLimit Giới hạn cứng số tác vụ chạy đồng thời
 */
async function runWorkerPool(
   tasks: number[],
   concurrencyLimit: number,
): Promise<ResizeResult[]> {
   const results: Promise<ResizeResult>[] = []; // Chứa tất cả các Promises kết quả để đợi ở cuối
   const executing: Promise<ResizeResult>[] = []; // Hộp chứa các Promises đang thực thi ngay tại thời điểm hiện tại

   console.log(
      `Starting Worker Pool with concurrency limit: ${concurrencyLimit}`,
   );

   for (const task of tasks) {
      // 1. Bắt đầu kích hoạt một tác vụ bất đồng bộ chạy
      const p = resizeImageTask(task).then((res) => {
         // 2. KHI TÁC VỤ NÀY HOÀN THÀNH:
         // Tự động xóa chính nó khỏi danh sách `executing` (đang chạy) để nhường chỗ cho tác vụ khác.
         executing.splice(executing.indexOf(p), 1);
         console.log(
            `[Worker Pool] Processed image_${res.imageId} in ${res.delay}ms`,
         );
         return res;
      });

      results.push(p); // Lưu lại để dùng Promise.all() cuối cùng thu thập kết quả
      executing.push(p); // Thêm vào danh sách đang chạy để giám sát số lượng đồng thời

      // 3. KIỂM SOÁT ĐỒNG THỜI (CONCURRENCY CONTROL):
      // Nếu số lượng tác vụ đang chạy đạt ngưỡng giới hạn cứng...
      if (executing.length >= concurrencyLimit) {
         // Dừng vòng lặp lại!
         // Sử dụng Promise.race(executing) để BLOCK vòng lặp cho đến khi **có ít nhất 1 tác vụ chạy xong**.
         // Khi tác vụ đó xong, hàm callback ở `.then` ở dòng 35 sẽ tự động chạy để xóa nó khỏi `executing`.
         // Vòng lặp được giải phóng (unblocked) và chuyển sang phần tử tiếp theo để nhét thêm 1 tác vụ mới.
         await Promise.race(executing);
      }
   }

   // Chờ tất cả toàn bộ các tác vụ (kể cả những tác vụ cuối cùng còn sót lại) hoàn thành và trả về
   return Promise.all(results);
}

async function main(): Promise<void> {
   const images = Array.from({ length: 20 }, (_, i) => i + 1); // 20 ảnh cần xử lý
   const limit = 5; // Chỉ cho phép xử lý tối đa 5 ảnh cùng lúc

   const start = Date.now();

   const results = await runWorkerPool(images, limit);

   console.log(
      `\nAll ${results.length} images processed in ${Date.now() - start}ms`,
   );
}

main();
