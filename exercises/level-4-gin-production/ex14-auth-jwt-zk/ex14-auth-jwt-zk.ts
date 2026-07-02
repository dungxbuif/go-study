/**
 * Ex14: Auth JWT + ZK Concept — TypeScript Version
 *
 * 🧠 So sánh key:
 * - Node.js: Dùng `passport.js` hoặc custom Express middlewares + thư viện `jsonwebtoken` và `bcrypt`
 *            để quản lý xác thực người dùng, băm mật khẩu, phân quyền role.
 * - Go:      Dùng các thư viện bên thứ ba như `golang-jwt` và `bcrypt`. Lưu trữ user context trực tiếp
 *            vào request context của Gin (`c.Set("userID", ...)`).
 *
 * 💡 Sự khác biệt lớn nhất:
 * 1. ZK (Zero-Knowledge) mock auth: Khi middleware kiểm tra và đọc trước Request Body (ví dụ để xác thực ZK proof),
 *    nếu không khéo léo sẽ làm mất stream body của các middleware/handlers tiếp theo.
 * 2. Trong Go, ta cần dùng kỹ thuật đọc và khôi phục body bằng `io.NopCloser` để đảm bảo stream dữ liệu vẫn đọc được bình thường ở tầng sau.
 */

import * as bcrypt from 'bcrypt';
import express, { NextFunction, Request, Response } from 'express';
import * as jwt from 'jsonwebtoken';

const app = express();

// 🧠 CƠ CHẾ ĐỌC BODY DƯỚI NẮP CAPO (Node.js express.json() Middleware):
// - Trong Node.js/Express, mặc định `req` là một Readable Stream kế thừa từ `http.IncomingMessage`.
//   Ta cũng chỉ có thể đọc stream này một lần để lấy body.
// - Tuy nhiên, Express giải quyết việc này bằng cách sử dụng middleware `express.json()` toàn cục.
//   Middleware này đọc toàn bộ stream body từ sớm, phân tích thành JSON Object và lưu vào thuộc tính `req.body`.
//   Từ đó, mọi middleware phía sau đều có thể truy cập `req.body` thoải mái mà không lo nghẽn stream hay EOF.
// - Đây là sự khác biệt cực kỳ quan trọng so với Go: Go Gin mặc định KHÔNG parse sẵn JSON body toàn cục vào RAM
//   để tiết kiệm tối đa hiệu năng bộ nhớ. Do đó trong Go, việc gọi binding `ShouldBindJSON` lần thứ hai
//   lên cùng một request sẽ gặp lỗi, trừ khi ta khôi phục thủ công `c.Request.Body`.
app.use(express.json());

const JWT_SECRET = 'super-secret-key';

interface User {
   username: string;
   password?: string;
   role: string;
}

// 🧠 MỞ RỘNG KIỂU DỮ LIỆU ĐỊNH DANH (TypeScript Request Extension vs Go Gin Context Keys):
// - Trong TypeScript, Express request (`Request`) là kiểu tĩnh cố định. Để gắn thêm thuộc tính `user`
//   hay `user_id` vào request nhằm chuyển tải dữ liệu giữa các middlewares, ta phải kế thừa và định nghĩa
//   một interface mới: `interface AuthenticatedRequest extends Request`.
// - Trái lại trong Go Gin, ta sử dụng `c.Set("key", value)` là một map động kiểu `map[string]any`.
//   Nhờ vậy, Go không cần định nghĩa lại struct Request, nhưng bù lại lập trình viên phải tự ép kiểu thủ công khi lấy dữ liệu ra.
interface AuthenticatedRequest extends Request {
   user?: any;
   user_id?: string;
}

const users = new Map<string, User>();

app.post('/auth/register', async (req: Request, res: Response) => {
   const { username, password } = req.body;
   if (!username || !password) {
      return res
         .status(400)
         .json({ error: 'Username and password are required' });
   }

   // Bcrypt.hash băm mật khẩu bất đồng bộ sử dụng luồng (Thread Pool - libuv) của Node.js,
   // tránh nghẽn Event Loop vì băm bcrypt cực kỳ tốn CPU.
   const hashedPassword = await bcrypt.hash(password, 10);
   users.set(username, { username, password: hashedPassword, role: 'user' });

   res.status(201).json({ message: 'User registered successfully' });
});

app.post('/auth/login', async (req: Request, res: Response) => {
   const { username, password } = req.body;
   const user = users.get(username);

   // bcrypt.compare đối chiếu bất đồng bộ
   if (
      !user ||
      !user.password ||
      !(await bcrypt.compare(password, user.password))
   ) {
      return res.status(401).json({ error: 'Invalid credentials' });
   }

   const token = jwt.sign(
      { username: user.username, role: user.role },
      JWT_SECRET,
      { expiresIn: '1h' },
   );
   res.json({ token });
});

function jwtMiddleware(
   req: AuthenticatedRequest,
   res: Response,
   next: NextFunction,
): void | Response {
   const authHeader = req.headers['authorization'];
   if (!authHeader || !authHeader.startsWith('Bearer ')) {
      return res.status(401).json({ error: 'Missing or malformed token' });
   }

   const token = authHeader.split(' ')[1];
   try {
      const claims = jwt.verify(token, JWT_SECRET);
      req.user = claims;
      next();
   } catch (err) {
      res.status(401).json({ error: 'Invalid or expired token' });
   }
}

function roleRequired(role: string) {
   return (
      req: AuthenticatedRequest,
      res: Response,
      next: NextFunction,
   ): void | Response => {
      if (!req.user || req.user.role !== role) {
         return res
            .status(403)
            .json({ error: 'Forbidden: Insufficient privileges' });
      }
      next();
   };
}

// mockZkAuthMiddleware: Xác thực Zero-Knowledge Proof mock.
// Nhờ có express.json() parse sẵn `req.body`, ta có thể đọc dữ liệu an toàn mà không sợ làm mất stream của các handler tiếp theo.
function mockZkAuthMiddleware(
   req: AuthenticatedRequest,
   res: Response,
   next: NextFunction,
): void | Response {
   const { user_id, proof, public_key } = req.body;

   if (proof === 'valid-proof') {
      req.user_id = user_id;
      next();
   } else {
      res.status(401).json({ error: 'Invalid ZK proof' });
   }
}

app.get(
   '/auth/profile',
   jwtMiddleware,
   (req: AuthenticatedRequest, res: Response) => {
      res.json({ user: req.user });
   },
);

app.post(
   '/auth/zk-verify',
   mockZkAuthMiddleware,
   (req: AuthenticatedRequest, res: Response) => {
      res.json({ message: 'ZK verified successfully', user_id: req.user_id });
   },
);

app.listen(8080, () => {
   console.log('Auth server running on port 8080');
});
