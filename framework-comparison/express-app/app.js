const express = require('express');
// Sử dụng zod để validate dữ liệu cho Express
const { z } = require('zod');

const app = express();

// Đăng ký middleware parse JSON body
app.use(express.json());

// 1. Định nghĩa Schema Validation với Zod
const UserSchema = z.object({
  name: z.string().min(2, { message: 'Tên phải dài tối thiểu 2 ký tự' }),
  email: z.string().email({ message: 'Email không hợp lệ' }),
});

// 2. Custom Middleware Ghi Log (Onion Model)
app.use((req, res, next) => {
  const start = Date.now();
  
  // Lắng nghe sự kiện finish khi response được gửi đi thành công
  res.on('finish', () => {
    const duration = Date.now() - start;
    console.log(`[${req.method}] ${res.statusCode} ${req.originalUrl} - ${duration}ms`);
  });
  
  next();
});

// 3. API Endpoints
app.post('/users', (req, res, next) => {
  try {
    // Validate dữ liệu từ body của request
    const parseResult = UserSchema.safeParse(req.body);
    
    if (!parseResult.success) {
      return res.status(422).json({
        success: false,
        errors: parseResult.error.flatten().fieldErrors,
      });
    }

    return res.status(201).json({
      success: true,
      data: parseResult.data,
    });
  } catch (error) {
    // Chuyển sang Centralized Error Handler
    next(error);
  }
});

// 4. Centralized Error Handler (Xử lý lỗi tập trung)
app.use((err, req, res, next) => {
  console.error(err.stack);
  res.status(500).json({
    success: false,
    error: 'Đã xảy ra lỗi hệ thống nghiêm trọng',
  });
});

const PORT = 8082;
app.listen(PORT, () => {
  console.log(`Express app đang chạy trên cổng ${PORT}...`);
});
