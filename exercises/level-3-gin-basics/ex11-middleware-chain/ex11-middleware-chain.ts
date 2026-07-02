/**
 * Ex11: Middleware Chain — TypeScript Version
 *
 * 🧠 So sánh key:
 * - Node.js: Express middleware chain hoạt động bằng cách gọi `next()` để chuyển sang middleware tiếp theo.
 *            Có thể lưu trữ/truyền trạng thái request qua object `req` (như `req.user`).
 * - Go:      Gin sử dụng mô hình "Hành tây" (Onion model). Middleware là một hàm nhận `*gin.Context`.
 *            Gọi `c.Next()` để đi sâu vào các lớp trong, và code sau `c.Next()` sẽ được chạy ngược
 *            ra ngoài sau khi handler xử lý xong (POST-processing).
 *            Để truyền dữ liệu giữa các middleware, Gin dùng `c.Set()` và `c.Get()`.
 */

import express, { NextFunction, Request, Response } from 'express';
const app = express();

// requestLogger: Ghi log thời gian phản hồi.
//
// 💡 CƠ CHẾ DƯỚI NẮP CAPO (Express Event Binding vs Gin Onion Model):
// - Trong Express, khi hàm `next()` được gọi, nó ngay lập tức kích hoạt phần middleware tiếp theo đồng bộ/bất đồng bộ.
// - Tuy nhiên, vì Javascript chạy bất đồng bộ (non-blocking I/O), hàm `next()` sẽ trả về trước khi response thực sự được gửi đi.
// - Do đó, để đo lường chính xác thời điểm hoàn tất HTTP request, Express bắt buộc phải đăng ký lắng nghe sự kiện
//   `res.on('finish', ...)` trên socket stream.
// - Trái lại trong Go, nhờ cơ chế đồng bộ hóa luồng chạy của Goroutine, Gin chỉ đơn giản là chặn tại `c.Next()`,
//   đợi tất cả xử lý xong, rồi chạy tiếp dòng code ngay bên dưới. Cực kỳ trực quan và không cần event listener!
function requestLogger(req: Request, res: Response, next: NextFunction): void {
   const start = Date.now();
   res.on('finish', () => {
      const duration = Date.now() - start;
      console.log(
         `${req.method} ${req.originalUrl} ${res.statusCode} - ${duration}ms`,
      );
   });
   next();
}

// corsMiddleware: Thêm CORS Headers.
// Thể hiện cơ chế ngắt mạch sớm trong Express. Bằng cách trả về `res.sendStatus(200)` mà không gọi `next()`,
// Express sẽ ngắt dòng đời request tại đây. Trong Go, ta gọi `c.Abort()` và cũng phải thêm `return` để đạt hiệu quả tương tự.
function corsMiddleware(
   req: Request,
   res: Response,
   next: NextFunction,
): void | Response {
   res.setHeader('Access-Control-Allow-Origin', '*');
   res.setHeader(
      'Access-Control-Allow-Methods',
      'GET, POST, PUT, DELETE, OPTIONS',
   );
   res.setHeader(
      'Access-Control-Allow-Headers',
      'Content-Type, Authorization, X-API-Key',
   );
   if (req.method === 'OPTIONS') {
      return res.sendStatus(200);
   }
   next();
}

function authMiddleware(
   req: Request,
   res: Response,
   next: NextFunction,
): void | Response {
   const apiKey = req.headers['x-api-key'];
   if (!apiKey || apiKey !== 'secret-key') {
      return res.status(401).json({ error: 'Unauthorized: Invalid API Key' });
   }
   next();
}

// errorHandler: Global Error Handler trong Express.
//
// 💡 SO SÁNH PHỤC HỒI LỖI (Express vs Go Recovery):
// - Express: Phát hiện middleware lỗi bằng cách đếm số lượng tham số của hàm (arity).
//   Hàm nào nhận đúng 4 tham số `(err, req, res, next)` sẽ được Express coi là Error Handler.
//   Mọi exception ném ra trong luồng đồng bộ sẽ được chuyển tới đây.
// - Go Gin: Sử dụng Recovery middleware dưới nắp capo với lệnh `defer recover()`.
//   Khi xảy ra panic (lỗi nghiêm trọng làm sập app), Go sẽ kích hoạt chuỗi defer. Recovery bắt panic đó,
//   ghi nhận log stack trace và gửi status 500, bảo vệ server hoạt động liên tục.
function errorHandler(
   err: any,
   req: Request,
   res: Response,
   next: NextFunction,
): void {
   console.error(err.stack);
   res.status(500).json({ error: 'Internal Server Error' });
}

app.use(requestLogger);
app.use(corsMiddleware);

app.get('/health', (req: Request, res: Response) => {
   res.json({ status: 'ok' });
});

const apiRouter = express.Router();
apiRouter.use(authMiddleware);

apiRouter.get('/data', (req: Request, res: Response) => {
   res.json({ data: 'Sensitive information' });
});

apiRouter.get('/panic', (req: Request, res: Response) => {
   throw new Error('Something went horribly wrong!');
});

app.use('/api/v1', apiRouter);
app.use(errorHandler);

app.listen(8080, () => {
   console.log('Middleware demo server running on port 8080');
});
