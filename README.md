# Go (Golang) Learning Journey 🚀

Chào mừng bạn đến với lộ trình học Go. Tài liệu này được tổng hợp từ nhiều nguồn uy tín và kinh nghiệm thực tế (bao gồm so sánh với Node.js) để giúp bạn nắm bắt Go một cách hệ thống và hiệu quả nhất.

---

## 🏁 1. Giới thiệu về Go

Go là ngôn ngữ lập trình mã nguồn mở được phát triển bởi Google. Nó được thiết kế để kết hợp hiệu năng của C/C++ với sự đơn giản và năng suất của các ngôn ngữ như Python/JavaScript.

### Đặc điểm cốt lõi:

- **Tĩnh & Mạnh (Static & Strong Type):** Bắt lỗi ngay lúc biên dịch.
- **Biên dịch (Compiled):** Chạy trực tiếp trên mã máy, không cần máy ảo (giống C++, khác Java/Node.js).
- **Đồng thời (Concurrency):** Hỗ trợ hàng triệu luồng đồng thời (Goroutines) cực kỳ nhẹ.
- **Đơn giản:** Chỉ có 25 từ khóa, không có lớp (class), không kế thừa phức tạp.
- **Hệ sinh thái Cloud-Native:** Là ngôn ngữ của Docker, Kubernetes, Terraform.

---

## 🛠 2. Tooling & Workspace

![Go Tools & Tools Ecosystem](./images/tools.png)

### Lệnh cơ bản

| Lệnh                 | Ý nghĩa                              |
| :------------------- | :----------------------------------- |
| `go run .`           | Biên dịch và chạy package hiện tại   |
| `go build`           | Xây dựng file thực thi (binary)      |
| `go mod init <name>` | Khởi tạo module mới                  |
| `go mod tidy`        | Dọn dẹp phụ thuộc (thêm/xóa package) |
| `go test ./...`      | Chạy toàn bộ test trong project      |
| `gofmt -w .`         | Tự động định dạng lại toàn bộ code   |

#### Quy tắc & Cú pháp quan trọng

- **Package main**: Mọi chương trình chạy được phải bắt đầu bằng `package main`.
- **Hàm main()**: Là điểm bắt đầu (entry point) của chương trình.
- **Dấu chấm phẩy (;)**: Go tự động thêm `;` khi biên dịch. Bạn **không nên** viết thủ công vào cuối dòng code.

### Cấu trúc Project khuyến nghị

```text
myapp/
├── go.mod        # Quản lý dependencies
├── main.go       # Entry point
├── internal/     # Code private, không cho package ngoài import
├── pkg/          # Code public, có thể tái sử dụng
└── cmd/          # Chứa các file thực thi khác nhau (Server, CLI...)
```

---

## 📜 3. Cú pháp cơ bản (Basic Syntax)

### 3. Khai báo Biến (Variables)

#### Các cách khai báo chính

Go hỗ trợ khai báo linh hoạt tùy vào phạm vi sử dụng.

- **Dùng từ khóa `var`:** Có thể dùng ở cả cấp độ package và function.
- **Sử dụng `:=` (Short declaration):** Chỉ dùng **trong hàm**. Loại này tự động nhận diện kiểu dữ liệu.

```go
var student1 string = "John" // Khai báo có kiểu
var student2 = "Jane"        // Tự động nhận diện kiểu
x := 2                       // Cách ngắn gọn (chỉ trong hàm)
```

#### Khai báo không khởi tạo (Zero Values)

Nếu không gán giá trị ngay, Go sẽ tự động gán giá trị mặc định cho biến:

- `int`: `0`
- `float32/64`: `0.0`
- `string`: `""`
- `bool`: `false`
- `pointer/interface/slice/map/channel`: `nil`

#### So sánh `var` và `:=`

