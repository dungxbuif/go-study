/**
 * Ex13: Clean Architecture — TypeScript Version
 * 
 * 🧠 So sánh key:
 * - TypeScript: Trong NestJS hoặc các cấu trúc TS nâng cao, ta dùng Decorators, Classes 
 *               và IoC Containers (như InversifyJS hay Nest DI) để tự động hóa Dependency Injection.
 * - Go:         Clean Architecture trong Go hoàn toàn dựa trên Interfaces và việc tự liên kết (wiring) 
 *               dependencies bằng tay trong file main.go (Manual Dependency Injection).
 * 
 * 💡 Sự khác biệt lớn nhất:
 * 1. Go không có reflection động dạng Decorators (metadata reflection của TS). Mọi thứ đều phải 
 *    tường minh, không có "phép thuật" ngầm (magic).
 * 2. Cấu trúc Clean Architecture giúp cô lập hoàn toàn business logic (Usecase) khỏi các tác nhân 
 *    bên ngoài như HTTP Framework (Gin) hay Database Driver (GORM).
 */

import express, { Request, Response } from 'express';

interface Todo {
  id: number;
  title: string;
  completed: boolean;
}

interface ITodoRepository {
  create(title: string): Promise<Todo>;
  findById(id: number): Promise<Todo | null>;
  findAll(): Promise<Todo[]>;
}

class InMemoryTodoRepository implements ITodoRepository {
  private todos: Map<number, Todo> = new Map();
  private nextId = 1;

  async create(title: string): Promise<Todo> {
    const todo: Todo = { id: this.nextId++, title, completed: false };
    this.todos.set(todo.id, todo);
    return todo;
  }

  async findById(id: number): Promise<Todo | null> {
    return this.todos.get(id) || null;
  }

  async findAll(): Promise<Todo[]> {
    return Array.from(this.todos.values());
  }
}

class TodoUsecase {
  constructor(private todoRepo: ITodoRepository) {}

  async createTodo(title: string): Promise<Todo> {
    if (!title || title.trim() === '') {
      throw new Error('Title cannot be empty');
    }
    return this.todoRepo.create(title);
  }

  async getTodo(id: number): Promise<Todo> {
    const todo = await this.todoRepo.findById(id);
    if (!todo) {
      throw new Error('Todo not found');
    }
    return todo;
  }

  async listTodos(): Promise<Todo[]> {
    return this.todoRepo.findAll();
  }
}

class TodoHandler {
  constructor(private todoUsecase: TodoUsecase) {}

  async handleCreate(req: Request, res: Response) {
    try {
      const { title } = req.body;
      const todo = await this.todoUsecase.createTodo(title);
      res.status(201).json({ success: true, data: todo });
    } catch (err) {
      res.status(400).json({ success: false, error: (err as Error).message });
    }
  }

  async handleGet(req: Request, res: Response) {
    try {
      const id = parseInt(req.params.id, 10);
      const todo = await this.todoUsecase.getTodo(id);
      res.json({ success: true, data: todo });
    } catch (err) {
      res.status(404).json({ success: false, error: (err as Error).message });
    }
  }

  async handleList(req: Request, res: Response) {
    const todos = await this.todoUsecase.listTodos();
    res.json({ success: true, data: todos });
  }
}

const app = express();
app.use(express.json());

const todoRepo = new InMemoryTodoRepository();
const todoUsecase = new TodoUsecase(todoRepo);
const todoHandler = new TodoHandler(todoUsecase);

app.post('/todos', (req, res) => todoHandler.handleCreate(req, res));
app.get('/todos/:id', (req, res) => todoHandler.handleGet(req, res));
app.get('/todos', (req, res) => todoHandler.handleList(req, res));

app.listen(8080, () => {
  console.log('Clean Architecture Todo server running on port 8080');
});
