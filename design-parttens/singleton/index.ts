class Database {
   private static instance: Database;
   private connection: string;

   private constructor() {
      this.connection = `Connected at ${new Date().toISOString()}`;
   }

   static getInstance(): Database {
      if (!Database.instance) {
         Database.instance = new Database();
      }
      return Database.instance;
   }

   query(sql: string): void {
      console.log(`[DB] Running: ${sql}`);
   }
}

const db1 = Database.getInstance();
const db2 = Database.getInstance();

console.log(db1 === db2); // true
db1.query('SELECT * FROM users');

export default Database;
