import { PGlite } from '@electric-sql/pglite';

// Initialize PGlite with file-based storage
export const db = new PGlite('./data');

// Initialize schema
export async function initDb() {
  await db.exec(`
    CREATE TABLE IF NOT EXISTS products (
      id SERIAL PRIMARY KEY,
      name TEXT,
      manufacture TEXT,
      output NUMERIC,
      price INTEGER,
      width NUMERIC,
      height NUMERIC
    );
  `);

  // Check if we need to seed initial data
  const result = await db.query('SELECT * FROM products;');
  const count = result.rows.length;

  if (count === 0) {
    await db.query(`
      INSERT INTO products (name, manufacture, output, price, width, height) VALUES
        ('Air Source Heat Pump', 'Generic Manufacturer', 12.5, 45000, 80, 120);
    `);
    console.log('Database seeded with initial data');
  }

  console.log('Database initialized');
}
