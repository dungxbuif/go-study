/**
 * Ex08: Context — TypeScript Version
 *
 * 🧠 So sánh key:
 * - Node.js: Quản lý việc hủy bỏ (cancellation) và giới hạn thời gian (timeout) bằng cách sử dụng
 *            `AbortController` và truyền `AbortSignal` vào các tác vụ I/O bất đồng bộ (như fetch, axios).
 * - Go:      Sử dụng `context.Context` (`context.WithTimeout`, `context.WithCancel`). Context trong Go
 *            được coi là "tiêu chuẩn vàng" để lan truyền tín hiệu huỷ bỏ, thời gian hết hạn và meta-data
 *            xuyên suốt tất cả các tầng kiến trúc trong hệ thống.
 */

import axios from 'axios';

interface FetchResult {
   url: string;
   data?: any;
   error?: string;
   success: boolean;
}

// fetchData thực hiện HTTP GET request và tích hợp cơ chế AbortSignal.
//
// 💡 CƠ CHẾ DƯỚI NẮP CAPO (Node.js & AbortSignal):
// 1. Khi ta truyền `signal` vào Axios (`axios.get(url, { signal })`), Axios sẽ đăng ký lắng nghe sự kiện 'abort' trên đối tượng `AbortSignal`.
// 2. Khi `controller.abort()` được gọi ở luồng chính, đối tượng `AbortSignal` sẽ chuyển trạng thái `aborted = true` và phát ra sự kiện.
// 3. Axios bắt được sự kiện này và ngay lập tức gọi phương thức `.abort()` trên underlying `http.ClientRequest` của Node.js.
// 4. Kết nối TCP socket sẽ bị hủy bỏ (destroyed) ngay tại tầng nhân, ngăn chặn việc truyền tải gói tin thêm và giải phóng tài nguyên.
async function fetchData(
   url: string,
   signal: AbortSignal,
): Promise<FetchResult> {
   try {
      const response = await axios.get(url, { signal });
      return { url, data: response.data, success: true };
   } catch (error: any) {
      // Phân biệt lỗi do chủ động hủy bỏ (Timeout) hay các lỗi mạng thông thường khác (500, Connection Refused...)
      if (
         axios.isCancel(error) ||
         error.name === 'AbortError' ||
         error.code === 'ERR_CANCELED'
      ) {
         return {
            url,
            error: 'Request was cancelled (Timeout)',
            success: false,
         };
      }
      return { url, error: error.message, success: false };
   }
}

// fetchMultiple thực hiện gửi nhiều request song song (Promise.all) với một giới hạn thời gian tối đa (timeoutMs).
//
// 🧠 CƠ CHẾ QUẢN LÝ TIMER LEAK (Rò rỉ bộ nhớ):
// - Trong Node.js, `setTimeout` tạo ra một đối tượng Timer được giữ lại trong hàng đợi sự kiện (Event Loop).
// - Nếu tất cả các request trong `Promise.all` hoàn thành XONG TRƯỚC thời gian timeoutMs, ta BẮT BUỘC phải gọi `clearTimeout(timeoutId)`.
// - Nếu không gọi `clearTimeout`, hàm callback hủy bỏ vẫn sẽ nằm trong bộ nhớ và tiếp tục chạy khi hết thời gian,
//   đồng thời ngăn cản garbage collector thu hồi bộ nhớ liên quan đến `AbortController`.
// - Điều này tương đương với quy tắc bắt buộc phải gọi `cancel()` khi dùng Context trong Go.
async function fetchMultiple(
   urls: string[],
   timeoutMs: number,
): Promise<FetchResult[]> {
   // Khởi tạo AbortController tương đương với context.WithCancel / context.WithTimeout trong Go
   const controller = new AbortController();
   const { signal } = controller;

   // Lên lịch để hủy bỏ request sau timeoutMs (Tương đương cơ chế timer của context.WithTimeout)
   const timeoutId = setTimeout(() => {
      controller.abort();
   }, timeoutMs);

   try {
      // Tạo mảng các Promise
      const promises = urls.map((url) => fetchData(url, signal));
      // Thực thi song song (Đợi toàn bộ phản hồi hoặc đến khi bị hủy bởi controller.abort())
      const results = await Promise.all(promises);
      return results;
   } finally {
      // Luôn dọn dẹp timer để tránh memory leak trong Node.js Event Loop
      clearTimeout(timeoutId);
   }
}

async function main(): Promise<void> {
   const urls: string[] = [
      'https://httpbin.org/delay/1',
      'https://httpbin.org/delay/5',
      'https://httpbin.org/status/200',
   ];

   console.log('=== Fetching with 3s Timeout ===');
   const results = await fetchMultiple(urls, 3000);
   console.table(results);
}

main();
