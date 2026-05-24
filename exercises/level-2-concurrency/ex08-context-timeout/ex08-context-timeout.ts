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

async function fetchData(url: string, signal: AbortSignal): Promise<FetchResult> {
  try {
    const response = await axios.get(url, { signal });
    return { url, data: response.data, success: true };
  } catch (error: any) {
    if (axios.isCancel(error) || error.name === 'AbortError' || error.code === 'ERR_CANCELED') {
      return { url, error: 'Request was cancelled (Timeout)', success: false };
    }
    return { url, error: error.message, success: false };
  }
}

async function fetchMultiple(urls: string[], timeoutMs: number): Promise<FetchResult[]> {
  const controller = new AbortController();
  const { signal } = controller;

  const timeoutId = setTimeout(() => {
    controller.abort();
  }, timeoutMs);

  try {
    const promises = urls.map(url => fetchData(url, signal));
    const results = await Promise.all(promises);
    return results;
  } finally {
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
