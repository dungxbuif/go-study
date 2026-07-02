package main

import (
	"context"
	"errors"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

// AttemptClaimLua là mã script Lua thực thi giật bao lì xì một cách nguyên tử.
//
// 🧠 TẠI SAO LUA SCRIPT LÀ "CHÌA KHÓA VÀNG" CHO BÀI TOÁN CONCURRENCY CỰC LỚN:
// - Redis hoạt động theo mô hình đơn luồng (Single-Threaded Event Loop) để thực thi các câu lệnh thực tế.
// - Khi một Lua script được chạy (qua lệnh EVAL hoặc EVALSHA), Redis sẽ khóa và chạy TOÀN BỘ script đó một cách liên tục và độc quyền.
// - Không có bất cứ client nào khác có thể chen ngang hay thực thi lệnh nào khác vào giữa quá trình script đang chạy.
// - Do đó, chuỗi hành động: (1) Đọc HGET kiểm tra -> (2) LPOP rút bao lì xì -> (3) HSET đánh dấu
//   được gộp lại thành 1 giao dịch có tính chất nguyên tử (Atomic Transaction) tuyệt đối.
// - Ta hoàn toàn loại bỏ được rủi ro race condition (1 bao lì xì bị giật bởi 2 người) mà KHÔNG cần dùng
//   tới các giải pháp khóa phân tán (Distributed Locks như Redlock) vốn tốn kém tài nguyên và latency cao.
const AttemptClaimLua = `
local poolKey = KEYS[1]
local reservedKey = KEYS[2]
local userID = ARGV[1]

local reserved = redis.call('HGET', reservedKey, userID)
if reserved then return {reserved, 'ALREADY_CLAIMED'} end

local amount = redis.call('LPOP', poolKey)
if not amount then return {0, 'SOLD_OUT'} end

redis.call('HSET', reservedKey, userID, amount)
return {amount, 'OK'}
`

type LuckyMoneyService struct {
	rdb *redis.Client
}

func NewLuckyMoneyService(rdb *redis.Client) *LuckyMoneyService {
	return &LuckyMoneyService{rdb: rdb}
}

func (s *LuckyMoneyService) CreateEnvelope(id string, amounts []int) error {
	// TODO: Đẩy danh sách amounts vào Redis List và dọn dẹp hash cũ của id
	return errors.New("not implemented")
}

// Claim thực hiện rút bao lì xì.
//
// 🧠 TỐI ƯU HÓA BĂNG THÔNG MẠNG (EVALSHA under the hood):
// - Nếu mỗi request giật bao lì xì ta đều gửi toàn bộ chuỗi text `AttemptClaimLua` qua mạng,
//   băng thông mạng của Server Redis sẽ nhanh chóng bị nghẽn (bandwidth bottleneck) khi có hàng triệu users giật cùng lúc.
// - Để giải quyết, driver `go-redis` dưới nắp capo sử dụng cơ chế **EVALSHA**:
//   1. Client tính toán mã băm SHA1 của chuỗi Lua Script (ví dụ: `23c8a9f...`).
//   2. Client gửi lệnh `EVALSHA 23c8a9f...` cùng các tham số.
//   3. Redis kiểm tra trong cache bộ nhớ xem có Script nào khớp với mã SHA1 này chưa.
//      - Nếu có, nó thực thi ngay lập tức mà không cần nhận lại thân code (tiết kiệm 99% băng thông).
//      - Nếu chưa (lỗi `NOSCRIPT`), Client mới tự động tải toàn bộ nội dung Script lên bằng lệnh `EVAL` để Redis cache lại cho các lần sau.
func (s *LuckyMoneyService) Claim(id string, userID string) (int, string, error) {
	// TODO: Sử dụng s.rdb.Eval để thực thi AttemptClaimLua một cách nguyên tử (atomic) trên Redis.
	// Nhận kết quả và trả về (amount, status, error).
	return 0, "", errors.New("not implemented")
}

func main() {
	// Demo Redis Lua script concurrency tại đây
}
