/**
 * Ex04: Pointers & Generics — TypeScript Version
 * 
 * 🧠 So sánh key:
 * - TypeScript: Object luôn được truyền bằng tham chiếu (reference), không có con trỏ (pointer) rõ ràng.
 *               Generics dùng cú pháp `<T>` tương tự hầu hết các ngôn ngữ OOP hiện đại.
 * - Go:         Quản lý vùng nhớ tường minh bằng con trỏ (`*` và `&`).
 *               Generics dùng cú pháp `[T any]` hoặc các constraints khác (như `[T comparable]`).
 * 
 * 💡 Sự khác biệt lớn nhất:
 * 1. Go phân biệt rõ truyền giá trị (truyền bản sao - copy) và truyền con trỏ (truyền địa chỉ ô nhớ gốc).
 * 2. Pointer có thể mang giá trị `nil`, dereference một nil pointer sẽ gây PANIC lập tức (không loose như null/undefined của JS).
 * 3. Cú pháp generic của Go sử dụng dấu ngoặc vuông thay vì ngoặc nhọn để tránh nhập nhèm ngữ nghĩa trong phân tích cú pháp.
 */

class ListNode<T> {
  value: T;
  next: ListNode<T> | null;

  constructor(value: T) {
    this.value = value;
    this.next = null;
  }
}

class LinkedList<T> {
  head: ListNode<T> | null;
  size: number;

  constructor() {
    this.head = null;
    this.size = 0;
  }

  push(val: T): void {
    const node = new ListNode(val);
    node.next = this.head;
    this.head = node;
    this.size++;
  }

  pop(): T {
    if (!this.head) {
      throw new Error("list is empty");
    }
    const val = this.head.value;
    this.head = this.head.next;
    this.size--;
    return val;
  }

  find(predicate: (val: T) => boolean): ListNode<T> | null {
    let current = this.head;
    while (current !== null) {
      if (predicate(current.value)) {
        return current;
      }
      current = current.next;
    }
    return null;
  }

  print(): void {
    const values: T[] = [];
    let current = this.head;
    while (current !== null) {
      values.push(current.value);
      current = current.next;
    }
    console.log(`LinkedList(${this.size}): [${values.join(" → ")}]`);
  }
}

const intList = new LinkedList<number>();
intList.push(10);
intList.push(20);
intList.push(30);
intList.print();

const found = intList.find((v) => v === 20);
console.log("Found:", found?.value);

const popped = intList.pop();
console.log("Popped:", popped);
intList.print();

const strList = new LinkedList<string>();
strList.push("Go");
strList.push("TypeScript");
strList.push("Rust");
strList.print();

const goNode = strList.find((v) => v.includes("Go"));
console.log("Found:", goNode?.value);
