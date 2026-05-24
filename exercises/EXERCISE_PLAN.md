# 🎯 Lộ Trình Bài Tập Go + Gin: Từ Node.js Dev → Senior Go Dev

> Thiết kế cho developer có nền tảng Node.js/Express/NestJS.
> Mỗi bài tập nhỏ, tập trung vào **1-2 kỹ năng cốt lõi**, xây dựng lên dần theo cấp độ.
> Các bài tập từ Level 3 trở đi tập trung vào các kỹ thuật nâng cao thường gặp trong production Go systems.

---

## 📐 Cấu Trúc Thư Mục

```
go-study/exercises/
├── EXERCISE_PLAN.md             
├── level-1-foundations/
│   ├── ex01-struct-methods/
│   ├── ex02-interfaces/
│   ├── ex03-error-handling/
│   └── ex04-pointers-generics/
├── level-2-concurrency/
│   ├── ex05-goroutine-waitgroup/
│   ├── ex06-channels-pipeline/
│   ├── ex07-worker-pool/
│   └── ex08-context-timeout/
├── level-3-gin-basics/
│   ├── ex09-hello-gin/
│   ├── ex10-crud-inmemory/
│   ├── ex11-middleware-chain/
│   └── ex12-validation-binding/
├── level-4-gin-production/
│   ├── ex13-clean-architecture/
│   ├── ex14-auth-jwt-zk/
│   ├── ex15-database-gorm/
│   └── ex16-testing-mocking/
└── level-5-senior/
    ├── ex17-semaphore-binary-retry/
    ├── ex18-redis-lua-concurrency/
    ├── ex19-websocket-room-hub/
    └── ex20-batch-pipeline-kafka/
```

---

## 📊 Checklist Tiến Độ

| # | Bài tập | Kỹ năng | Status |
|:--|:--------|:--------|:-------|
| 01 | Struct & Methods | Pointer/Value receiver | ☐ |
| 02 | Interfaces | Implicit satisfaction, DI | ☐ |
| 03 | Error Handling | Custom errors, wrapping | ☐ |
| 04 | Pointers & Generics | `*T`, `[T any]` | ☐ |
| 05 | Goroutine + WaitGroup | `sync.WaitGroup`, `Mutex` | ☐ |
| 06 | Channels Pipeline | fan-out, fan-in, `close()` | ☐ |
| 07 | Worker Pool | Buffered channel job queue | ☐ |
| 08 | Context Timeout | `WithTimeout`, `WithCancel` | ☐ |
| 09 | Hello Gin | Routes, groups, params | ☐ |
| 10 | CRUD In-Memory | JSON binding, `sync.RWMutex` | ☐ |
| 11 | Middleware Chain | Logger, Auth, Onion model | ☐ |
| 12 | Validation & Binding | Struct tags, custom validator | ☐ |
| 13 | Clean Architecture | Interface layers, manual DI | ☐ |
| 14 | Auth JWT + ZK concept | JWT middleware, c.Set/Get | ☐ |
| 15 | Database + SQL | UPSERT, Preload, pagination | ☐ |
| 16 | Testing & Mocking | Table-driven, `httptest` | ☐ |
| 17 | Semaphore + Binary Retry | Channel semaphore, đệ quy | ☐ |
| 18 | Redis Lua Concurrency | Atomic claim, LPOP/HSET | ☐ |
| 19 | WebSocket Room Hub | Gorilla WS, Hub pattern | ☐ |
| 20 | Batch Pipeline + Event | Batcher, CTE SQL, Kafka | ☐ |

---

## 🟢 Level 1: Go Foundations (Tư duy khác Node.js)

> **Mục tiêu**: Bỏ thói quen OOP class-based của Node.js, làm quen tư duy Go.
> **Thời gian ước tính**: 1-2 giờ/bài

---

### Ex01: Struct & Methods — "User Service"

**🧠 Node.js → Go**: `class UserService` → `type UserService struct` + methods with receivers

**Yêu cầu:**
1. Tạo struct `User` với các trường: `ID`, `Name`, `Email`, `CreatedAt`
2. Tạo struct `UserService` chứa field `users []User` (in-memory store)
3. Implement các methods:
   - `(s *UserService) Create(name, email string) User` — pointer receiver vì sửa state
   - `(s *UserService) FindByID(id int) (User, error)` — trả error nếu không tìm thấy
   - `(s *UserService) FindByEmail(email string) (User, error)`
   - `(s UserService) Count() int` — **value receiver** vì chỉ đọc, không sửa state
4. Viết `main()` demo tạo 3 users, tìm theo ID, tìm theo email

**🎯 Kỹ năng rèn:**
- Pointer receiver (`*T`) vs Value receiver (`T`) — khi nào dùng `*` khi nào không
- Struct initialization, zero values
- Multiple return values — Go trả `(result, error)`, không `throw`
- `time.Now()` thay cho `new Date()`

