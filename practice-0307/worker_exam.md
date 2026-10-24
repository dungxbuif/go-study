```go
// worker nhận job từ channel jobs (chỉ đọc: <-chan) và gửi kết quả vào channel results (chỉ ghi: chan<-)
func worker(id int, jobs <-chan int, results chan<- int) {
	// Vòng lặp range trên channel sẽ chạy liên tục cho đến khi channel bị close
	for j := range jobs {
		fmt.Println("worker", id, "started  job", j)
		time.Sleep(time.Second) // Giả lập công việc tốn 1 giây để xử lý
		fmt.Println("worker", id, "finished job", j)
		results <- j * 2 // Gửi kết quả sau khi xử lý thành công vào channel results
	}
}

func main() {
	const numJobs = 5
	// Tạo các buffered channels với dung lượng 5 để không bị block khi gửi dữ liệu
	jobs := make(chan int, numJobs)
	results := make(chan int, numJobs)

	// 1. Khởi chạy 3 worker goroutines chạy song song (Concurrent Workers)
	for w := 1; w <= 3; w++ {
		go worker(w, jobs, results)
	}

	// 2. Đẩy 5 công việc vào channel jobs
	for j := 1; j <= numJobs; j++ {
		jobs <- j
	}
	// Đóng channel jobs để báo cho các worker biết là không còn job mới nào nữa.
	// Điều này giúp vòng lặp `for range jobs` ở các worker tự động thoát ra khi đã xử lý hết job.
	close(jobs)

	// 3. Cơ chế đồng bộ hóa (Synchronization) không dùng WaitGroup:
	// Vì ta biết chính xác có 5 jobs -> sẽ sinh ra đúng 5 kết quả trả về trong channel `results`.
	// Vòng lặp dưới đây sẽ lấy ra (đọc) 5 kết quả từ channel `results`.
	// Dòng code `<-results` sẽ block (treo/chờ) cho đến khi có dữ liệu được gửi từ worker.
	// Nhờ việc lặp đúng 5 lần, hàm main sẽ đợi cho đến khi nhận đủ cả 5 kết quả rồi mới kết thúc chương trình.
	for a := 1; a <= numJobs; a++ {
		<-results // Đọc kết quả và block chờ nếu kết quả chưa tới
	}
}
```

---

```go

```
