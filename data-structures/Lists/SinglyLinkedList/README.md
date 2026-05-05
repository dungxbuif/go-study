# Linked List (Danh sách liên kết)
- Các phần tử không nối tiếp nhau trong bộ nhớ.
- Mỗi phần tử chứa data và con trỏ tới phần tử tiếp theo.

## 1. Singly Linked List (Danh sách liên kết đơn)
Một SinglyLinkedList bao gồm:
- head: Con trỏ tới Node đầu tiên.
- tail: Con trỏ tới Node cuối cùng.
- size: Số lượng phần tử trong danh sách.

Mỗi Node chứa:
- data: Giá trị cần lưu trữ.
- next: Con trỏ tới Node tiếp theo. Node cuối cùng có next = nil.

### 🛠️ Các phương thức chính (Methods)

| Phương thức | Mô tả | Độ phức tạp |
| :--- | :--- | :--- |
| `New()` | Khởi tạo một danh sách mới. | O(1) |
| `PushFront(v)` | Thêm giá trị `v` vào đầu danh sách. | O(1) |
| `PushBack(v)` | Thêm giá trị `v` vào cuối danh sách. | O(1) |
| `InsertAfter(v, mark)` | Chèn giá trị `v` vào sau phần tử `mark`. | O(1) |
| `Remove(e)` | Xóa phần tử `e` khỏi danh sách. | O(n)* |
| `MoveToFront(e)` | Di chuyển phần tử `e` lên đầu danh sách. | O(n) |
| `MoveToBack(e)` | Di chuyển phần tử `e` xuống cuối danh sách. | O(n) |