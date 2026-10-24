- Go chỉ có for không có While
- For Range tạo ra bản sao của phần tử
- Switch Go tự dừng ko cần break

# Why Go efficient?

- **Code đơn giản, tường minh (Ít abstractions):** Dễ code, dễ duy trì (maintain). Mọi thứ xoay quanh function và struct. Thiết kế ngôn ngữ giữ sự cân bằng: hạn chế "phép thuật" (magic) giúp code dễ đọc, dễ học đối với mọi thành viên trong team.
- **Nhiều tools tích hợp sẵn:** Định dạng code (`gofmt`), testing, profiling hiệu năng (`pprof`), phát hiện race condition (`-race`),... giúp nhà phát triển tối ưu hóa ứng dụng cực nhanh mà không cần cài thêm thư viện bên thứ ba.

---

## 1. Low level control (Kiểm soát cấp thấp)

Go cho phép bạn kiểm soát dữ liệu được cấp phát và hoạt động như thế nào trong bộ nhớ (Memory):

### Value type vs Pointer type

Bạn có quyền quyết định truyền giá trị hay truyền tham chiếu:

```go
// Value type: copy toàn bộ dữ liệu, thường nằm trên STACK
a := User{Name: "Antigravity"}
b := a  // b sao chép toàn bộ data của a, sửa b không ảnh hưởng đến a

// Pointer type: chỉ sao chép địa chỉ vùng nhớ, dữ liệu nằm trên HEAP
a := &User{Name: "Antigravity"}
b := a  // b trỏ tới cùng vùng nhớ với a, sửa b sẽ làm thay đổi a
```

### Escape Analysis (Phân tích thoát)

Compiler của Go tự động phân tích tĩnh (static analysis) lúc build để quyết định một biến nằm ở đâu:

- **Stack:** Truy cập cực nhanh, tự động giải phóng khi hàm kết thúc (không cần GC dọn dẹp).
- **Heap:** Chậm hơn Stack, cần GC theo dõi và thu gom.

```go
func newPointer() *User {
    u := User{Name: "A"}
    return &u // Trả về con trỏ ra ngoài hàm -> Biến u "escape" (thoát) lên HEAP
}

func newVal() {
    u := User{Name: "A"}
    fmt.Println(u.Name) // Biến u chỉ dùng nội bộ trong hàm -> Giữ lại ở STACK
}
```

> _Mẹo:_ Bạn có thể xem kết quả phân tích bằng lệnh `go build -gcflags="-m" main.go`.

---

## 2. M:N Scheduling (Bộ lập lịch M:N)

---

## 3. Memory efficiency (Hiệu quả bộ nhớ)

### Bộ nhớ liên tiếp (Contiguous Memory Layout)

Trong Go, các phần tử của Array hoặc Struct nằm liền kề nhau trong RAM:

```go
type Point struct {
    X int64 // 8 bytes
    Y int64 // 8 bytes
}
points := [3]Point{{1,2}, {3,4}, {5,6}}
```

Vùng nhớ của array `points` sẽ trải dài liên tiếp **48 bytes** trong RAM:
`[ 1 ][ 2 ][ 3 ][ 4 ][ 5 ][ 6 ]` (mỗi ô 8 bytes).

**Lợi ích:** CPU đọc dữ liệu theo các đường truyền gọi là _Cache Line_ (thường là 64 bytes). Khi bạn truy cập `points[0]`, CPU sẽ tiện tay nạp luôn `points[1]` vào L1/L2 Cache của CPU. Việc truy cập phần tử tiếp theo sẽ nhanh hơn hàng trăm lần so với việc phải đọc từ RAM chính.

### Đệm cấu trúc (Struct Padding)

Go sắp xếp các trường trong struct dựa theo nguyên tắc căn lề (alignment rules). Thứ tự khai báo trường có thể làm thay đổi dung lượng RAM tiêu thụ:

```go
// Tốn 24 bytes (Do bị chèn padding)
type Bad struct {
    a bool   // 1 byte + 7 bytes padding để căn lề cho b
    b int64  // 8 bytes
    c bool   // 1 byte + 7 bytes padding ở cuối
}

// Tốn 16 bytes (Đã tối ưu hóa)
type Good struct {
    b int64  // 8 bytes
    a bool   // 1 byte
    c bool   // 1 byte + 6 bytes padding ở cuối
}
```

