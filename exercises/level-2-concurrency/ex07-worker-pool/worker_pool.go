package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// ============================================================================
// 💡 PHÂN TÍCH CHUYÊN SÂU CƠ CHẾ GO WORKER POOL DƯỚI NẮP CAPO
// ============================================================================
//
// 1. Luồng chạy của Go Worker Pool hoạt động thế nào?
//    - Chúng ta có 2 kênh (channels):
//      - `jobs` (đầu vào): Kênh chứa danh sách công việc cần làm.
//      - `results` (đầu ra): Kênh chứa kết quả của các công việc đã làm xong.
//    - Chúng ta dùng **Buffered Channels** để tránh việc Goroutine gửi bị chặn ngay lập tức,
//      cho phép chuẩn bị sẵn hàng đợi công việc.
//
// 2. Cơ chế phân phối công việc tự động (Implicit Load Balancing):
//    - Chúng ta sinh ra `numWorkers` Goroutines độc lập cùng chạy hàm `Worker`.
//    - Cả `numWorkers` Goroutines này cùng chạy vòng lặp `for job := range jobs` trên **CÙNG MỘT channel**.
//    - **Dưới mui xe:** Go Runtime scheduler quản lý việc xếp hàng đọc từ channel cực kỳ thông minh.
//      Khi có công việc mới đẩy vào `jobs`, Go sẽ tự động chọn một Goroutine (Worker) đang ở trạng thái
//      rảnh (idle) để trao việc. Phép phân phối này hoàn toàn an toàn luồng (thread-safe),
//      không bao giờ lo tranh chấp tài nguyên (race condition) hay bị trùng lặp công việc!
//
// 3. sync.WaitGroup đóng vai trò gì?
//    - Để biết khi nào **tất cả** các Workers đã hoàn thành và dọn dẹp sạch sẽ tài nguyên.
//    - Mỗi khi tạo 1 Worker, ta gọi `wg.Add(1)`. Khi Worker chạy xong và thoát vòng lặp, ta gọi `wg.Done()`.
//    - Hàm chính sẽ đợi qua `wg.Wait()` trước khi đóng channel kết quả `results`.
//
// 4. Cơ chế Đóng Kênh Domino:
//    - Bước 1: Sau khi nạp hết việc, ta gọi `close(jobs)`.
//    - Bước 2: Các workers khi đọc hết việc cũ trong channel sẽ phát hiện channel đã đóng ➡️ tự động
//      thoát khỏi vòng lặp `for range jobs` ➡️ gọi `wg.Done()`.
//    - Bước 3: `wg.Wait()` ở main goroutine được mở khóa ➡️ gọi `close(results)` để báo kết thúc kết quả.
//    - Bước 4: Vòng lặp thu thập kết quả `for res := range results` kết thúc, in ra màn hình.
// ============================================================================

type Job struct {
	ID        int
	ImageName string
}

type Result struct {
	JobID    int
	Duration int
}

/*
Edited implement.go
Viewed implement.go:1-16

Dòng code `jobChan := make(chan Job, len(jobs))` đang làm nhiệm vụ **khởi tạo một Buffered Channel** (Channel có bộ đệm). 

Dưới đây là giải thích chi tiết về lý do tại sao phải dùng `make` và cơ chế hoạt động "dưới nắp capo" (under the hood) của nó trong Go runtime.

### 1. Tại sao phải dùng `make`? (So sánh với Node.js)
Trong Go, có 3 kiểu dữ liệu đặc biệt là **Channel, Map, và Slice**. Chúng là các kiểu tham chiếu (reference types).
- Nếu bạn chỉ khai báo `var jobChan chan Job`, Go sẽ cấp phát nó với giá trị zero-value là `nil`. Việc gửi hoặc nhận dữ liệu từ một `nil` channel sẽ làm goroutine bị **block vĩnh viễn (deadlock)**.
- Từ khóa `make` được sinh ra để vừa cấp phát bộ nhớ (allocate memory) vừa **khởi tạo** cấu trúc dữ liệu nền tảng bên dưới cho 3 kiểu này. (Khác với từ khóa `new` chỉ cấp phát bộ nhớ và trả về con trỏ, nhưng không khởi tạo).

### 2. Cơ chế "dưới nắp capo" (Under the hood)

Khi bạn gọi `make(chan Job, len(jobs))`, trình biên dịch của Go sẽ biên dịch nó thành một lệnh gọi hàm nội bộ trong Go runtime là `makechan(t *chantype, size int) *hchan`.

Hàm này sẽ tạo ra và trả về một con trỏ trỏ tới một struct tên là **`hchan`** (nằm trong thư mục `runtime/chan.go` của mã nguồn Go). `jobChan` bản chất chính là một con trỏ trỏ tới struct `hchan` này.

Cấu trúc của `hchan` chứa các thành phần cốt lõi sau để quản lý concurrency:

```go
type hchan struct {
    qcount   uint           // Số lượng phần tử đang có trong buffer
    dataqsiz uint           // Kích thước tối đa của buffer (ở đây là len(jobs))
    buf      unsafe.Pointer // Con trỏ trỏ tới mảng thực tế lưu trữ dữ liệu (Circular Queue)
    sendx    uint           // Vị trí (index) để ghi phần tử tiếp theo vào mảng
    recvx    uint           // Vị trí (index) để đọc phần tử tiếp theo từ mảng
    recvq    waitq          // Hàng đợi các Goroutines đang bị block chờ đọc (chờ nhận job)
    sendq    waitq          // Hàng đợi các Goroutines đang bị block chờ ghi (chờ gửi job)
    lock     mutex          // Khóa Mutex để đảm bảo an toàn (Thread-safe)
}
```

### 3. Phân tích cụ thể ở dòng lệnh của bạn: `make(chan Job, len(jobs))`

Vì bạn truyền vào tham số thứ hai là `len(jobs)` (giả sử = 50), đây là một **Buffered Channel**:
1. **Cấp phát mảng (Ring Buffer):** Go sẽ tạo ra một mảng vòng (Circular Queue) nằm trên Heap memory có khả năng chứa đúng 50 struct `Job`. `hchan.buf` sẽ trỏ tới mảng này, và `hchan.dataqsiz` sẽ bằng 50.
2. **Hành vi Non-blocking (Không nghẽn):** Khi hàm Main (hoặc goroutine đẩy job) gọi `jobChan <- job`, Go sẽ khóa `mutex` lại, copy giá trị `job` vào mảng tại vị trí `sendx`, tăng `qcount` lên 1, sau đó mở khóa `mutex`. Vì buffer có sức chứa 50, vòng lặp đẩy 50 job vào channel sẽ chạy **cực kỳ nhanh mà không bị khựng lại (block)** chờ các worker xử lý.
3. **Goroutine Scheduler:** Nếu bạn không cấp `len(jobs)` (tức là tạo Unbuffered Channel với `make(chan Job)`), mảng `buf` sẽ không được tạo ra. Khi đó, mỗi lần `jobChan <- job`, Goroutine chính gửi dữ liệu sẽ bị "tạm giam" (vào hàng đợi `sendq`) và nhường CPU lại cho đến khi có một Worker Goroutine tới nhận `<-jobChan`. Bằng việc cấp buffer bằng đúng số lượng jobs, bạn đảm bảo luồng chính phân phối xong việc là có thể đi làm việc khác ngay lập tức!

**Tóm lại:** Dòng `make` đó thiết lập một cấu trúc Thread-safe (có Mutex bảo vệ) dùng mảng vòng (Ring buffer) trên RAM, giúp các Goroutines gửi và nhận `Job` giao tiếp với nhau mà không lo Race Condition hay làm sập (Crash) chương trình.

*/