| Đặc điểm      | `var`                                | `:=`                               |
| :------------ | :----------------------------------- | :--------------------------------- |
| **Phạm vi**   | 🌍 Trong & Ngoài hàm (Package level) | 🏠 Chỉ **Trong** hàm (Local level) |
| **Tách biệt** | ✅ Có thể khai báo trước, gán sau    | ❌ Phải khai báo và gán cùng lúc   |
| **Cú pháp**   | `var name type = val`                | `name := val`                      |

> [!TIP]
> Trong thực tế, hãy ưu tiên dùng `:=` bên trong hàm để code ngắn gọn. Chỉ dùng `var` khi cần khai báo biến toàn cục (package level) hoặc khi chưa biết giá trị khởi tạo ngay lập tức.

#### Khai báo nhiều biến

```go
// Trên cùng 1 dòng
var a, b, c, d int = 1, 3, 5, 7
var e, f = 6, "Hello" // Khác kiểu nếu không chỉ định type
g, h := 7, "World"

// Khai báo theo khối (Group block) - Tăng tính thẩm mỹ
var (
    userID   int
    userName string = "Guest"
    isActive bool   = true
)
```

> [!NOTE]
> Sử dụng khối `var (...)` giúp nhóm các biến có liên quan lại với nhau, giúp code sạch sẽ và dễ quản lý hơn.

#### Quy tắc đặt tên biến (Naming Rules)

- Phải bắt đầu bằng chữ cái hoặc dấu gạch dưới `_`.
- Không được bắt đầu bằng số.
- Chỉ chứa chữ cái, số và `_`.
- **Case-sensitive:** `age` và `Age` là hai biến khác nhau.
- Không chứa khoảng trắng và không trùng từ khóa Go.

**Conventions:**

- **Camel Case:** `myVariableName` (Dùng cho biến local).
- **Pascal Case:** `MyVariableName` (Dùng để **Export** biến ra ngoài package).
- **Snake Case:** `my_variable_name` (Ít dùng trong Go).

#### Hằng số (Constants) & Iota

Hằng số là các giá trị không đổi sau khi khai báo. `iota` giúp tạo các số tăng dần tự động.

```go
const Pi = 3.14
const (
    Read = 1 << iota  // 1 (1 << 0)
    Write             // 2 (1 << 1)
    Execute           // 4 (1 << 2)
)
```

> [!NOTE]
> **Lưu ý về Hằng số:**
>
> - Hằng số phải được gán giá trị ngay khi khai báo.
> - Không thể dùng cú pháp `:=` cho hằng số.
> - `iota` sẽ tự động reset về 0 ở mỗi khối `const (...)` mới.

#### Hàm Xuất dữ liệu (Output Functions)

Sử dụng package `fmt` để in dữ liệu ra màn hình:

- **`Print()`**: In các đối mục sát nhau. Không tự thêm dòng mới.
- **`Println()`**: Thêm dấu cách giữa các đối mục và tự động xuống dòng (`\n`).
- **`Printf()`**: In định dạng (Format) bằng các **Verbs**:
   - `%v`: Giá trị mặc định.
   - `%T`: Kiểu dữ liệu của biến.
   - `%d`: Số nguyên (decimal).
   - `%s`: Chuỗi (string).
   - `%f`: Số thực (float). Ví dụ `%.2f` lấy 2 chữ số thập phân.
   - `%t`: Boolean (true/false).

#### Hệ thống Toán tử (Operators)

| Loại        | Toán tử                                                     |
| :---------- | :---------------------------------------------------------- | ------------------ | ------ |
| **Số học**  | `+`, `-`, `*`, `/`, `%`, `++`, `--`                         |
| **Gán**     | `=`, `+=`, `-=`, `*=`, `/=`, `%=`, `&=`, `^=`, `<<=`, `>>=` |
| **So sánh** | `==`, `!=`, `<`, `>`, `<=`, `>=`                            |
| **Logic**   | `&&`, `                                                     |                    | `, `!` |
| **Bitwise** | `&`, `                                                      | `, `^`, `<<`, `>>` |

### Kiểu dữ liệu cơ bản