---

## 4. Fast compilation (Biên dịch nhanh)

Tốc độ biên dịch của Go nhanh vượt trội do những quyết định thiết kế ngôn ngữ cực kỳ thực tế:

- **Không cho phép Import vòng quanh (No Circular Imports):** Go cấm Package A import Package B và ngược lại. Compiler duyệt đồ thị phụ thuộc (dependency graph) theo một chiều duy nhất, không phải quay đầu tính toán lại.
- **Khai báo import rõ ràng:** Compiler chỉ cần đọc dòng đầu các file để biết chính xác cần compile những gì, không cần quét toàn bộ source code của project.
- **Import không dùng là lỗi compile:** Giúp loại bỏ hoàn toàn code thừa ngay lúc viết.
- **Không có cơ chế Generics quá phức tạp:** Hạn chế tối đa việc tạo thêm code giả trong lúc build (như template metaprogramming của C++).

---

## 5. Garbage Collector (GC) low-latency

Go GC được thiết kế ưu tiên giảm thiểu tối đa độ trễ (latency) hơn là tối ưu hóa băng thông (throughput).

### Thuật toán Đánh dấu 3 màu (Tri-color concurrent mark-and-sweep)

1. **Trắng (White):** Các đối tượng chưa được kiểm tra (ứng viên tiềm năng làm rác).
2. **Xám (Grey):** Các đối tượng đã được tìm thấy, nhưng các đối tượng con mà nó trỏ tới chưa được kiểm tra.
3. **Đen (Black):** Các đối tượng đang được sử dụng và các đối tượng con của nó đã được quét xong.

```
Ban đầu: Tất cả objects đều màu Trắng
   │
   ▼
Mark phase: Quét từ các đối tượng gốc (Stack, Globals...) -> Đánh dấu Xám
   │
   ▼
Quét đệ quy các pointer bên trong Grey object -> Chuyển object đó sang Đen, chuyển các con của nó sang Xám
   │
   ▼
Sweep phase: Khi không còn object Xám, mọi object còn Trắng chính là rác -> Giải phóng.
```

### Tại sao thời gian dừng hệ thống (Stop-The-World - STW) lại dưới 1ms?

- **Chạy song song (Concurrent):** Quá trình GC chạy song song cùng lúc với các Goroutine của ứng dụng đang chạy.
- **Rào cản ghi (Write Barrier):** Khi ứng dụng thay đổi các liên kết pointer trong lúc GC đang hoạt động, Write Barrier sẽ can thiệp để đảm bảo các object mới/bị thay đổi không bị GC đánh dấu sai màu hoặc bỏ sót.
- **STW cực ngắn:** Go chỉ dừng ứng dụng hoàn toàn ở 2 thời điểm siêu ngắn: lúc bật Write Barrier ban đầu và lúc tắt Write Barrier/hoàn thành Mark. Mỗi lần dừng chỉ mất từ vài chục micro giây ($\mu s$).

---

# Go Scheduler (GMP Model)

Go runtime sở hữu một bộ lập lịch (Scheduler) tùy biến chạy trên không gian người dùng (user space), quản lý luồng song song thông qua mô hình GMP.

```
          ┌─────────────────────────────────────────┐
          │            Go Runtime Scheduler          │
          └─────────────────────────────────────────┘
                            │
           ┌────────────────┼────────────────┐
           ▼                ▼                ▼
        ┌─────┐          ┌─────┐          ┌─────┐
        │  P0 │          │  P1 │          │  P2 │     (Logical Processors - Mặc định bằng số CPU Cores)
        └──┬──┘          └──┬──┘          └──┬──┘
           │                │                │
     ┌─────┴─────┐    ┌─────┴─────┐    ┌─────┴─────┐
     │ Local Run │    │ Local Run │    │ Local Run │  (Hàng đợi chứa các Goroutine đang chờ chạy)
     │   Queue   │    │   Queue   │    │   Queue   │
     │ G1, G2    │    │ G3        │    │ Rỗng      │
     └─────┬─────┘    └─────┬─────┘    └─────┬─────┘
           │                │                │
           ▼                ▼                ▼
        ┌─────┐          ┌─────┐          ┌─────┐
        │  M0 │          │  M1 │          │  M2 │     (OS Threads thật sự do nhân OS quản lý)
        └─────┘          └─────┘          └─────┘
```

