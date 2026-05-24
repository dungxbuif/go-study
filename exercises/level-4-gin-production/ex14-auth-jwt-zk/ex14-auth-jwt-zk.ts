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

import express, { Request, Response, NextFunction } from 'express';
import * as jwt from 'jsonwebtoken';
import * as bcrypt from 'bcrypt';

const app = express();
app.use(express.json());

const JWT_SECRET = 'super-secret-key';

interface User {
  username: string;
  password?: string;
  role: string;
}

interface AuthenticatedRequest extends Request {
  user?: any;
  user_id?: string;
}

const users = new Map<string, User>();

app.post('/auth/register', async (req: Request, res: Response) => {
  const { username, password } = req.body;
  if (!username || !password) {
    return res.status(400).json({ error: 'Username and password are required' });
  }

  const hashedPassword = await bcrypt.hash(password, 10);
  users.set(username, { username, password: hashedPassword, role: 'user' });

  res.status(201).json({ message: 'User registered successfully' });
});

app.post('/auth/login', async (req: Request, res: Response) => {
  const { username, password } = req.body;
  const user = users.get(username);

  if (!user || !user.password || !(await bcrypt.compare(password, user.password))) {
    return res.status(401).json({ error: 'Invalid credentials' });
  }

  const token = jwt.sign({ username: user.username, role: user.role }, JWT_SECRET, { expiresIn: '1h' });
  res.json({ token });
});

function jwtMiddleware(req: AuthenticatedRequest, res: Response, next: NextFunction): void | Response {
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
  return (req: AuthenticatedRequest, res: Response, next: NextFunction): void | Response => {
    if (!req.user || req.user.role !== role) {
      return res.status(403).json({ error: 'Forbidden: Insufficient privileges' });
    }
    next();
  };
}

function mockZkAuthMiddleware(req: AuthenticatedRequest, res: Response, next: NextFunction): void | Response {
  const { user_id, proof, public_key } = req.body;
  
  if (proof === 'valid-proof') {
    req.user_id = user_id;
    next();
  } else {
    res.status(401).json({ error: 'Invalid ZK proof' });
  }
}

app.get('/auth/profile', jwtMiddleware, (req: AuthenticatedRequest, res: Response) => {
  res.json({ user: req.user });
});

app.post('/auth/zk-verify', mockZkAuthMiddleware, (req: AuthenticatedRequest, res: Response) => {
  res.json({ message: 'ZK verified successfully', user_id: req.user_id });
});

app.listen(8080, () => {
  console.log('Auth server running on port 8080');
});
