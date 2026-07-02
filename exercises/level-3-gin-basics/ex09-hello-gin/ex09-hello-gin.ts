/**
 * Ex09: Hello Gin — TypeScript Version
 *
 * 🧠 So sánh key:
 * - Node.js: Dùng Express.js để tạo HTTP server, quản lý routing, params, query và router group.
 * - Go:      Dùng Gin framework. Gin được viết trên nền tảng httprouter cực kỳ nhanh (sử dụng Radix Tree),
 *            cho hiệu năng xử lý request vượt trội, lượng cấp phát RAM (allocation) cực thấp.
 *
 * 💡 Sự khác biệt lớn nhất:
 * 1. Express có cơ chế linh động (dynamic typing), trong khi Gin yêu cầu kiểu dữ liệu tường minh (static typing) cho response.
 * 2. Cú pháp khai báo routing và router group của Gin (`r.Group`) tương đối giống và dễ map từ Express.
 */

import express, { Request, Response } from 'express';
const app = express();

const PORT = process.env.PORT || 8080;

// Root Route
//
// 💡 CƠ CHẾ DƯỚI NẮP CAPO (Node.js Express vs Go Gin):
// - Express sử dụng mô hình đơn luồng (Single-Threaded) dựa trên Event Loop.
//   Tất cả các HTTP requests đi vào đều được xử lý trên một luồng chính duy nhất (Main Thread).
//   Nếu callback này thực hiện một phép tính toán CPU nặng (như mã hóa mật khẩu, lặp dữ liệu lớn),
//   nó sẽ BLOCK hoàn toàn Event Loop, khiến tất cả các request khác đến sau bị treo (lag).
// - Go Gin chạy đa luồng cực kỳ hiệu quả thông qua scheduler của Go. Mỗi request là một Goroutine riêng.
//   Mặc dù một request tính toán nặng cũng có thể chiếm CPU, nhưng scheduler sẽ tự động phân phối (preemption)
//   để các Goroutine khác trên các luồng CPU (Thread) khác vẫn chạy trơn tru mà không lo bị nghẽn mạng!
app.get('/', (req: Request, res: Response) => {
   res.json({ message: 'Hello, Go!' });
});

app.get('/health', (req: Request, res: Response) => {
   res.json({ status: 'ok', uptime: process.uptime() });
});

// Route Parameters
// Express giải quyết động qua regex compiler (path-to-regexp package).
// Khi khởi động, Express dịch '/users/:id' thành một RegExp phức tạp.
// Khi có request, nó chạy hàm `.test()` trên URL path, việc này tốn CPU hơn nhiều so với cấu trúc Radix Tree của Gin.
app.get('/users/:id', (req: Request, res: Response) => {
   const id = req.params.id;
   res.json({ user_id: id });
});

// Query Parameters
// Trong Express, `req.query` được phân tích tự động bằng thư viện `qs` hoặc `querystring` dưới dạng Object.
// Lợi thế của JS là kiểu dữ liệu động, `parseInt(page, 10)` được sử dụng để chuyển đổi tường minh an toàn.
app.get('/search', (req: Request, res: Response) => {
   const q = (req.query.q as string) || '';
   const page = (req.query.page as string) || '1';
   res.json({ query: q, page: parseInt(page, 10) });
});

// Express Router Grouping
// express.Router() hoạt động tương đương như `r.Group` trong Gin.
// Giúp gom nhóm các API `/api/v1` thành một module độc lập, dễ quản lý và dễ viết middleware chuyên biệt.
const apiRouter = express.Router();
apiRouter.get('/ping', (req: Request, res: Response) => {
   res.json({ message: 'pong' });
});

app.use('/api/v1', apiRouter);

app.listen(PORT, () => {
   console.log(`Server is running on port ${PORT}`);
});
