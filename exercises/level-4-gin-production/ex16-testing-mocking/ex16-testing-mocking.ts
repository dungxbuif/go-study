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
