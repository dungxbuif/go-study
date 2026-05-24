/**
 * Ex17: Channel Semaphore + Binary Split Retry — TypeScript Version
 * 
 * 🧠 So sánh key:
 * - Node.js: Để giới hạn concurrency, ta tự xây dựng cấu trúc hàng đợi hoặc dùng thư viện.
 *            Để thực hiện thuật toán chia đôi thử lại (Binary Split Retry) nhằm cô lập lỗi trong mảng/chunk,
 *            ta thực hiện chia đôi mảng bằng đệ quy và chạy song song bằng `Promise.all()`.
 * - Go:      Sử dụng buffer channel làm Semaphore giới hạn luồng ghi cực kỳ tinh gọn.
 *            Khi một chunk dữ liệu ghi bị lỗi (ví dụ lỗi mạng/dữ liệu sai), Go áp dụng thuật toán đệ quy 
 *            chia đôi (Binary Split Retry) để chia nhỏ chunk đó thành các Goroutine chạy song song.
 *            Cơ chế này giúp định vị và cách ly chính xác dữ liệu lỗi (như lỗi block #37) 
 *            mà không làm tắc nghẽn toàn bộ hệ thống hay bị lock DB.
 */

async function sleep(ms: number): Promise<void> {
  return new Promise(resolve => setTimeout(resolve, ms));
}

async function processItem(item: number): Promise<number> {
  await sleep(50);
  if (item === 37 || item === 82) {
    throw new Error(`Failed to process item ${item}`);
  }
  return item;
}

class PromiseSemaphore {
  private maxConcurrent: number;
  private currentConcurrent: number;
  private queue: (() => void)[];

  constructor(maxConcurrent: number) {
    this.maxConcurrent = maxConcurrent;
    this.currentConcurrent = 0;
    this.queue = [];
  }

  public async acquire(): Promise<void> {
    if (this.currentConcurrent >= this.maxConcurrent) {
      await new Promise<void>(resolve => this.queue.push(resolve));
    }
    this.currentConcurrent++;
  }

  public release(): void {
    this.currentConcurrent--;
    if (this.queue.length > 0) {
      const next = this.queue.shift();
      if (next) next();
    }
  }
}

interface ChunkResult {
  results: number[];
  failedItems: number[];
}

async function processChunk(chunk: number[], sem: PromiseSemaphore): Promise<ChunkResult> {
  const results: number[] = [];
  const failedItems: number[] = [];

  const promises = chunk.map(async (item) => {
    await sem.acquire();
    try {
      const res = await processItem(item);
      results.push(res);
    } catch (err) {
      failedItems.push(item);
    } finally {
      sem.release();
    }
  });

  await Promise.all(promises);
  return { results, failedItems };
}

async function processChunkWithBinaryRetry(chunk: number[], sem: PromiseSemaphore): Promise<void> {
  const { failedItems } = await processChunk(chunk, sem);
  if (failedItems.length === 0) {
    return;
  }

  if (failedItems.length === 1) {
    console.log(`[FAILED ITEM ISOLATED] Item ${failedItems[0]} permanently failed`);
    return;
  }

  console.log(`[CHUNK FAILED] Chunk failed. Splitting failed items (${failedItems.length}) in half...`);
  const mid = Math.floor(failedItems.length / 2);
  const leftChunk = failedItems.slice(0, mid);
  const rightChunk = failedItems.slice(mid);

  await Promise.all([
    processChunkWithBinaryRetry(leftChunk, sem),
    processChunkWithBinaryRetry(rightChunk, sem),
  ]);
}

async function main(): Promise<void> {
  const items = Array.from({ length: 100 }, (_, i) => i + 1);
  const chunkSize = 10;
  const chunks: number[][] = [];
  for (let i = 0; i < items.length; i += chunkSize) {
    chunks.push(items.slice(i, i + chunkSize));
  }

  const sem = new PromiseSemaphore(3);
  const start = Date.now();

  const promises = chunks.map(chunk => processChunkWithBinaryRetry(chunk, sem));
  await Promise.all(promises);

  console.log(`All operations completed in ${Date.now() - start}ms`);
}

main();
