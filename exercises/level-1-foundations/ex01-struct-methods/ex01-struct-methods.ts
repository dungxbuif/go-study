/**
 * Ex01: Struct & Methods — TypeScript Version
 * 
 * 🧠 So sánh key:
 * - Node.js: Dùng class, constructor, và cơ chế throw Error khi xảy ra lỗi.
 * - Go:      Không có class, dùng struct và các methods có receiver (pointer/value receiver),
 *            lỗi được trả về như một giá trị (result, error) thay vì throw.
 * 
 * 💡 Sự khác biệt lớn nhất:
 * 1. Go không có class hay constructor thừa kế, thay vào đó dùng struct và NewXxx() factory function.
 * 2. Go không throw exception mà trả về multiple values (value, error), buộc caller phải kiểm tra lỗi ngay lập tức.
 * 3. Go phân biệt pointer receiver (*T - có thể thay đổi state của struct) và value receiver (T - chỉ đọc, nhận bản sao).
 */

class User {
  public id: number;
  public name: string;
  public email: string;
  public createdAt: Date;

  constructor(id: number, name: string, email: string) {
    this.id = id;
    this.name = name;
    this.email = email;
    this.createdAt = new Date();
  }
}

class UserService {
  private users: User[];
  private nextID: number;

  constructor() {
    this.users = [];
    this.nextID = 1;
  }

  public create(name: string, email: string): User {
    const user = new User(this.nextID++, name, email);
    this.users.push(user);
    return user;
  }

  public findById(id: number): User {
    const user = this.users.find((u) => u.id === id);
    if (!user) {
      throw new Error(`user with id ${id} not found`);
    }
    return user;
  }

  public findByEmail(email: string): User {
    const user = this.users.find((u) => u.email === email);
    if (!user) {
      throw new Error(`user with email ${email} not found`);
    }
    return user;
  }

  public count(): number {
    return this.users.length;
  }
}

const svc = new UserService();

svc.create("Alice", "alice@example.com");
svc.create("Bob", "bob@example.com");
svc.create("Charlie", "charlie@example.com");

console.log(`Total users: ${svc.count()}`);
console.log("Found by ID:", svc.findById(2));
console.log("Found by email:", svc.findByEmail("alice@example.com"));

try {
  svc.findById(999);
} catch (err) {
  console.error("Error:", (err as Error).message);
}