**💡 So sánh Express:**
```javascript
// Node.js: class-based
class UserService {
  constructor() { this.users = []; }
  create(name, email) { /* ... */ }
  findById(id) { /* throw nếu không tìm thấy */ }
}
```
```go
// Go: Struct + Methods
type UserService struct { users []User }
func (s *UserService) Create(name, email string) User { /* ... */ }
func (s *UserService) FindByID(id int) (User, error) { /* trả error, không throw */ }
```

---

### Ex02: Interfaces — "Payment Gateway"

**🧠 Node.js → Go**: `implements Interface` (explicit) → Duck Typing (implicit)

**Yêu cầu:**
1. Định nghĩa interface `PaymentGateway`:
   ```go
   type PaymentGateway interface {
       Charge(amount float64, currency string) (txID string, err error)
       Refund(txID string) error
   }
   ```
2. Implement 2 struct: `StripeGateway` và `MomoGateway` — mỗi cái logic khác nhau
3. Viết function `ProcessOrder(gw PaymentGateway, amount float64)` — nhận interface, không biết implementation cụ thể
4. Trong `main()`, truyền lần lượt Stripe và Momo vào `ProcessOrder`

**🎯 Kỹ năng rèn:**
- Implicit interface satisfaction (không cần `implements`)
- Dependency Injection kiểu Go (truyền interface vào function)
- Interface-based architecture — tách biệt implementations khỏi business logic

---

### Ex03: Error Handling — "Config Parser"

**🧠 Node.js → Go**: `try/catch/throw` → `if err != nil` + error wrapping

**Yêu cầu:**
1. Tạo custom error types:
   ```go
   type ValidationError struct { Field, Message string }
   func (e *ValidationError) Error() string { ... }
   
   type NotFoundError struct { Resource string; ID int }
   func (e *NotFoundError) Error() string { ... }
   ```
2. Implement `ParseConfig(filename string) (Config, error)`:
   - File không tồn tại → `NotFoundError`
   - Nội dung sai format → `ValidationError`
   - Wrap error: `fmt.Errorf("parsing config: %w", err)`
3. Ở caller, dùng `errors.Is()` và `errors.As()` để xử lý từng loại lỗi khác nhau

**🎯 Kỹ năng rèn:**
- Custom error types (implement `Error() string`)
- Error wrapping với `%w` — giữ nguyên error chain
- `errors.Is()` vs `errors.As()` — tương đương `instanceof` trong JS
- Sentinel errors: `var ErrNotFound = errors.New("not found")`

---

### Ex04: Pointers & Generics — "LinkedList"

**🧠 Node.js → Go**: JS tự quản reference → Go phải hiểu `*` và `&`

**Yêu cầu:**
1. Implement `LinkedList[T any]` (dùng Generics):
   ```go
   type Node[T any] struct {
       Value T
       Next  *Node[T]
   }
   type LinkedList[T any] struct {
       Head *Node[T]
       Size int
   }
   ```
2. Methods: `Push(val T)`, `Pop() (T, error)`, `Find(predicate func(T) bool) *Node[T]`, `Print()`
3. Trong `main()`, tạo `LinkedList[int]` và `LinkedList[string]`

**🎯 Kỹ năng rèn:**
- Pointers: `&Node[T]{...}` tạo node, `current = current.Next` traverse
- Generics `[T any]` — tương đương `<T>` trong TypeScript
- Nil checking — pointer có thể nil, phải check trước khi dereference
- `func(T) bool` — function as parameter, tương đương callback trong JS

---

## 🔵 Level 2: Concurrency (Vũ khí tối thượng của Go)

> **Mục tiêu**: Thành thạo Goroutine, Channel, WaitGroup, Context.
> Node.js **không có tương đương trực tiếp** cho phần này.
> **Thời gian ước tính**: 2-3 giờ/bài

---

### Ex05: Goroutine + WaitGroup — "Parallel URL Checker"

**🧠 Node.js → Go**: `Promise.all([fetch(...)])` → `go func() {}` + `sync.WaitGroup`

**Yêu cầu:**
1. Input: danh sách 10 URLs hardcoded
2. Mỗi URL = 1 goroutine, gửi HTTP GET kiểm tra status
3. Dùng `sync.WaitGroup` đợi tất cả hoàn thành
4. Dùng `sync.Mutex` bảo vệ shared slice kết quả (tránh race condition)
5. In bảng kết quả: `URL | Status Code | Response Time`
6. **Bonus**: Chạy `go test -race ./...` phát hiện race condition

**🎯 Kỹ năng rèn:**
- `go func(){}()` — anonymous goroutine
- `sync.WaitGroup`: `Add()`, `Done()`, `Wait()`
- `sync.Mutex`: `Lock()`, `Unlock()` — bảo vệ shared data
- So sánh: Express chạy single-thread → Go chạy multi-goroutine thực sự

---

### Ex06: Channels + Pipeline — "Data Processing Pipeline"

