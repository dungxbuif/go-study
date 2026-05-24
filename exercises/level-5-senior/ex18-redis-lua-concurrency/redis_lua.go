package main

import (
	"context"
	"errors"

	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

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

func (s *LuckyMoneyService) Claim(id string, userID string) (int, string, error) {
	// TODO: Sử dụng s.rdb.Eval để thực thi AttemptClaimLua một cách nguyên tử (atomic) trên Redis.
	// Nhận kết quả và trả về (amount, status, error).
	return 0, "", errors.New("not implemented")
}

func main() {
	// Demo Redis Lua script concurrency tại đây
}