- **Số nguyên:** `int`, `int8`, `int64`, `uint`, `byte` (alias của `uint8`).
- **Số thực:** `float32`, `float64`.
- **Logic:** `bool`.
- **Chuỗi:** `string` (immutable UTF-8).
- **Ký tự:** `rune` (alias của `int32`, đại diện cho 1 Unicode code point).

---

## ⚙️ 4. Control Flow

### If / Else

Có thể khai báo biến tạm ngay trong `if`.

```go
if v := compute(); v > 10 {
    fmt.Println(v)
} else {
    fmt.Println("Small")
}
```

### Loops (Chỉ có `for`)

Go không có `while` hay `do-while`.

```go
// 1. Vòng lặp cơ bản
for i := 0; i < 5; i++ {
    fmt.Println(i)
}

// 2. Vòng lặp dạng while
i := 1
for i < 5 {
    i *= 2
}

// 3. Lặp qua Arrays/Slices/Maps (Range)
for index, value := range fruits {
    fmt.Printf("Vị trí %d là %s\n", index, value)
}
```

> [!TIP]
> **Lưu ý về Range:**
>
> - Nếu bạn không cần `index`, hãy dùng dấu gạch dưới `_`: `for _, value := range fruits`.
> - Range tạo ra một **bản sao (copy)** của phần tử, không phải chính phần tử đó.

### Switch

Không cần lệnh `break` (Go tự động dừng). Hỗ trợ **Multi-case** (nhiều giá trị trong 1 case).

```go
switch day {
case "Sat", "Sun":
    fmt.Println("Weekend") // Chạy cho cả 2 trường hợp
case "Mon":
    fmt.Println("Start work")
default:
    fmt.Println("Working...")
}
```

---

## 📦 5. Data Structures

### Arrays & Slices

Đặc điểm,Array (Mảng),Slice (Lát cắt)
Kích thước,Cố định khi khai báo.,"Linh hoạt, có thể co giãn."
Kiểu dữ liệu,Kích thước là một phần của kiểu dữ liệu (ví dụ [3]int khác [5]int).,Không chứa kích thước trong kiểu dữ liệu (ví dụ []int).
Vùng nhớ,Lưu trữ giá trị trực tiếp.,Lưu tham chiếu đến một Array ngầm định (Underlying Array).
Cách truyền (Pass by),Truyền giá trị (copy toàn bộ mảng).,Truyền tham chiếu (chỉ copy cấu trúc Slice Header).
Tính phổ biến,Ít dùng trực tiếp trong logic nghiệp vụ.,Sử dụng thường xuyên trong hầu hết mọi trường hợp.

- **Array:** Kích thước cố định. `var a [3]int = [3]int{1, 2, 3}`.
   - Fixed length
   - Same type
   - Indexable
   - Contigous in Mem
- **Slice:** Linh hoạt hơn, là "view" của một array.

```go
// Tạo slice bằng make
s := make([]int, 5, 10) // len=5, cap=10

// Khởi tạo nhanh
fruits := []string{"apple", "orange"}

// Cắt slice từ array
arr := [5]int{1, 2, 3, 4, 5}
s1 := arr[1:3] // [2, 3]

// Copy slice
dst := make([]int, len(s1))
copy(dst, s1)
```

> [!IMPORTANT]
> **Lưu ý về Slice:**
>
> - Slice là một **reference type**. Khi bạn gán `s1 = s2`, cả hai cùng trỏ vào một vùng nhớ.
> - Nếu bạn thay đổi phần tử trong slice được cắt từ array, array gốc cũng bị thay đổi theo.
> - Dùng `cap()` để kiểm tra sức chứa tối đa trước khi Go phải cấp phát lại bộ nhớ mới.

### Maps (Key-Value)

```go
m := make(map[string]int)
m["gold"] = 100
delete(m, "gold")      // Xóa phần tử theo key
val, ok := m["gold"]   // ok = false nếu không tìm thấy
```

### Structs (Thay thế Class)