**🧠 Node.js → Go**: `stream.pipe()` / RxJS → Go channels pipeline

**Yêu cầu:**
1. Xây pipeline 3 stage chạy **đồng thời**:
   - **Stage 1 - Generator**: Sinh 100 số nguyên, đẩy vào `chan int`
   - **Stage 2 - Filter**: Chỉ giữ số chẵn, đẩy vào channel khác
   - **Stage 3 - Squarer**: Bình phương giá trị, đẩy vào channel kết quả
2. Main goroutine đọc kết quả cuối cùng
3. Đóng channel đúng cách: `close(ch)` + `for val := range ch`

**🎯 Kỹ năng rèn:**
- Unbuffered vs Buffered channel — blocking behavior khác nhau
- `close()` + `for range ch` pattern
- Channel direction: `chan<- int` (send-only), `<-chan int` (receive-only)
- Pipeline pattern — nền tảng cho các hệ thống data processing thực tế

---

### Ex07: Worker Pool — "Image Resizer Simulator"

**🧠 Node.js → Go**: `p-limit` / Bull Queue → native Worker Pool

**Yêu cầu:**
1. Giả lập 50 "ảnh" cần xử lý (mỗi ảnh mất random 100-500ms)
2. Tạo Worker Pool với **5 workers** (5 goroutines)
3. 1 `jobs` channel phân phối, 1 `results` channel thu kết quả
4. Log: `Worker #3 processed image_42.jpg in 234ms`
5. Đo tổng thời gian — so sánh tuần tự vs parallel

**🎯 Kỹ năng rèn:**
- Worker Pool pattern — pattern production Go phổ biến nhất
- Buffered channels làm job queue
- Tái sử dụng goroutine (5 worker, không tạo 50 goroutine)
- Chuẩn bị cho Ex17 (Semaphore nâng cao)

---

### Ex08: Context — "Timeout & Cancellation"

**🧠 Node.js → Go**: `AbortController` → `context.WithTimeout` / `context.WithCancel`

**Yêu cầu:**
1. `FetchData(ctx context.Context, url string) (string, error)` — gửi HTTP với context
2. `FetchMultiple(urls []string, timeout time.Duration)`:
   - Tạo `context.WithTimeout` tổng 3 giây
   - Fetch tất cả URLs song song
   - Nếu quá 3s → cancel tất cả requests còn lại
3. Demo: mix URLs nhanh + chậm (`httpbin.org/delay/5`)
4. **Bonus**: Viết function chạy nền dùng `context.WithCancel`, bắt signal `SIGINT` để dừng

**🎯 Kỹ năng rèn:**
- `context.WithTimeout()`, `context.WithCancel()`
- `defer cancel()` — **bắt buộc** tránh goroutine leak
- `select` + `ctx.Done()` pattern
- Context propagation xuống mọi tầng

---

## 🟡 Level 3: Gin Framework Basics

> **Mục tiêu**: Làm quen Gin framework, mapping từ Express/NestJS.
> **Thời gian ước tính**: 2-3 giờ/bài

---

### Ex09: Hello Gin — "First Gin Server"

**🧠 Node.js → Go**: `express()` → `gin.Default()`, `app.get()` → `r.GET()`

**Yêu cầu:**
1. Gin server với routes:
   - `GET /` → `{ "message": "Hello, Go!" }`
   - `GET /health` → `{ "status": "ok", "uptime": "..." }`
   - `GET /users/:id` → lấy param `id`
   - `GET /search?q=xxx&page=1` → lấy query params
2. Route group: `/api/v1` prefix
3. Custom port từ `os.Getenv("PORT")`

**💡 So sánh Express:**
```javascript
app.get('/users/:id', (req, res) => {
  res.json({ user_id: req.params.id });
});
```
```go
r.GET("/users/:id", func(c *gin.Context) {
    c.JSON(200, gin.H{"user_id": c.Param("id")})
})
```

**🎯 Kỹ năng rèn:**
- `gin.Default()` vs `gin.New()` (Default có Logger + Recovery)
- `c.Param()`, `c.Query()`, `c.DefaultQuery()`
- `c.JSON()`, `c.String()`
- Route group `r.Group("/api/v1")`

---

### Ex10: CRUD In-Memory — "Todo API"

**🧠 Node.js → Go**: Dynamic JS object → Typed Go struct + JSON tags

**Yêu cầu:**
1. Struct `Todo` với JSON tags:
   ```go
   type Todo struct {
       ID        int       `json:"id"`
       Title     string    `json:"title"`
       Completed bool      `json:"completed"`
       CreatedAt time.Time `json:"created_at"`
   }
   ```
2. CRUD endpoints: `POST /todos`, `GET /todos`, `GET /todos/:id`, `PUT /todos/:id`, `DELETE /todos/:id`
3. Lưu trữ bằng `map[int]Todo` + `sync.RWMutex` (thread-safe vì Gin xử lý concurrent)
4. Response format chuẩn: `{ "success": true, "data": {...} }`
5. Filter: `GET /todos?completed=true`

