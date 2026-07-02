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

// ============================================================================
// 💡 GIẢI THÍCH CHI TIẾT VỀ GENERATOR FUNCTION VÀ TỪ KHÓA "yield" TRONG JS/TS
// ============================================================================
// 1. Hàm thông thường (Regular Function):
//    - Khi gọi, nó chạy một mạch từ dòng đầu tiên đến khi gặp "return" hoặc hết hàm.
//    - Người gọi (Caller) không thể tạm dừng hàm giữa chừng hoặc yêu cầu hàm trả về từng phần dữ liệu.
//
// 2. Hàm Generator (Generator Function) - Khai báo bằng cú pháp `function*`:
//    - Là hàm đặc biệt có khả năng **TẠM DỪNG (PAUSE)** thực thi tại bất kỳ thời điểm nào và
//      sau đó **TIẾP TỤC (RESUME)** chạy lại từ chính điểm đã dừng.
//
// 3. Từ khóa `yield`:
//    - Được ví như một nút **"TẠM DỪNG VÀ GỬI DỮ LIỆU RA NGOÀI"**.
//    - Khi hàm chạy đến dòng `yield <giá trị>`, nó sẽ dừng chương trình tại đó, đóng gói giá trị và
//      trả ra ngoài cho Caller. Trạng thái của hàm (biến số cục bộ, dòng lệnh đang chạy) được đóng băng.
//    - Khi Caller gọi tiếp (ví dụ qua `.next()` hoặc vòng lặp `for await`), hàm sẽ "rã đông" và
//      chạy tiếp từ dòng ngay sau lệnh `yield` đó.
//
// 4. Async Generator (`async function*`):
//    - Kết hợp sức mạnh của Promise (Async/Await) và Generator.
//    - Trả về một đối tượng `AsyncGenerator` chứa các phần tử sẽ được sinh ra bất đồng bộ theo thời gian.
//    - Giúp chúng ta duyệt qua các phần tử bằng cú pháp cực kỳ sạch sẽ: `for await (const x of generator)`.
//
// 5. Tại sao dùng Generator để làm Pipeline?
//    - Cơ chế này gọi là **Lazy Evaluation** (Đánh giá lười biếng / Xử lý theo yêu cầu).
//    - Dữ liệu KHÔNG cần nạp toàn bộ vào mảng gây tốn RAM. Dữ liệu chỉ được sinh ra từng phần tử một,
//      chảy qua các đường ống filter, transform theo lượt. Cơ chế này mô phỏng rất giống cách Go
//      truyền nhận dữ liệu qua Channel!
// ============================================================================

/**
 * Stage 1: Generator Stage (Sinh dữ liệu)
 * Hàm này sinh ra các số nguyên từ 1 đến limit.
 * Dấu `*` sau `async function` chỉ định đây là một Async Generator.
 */
async function* generatorStage(limit: number): AsyncGenerator<number> {
   for (let i = 1; i <= limit; i++) {
      // Tạm dừng ở đây, trả số `i` ra ngoài cho Stage tiếp theo xử lý.
      // Lần lặp tiếp theo sẽ chạy tiếp từ dòng này (tăng `i` lên và tiếp tục vòng lặp).
      yield i;
   }
}

/**
 * Stage 2: Filter Stage (Lọc dữ liệu)
 * Nhận đầu vào là một nguồn AsyncIterable (nguồn sinh dữ liệu bất đồng bộ).
 */
async function* filterStage(
   source: AsyncIterable<number>,
): AsyncGenerator<number> {
   // Cú pháp `for await` liên tục kéo (pull) từng phần tử từ Stage 1 khi nó được sinh ra
   for await (const val of source) {
      // Chỉ giữ lại các số chẵn
      if (val % 2 === 0) {
         // Nếu là số chẵn, tạm dừng và gửi nó sang Stage 3.
         // Nếu là số lẻ, bỏ qua (không yield) và tiếp tục chờ phần tử tiếp theo từ nguồn.
         yield val;
      }
   }
}

/**
 * Stage 3: Square Stage (Bình phương dữ liệu)
 * Nhận đầu vào từ Stage 2 (chỉ chứa các số chẵn).
 */
async function* squareStage(
   source: AsyncIterable<number>,
): AsyncGenerator<number> {
   for await (const val of source) {
      // Tạm dừng và gửi giá trị đã được bình phương ra ngoài cùng (Caller nhận được)
      yield val * val;
   }
}

/**
 * Hàm điều phối chạy thử Generator Pipeline
 */
async function runGeneratorPipeline(): Promise<void> {
   console.log('=== Running Pipeline via Async Generators ===');

   // Thiết lập đường ống (Chưa có code nào thực sự chạy ở đây - Lazy!)
   const generatorStream = generatorStage(20); // Tạo luồng sinh số từ 1 đến 20
   const filteredStream = filterStage(generatorStream); // Lồng luồng filter vào
   const squaredStream = squareStage(filteredStream); // Lồng luồng bình phương vào

   // Bắt đầu kéo (pull) kết quả từ Stage cuối cùng.
   // Lúc này, squaredStream sẽ yêu cầu filteredStream gửi phần tử, filteredStream
   // lại yêu cầu generatorStream gửi phần tử. Các Stage kích hoạt chạy theo phản ứng dây chuyền!
   for await (const result of squaredStream) {
      console.log(`Pipeline Output (yield): ${result}`);
   }
}

/**
 * Cách làm truyền thống sử dụng Node.js Streams
 * Dành cho bạn đối chiếu: Stream dùng event-driven thay vì Generator.
 */
function runStreamsPipeline(): void {
   console.log('\n=== Running Pipeline via Node.js Streams ===');

   let count = 1;
   const numbersSource = new Readable({
      objectMode: true,
      read() {
         if (count > 20) {
            this.push(null); // Kết thúc stream
         } else {
            this.push(count++); // Đẩy số tiếp theo vào stream
         }
      },
   });

   const filterEven = new Transform({
      objectMode: true,
      transform(chunk: number, encoding, callback) {
         if (chunk % 2 === 0) {
            callback(null, chunk); // Giữ lại phần tử chẵn
         } else {
            callback(); // Bỏ qua phần tử lẻ
         }
      },
   });

   const square = new Transform({
      objectMode: true,
      transform(chunk: number, encoding, callback) {
         callback(null, chunk * chunk); // Nhân đôi/Bình phương phần tử
      },
   });

   const printer = new Transform({
      objectMode: true,
      transform(chunk: number, encoding, callback) {
         console.log(`Stream Output: ${chunk}`);
         callback(null, chunk);
      },
   });

   // Nối các ống stream lại với nhau
   pipeline(numbersSource, filterEven, square, printer, (err) => {
      if (err) console.error('Pipeline failed:', err);
      else console.log('Streams Pipeline completed.');
   });
}

async function main(): Promise<void> {
   await runGeneratorPipeline();
   runStreamsPipeline();
}

main();
