/**
 * Ex10: CRUD In-Memory — TypeScript Version
 * 
 * 🧠 So sánh key:
 * - Node.js: Lưu dữ liệu trong cấu trúc Map hoặc Array toàn cục. Không cần quan tâm đến Lock 
 *            cho các thao tác ghi đồng thời (concurrent write) vì Event Loop là đơn luồng (single-thread).
 * - Go:      HTTP server của Gin chạy đa luồng đồng thời (mỗi request là một goroutine riêng).
 *            Do đó khi đọc/ghi vào shared memory (map, slice) bắt buộc phải sử dụng `sync.RWMutex` 
 *            hoặc `sync.Mutex` để tránh lỗi race condition/panic do ghi đè vùng nhớ song song.
 */

import express, { Request, Response } from 'express';
const app = express();

app.use(express.json());

interface Todo {
  id: number;
  title: string;
  completed: boolean;
  created_at: Date;
}

const todos = new Map<number, Todo>();
let nextId = 1;

app.post('/todos', (req: Request, res: Response) => {
  const { title } = req.body;
  if (!title) {
    return res.status(400).json({ success: false, error: 'Title is required' });
  }

  const todo: Todo = {
    id: nextId++,
    title,
    completed: false,
    created_at: new Date()
  };

  todos.set(todo.id, todo);
  res.status(201).json({ success: true, data: todo });
});

app.get('/todos', (req: Request, res: Response) => {
  const completed = req.query.completed as string | undefined;
  let list = Array.from(todos.values());

  if (completed !== undefined) {
    const isCompleted = completed === 'true';
    list = list.filter(t => t.completed === isCompleted);
  }

  res.json({ success: true, data: list });
});

app.get('/todos/:id', (req: Request, res: Response) => {
  const id = parseInt(req.params.id, 10);
  const todo = todos.get(id);

  if (!todo) {
    return res.status(404).json({ success: false, error: 'Todo not found' });
  }

  res.json({ success: true, data: todo });
});

app.put('/todos/:id', (req: Request, res: Response) => {
  const id = parseInt(req.params.id, 10);
  const todo = todos.get(id);

  if (!todo) {
    return res.status(404).json({ success: false, error: 'Todo not found' });
  }

  const { title, completed } = req.body;
  if (title !== undefined) todo.title = title;
  if (completed !== undefined) todo.completed = completed;

  res.json({ success: true, data: todo });
});

app.delete('/todos/:id', (req: Request, res: Response) => {
  const id = parseInt(req.params.id, 10);
  if (!todos.has(id)) {
    return res.status(404).json({ success: false, error: 'Todo not found' });
  }

  todos.delete(id);
  res.json({ success: true, message: 'Todo deleted successfully' });
});

app.listen(8080, () => {
  console.log('Todo server running on port 8080');
});