**🎯 Kỹ năng rèn:**
- `c.ShouldBindJSON(&input)` — tương đương `req.body`
- JSON struct tags: `json:"field"`, `json:"-"`, `json:",omitempty"`
- `sync.RWMutex`: `RLock/RUnlock` cho đọc, `Lock/Unlock` cho ghi
- `map[K]struct{}` trick — Go set pattern (tiết kiệm RAM)

---

### Ex11: Middleware Chain — "Logger, Auth, CORS"

**🧠 Node.js → Go**: `app.use(middleware)` → `r.Use(middleware)` + Onion model

**Yêu cầu:**
1. Viết 3 custom middlewares:
   - **RequestLogger**: Log method, path, status code, response time (tương tự `morgan`)
   - **AuthMiddleware**: Kiểm tra header `X-API-Key`, nếu sai → `c.AbortWithStatusJSON(401, ...)`
   - **CORSMiddleware**: Set CORS headers
2. Apply middleware theo scope:
   - Global: Logger + CORS (tất cả routes)
   - Group level: Auth chỉ cho `/api/v1` (không cho `/health`)
3. Custom Recovery middleware: catch panic → trả JSON error

**🎯 Kỹ năng rèn:**
- Gin middleware signature: `func(c *gin.Context)` + `c.Next()`
- `c.Abort()` vs `c.Next()` — dừng chain vs tiếp tục
- `c.Set()` / `c.Get()` — truyền data qua middleware (tương đương `req.locals`)
- **Onion model**: code trước `c.Next()` = PRE, code sau `c.Next()` = POST

**💡 So sánh Express:**
```javascript
// Express: next() để tiếp tục
app.use((req, res, next) => { /* ... */ next(); });
```
```go
// Gin: c.Next() + Onion model
func Logger() gin.HandlerFunc {
    return func(c *gin.Context) {
        start := time.Now()   // PRE-handler
        c.Next()               // → xử lý handler → quay lại đây
        latency := time.Since(start) // POST-handler
        log.Printf("%s %s %d %v", c.Request.Method, c.Request.URL, c.Writer.Status(), latency)
    }
}
```

---

### Ex12: Validation & Binding — "User Registration"

**🧠 Node.js → Go**: `class-validator` / `zod` → Struct tags `binding:"required"`

**Yêu cầu:**
1. DTO struct với validation:
   ```go
   type CreateUserDTO struct {
       Username string `json:"username" binding:"required,min=3,max=20,alphanum"`
       Email    string `json:"email" binding:"required,email"`
       Password string `json:"password" binding:"required,min=8"`
       Age      int    `json:"age" binding:"required,gte=18,lte=120"`
   }
   ```
2. `POST /register`: bind + validate, trả error message tiếng Việt
3. `UserResponse` struct riêng (không chứa password — tương đương `@Exclude()` NestJS)
4. Custom validator: `ValidateVietnamPhone` (bắt đầu `+84` hoặc `0`)

**🎯 Kỹ năng rèn:**
- `c.ShouldBindJSON()` + `binding` tags
- DTO pattern: Input DTO ≠ Response DTO (tách biệt dữ liệu vào/ra)
- Custom validator: `v.RegisterValidation(...)`

---

## 🟠 Level 4: Gin Production Patterns

> **Mục tiêu**: Clean architecture, auth, database, testing — đủ trình độ xây dựng project thực tế.
> **Thời gian ước tính**: 3-4 giờ/bài

---

### Ex13: Clean Architecture — "Refactor Todo API"

**🧠 Node.js → Go**: NestJS Module/Service/Controller → Go interface-based layers

**Yêu cầu:**
Refactor bài Ex10:
```
ex13-clean-architecture/
├── cmd/main.go                    # Bootstrap: wiring dependencies thủ công
├── internal/
│   ├── domain/todo.go             # Entity + TodoRepository interface
│   ├── usecase/todo_usecase.go    # Business logic, nhận interface
│   ├── repository/todo_memory.go  # In-memory implementation
│   └── delivery/http/
│       ├── handler.go             # Gin handlers
│       └── router.go             # Route setup
└── go.mod
```

1. **Domain**: `Todo` entity + `TodoRepository` interface
2. **Repository**: `InMemoryTodoRepository` satisfy interface
3. **Usecase**: `TodoUsecase` nhận `TodoRepository` qua constructor
4. **Delivery**: Gin handlers gọi usecase
5. **main.go**: Wire dependencies thủ công (DI kiểu Go)

**🎯 Kỹ năng rèn:**
- Interface-based architecture
- Constructor injection (không magic, không decorator)
- `internal/` package — code private
- Tương đương NestJS modules nhưng explicit

---

### Ex14: Auth JWT + ZK Concept — "Login & Protected Routes"