## Định nghĩa các thành phần:

- **G (Goroutine):** Đơn vị thực thi nhỏ nhất. Cực nhẹ, stack khởi tạo ban đầu chỉ khoảng 2KB và có khả năng tự động tăng/giảm dung lượng linh hoạt.
- **M (Machine/OS Thread):** Luồng xử lý thực tế của hệ điều hành.
- **P (Processor):** Đại diện cho tài nguyên xử lý logic. Số lượng `P` thường bằng số CPU Cores (cấu hình qua biến môi trường `GOMAXPROCS`). Mỗi `P` sở hữu một hàng đợi cục bộ (Local Run Queue) chứa các `G` đang sẵn sàng chạy.

## Các cơ chế vận hành chính:

### 1. Work Stealing (Ăn cắp việc)

- Khi hàng đợi của `P2` bị rỗng (không còn `G` nào để xử lý), bộ lập lịch sẽ điều phối `P2` đi quét:
   1. Kiểm tra hàng đợi toàn cục (Global Run Queue) xem có việc không.
   2. Nếu Global cũng rỗng, nó sẽ tìm đến các `P` khác (như `P0` hoặc `P1`) để **"ăn cắp" một nửa** số lượng `G` đang chờ trong Local Run Queue của họ về xử lý. Điều này giúp tối ưu hóa hiệu suất CPU, tránh tình trạng Core chơi Core làm.

### 2. Handoff (Bàn giao Thread khi bị block)

- Giả sử `G1` đang chạy trên `M0` (được liên kết bởi `P0`) thực hiện một thao tác đọc file dung lượng lớn hoặc truy vấn mạng (Blocking System Call):
   1. `M0` sẽ bị khóa cứng (blocked) bởi hệ điều hành để chờ dữ liệu.
   2. Bộ lập lịch Go phát hiện ra và sẽ **ngay lập tức tách** `P0` ra khỏi `M0`.
   3. `P0` sẽ đi tìm một Thread `M` khác đang rỗi (hoặc tạo mới một `M`) để liên kết và tiếp tục chạy các `G2`, `G3` trong hàng đợi của nó.
   4. Khi `M0` hoàn thành chặn hệ thống, `G1` sẽ tìm một `P` trống để chạy tiếp, nếu không có nó sẽ được đẩy vào Global Run Queue, còn `M0` được đưa vào danh sách thread chờ.

### 3. Cooperative & Asynchronous Preemption (Nhường quyền & Ngắt không đồng bộ)

- **Cooperative (Nhường quyền tự nguyện):** Bản thân Goroutine khi biên dịch sẽ được chèn các đoạn mã kiểm tra (preemption points). Khi thực hiện gọi hàm, thao tác channel (`ch <- val`), hoặc gọi syscall, `G` sẽ tự động nhường sân khấu cho các `G` khác.
- **Asynchronous Preemption (Ngắt không đồng bộ):** Từ Go 1.14+, để tránh trường hợp một `G` chứa vòng lặp vô hạn (không gọi hàm, không thao tác I/O) chiếm dụng CPU mãi mãi, Go Runtime sẽ gửi một tín hiệu (Signal `SIGURG`) trực tiếp từ OS đến Thread `M` đang chạy để tạm dừng `G` cứng đầu đó lại và nhường chỗ cho các `G` khác.

---

# `sync.WaitGroup` (Chuẩn Go idiomatic)

- Nếu không dùng `sync.WaitGroup`, ta cần phải biết chính xác số lượng job/kết quả và đọc bằng vòng lặp đếm số lượng. Nếu số lượng job là động (dynamic) hoặc một số job bị lỗi/bị lọc bỏ giữa chừng (không gửi kết quả về `results`), hàm `main` sẽ đợi mãi mãi $\rightarrow$ **Deadlock**.

- Nếu dùng `sync.WaitGroup`, ta chỉ cần cộng `wg.Add(1)` khi bắt đầu worker goroutine, và gọi `wg.Done()` khi worker hoàn thành. Khi tất cả worker đã xong, ta gọi `wg.Wait()` để đóng channel kết quả.
