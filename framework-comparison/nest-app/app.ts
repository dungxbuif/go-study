import { 
  Controller, 
  Post, 
  Body, 
  Injectable, 
  Module, 
  ValidationPipe, 
  UsePipes, 
  NestInterceptor, 
  ExecutionContext, 
  CallHandler,
  HttpException,
  HttpStatus
} from '@nestjs/common';
import { NestFactory } from '@nestjs/core';
import { IsString, IsEmail, MinLength } from 'class-validator';
import { Observable } from 'rxjs';
import { tap } from 'rxjs/operators';

// ==========================================
// 1. DTO (Data Transfer Object) & Validation Rules
// ==========================================
export class CreateUserDto {
  @IsString()
  @MinLength(2, { message: 'Tên phải dài tối thiểu 2 ký tự' })
  name: string;

  @IsEmail({}, { message: 'Email không hợp lệ' })
  email: string;
}

// ==========================================
// 2. Custom Interceptor (Logger)
// ==========================================
@Injectable()
export class LoggingInterceptor implements NestInterceptor {
  intercept(context: ExecutionContext, next: CallHandler): Observable<any> {
    const req = context.switchToHttp().getRequest();
    const start = Date.now();
    
    return next.handle().pipe(
      tap(() => {
        const duration = Date.now() - start;
        const res = context.switchToHttp().getResponse();
        console.log(`[${req.method}] ${res.statusCode} ${req.url} - ${duration}ms`);
      }),
    );
  }
}

// ==========================================
// 3. Service Layer (Dependency Injection Provider)
// ==========================================
@Injectable()
export class UserService {
  create(dto: CreateUserDto) {
    return {
      success: true,
      data: dto
    };
  }
}

// ==========================================
// 4. Controller Layer (Xử lý Routing & Request)
// ==========================================
@Controller('users')
export class UserController {
  // Inject UserService tự động thông qua Constructor DI
  constructor(private readonly userService: UserService) {}

  @Post()
  // Sử dụng ValidationPipe tự động validate theo rules của DTO class
  @UsePipes(new ValidationPipe({ whitelist: true }))
  createUser(@Body() createUserDto: CreateUserDto) {
    return this.userService.create(createUserDto);
  }
}

// ==========================================
// 5. Main Module Config (Đóng gói ứng dụng)
// ==========================================
@Module({
  controllers: [UserController],
  providers: [UserService],
})
export class AppModule {}

// ==========================================
// 6. Bootstrap Server (Khởi chạy ứng dụng)
// ==========================================
async function bootstrap() {
  const app = await NestFactory.create(AppModule);
  
  // Áp dụng Logger Interceptor toàn cục
  app.useGlobalInterceptors(new LoggingInterceptor());
  
  const PORT = 8083;
  await app.listen(PORT);
  console.log(`NestJS app đang chạy trên cổng ${PORT}...`);
}
bootstrap();