**🧠 Node.js → Go**: `passport.js` / `@UseGuards(AuthGuard)` → Gin middleware + JWT

**Yêu cầu:**
1. Endpoints:
   - `POST /auth/register` — hash password bằng `bcrypt`
   - `POST /auth/login` — verify password, trả JWT token
   - `GET /auth/profile` — **protected**, cần JWT header `Authorization: Bearer <token>`
2. JWT middleware:
   - Parse token → verify signature → extract claims
   - `c.Set("userID", claims.UserID)` — truyền user context
   - Invalid/expired → `c.AbortWithStatusJSON(401, ...)`
3. Role-based: middleware `RoleRequired("admin")`
4. **Bonus ZK concept**: Viết một middleware giả lập ZK-Auth:
   - Client gửi `{ "user_id": "...", "proof": "...", "public_key": "..." }`
   - Middleware đọc body, verify (giả lập), `c.Set("user_id")`, rồi **restore body** cho handler dùng tiếp (dùng `io.NopCloser`)

**🎯 Kỹ năng rèn:**
- JWT: `github.com/golang-jwt/jwt/v5`
- bcrypt: `golang.org/x/crypto/bcrypt`
- `c.Set()` / `c.MustGet()` — truyền user context
- Body reading + restoring trong middleware (kỹ thuật từ `zk_auth.go`)

---

### Ex15: Database + SQL — "Blog API"

**🧠 Node.js → Go**: `TypeORM` / `Prisma` → `GORM` + raw SQL khi cần

**Yêu cầu:**
1. Models:
   ```go
   type User struct {
       gorm.Model
       Username string `gorm:"uniqueIndex;size:50"`
       Posts    []Post `gorm:"foreignKey:AuthorID"`
   }
   type Post struct {
       gorm.Model
       Title    string `gorm:"size:200"`
       Content  string `gorm:"type:text"`
       AuthorID uint
       Author   User
   }
   ```
2. CRUD cho Posts: list (pagination `?page=1&limit=10`), get (`Preload("Author")`), create, update (chỉ tác giả), soft delete
3. Database: **SQLite** (đơn giản, dùng `gorm.io/driver/sqlite`)
4. **Nâng cao — Viết raw SQL UPSERT**:
   ```go
   // UPSERT pattern: chèn hoặc cập nhật nếu đã tồn tại
   db.Exec(`
       INSERT INTO posts (title, content, author_id)
       VALUES (?, ?, ?)
       ON CONFLICT (id) DO UPDATE SET
           title = EXCLUDED.title,
           updated_at = CURRENT_TIMESTAMP
   `, title, content, authorID)
   ```
5. Connection pool: `SetMaxOpenConns`, `SetMaxIdleConns`

**🎯 Kỹ năng rèn:**
- GORM basics: associations, preloading, pagination, soft delete
- Raw SQL: UPSERT pattern (`ON CONFLICT DO UPDATE`)
- Connection pool management
- Hiểu khi nào dùng ORM vs `database/sql` thuần (write-heavy vs read-heavy)

---

### Ex16: Testing & Mocking — "Test Todo Usecase"

**🧠 Node.js → Go**: `jest.mock()` → Interface-based mocking (không magic)

**Yêu cầu:**
1. **Unit Tests** cho `TodoUsecase`:
   - Tạo `MockTodoRepository` implement `TodoRepository` interface
   - Test: Create thành công, FindByID not found, Update đã completed
2. **HTTP Integration Tests**:
   - `httptest.NewRecorder()` + `gin.CreateTestContext()`
   - `POST /todos` body hợp lệ → 201
   - `POST /todos` thiếu title → 400
   - `GET /todos/999` → 404
3. **Table-Driven Tests** (chuẩn Go):
   ```go
   tests := []struct {
       name     string
       input    CreateTodoDTO
       wantCode int
       wantErr  bool
   }{
       {"valid todo", CreateTodoDTO{Title: "Buy milk"}, 201, false},
       {"empty title", CreateTodoDTO{Title: ""}, 400, true},
   }
   for _, tt := range tests {
       t.Run(tt.name, func(t *testing.T) { /* ... */ })
   }
   ```
4. Coverage: `go test -cover -race ./...`

**🎯 Kỹ năng rèn:**
- Table-driven tests — phong cách chuẩn Go
- Interface-based mocking — không cần `jest.mock()` magic
- `httptest` package — test handler không cần start server
- `go test -race` — phát hiện race condition

---

## 🔴 Level 5: Senior Production Techniques

> **Mục tiêu**: Làm chủ các kỹ thuật nâng cao thường dùng trong các hệ thống Go quy mô lớn.
> Hoàn thành level này = bạn tự tin mổ xẻ code và kiến trúc khi phỏng vấn vị trí Senior.
> **Thời gian ước tính**: 4-6 giờ/bài

---

### Ex17: Channel Semaphore + Binary Split Retry

