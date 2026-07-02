/**
 * Ex21: Circuit Breaker + Fallback — TypeScript Version
 * 
 * 🧠 So sánh key:
 * - TypeScript: Trong Node.js, ta có thể xây dựng Circuit Breaker bằng cách dùng lớp/hàm bất đồng bộ
 *               kết hợp với quản lý trạng thái tĩnh. Vì chạy đơn luồng, ta không lo ngại tranh chấp bộ nhớ 
 *               khi thay đổi trạng thái mạch giữa các requests đồng thời.
 * - Go:         Đòi hỏi sử dụng `sync.Mutex` để đồng bộ hóa tuyệt đối việc cập nhật trạng thái mạch 
 *               (Closed, Open, Half-Open) giữa hàng ngàn Goroutines song song.
 */

enum State {
  Closed,
  Open,
  HalfOpen,
}

class CircuitBreaker {
  private state: State;
  private failureCount: number;
  private failureThreshold: number;
  private cooldownWindowMs: number;
  private lastStateChangeTime: number;

  constructor(failureThreshold: number, cooldownWindowMs: number) {
    this.state = State.Closed;
    this.failureCount = 0;
    this.failureThreshold = failureThreshold;
    this.cooldownWindowMs = cooldownWindowMs;
    this.lastStateChangeTime = Date.now();
  }

  // execute: Thực thi tác vụ bất đồng bộ được bảo vệ bởi Circuit Breaker.
  //
  // 🧠 CƠ CHẾ KHÔNG KHÓA CỦA NODE.JS (Lock-free asynchronous state transition):
  // - Khác với Go cần `sync.Mutex` để khóa đồng bộ hóa ô nhớ trạng thái, Node.js giải quyết cực kỳ an toàn
  //   nhờ cơ chế Event Loop chạy đơn luồng.
  // - Quá trình chuyển mạch từ Open sang Half-Open bằng thời gian (`Date.now() - lastStateChangeTime`)
  //   diễn ra đồng bộ và liên tục mà không bao giờ bị chen ngang hay bị tranh chấp tài nguyên bộ nhớ bởi các luồng khác.
  public async execute<T>(action: () => Promise<T>, fallback: () => Promise<T>): Promise<T> {
    // Lazy state transition from Open to HalfOpen
    if (this.state === State.Open && Date.now() - this.lastStateChangeTime > this.cooldownWindowMs) {
      console.log('[CircuitBreaker] Cooldown elapsed. Transitioning to HALF-OPEN');
      this.state = State.HalfOpen;
      this.lastStateChangeTime = Date.now();
    }

    if (this.state === State.Open) {
      console.log('[CircuitBreaker] Circuit is OPEN. Triggering fallback fast-fail.');
      return fallback();
    }

    try {
      const result = await action();
      this.onSuccess();
      return result;
    } catch (err) {
      this.onFailure();
      return fallback(); // Chạy fallback cứu hộ
    }
  }

  private onFailure(): void {
    this.failureCount++;
    console.log(`[CircuitBreaker] Failure recorded. Count: ${this.failureCount}, State: ${State[this.state]}`);

    if (this.state === State.Closed && this.failureCount >= this.failureThreshold) {
      console.log('[CircuitBreaker] Threshold exceeded. Tripping circuit to OPEN!');
      this.state = State.Open;
      this.lastStateChangeTime = Date.now();
    } else if (this.state === State.HalfOpen) {
      console.log('[CircuitBreaker] Test request failed in HALF-OPEN. Tripping circuit back to OPEN!');
      this.state = State.Open;
      this.lastStateChangeTime = Date.now();
    }
  }

  private onSuccess(): void {
    console.log(`[CircuitBreaker] Success recorded. State: ${State[this.state]}`);
    if (this.state === State.HalfOpen) {
      console.log('[CircuitBreaker] Test request succeeded in HALF-OPEN. Closing circuit back to CLOSED!');
      this.state = State.Closed;
      this.failureCount = 0;
    } else if (this.state === State.Closed) {
      this.failureCount = 0; // Reset bộ đếm lỗi khi thành công
    }
  }
}

async function sleep(ms: number): Promise<void> {
  return new Promise(resolve => setTimeout(resolve, ms));
}

async function main(): Promise<void> {
  const cb = new CircuitBreaker(3, 2000);

  const unstableAPI = (shouldFail: boolean) => {
    return async (): Promise<string> => {
      if (shouldFail) {
        throw new Error('Third-party API failure');
      }
      return 'API success';
    };
  };

  const fallbackData = async (): Promise<string> => {
    console.log('[Fallback] Returning cached static mock data.');
    return 'Fallback response';
  };

  console.log('--- 1. Testing normal behavior (CLOSED) ---');
  await cb.execute(unstableAPI(false), fallbackData);

  console.log('\n--- 2. Inducing errors to trip the circuit ---');
  for (let i = 0; i < 4; i++) {
    await cb.execute(unstableAPI(true), fallbackData);
  }

  console.log('\n--- 3. Circuit is now OPEN. Requests should fail-fast instantly ---');
  await cb.execute(unstableAPI(false), fallbackData);

  console.log('\n--- 4. Waiting for Cooldown window to expire (2 seconds) ---');
  await sleep(2100);

  console.log('\n--- 5. Making a request now. Circuit should transition to Half-Open and close ---');
  await cb.execute(unstableAPI(false), fallbackData);

  console.log('\n--- 6. Circuit is CLOSED again. Back to normal ---');
  await cb.execute(unstableAPI(false), fallbackData);
}

main();
