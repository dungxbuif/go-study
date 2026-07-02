/**
 * Ex20: Batch Pipeline + Event Streaming — TypeScript Version
 *
 * 🧠 So sánh key:
 * - Node.js: Tự động gom lô (batch) các thao tác bất đồng bộ thông qua cấu trúc hàng đợi mảng
 *            và bộ đếm thời gian `setTimeout`.
 *            Sử dụng các transaction của cơ sở dữ liệu để bảo đảm tính toàn vẹn dữ liệu.
 * - Go:      Sử dụng WalletUpdateBatcher pattern kết hợp buffered channel làm queue và `time.Ticker`/`time.Timer`
 *            trực tiếp trong vòng lặp `select` để kiểm soát cả hai điều kiện: kích thước lô (batch size)
 *            hoặc thời gian chờ tối đa (timeout).
 *
 * 💡 Sự khác biệt lớn nhất:
 * 1. Go select block một cách tự nhiên và có hiệu năng cực cao khi điều phối tín hiệu timeout và dữ liệu đến song song.
 * 2. Cơ chế Graceful Shutdown sử dụng `signal.Notify` bắt SIGINT/SIGTERM, đợi flush hết lô dữ liệu hiện tại
 *    rồi mới ngắt kết nối DB nhằm không làm mất mát bất kỳ sự kiện hay dữ liệu nào.
 */

// Batcher: Hỗ trợ gom lô các sự kiện hoặc câu lệnh ghi để tối ưu I/O.
//
// 🧠 CƠ CHẾ GOM LÔ TRONG NODE.JS (Array Queue & setTimeout under the hood):
// - Trong Node.js, `queue` đơn giản là một mảng JS động (`string[]`). Mảng này tự co giãn kích thước trên V8 Heap.
// - **Push (Gửi dữ liệu)**:
//   - Ta đẩy dữ liệu vào `queue` bằng `push()`.
//   - Nếu kích thước mảng đạt `batchSize`, ta chủ động kích hoạt `flush()` ngay lập tức.
//   - Nếu chưa đạt giới hạn, nhưng đây là item đầu tiên của lô mới, ta bắt đầu đếm ngược thời gian bằng `setTimeout`.
//     Timer này đảm bảo nếu lưu lượng tin nhắn bị ngắt quãng (sparse traffic), ta vẫn không giữ dữ liệu bị "kẹt" trong RAM quá lâu.
// - **Flush (Xóa hàng đợi & Xử lý)**:
//   - Ta dọn dẹp `setTimeout` cũ bằng `clearTimeout` để tránh rò rỉ timer của Event Loop.
//   - Sao chép nông (shallow copy) mảng `queue` ra một biến tạm `batch = [...this.queue]` rồi xóa sạch `queue` gốc.
//     Việc sao chép này cực kỳ quan trọng để giải phóng nhanh hàng đợi chính, cho phép các thao tác `push` tiếp theo tiếp tục nhận dữ liệu
//     mà không bị ảnh hưởng bởi quá trình xử lý bất đồng bộ `processBatchFn` đang diễn ra!
class Batcher {
   private batchSize: number;
   private timeoutMs: number;
   private processBatchFn: (batch: string[]) => Promise<void>;
   private queue: string[];
   private timer: NodeJS.Timeout | null;

   constructor(
      batchSize: number,
      timeoutMs: number,
      processBatchFn: (batch: string[]) => Promise<void>,
   ) {
      this.batchSize = batchSize;
      this.timeoutMs = timeoutMs;
      this.processBatchFn = processBatchFn;
      this.queue = [];
      this.timer = null;
   }

   public push(item: string): void {
      this.queue.push(item);
      if (this.queue.length >= this.batchSize) {
         this.flush();
      } else if (!this.timer) {
         this.timer = setTimeout(() => {
            this.flush();
         }, this.timeoutMs);
      }
   }

   public async flush(): Promise<void> {
      if (this.timer) {
         clearTimeout(this.timer);
         this.timer = null;
      }

      if (this.queue.length === 0) {
         return;
      }

      const batch = [...this.queue];
      this.queue = [];

      try {
         await this.processBatchFn(batch);
      } catch (err: any) {
         console.error('Failed to process batch:', err.message);
      }
   }
}

async function simulateDatabaseTransactionAndEventPublishing(
   batch: string[],
): Promise<void> {
   console.log(`[Batcher] Processing batch of size ${batch.length}...`);
   await new Promise((resolve) => setTimeout(resolve, 100));
   console.log(
      `[Database] Inserted and committed ${batch.length} records successfully.`,
   );
   batch.forEach((item) => {
      console.log(
         `[Kafka] Published event: { "status": "committed", "id": "${item}" }`,
      );
   });
}

async function main(): Promise<void> {
   const batcher = new Batcher(
      5,
      1000,
      simulateDatabaseTransactionAndEventPublishing,
   );

   console.log('--- Pushing 12 items into Batcher ---');
   for (let i = 1; i <= 12; i++) {
      batcher.push(`tx_${i}`);
      await new Promise((resolve) => setTimeout(resolve, 50));
   }

   await new Promise((resolve) => setTimeout(resolve, 1500));
   await batcher.flush();
}

main();
