/**
 * Ex16: Testing & Mocking — TypeScript Version
 *
 * 🧠 So sánh key:
 * - Node.js: Dùng Jest làm test runner, sử dụng các thư viện Mock mạnh mẽ (`jest.mock()`)
 *            cho phép "hack" và ghi đè bất cứ module/file nào trong hệ thống lúc chạy test.
 * - Go:      Không dùng "magic mock" kiểu JS/Python được. Testing trong Go hoàn toàn dựa trên
 *            Interface-based Mocking. Code bắt buộc phải được thiết kế decoupled bằng Interface,
 *            từ đó tự viết Struct Mock khớp với Interface để inject vào lúc chạy test.
 *
 * 💡 Sự khác biệt lớn nhất:
 * 1. Go khuyến khích viết Table-Driven Tests (một mảng các testcases chứa input, target và expected).
 * 2. Cực kỳ dễ dàng phát hiện Race Condition khi viết test nhờ cờ `-race` mặc định (`go test -race ./...`).
 */

import * as assert from 'assert';

interface Todo {
   id: number;
   title: string;
   completed: boolean;
}

interface ITodoRepository {
   create(title: string): Promise<Todo>;
}

// MockTodoRepository: Triển khai thủ công Mock Interface.
//
// 🧠 CƠ CHẾ MOCK TRONG JAVASCRIPT VS GOLANG STATIC TYPING:
// - Trong Javascript, nhờ tính chất dynamic typing, ta có thể "monkey patch" bất kỳ đối tượng nào lúc chạy test,
//   thậm chí không cần định nghĩa Class Mock đầy đủ. Thư viện Jest có thể tự động thay thế `ITodoRepository` bằng
//   một đối tượng rỗng `{ create: jest.fn() }` và V8 Engine vẫn chạy trơn tru.
// - Tuy nhiên, trong TypeScript và Go, để đảm bảo tính an toàn kiểu dữ liệu (type-safety) lúc compile,
//   ta thường tự viết Class Mock implement đầy đủ Interface như thế này.
// - Cách viết Class Mock thủ công này giúp kiểm thử cực kỳ trực quan, độc lập với các thư viện Mock ma thuật
//   và dễ dàng chuyển dịch sang Go mà không gặp bất cứ rào cản tư duy nào.
class MockTodoRepository implements ITodoRepository {
   public todos: Todo[];
   public shouldFail: boolean;

   constructor() {
      this.todos = [];
      this.shouldFail = false;
   }

   public async create(title: string): Promise<Todo> {
      if (this.shouldFail) {
         throw new Error('Database connection failed');
      }
      const todo: Todo = { id: this.todos.length + 1, title, completed: false };
      this.todos.push(todo);
      return todo;
   }
}

class TodoUsecase {
   private todoRepo: ITodoRepository;

   constructor(todoRepo: ITodoRepository) {
      this.todoRepo = todoRepo;
   }

   public async createTodo(title: string): Promise<Todo> {
      if (!title || title.trim() === '') {
         throw new Error('Title is required');
      }
      return this.todoRepo.create(title);
   }
}

interface TestCase {
   name: string;
   title: string;
   shouldFailDb: boolean;
   expectedError: string | null;
   expectTodo: boolean;
}

// runTableDrivenTests: Mô phỏng mô hình Table-Driven Test kinh điển của Go trong TypeScript.
//
// 🧠 LỢI ÍCH CỦA TABLE-DRIVEN TESTS:
// - Giúp code kiểm thử cực kỳ sạch sẽ và dễ mở rộng.
// - Khi muốn kiểm tra thêm một ca biên mới (Edge Case), ta chỉ việc khai báo thêm 1 dòng Object vào mảng `tests`
//   mà không cần phải copy-paste hàng chục dòng code `describe` hay `it` lặp đi lặp lại.
// - Việc chạy lặp trong khối `try/catch` kết hợp các câu lệnh khẳng định `assert` cung cấp một khung kiểm thử chuẩn chỉ,
//   dễ dàng map trực tiếp sang framework `testing` của Golang.
async function runTableDrivenTests(): Promise<void> {
   const tests: TestCase[] = [
      {
         name: 'successful creation',
         title: 'Buy groceries',
         shouldFailDb: false,
         expectedError: null,
         expectTodo: true,
      },
      {
         name: 'empty title validation failure',
         title: '',
         shouldFailDb: false,
         expectedError: 'Title is required',
         expectTodo: false,
      },
      {
         name: 'database connection error',
         title: 'Valid title',
         shouldFailDb: true,
         expectedError: 'Database connection failed',
         expectTodo: false,
      },
   ];

   for (const tt of tests) {
      console.log(`Running test: ${tt.name}`);
      const mockRepo = new MockTodoRepository();
      mockRepo.shouldFail = tt.shouldFailDb;
      const usecase = new TodoUsecase(mockRepo);

      try {
         const todo = await usecase.createTodo(tt.title);
         assert.strictEqual(tt.expectedError, null);
         assert.ok(todo);
         assert.strictEqual(todo.title, tt.title);
      } catch (err: any) {
         assert.ok(tt.expectedError);
         assert.strictEqual(err.message, tt.expectedError);
      }
   }
   console.log('All tests completed successfully!');
}

runTableDrivenTests();
