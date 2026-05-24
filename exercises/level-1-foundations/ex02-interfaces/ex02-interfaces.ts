/**
 * Ex02: Interfaces — TypeScript Version
 * 
 * 🧠 So sánh key:
 * - TypeScript: Sử dụng Explicit Interface (implements) — class bắt buộc phải khai báo tường minh.
 * - Go:         Sử dụng Implicit Interface (Duck Typing) — struct chỉ cần định nghĩa đủ methods
 *               là tự động thoả mãn interface mà không cần khai báo từ khoá implements.
 * 
 * 💡 Sự khác biệt lớn nhất:
 * 1. Go interface cho phép định nghĩa bởi CONSUMER (bên nhận) thay vì PRODUCER (bên cung cấp).
 * 2. Giúp giảm thiểu sự phụ thuộc (decoupling) giữa các package cực kỳ mạnh mẽ.
 * 3. Go Dependency Injection được thực hiện tự nhiên bằng cách truyền interface vào struct/function.
 */

interface PaymentGateway {
  charge(amount: number, currency: string): Promise<{ txID: string }>;
  refund(txID: string): Promise<void>;
}

class StripeGateway implements PaymentGateway {
  async charge(
    amount: number,
    currency: string
  ): Promise<{ txID: string }> {
    if (amount > 10000) {
      throw new Error("Stripe: amount exceeds limit");
    }
    const txID = `stripe_${Date.now()}`;
    console.log(`Stripe: Charged ${amount} ${currency} → tx: ${txID}`);
    return { txID };
  }

  async refund(txID: string): Promise<void> {
    console.log(`Stripe: Refunded tx ${txID}`);
  }
}

class MomoGateway implements PaymentGateway {
  async charge(
    amount: number,
    currency: string
  ): Promise<{ txID: string }> {
    if (currency !== "VND") {
      throw new Error("Momo: only supports VND");
    }
    const txID = `momo_${Date.now()}`;
    console.log(`Momo: Charged ${amount} ${currency} → tx: ${txID}`);
    return { txID };
  }

  async refund(txID: string): Promise<void> {
    if (!txID) {
      throw new Error("Momo: txID is required for refund");
    }
    console.log(`Momo: Refunded tx ${txID}`);
  }
}

async function processOrder(
  gw: PaymentGateway,
  amount: number
): Promise<void> {
  try {
    const { txID } = await gw.charge(amount, "VND");
    console.log(`Order processed successfully: ${txID}`);
  } catch (err) {
    console.error(`Order failed: ${(err as Error).message}`);
  }
}

async function main() {
  const stripe = new StripeGateway();
  const momo = new MomoGateway();

  console.log("=== Using Stripe ===");
  await processOrder(stripe, 5000);
  await processOrder(stripe, 15000);

  console.log("\n=== Using Momo ===");
  await processOrder(momo, 50000);
}

main();
