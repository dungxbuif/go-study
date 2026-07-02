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

// CreateUserSchema định nghĩa hình dạng (shape) của dữ liệu lúc runtime.
//
// 🧠 CƠ CHẾ CỦA ZOD VS GOLANG REFLECTION:
// - TypeScript: Các Type (`type`, `interface`) sẽ bị BỐC HƠI hoàn toàn (erased) sau khi biên dịch sang Javascript.
//   Do đó tại runtime, Node.js không thể biết một đối tượng có thỏa mãn kiểu `CreateUserDTO` hay không.
//   Zod giải quyết vấn đề này bằng cách tạo ra một đối tượng Schema thực tế lúc runtime để kiểm tra (safeParse).
// - `z.infer<typeof CreateUserSchema>` là một cầu nối ngược: Nó tự động suy luận ra kiểu TypeScript tĩnh từ đối tượng Zod runtime.
// - Trái lại trong Go, struct và struct tags luôn tồn tại song song cả lúc compile-time lẫn runtime dưới dạng nhị phân,
//   giúp ta không cần duy trì 2 định nghĩa riêng biệt (schema và type/interface).
const CreateUserSchema = z.object({
   username: z
      .string()
      .min(3)
      .max(20)
      .regex(/^[a-zA-Z0-9]+$/),
   email: z.string().email(),
   password: z.string().min(8),
   age: z.number().int().gte(18).lte(120),
   phone: z
      .string()
      .refine((val) => val.startsWith('+84') || val.startsWith('0'), {
         message: 'Phone number must start with +84 or 0',
      }),
});

type CreateUserDTO = z.infer<typeof CreateUserSchema>;

// Route Register
//
// 🧠 PHÂN TÍCH RỦI RO BỘ NHỚ (Dynamic Payload Parsing):
// - Trong Node.js, `req.body` là một đối tượng Javascript động có thể chứa hàng triệu thuộc tính ngẫu nhiên
//   do hacker truyền vào, chiếm dụng RAM và CPU của tiến trình đơn luồng.
// - Bằng việc sử dụng `CreateUserSchema.safeParse(req.body)`, Zod lọc bỏ hoàn toàn các trường dư thừa (strip by default),
//   chỉ trả về một đối tượng sạch sẽ chứa các thuộc tính đã được khai báo trong `result.data`.
// - Điều này cung cấp tính năng an toàn bảo mật tương tự như cơ chế đón nhận DTO tĩnh trong Go Gin!
app.post('/register', (req: Request, res: Response) => {
   const result = CreateUserSchema.safeParse(req.body);

   if (!result.success) {
      return res.status(400).json({
         success: false,
         errors: result.error.errors.map((err) => ({
            field: err.path.join('.'),
            message: err.message,
         })),
      });
   }

   // Nhận dữ liệu sạch đã được validate và lọc sạch (safe payload)
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
