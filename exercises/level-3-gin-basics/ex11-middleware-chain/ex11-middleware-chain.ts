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

import express, { Request, Response, NextFunction } from 'express';
const app = express();

function requestLogger(req: Request, res: Response, next: NextFunction): void {
  const start = Date.now();
  res.on('finish', () => {
    const duration = Date.now() - start;
    console.log(`${req.method} ${req.originalUrl} ${res.statusCode} - ${duration}ms`);
  });
  next();
}

function corsMiddleware(req: Request, res: Response, next: NextFunction): void | Response {
  res.setHeader('Access-Control-Allow-Origin', '*');
  res.setHeader('Access-Control-Allow-Methods', 'GET, POST, PUT, DELETE, OPTIONS');
  res.setHeader('Access-Control-Allow-Headers', 'Content-Type, Authorization, X-API-Key');
  if (req.method === 'OPTIONS') {
    return res.sendStatus(200);
  }
  next();
}

function authMiddleware(req: Request, res: Response, next: NextFunction): void | Response {
  const apiKey = req.headers['x-api-key'];
  if (!apiKey || apiKey !== 'secret-key') {
    return res.status(401).json({ error: 'Unauthorized: Invalid API Key' });
  }
  next();
}

function errorHandler(err: any, req: Request, res: Response, next: NextFunction): void {
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
