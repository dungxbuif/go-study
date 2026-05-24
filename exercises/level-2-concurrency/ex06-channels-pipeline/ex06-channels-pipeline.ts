/**
 * Ex06: Channels + Pipeline — TypeScript Version
 * 
 * 🧠 So sánh key:
 * - Node.js: Xây dựng luồng xử lý liên tục (pipeline) thông qua Node.js Streams hoặc Async Generators.
 *            Bản chất Event Loop vẫn chuyển đổi tuần tự giữa các microtask phát sinh trong pipeline.
 * - Go:      Các Stage trong pipeline chạy độc lập hoàn toàn trên các Goroutines khác nhau.
 *            Dữ liệu truyền qua lại bằng Channels. Lợi ích lớn nhất là tính năng Blocking mặc định của 
 *            Channel giúp tự động kiểm soát tốc độ giữa các Stage (Backpressure) cực kỳ mượt mà.
 */

import { Readable, Transform, pipeline } from 'stream';

async function* generatorStage(limit: number): AsyncGenerator<number> {
  for (let i = 1; i <= limit; i++) {
    yield i;
  }
}

async function* filterStage(source: AsyncIterable<number>): AsyncGenerator<number> {
  for await (const val of source) {
    if (val % 2 === 0) {
      yield val;
    }
  }
}

async function* squareStage(source: AsyncIterable<number>): AsyncGenerator<number> {
  for await (const val of source) {
    yield val * val;
  }
}

async function runGeneratorPipeline(): Promise<void> {
  console.log('=== Running Pipeline via Async Generators ===');
  const generatorStream = generatorStage(20);
  const filteredStream = filterStage(generatorStream);
  const squaredStream = squareStage(filteredStream);

  for await (const result of squaredStream) {
    console.log(`Pipeline Output: ${result}`);
  }
}

function runStreamsPipeline(): void {
  console.log('\n=== Running Pipeline via Node.js Streams ===');

  let count = 1;
  const numbersSource = new Readable({
    objectMode: true,
    read() {
      if (count > 20) {
        this.push(null);
      } else {
        this.push(count++);
      }
    }
  });

  const filterEven = new Transform({
    objectMode: true,
    transform(chunk: number, encoding, callback) {
      if (chunk % 2 === 0) {
        callback(null, chunk);
      } else {
        callback();
      }
    }
  });

  const square = new Transform({
    objectMode: true,
    transform(chunk: number, encoding, callback) {
      callback(null, chunk * chunk);
    }
  });

  const printer = new Transform({
    objectMode: true,
    transform(chunk: number, encoding, callback) {
      console.log(`Stream Output: ${chunk}`);
      callback(null, chunk);
    }
  });

  pipeline(
    numbersSource,
    filterEven,
    square,
    printer,
    (err) => {
      if (err) console.error('Pipeline failed:', err);
      else console.log('Streams Pipeline completed.');
    }
  );
}

async function main(): Promise<void> {
  await runGeneratorPipeline();
  runStreamsPipeline();
}

main();