```go
type Person struct {
    name   string
    age    int
    job    string
    salary int
}

// Khởi tạo
var p1 Person
p1.name = "John" // Gán giá trị sau khi khai báo

p2 := Person{name: "Jane", age: 25} // Khởi tạo nhanh
```

### Methods (Hàm dành riêng cho Struct)

Methods không thuộc về Struct nhưng "gắn" vào nó thông qua **receiver**.

```go
func (p Person) Describe() {
    fmt.Printf("%s is %d years old\n", p.name, p.age)
}
```

> [!IMPORTANT]
> **Lưu ý về Struct & Methods:**
>
> - **Pointer Receiver:** Nếu muốn thay đổi giá trị của struct bên trong method, bạn phải dùng con trỏ: `func (p *Person) UpdateAge(newAge int)`.
> - **Value Receiver:** Nếu dùng `(p Person)`, Go sẽ copy toàn bộ struct vào hàm, mọi thay đổi bên trong sẽ không ảnh hưởng gốc.

---

## 🧬 6. Concurrency (Goroutines & Channels)

Đây là "vũ khí" mạnh nhất của Go.

- **Goroutine:** Dùng từ khóa `go` trước một hàm để chạy song song. Chỉ tốn khoảng 2KB RAM.
- **Channel:** Dùng để giao tiếp giữa các goroutines (Don't communicate by sharing memory; share memory by communicating).

```go
ch := make(chan string)
go func() { ch <- "Done" }()
msg := <-ch // Chờ nhận dữ liệu
```

---

## 7. Các từ khóa đặc biệt (`defer`, `panic`, `recover`)

### Defer

Dùng để trì hoãn việc thực thi một hàm cho đến khi hàm bao quanh nó kết thúc (thường dùng để đóng file, kết nối DB).

```go
func main() {
    defer fmt.Println("Thế giới!")
    fmt.Println("Chào")
}
// Kết quả: Chào Thế giới!
```

> [!TIP]
> **Lưu ý về Defer:**
>
> - Nếu có nhiều lệnh `defer`, chúng sẽ chạy theo thứ tự **LIFO** (Last In, First Out) - lệnh nào khai báo sau cùng sẽ chạy trước.

### Panic & Recover

- `panic`: Dùng để dừng chương trình ngay lập tức khi gặp lỗi nghiêm trọng.
- `recover`: Dùng trong hàm `defer` để "cứu" chương trình khỏi bị crash sau khi `panic`.

---

## ⚠️ 8. Xử lý lỗi (Error Handling)

Go không dùng `try/catch`. Lỗi là một giá trị trả về.

```go
func doWork() (int, error) {
    if fail { return 0, errors.New("failed") }
    return 100, nil
}

// Cách dùng phổ biến
val, err := doWork()
if err != nil {
    log.Fatal(err)
}
```

---

## 🔗 8. So sánh Go với Node.js (Cho Developer)

Để xem chi tiết bản so sánh tư duy lập trình và báo cáo hiệu năng thực tế giữa Go và Node.js/TypeScript (Bun), vui lòng xem tài liệu riêng biệt:

👉 **[Báo cáo So sánh Go vs Node.js](file:///Users/dungxbuif/workspace/go-study/GO_NODEJS.md)**

---

## ❓ Q&A Quan trọng

### 1. Số thực (Float) được lưu thế nào?

Theo tiêu chuẩn **IEEE 754**. Ví dụ `float64` có 1 bit dấu, 11 bit mũ và 52 bit phần lẻ. Tránh dùng float cho tiền tệ vì sai số nhị phân.

### 2. Interface trong Go có gì đặc biệt?

Interface trong Go là **Implicit**. Một Struct không cần khai báo `implements Interface`. Chỉ cần nó có đủ các phương thức mà Interface yêu cầu, nó sẽ tự động được coi là implement Interface đó (Duck Typing).

---

_Tài liệu này được tổng hợp để hỗ trợ quá trình học tập. Hãy thực hành viết code thường xuyên để nắm vững kiến thức!_

ASCII

- use 7 bits bieu dien 128 ky tu
