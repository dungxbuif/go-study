/**
 * Ex18: Redis Lua Script — TypeScript Version
 * 
 * 🧠 So sánh key:
 * - Node.js: Thường sử dụng thư viện `ioredis` để tải và chạy các mã script Lua nguyên tử (atomic).
 * - Go:      Dùng các thư viện driver như `go-redis` để quản lý các Redis script, băm SHA, 
 *            và gọi lệnh atomic trên RAM Redis.
 * 
 * 💡 Sự khác biệt lớn nhất:
 * 1. Với những bài toán concurrency cực lớn (như giật bao lì xì, flash-sale), các lệnh khóa DB truyền thống 
 *    như `SELECT FOR UPDATE` gây nghẽn cổ chai DB cực kỳ nghiêm trọng (Lock Contention).
 * 2. Giải pháp hoàn hảo là dùng Redis Lua Script chạy đơn luồng cực nhanh trên RAM, đảm bảo 
 *    tính nguyên tử (atomic), sau đó đưa kết quả "giữ chỗ" đó cho một Async Worker lưu bền vững (durability) xuống PostgreSQL.
 */

import Redis from 'ioredis';
import express, { Request, Response } from 'express';

const app = express();
app.use(express.json());

const redis = new Redis() as any;

const attemptClaimLua = `
local poolKey = KEYS[1]
local reservedKey = KEYS[2]
local userID = ARGV[1]

local reserved = redis.call('HGET', reservedKey, userID)
if reserved then return {reserved, 'ALREADY_CLAIMED'} end

local amount = redis.call('LPOP', poolKey)
if not amount then return {0, 'SOLD_OUT'} end

redis.call('HSET', reservedKey, userID, amount)
return {amount, 'OK'}
`;

redis.defineCommand('attemptClaim', {
  numberOfKeys: 2,
  lua: attemptClaimLua,
});

app.post('/red-envelope/create', async (req: Request, res: Response) => {
  const { id, amounts } = req.body;
  const poolKey = `red_envelope:${id}:pool`;

  await redis.del(poolKey);
  await redis.del(`red_envelope:${id}:reserved`);

  await redis.rpush(poolKey, ...amounts);
  res.status(201).json({ success: true, message: 'Envelope created' });
});

app.post('/red-envelope/:id/claim', async (req: Request, res: Response) => {
  const id = req.params.id;
  const { user_id } = req.body;

  const poolKey = `red_envelope:${id}:pool`;
  const reservedKey = `red_envelope:${id}:reserved`;

  try {
    const [amount, status] = await redis.attemptClaim(poolKey, reservedKey, user_id);
    res.json({ success: true, amount: parseInt(amount, 10), status });
  } catch (err: any) {
    res.status(500).json({ success: false, error: err.message });
  }
});

app.get('/red-envelope/:id/status', async (req: Request, res: Response) => {
  const id = req.params.id;
  const reservedKey = `red_envelope:${id}:reserved`;

  const claims = await redis.hgetall(reservedKey);
  res.json({ success: true, data: claims });
});

app.listen(8080, () => {
  console.log('Redis Lucky Money Express server running on port 8080');
});