// Worker là hàm xử lý công việc chạy ngầm
// - id: Mã định danh của worker
// - jobs: Kênh chỉ đọc (receive-only channel: `<-chan Job`)
// - results: Kênh chỉ ghi (send-only channel: `chan<- Result`)
// - wg: Con trỏ WaitGroup để thông báo hoàn thành
func Worker(id int, jobs <-chan Job, results chan<- Result, wg *sync.WaitGroup) {
	// Đảm bảo thông báo WaitGroup giảm đi 1 khi hàm này kết thúc
	defer wg.Done()
	
	// Vòng lặp liên tục kéo việc từ channel.
	// Vòng lặp này sẽ chặn (Block) worker nếu channel chưa có việc mới,
	// và tự động thoát (break) khi channel `jobs` bị đóng.
	for job := range jobs {
		// Giả lập xử lý ảnh mất 100-500ms
		duration := rand.Intn(400) + 100
		time.Sleep(time.Duration(duration) * time.Millisecond)

		// Đẩy kết quả vào channel kết quả
		results <- Result{
			JobID:    job.ID,
			Duration: duration,
		}
		fmt.Printf("Worker #%d processed image_%d.jpg in %dms\n", id, job.ID, duration)
	}
}

// RunWorkerPool khởi tạo và quản lý toàn bộ vòng đời của Worker Pool
func RunWorkerPool(numWorkers int, jobsList []Job) []Result {
	// Khởi tạo các buffered channels chứa đủ dung lượng công việc
	jobs := make(chan Job, len(jobsList))
	results := make(chan Result, len(jobsList))

	var wg sync.WaitGroup

	// 1. Sinh ra N Workers chạy đồng thời
	for w := 1; w <= numWorkers; w++ {
		wg.Add(1)
		go Worker(w, jobs, results, &wg)
	}

	// 2. Nạp toàn bộ công việc vào channel đầu vào
	for _, j := range jobsList {
		jobs <- j
	}
	// Đóng channel công việc ngay sau khi nạp xong.
	// Lưu ý: Đóng channel chỉ báo hiệu không nhận việc mới, các worker vẫn đọc nốt việc cũ bình thường.
	close(jobs)

	// 3. Đợi tất cả các worker xử lý xong toàn bộ hàng đợi việc
	wg.Wait()
	
	// 4. Đóng channel kết quả sau khi không còn worker nào ghi kết quả nữa
	close(results)

	// 5. Thu thập kết quả từ channel đầu ra
	var output []Result
	for res := range results {
		output = append(output, res)
	}

	return output
}

func main() {
	rand.Seed(time.Now().UnixNano())

	// Tạo danh sách 20 ảnh cần xử lý
	jobs := make([]Job, 20)
	for i := 0; i < 20; i++ {
		jobs[i] = Job{ID: i + 1, ImageName: fmt.Sprintf("image_%d.jpg", i+1)}
	}

	start := time.Now()
	// Khởi chạy Worker Pool với 5 Workers song song để xử lý 20 công việc
	results := RunWorkerPool(5, jobs)
	duration := time.Since(start)

	fmt.Printf("\nAll %d images processed in %v\n", len(results), duration)
}
