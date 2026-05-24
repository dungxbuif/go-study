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

app.get('/', (req: Request, res: Response) => {
  res.json({ message: 'Hello, Go!' });
});

app.get('/health', (req: Request, res: Response) => {
  res.json({ status: 'ok', uptime: process.uptime() });
});

app.get('/users/:id', (req: Request, res: Response) => {
  const id = req.params.id;
  res.json({ user_id: id });
});

app.get('/search', (req: Request, res: Response) => {
  const q = (req.query.q as string) || "";
  const page = (req.query.page as string) || "1";
  res.json({ query: q, page: parseInt(page, 10) });
});

const apiRouter = express.Router();
apiRouter.get('/ping', (req: Request, res: Response) => {
  res.json({ message: 'pong' });
});

app.use('/api/v1', apiRouter);

app.listen(PORT, () => {
  console.log(`Server is running on port ${PORT}`);
});
