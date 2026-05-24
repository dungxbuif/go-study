package main

import (
	"fmt"
)

// ============================================================================
// 🧠 GIẢI THÍCH VỀ CON TRỎ (POINTER) TRONG GO CHO NODE.JS DEV:
//
// 1. Trong JavaScript:
//    - Mọi Object, Array, Class instance đều được truyền dưới dạng THAM CHIẾU (Reference).
//    - Khi bạn gán `obj2 = obj1`, cả hai biến đều trỏ chung vào cùng một vùng nhớ.
//
// 2. Trong Go:
//    - Mọi kiểu dữ liệu (kể cả struct) mặc định đều được truyền dưới dạng THAM TRỊ (Value - copy giá trị).
//    - Nếu bạn gán `struct2 = struct1`, Go sẽ COPY toàn bộ dữ liệu của struct1 sang ô nhớ mới cho struct2.
//    - Để thay đổi dữ liệu của đối tượng gốc hoặc tránh việc copy cấu trúc dữ liệu lớn gây tốn RAM,
//      chúng ta sử dụng Con Trỏ (Pointer).
//
// 3. Ký tự cần nhớ:
//    - `*T` (Ví dụ `*Node`): Đọc là "Con trỏ trỏ tới kiểu T" (biến này chỉ lưu địa chỉ ô nhớ, không lưu giá trị trực tiếp).
//    - `&`  (Ví dụ `&Node{}`): Đọc là "Lấy địa chỉ ô nhớ của...". Dùng để chuyển một giá trị thường thành con trỏ.
//    - `*ptr` (Dereferencing): Đọc là "Giá trị tại địa chỉ mà ptr đang trỏ tới".
// ============================================================================

// Node đại diện cho một mắt xích trong danh sách liên kết.
type Node[T any] struct {
	Value T
	Next  *Node[T] // Next là một CON TRỎ. Nó không lưu Node tiếp theo trực tiếp, mà lưu "Địa chỉ ô nhớ" của Node tiếp theo.
}

// NewNode là hàm khởi tạo một Node mới. Nó trả về *Node[T] (tức là địa chỉ vùng nhớ chứa Node đó).
func NewNode[T any](value T) *Node[T] {
	// Ký tự `&` ở đây có nghĩa là: Khởi tạo struct Node trong bộ nhớ (RAM), sau đó lấy ĐỊA CHỈ ô nhớ của nó trả về.
	return &Node[T]{
		Value: value,
		Next:  nil, // nil tương đương null/undefined trong JS, thể hiện con trỏ chưa trỏ vào ô nhớ nào.
	}
}

type LinkedList[T any] struct {
	Head *Node[T] // Head là con trỏ, giữ địa chỉ của Node đầu tiên.
	Size int      // Size là giá trị thường (kiểu int), lưu số lượng phần tử.
}

// NewLinkedList khởi tạo danh sách và trả về con trỏ quản lý danh sách đó.
func NewLinkedList[T any]() *LinkedList[T] {
	return &LinkedList[T]{
		Head: nil,
		Size: 0,
	}
}

// Receiver `(l *LinkedList[T])` bắt buộc phải là CON TRỎ `*`
// Vì khi thêm phần tử (Push), ta cần SỬA ĐỔI thuộc tính Head và Size của đối tượng gốc l.
// Nếu không dùng `*`, Go sẽ tạo một bản copy của LinkedList khi gọi hàm, và thuộc tính Head gốc sẽ không bao giờ thay đổi.
func (l *LinkedList[T]) Push(val T) {
	node := NewNode[T](val) // node lúc này là con trỏ (*Node[T])

	// Gán địa chỉ Node đầu tiên hiện tại của List vào trường Next của Node mới.
	// Lưu ý: Trong Go, thay vì dùng toán tử `->` như C/C++, ta dùng thẳng dấu chấm `.`
	// Go compiler sẽ tự động ngầm hiểu (auto-dereference) `(*node).Next = l.Head`
	node.Next = l.Head

	// Cập nhật Head của danh sách trỏ vào địa chỉ của Node mới tạo.
	l.Head = node
	l.Size++
}

// Pop lấy phần tử đầu tiên ra. Trả về giá trị của phần tử và trạng thái thành công (bool)
func (l *LinkedList[T]) Pop() (T, bool) {
	if l.Head == nil { // Nếu Head đang trỏ vào nil -> Danh sách trống.
		var zero T
		return zero, false
	}

	// l.Head đang lưu địa chỉ của Node đầu tiên.
	// l.Head.Value lấy giá trị nằm bên trong Node đó.
	val := l.Head.Value

	// Dịch chuyển Head trỏ sang Node tiếp theo (Next).
	// Next cũng là một con trỏ lưu địa chỉ Node tiếp theo.
	l.Head = l.Head.Next
	l.Size--

	return val, true
}

// Find duyệt qua danh sách tìm Node đầu tiên thoả mãn điều kiện.
func (l *LinkedList[T]) Find(predicate func(T) bool) *Node[T] {
	current := l.Head // Gán địa chỉ của Head cho biến chạy current.

	for current != nil { // Vòng lặp chạy đến khi current gặp nil (kết thúc danh sách).
		if predicate(current.Value) {
			return current // Trả về địa chỉ của Node tìm thấy.
		}
		current = current.Next // Dịch chuyển biến chạy sang địa chỉ của Node tiếp theo.
	}
	return nil
}

func (l *LinkedList[T]) Print() {
	values := []T{}
	current := l.Head
	for current != nil {
		values = append(values, current.Value)
		current = current.Next
	}
	fmt.Printf("LinkedList(%d): %v\n", l.Size, values)
}

func main() {
	// intList là con trỏ trỏ tới struct LinkedList trong RAM (*LinkedList[int])
	intList := NewLinkedList[int]()

	intList.Push(10)
	intList.Push(20)
	intList.Push(30)

	intList.Print() // LinkedList(3): [30 20 10]

	fmt.Printf("Int List Size: %d\n", intList.Size)

	// found nhận được con trỏ *Node[int] trỏ thẳng tới ô nhớ chứa giá trị 20 trong Linked List.
	found := intList.Find(func(v int) bool { return v == 20 })
	if found != nil {
		fmt.Printf("Found: %v (địa chỉ ô nhớ: %p)\n", found.Value, found)
	}

	val, _ := intList.Pop()
	fmt.Printf("Popped: %d\n", val)
	intList.Print() // LinkedList(2): [20 10]
}
