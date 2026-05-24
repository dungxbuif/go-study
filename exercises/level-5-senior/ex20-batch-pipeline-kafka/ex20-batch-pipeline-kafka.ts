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

class Batcher {
  private batchSize: number;
  private timeoutMs: number;
  private processBatchFn: (batch: string[]) => Promise<void>;
  private queue: string[];
  private timer: NodeJS.Timeout | null;

  constructor(batchSize: number, timeoutMs: number, processBatchFn: (batch: string[]) => Promise<void>) {
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

async function simulateDatabaseTransactionAndEventPublishing(batch: string[]): Promise<void> {
  console.log(`[Batcher] Processing batch of size ${batch.length}...`);
  await new Promise(resolve => setTimeout(resolve, 100));
  console.log(`[Database] Inserted and committed ${batch.length} records successfully.`);
  batch.forEach((item) => {
    console.log(`[Kafka] Published event: { "status": "committed", "id": "${item}" }`);
  });
}

async function main(): Promise<void> {
  const batcher = new Batcher(5, 1000, simulateDatabaseTransactionAndEventPublishing);

  console.log('--- Pushing 12 items into Batcher ---');
  for (let i = 1; i <= 12; i++) {
    batcher.push(`tx_${i}`);
    await new Promise(resolve => setTimeout(resolve, 50)); 
  }

  await new Promise(resolve => setTimeout(resolve, 1500));
  await batcher.flush();
}

main();
