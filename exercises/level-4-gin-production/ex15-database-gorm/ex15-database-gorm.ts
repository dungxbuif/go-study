/**
 * Ex15: Database + SQL — TypeScript Version
 *
 * 🧠 So sánh key:
 * - TypeScript: Dùng Sequelize, Prisma hay TypeORM để viết các câu lệnh truy vấn,
 *               preload quan hệ (join/include), phân trang và soft delete.
 * - Go:         Dùng GORM framework để mapping struct sang database table, preload quan hệ,
 *               phân trang và quản lý Connection Pool (`SetMaxOpenConns`, `SetMaxIdleConns`).
 *
 * 💡 Sự khác biệt lớn nhất:
 * 1. Mặc dù ORM tiện lợi, các hệ thống Go lớn ưu tiên hiệu năng cực cao thường dùng SQL thuần (`database/sql` hoặc `sqlx`, `pgx`)
 *    cho các tác vụ phức tạp hoặc ghi dữ liệu lớn.
 * 2. UPSERT pattern (`ON CONFLICT DO UPDATE`) giúp tránh lỗi ghi đè dữ liệu bất đồng bộ đồng thời cực kỳ hiệu quả.
 */

import { DataTypes, Model, Sequelize } from 'sequelize';

// 🧠 CONNECTION POOL TRONG NODE.JS (Sequelize vs Go):
// - Sequelize sử dụng một connection pool nội bộ (thông qua thư viện `generic-pool`) để quản lý các kết nối tới Database.
// - Vì Node.js chạy đơn luồng, các kết nối trong pool được chia sẻ bất đồng bộ thông qua Event Loop.
// - Khi một truy vấn được gửi đi (ví dụ: `Post.create()`), Sequelize sẽ rút (borrow) một kết nối vật lý, gửi lệnh SQL
//   qua socket, và đăng ký một callback chờ kết quả. Nhờ non-blocking I/O, luồng chính không bị đứng im chờ database phản hồi.
// - Go cũng sử dụng cơ chế pool tương tự (`database/sql` pool), nhưng mỗi kết nối được sử dụng độc quyền bởi một Goroutine đồng bộ.
//   Mô hình này giúp Go loại bỏ hoàn toàn chi phí quản lý bất đồng bộ phức tạp (Event Loop callbacks) và có hiệu năng thực thi gần như bare-metal.
const sequelize = new Sequelize({
   dialect: 'sqlite',
   storage: ':memory:',
   logging: false,
});

class User extends Model {
   public id!: number;
   public username!: string;
}

User.init(
   {
      id: {
         type: DataTypes.INTEGER,
         autoIncrement: true,
         primaryKey: true,
      },
      username: {
         type: DataTypes.STRING,
         allowNull: false,
         unique: true,
      },
   },
   { sequelize, modelName: 'user' },
);

class Post extends Model {
   public id!: number;
   public title!: string;
   public content!: string;
   public authorId!: number;
}

Post.init(
   {
      id: {
         type: DataTypes.INTEGER,
         autoIncrement: true,
         primaryKey: true,
      },
      title: {
         type: DataTypes.STRING,
         allowNull: false,
      },
      content: {
         type: DataTypes.TEXT,
         allowNull: false,
      },
   },
   { sequelize, modelName: 'post' },
);

User.hasMany(Post, { foreignKey: 'authorId', as: 'posts' });
Post.belongsTo(User, { foreignKey: 'authorId', as: 'author' });

// upsertPost: Thực hiện UPSERT nguyên tử bằng SQL thô.
// Bằng cách sử dụng câu lệnh `ON CONFLICT` tại mức cơ sở dữ liệu, ta đảm bảo tính nguyên tử tuyệt đối (Atomicity).
// Sequelize cũng hỗ trợ `.upsert()` tự động chuyển dịch cú pháp này dựa trên hệ quản trị cơ sở dữ liệu (PostgreSQL, SQLite, MySQL...).
async function upsertPost(title: string, content: string, authorId: number) {
   await sequelize.query(
      `INSERT INTO posts (title, content, authorId, createdAt, updatedAt)
     VALUES (:title, :content, :authorId, CURRENT_TIMESTAMP, CURRENT_TIMESTAMP)
     ON CONFLICT(id) DO UPDATE SET
       title = EXCLUDED.title,
       content = EXCLUDED.content,
       updatedAt = CURRENT_TIMESTAMP`,
      {
         replacements: { title, content, authorId },
      },
   );
}

// getPostsWithPagination: Phân trang và liên kết dữ liệu (Eager Loading).
//
// 🧠 CƠ CHẾ EAGER LOADING (Sequelize include vs GORM Preload):
// - Trong Sequelize, khi sử dụng `include: [{ model: User }]`, mặc định Sequelize sẽ tạo ra một câu lệnh `LEFT OUTER JOIN`
//   để gộp 2 bảng lại làm một trong một lần truy vấn duy nhất.
// - Trái lại trong Go GORM, `Preload` mặc định tách thành 2 câu truy vấn tuần tự cực kỳ sạch sẽ (Select * from posts -> Select * from users where id IN (...)).
// - Tải bằng JOIN (Sequelize) có ưu điểm là chỉ tốn 1 round-trip tới DB, nhưng nếu dữ liệu có mối quan hệ 1-N lớn (ví dụ: N bản ghi con),
//   nó sẽ dẫn đến hiện tượng trùng lặp dữ liệu dư thừa khổng lồ ở các cột cha (Cartesian Product), gây tốn RAM để truyền tải và parse.
// - Việc GORM tách thành các câu lệnh `IN` tuần tự đôi khi đem lại hiệu năng tối ưu hơn rất nhiều cho các tập dữ liệu phức tạp.
async function getPostsWithPagination(page: number, limit: number) {
   const offset = (page - 1) * limit;
   const { rows, count } = await Post.findAndCountAll({
      limit,
      offset,
      include: [{ model: User, as: 'author', attributes: ['username'] }],
   });

   return {
      total: count,
      page,
      limit,
      data: rows,
   };
}

async function main() {
   await sequelize.sync({ force: true });

   const alice = await User.create({ username: 'alice' });

   await Post.create({
      title: 'First Post',
      content: 'Hello World',
      authorId: alice.id,
   });
   await Post.create({
      title: 'Second Post',
      content: 'Learning SQL',
      authorId: alice.id,
   });

   await upsertPost('Upserted Post', 'This is upserted content', alice.id);

   const paginatedResults = await getPostsWithPagination(1, 10);
   console.log(JSON.stringify(paginatedResults, null, 2));

   await sequelize.close();
}

main();
