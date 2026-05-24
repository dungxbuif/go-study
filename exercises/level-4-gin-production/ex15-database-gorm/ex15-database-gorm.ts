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

import { Sequelize, DataTypes, Model } from 'sequelize';

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
  { sequelize, modelName: 'user' }
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
  { sequelize, modelName: 'post' }
);

User.hasMany(Post, { foreignKey: 'authorId', as: 'posts' });
Post.belongsTo(User, { foreignKey: 'authorId', as: 'author' });

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
    }
  );
}

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

  await Post.create({ title: 'First Post', content: 'Hello World', authorId: alice.id });
  await Post.create({ title: 'Second Post', content: 'Learning SQL', authorId: alice.id });

  await upsertPost('Upserted Post', 'This is upserted content', alice.id);

  const paginatedResults = await getPostsWithPagination(1, 10);
  console.log(JSON.stringify(paginatedResults, null, 2));

  await sequelize.close();
}

main();