**Mục tiêu phỏng vấn**: Trả lời câu hỏi _"Làm sao bạn giới hạn số lượng goroutines xử lý đồng thời để tránh làm quá tải downstream services?"_

**Yêu cầu:**
1. **Channel-Based Semaphore** — giới hạn goroutine concurrent:
   ```go
   sem := make(chan struct{}, maxConcurrent) // ví dụ 3
   
   for _, chunk := range chunks {
       go func(c []int) {
           sem <- struct{}{}          // ACQUIRE: block nếu đầy 3
           defer func() { <-sem }()   // RELEASE: nhả slot
           processChunk(c)
       }(chunk)
   }
   ```
2. **Binary Split Retry** — khi 1 chunk lỗi, chia đôi thử lại:
   ```go
   func processChunkWithRetry(chunk []int, sem chan struct{}) {
       results, failedItems := processChunk(chunk)
       if len(failedItems) == 0 { return } // OK
       if len(failedItems) == 1 { 
           markAsFailed(failedItems[0])     // Cô lập block lỗi
           return 
       }
       // Chia đôi và thử lại song song
       mid := len(failedItems) / 2
       go processChunkWithRetry(failedItems[:mid], sem)
       go processChunkWithRetry(failedItems[mid:], sem)
   }
   ```
3. **Kịch bản test**: 100 items chia thành 10 chunks, 3 concurrent, item #37 và #82 luôn fail
4. Log chi tiết: chunk nào chạy, chia đôi ra sao, kết quả cuối cùng
5. Đo tổng thời gian xử lý

**🎯 Kỹ năng rèn:**
- **Channel Semaphore** — kỹ thuật senior Go #1 (zero-overhead, Go Scheduler park goroutines)
- **Binary Split Retry** — thuật toán đệ quy chia đôi cô lập lỗi
- Rate limiting / Backpressure (chống sốc hệ thống)
- DB Connection Protection (quản lý kết nối DB hiệu quả)

**🗣️ Câu trả lời phỏng vấn mẫu:**
> _"Tôi sử dụng Channel-Based Semaphore để giới hạn cứng ConcurrentRequests. Khi buffer channel đầy, goroutine sẽ ở trạng thái Parked bởi Go Scheduler — hoàn toàn zero CPU overhead. Khi một chunk bị lỗi, thay vì retry toàn bộ, tôi áp dụng Binary Split Retry để cô lập chính xác block nào lỗi, giảm thiểu tải lên downstream service/database."_

---

### Ex18: Redis Lua Script — "Lucky Money Concurrency"

**Mục tiêu phỏng vấn**: Trả lời câu hỏi _"Tại sao không dùng SELECT FOR UPDATE để xử lý giao dịch concurrency cao?"_

**Yêu cầu:**
1. **Redis Lua Script** cho Atomic Claim:
   ```lua
   -- AttemptClaim: chạy đơn luồng trong Redis, zero race condition
   local poolKey = KEYS[1]       -- List các phần tiền lì xì
   local reservedKey = KEYS[2]   -- Hash: userID → amount đã nhận
   local userID = ARGV[1]
   
   -- 1. Kiểm tra user đã giật chưa
   local reserved = redis.call('HGET', reservedKey, userID)
   if reserved then return {reserved, 'ALREADY_CLAIMED'} end
   
   -- 2. LPOP lấy 1 phần tiền (Atomic)
   local amount = redis.call('LPOP', poolKey)
   if not amount then return {0, 'SOLD_OUT'} end
   
   -- 3. Đánh dấu đã giữ chỗ
   redis.call('HSET', reservedKey, userID, amount)
   return {amount, 'OK'}
   ```
2. **Gin API**:
   - `POST /red-envelope/create` — tạo bao lì xì, chia tiền random push vào Redis List
   - `POST /red-envelope/:id/claim` — giật lì xì (gọi Lua Script)
   - `GET /red-envelope/:id/status` — xem ai đã giật, còn bao nhiêu
3. **Concurrency test**: Dùng `sync.WaitGroup` tạo 100 goroutines cùng giật 1 bao chỉ có 10 phần
4. Verify: đúng 10 người nhận, 90 bị từ chối, không ai nhận 2 lần

**🎯 Kỹ năng rèn:**
- Redis Lua Script — atomic operations trên RAM
- `LPOP` — atomic pop, 100 người cùng gọi chỉ 10 nhận
- 2-layer validation: Redis (speed) → PostgreSQL (durability)
- Concurrency testing với goroutines

**🗣️ Câu trả lời phỏng vấn mẫu:**
> _"SELECT FOR UPDATE sẽ khóa dòng dữ liệu trong DB, gây Lock Contention nghiêm trọng với burst traffic. Tôi đưa bài toán về Redis Lua Script — chạy đơn luồng nguyên tử trên RAM, nhanh hơn hàng trăm lần. LPOP đảm bảo chỉ có đúng N người nhận quà. Sau khi Redis 'giữ chỗ', Async Worker mới ghi bền vững vào Postgres."_

