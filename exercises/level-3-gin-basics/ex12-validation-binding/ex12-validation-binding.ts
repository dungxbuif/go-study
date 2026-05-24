/**
 * Ex12: Validation & Binding — TypeScript Version
 * 
 * 🧠 So sánh key:
 * - TypeScript: Sử dụng thư viện `zod` để định nghĩa schema run-time và validate động dữ liệu
 *               được gửi lên từ request body (`safeParse()`).
 * - Go:         Sử dụng cơ chế struct tags và thư viện validator mặc định của Gin (`binding:"required,email"`).
 *               Quá trình binding diễn ra tĩnh (static) và tự động ánh xạ (marshal/unmarshal) 
 *               thông qua reflection.
 * 
 * 💡 Sự khác biệt lớn nhất:
 * 1. Zod schema được định nghĩa tách biệt dưới dạng runtime-code. Còn Go định nghĩa trực tiếp 
 *    metadata (struct tags) trên chính struct nhận dữ liệu.
 * 2. Go hỗ trợ custom validator bằng cách đăng ký callback với validator engine của Gin.
 */

import express, { Request, Response } from 'express';
import { z } from 'zod';

const app = express();
app.use(express.json());

const CreateUserSchema = z.object({
  username: z.string().min(3).max(20).regex(/^[a-zA-Z0-9]+$/),
  email: z.string().email(),
  password: z.string().min(8),
  age: z.number().int().gte(18).lte(120),
  phone: z.string().refine(val => val.startsWith('+84') || val.startsWith('0'), {
    message: 'Phone number must start with +84 or 0',
  }),
});

type CreateUserDTO = z.infer<typeof CreateUserSchema>;

app.post('/register', (req: Request, res: Response) => {
  const result = CreateUserSchema.safeParse(req.body);

  if (!result.success) {
    return res.status(400).json({
      success: false,
      errors: result.error.errors.map(err => ({
        field: err.path.join('.'),
        message: err.message,
      })),
    });
  }

  const { username, email, age, phone } = result.data;

  const userResponse = {
    username,
    email,
    age,
    phone,
  };

  res.status(201).json({
    success: true,
    data: userResponse,
  });
});

app.listen(8080, () => {
  console.log('Validation server running on port 8080');
});
