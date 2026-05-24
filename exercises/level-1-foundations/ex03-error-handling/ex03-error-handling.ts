/**
 * Ex03: Error Handling — TypeScript Version
 * 
 * 🧠 So sánh key:
 * - Node.js: Sử dụng try/catch/throw để quản lý lỗi động, có thể bọc lỗi bằng cause.
 * - Go:         Sử dụng cơ chế kiểm tra tường minh `if err != nil`. Bọc lỗi bằng error wrapping
 *               thông qua %w để giữ nguyên chuỗi lỗi (error chain) và unwrap khi cần.
 * 
 * 💡 Sự khác biệt lớn nhất:
 * 1. Go không có try/catch truyền thống. Mọi lỗi đều phải được xử lý ngay tại nơi xảy ra hoặc chuyển tiếp có chủ đích.
 * 2. Để kiểm tra kiểu lỗi (custom error types), Go sử dụng `errors.As()` (tương đương `instanceof` trong JS).
 * 3. Để kiểm tra giá trị lỗi cụ thể (sentinel errors), Go sử dụng `errors.Is()` (tương đương `===` trong JS).
 */

import * as fs from 'fs';

class ValidationError extends Error {
  public field: string;

  constructor(field: string, message: string) {
    super(`validation error: field '${field}' — ${message}`);
    this.name = "ValidationError";
    this.field = field;
  }
}

class NotFoundError extends Error {
  public resource: string;
  public id: string;

  constructor(resource: string, id: string) {
    super(`${resource} with id '${id}' not found`);
    this.name = "NotFoundError";
    this.resource = resource;
    this.id = id;
  }
}

interface Config {
  port: number;
}

function parseConfig(filename: string): Config {
  if (!fs.existsSync(filename)) {
    throw new NotFoundError("config file", filename);
  }

  const content = fs.readFileSync(filename, "utf-8");

  let config: any;
  try {
    config = JSON.parse(content);
  } catch {
    throw new ValidationError("content", "invalid JSON format");
  }

  if (!config.port) {
    throw new ValidationError("port", "is required");
  }

  if (typeof config.port !== "number") {
    throw new ValidationError("port", "must be a number");
  }

  return config as Config;
}

function main() {
  const testFiles = [
    "nonexistent.json",
    "/dev/null",
  ];

  for (const file of testFiles) {
    try {
      const config = parseConfig(file);
      console.log("Config loaded:", config);
    } catch (err) {
      if (err instanceof NotFoundError) {
        console.error(`[NOT FOUND] ${err.resource}: ${err.id}`);
      } else if (err instanceof ValidationError) {
        console.error(`[VALIDATION] field '${err.field}': ${err.message}`);
      } else {
        console.error(`[UNKNOWN] ${(err as Error).message}`);
      }
    }
  }

  console.log("\n--- Error Wrapping Demo ---");

  function loadAppConfig(): Config {
    try {
      return parseConfig("app.json");
    } catch (err) {
      throw new Error("failed to load app config", { cause: err });
    }
  }

  try {
    loadAppConfig();
  } catch (err) {
    console.error("Top-level error:", (err as Error).message);
    console.error("Original cause:", ((err as Error).cause as Error)?.message);
  }
}

main();