---

### Ex19: WebSocket Room + Hub Pattern — "Real-time Chat"

**Mục tiêu phỏng vấn**: Trả lời câu hỏi _"Điều gì xảy ra nếu không dùng Mutex trong WSService?"_ và _"Tại sao dùng Ping/Pong?"_

**Yêu cầu:**
1. **WSService (Hub)** — quản lý connections thread-safe:
   ```go
   type WSService struct {
       mu          sync.RWMutex
       connections map[string][]*websocket.Conn            // address → list conns
       rooms       map[string]map[*websocket.Conn]struct{} // room → set of conns
   }
   ```
2. **Methods**:
   - `AddConnection(address string, conn *websocket.Conn)`
   - `RemoveConnection(address string, conn *websocket.Conn)` — xóa khỏi map + tất cả rooms
   - `SendToAddress(address string, msg interface{})` — targeted notification
   - `BroadcastToRoom(room string, msg interface{})` — room broadcast
   - `AddToRoom(room string, conn *websocket.Conn)` / `RemoveFromRoom`
3. **WebSocket Handler**:
   - `GET /ws?token=<jwt>` — upgrade connection
   - Goroutine **Ping ticker** phát hiện zombie:
     ```go
     go func() {
         ticker := time.NewTicker(30 * time.Second)
         for range ticker.C {
             if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
                 wsService.RemoveConnection(addr, conn)
                 return
             }
         }
     }()
     ```
   - `conn.SetPongHandler()` + `conn.SetReadDeadline()` — gia hạn khi nhận Pong
4. **Offline Inbox Pattern**:
   - Nếu address offline → lưu event vào DB
   - Khi connect lại → flush tất cả pending events qua WebSocket → xóa khỏi DB
5. **HTTP Bridge** (cho các service khác bắn event):
   - `POST /api/event` (kèm `X-API-Key`) → tìm connection → send hoặc save offline

**🎯 Kỹ năng rèn:**
- `sync.RWMutex`: `RLock` cho broadcast (nhiều goroutine đọc), `Lock` cho add/remove
- `map[*websocket.Conn]struct{}` — Go Set pattern, tiết kiệm RAM
- Goroutine per connection (2: read + ping)
- Heartbeat Ping/Pong — phát hiện TCP dead connections
- `defer RemoveConnection()` — cleanup khi disconnect
- Offline Inbox Pattern — không mất thông báo khi user offline

**🗣️ Câu trả lời phỏng vấn mẫu:**
> _"Go Map không thread-safe. Nếu AddConnection và RemoveConnection chạy đồng thời sẽ panic: concurrent map write. Tôi dùng sync.RWMutex — cho phép hàng nghìn goroutine đọc đồng thời (RLock) khi broadcast, nhưng khóa độc quyền (Lock) khi thêm/xóa connection."_
>
> _"Ping/Pong cần thiết vì TCP socket có thể bị đứt ngầm ở tầng mạng mà Read không biết. Server chủ động gửi Ping, nếu không nhận Pong trong PongWait thì coi là zombie — dọn dẹp RAM ngay."_

---

### Ex20: Batch Pipeline + Event Streaming

**Mục tiêu phỏng vấn**: Trả lời câu hỏi _"Hệ thống đảm bảo dữ liệu ghi vào DB luôn đúng đắn thế nào?"_

**Yêu cầu:**
1. **3-Layer Storage** (Bộ nhớ đệm 3 tầng):
   - **Staging** (in-memory): nhận data thô nhanh nhất có thể
   - **Main DB** (SQLite): dữ liệu chuẩn hóa, có index
   - **Orchestrator** (failure tracking): lưu vết item lỗi
2. **Batcher** pattern:
   ```go
   type Batcher struct {
       queue     chan string          // Buffered channel 10000
       batchSize int                  // 50
       timeout   time.Duration       // 2 seconds
   }
   
   func (b *Batcher) Run(ctx context.Context) {
       batch := make([]string, 0, b.batchSize)
       timer := time.NewTimer(b.timeout)
       for {
           select {
           case item := <-b.queue:
               batch = append(batch, item)
               if len(batch) >= b.batchSize {
                   b.processBatch(batch)
                   batch = batch[:0]
                   timer.Reset(b.timeout)
               }
           case <-timer.C:
               if len(batch) > 0 {
                   b.processBatch(batch) // Timeout flush
                   batch = batch[:0]
               }
               timer.Reset(b.timeout)
           case <-ctx.Done():
               return
           }
       }
   }
   ```
3. **SQL Transaction (ACID)**:
   - `BeginTx()` → insert block + insert transactions → `Commit()` hoặc `Rollback()`
   - UPSERT: `ON CONFLICT DO UPDATE`
4. **Event Publishing** (giả lập Message Broker):
   - Sau commit thành công → publish event `{ status: "new", data: {...} }`
   - Khi phát hiện conflict → publish `{ status: "reverted" }`
5. **Concurrency Control**:
   - Semaphore giới hạn goroutine ghi DB (`maxOpen / 2` — dự phòng cho API đọc)
6. **Graceful Shutdown**:
   - Bắt `SIGINT/SIGTERM` → flush batcher → đóng DB → exit

**🎯 Kỹ năng rèn:**
- Batcher pattern: gom lô theo size HOẶC timeout (whichever first)
- SQL Transaction: `BeginTx`, `Commit`, `Rollback` — đảm bảo ACID
- UPSERT pattern — xử lý trùng lặp
- Event-driven architecture (publish sau commit)
- Semaphore DB connection protection
- Graceful shutdown: `signal.Notify` + `context.WithTimeout`

**🗣️ Ba luận điểm vàng khi phỏng vấn:**
> 1. _"Giao dịch nguyên tử & UPSERT: Toàn bộ dữ liệu ghi → nghiệp vụ bọc trong 1 SQL Transaction. Lỗi giữa chừng → Rollback sạch sẽ, không rác."_
> 2. _"Cứu hộ bất đồng bộ: Dữ liệu lỗi mạng → đẩy vào failure tracking. Tiến trình nền quét định kỳ, cứu hộ tự động, không nghẽn luồng chính."_
> 3. _"Reorg Handler: Phát hiện dữ liệu không đồng nhất → tải lại dữ liệu chuẩn → ghi đè nguyên tử → publish 'reverted' cho downstream."_

---

## ⚡ Hướng Dẫn Thực Hành

### Thời gian đề xuất
| Level | Bài tập | Thời gian/bài | Tổng |
|:------|:--------|:-------------|:-----|
| Level 1 | Ex01-04 | 1-2 giờ | ~6 giờ |
| Level 2 | Ex05-08 | 2-3 giờ | ~10 giờ |
| Level 3 | Ex09-12 | 2-3 giờ | ~10 giờ |
| Level 4 | Ex13-16 | 3-4 giờ | ~14 giờ |
| Level 5 | Ex17-20 | 4-6 giờ | ~20 giờ |
| **Tổng** | **20 bài** | | **~60 giờ** |

### Quy tắc vàng
1. **Làm tuần tự** — mỗi bài xây trên kiến thức bài trước
2. **Gõ tay từng dòng** — không copy-paste, phải nhớ cú pháp
3. **Viết test từ Level 2** — `go test -race ./...`
4. **So sánh Node.js** — mỗi bài viết cùng logic bằng cả 2 ngôn ngữ
5. **Nếu stuck > 2 giờ** — hỏi ngay, đừng để mất momentum

### Dependencies cần cài
```bash
# Level 3+: Gin framework
go get -u github.com/gin-gonic/gin

# Level 4: JWT + bcrypt
go get -u github.com/golang-jwt/jwt/v5
go get -u golang.org/x/crypto/bcrypt

# Level 4: GORM + SQLite
go get -u gorm.io/gorm
go get -u gorm.io/driver/sqlite

# Level 5: WebSocket
go get -u github.com/gorilla/websocket

# Level 5: Redis
go get -u github.com/redis/go-redis/v9
```

---

## 🗺️ Mapping: Bài Tập ↔ Câu Hỏi Phỏng Vấn Senior

| Câu hỏi phỏng vấn dự kiến | Bài tập cover |
|:--------------------------|:-------------|
| _"Giải thích cơ chế Semaphore và giới hạn concurrency trong Go?"_ | **Ex17** |
| _"Tại sao không dùng SELECT FOR UPDATE cho các giao dịch concurrency cao?"_ | **Ex18** |
| _"Làm sao quản lý hàng vạn WebSocket connections thread-safe?"_ | **Ex19** |
| _"Làm thế nào để đảm bảo tính đúng đắn dữ liệu khi ghi database hàng loạt (Batching)?"_ | **Ex20** |
| _"Giải thích Channel-based concurrency trong Go?"_ | **Ex06, Ex07** |
| _"Context dùng để làm gì? Quên cancel thì sao?"_ | **Ex08** |
| _"Gin middleware hoạt động thế nào? So sánh với Express?"_ | **Ex11** |
| _"Clean Architecture trong Go khác NestJS thế nào?"_ | **Ex13** |
| _"Stateless authentication với JWT và các phương án nâng cao?"_ | **Ex14** |
| _"UPSERT và tối ưu hóa ghi DB hàng loạt?"_ | **Ex15, Ex20** |
| _"Table-driven tests là gì? Tại sao Go community ưa chuộng?"_ | **Ex16** |
| _"Graceful Shutdown xử lý thế nào?"_ | **Ex20** |

---

> _Hoàn thành 20 bài tập này = bạn tự tin code Go + Gin như một Senior thực thụ trước mặt interviewer._
>
> **Let's Go! 🚀**
