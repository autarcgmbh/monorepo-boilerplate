import { Hono } from 'hono';
import { serve } from '@hono/node-server';
import { db, initDb } from './db';

const app = new Hono();

type Product = {
  id: number;
  name: string | null;
  manufacture: string | null;
  output: number | null;
  price: number | null;
  width: number | null;
  height: number | null;
};

app.get('/products', async (c) => {
  const result = await db.query<Product>('SELECT * FROM products ORDER BY id;');
  return c.json(result.rows);
});

app.get('/products/:id', async (c) => {
  const id = parseInt(c.req.param('id'));
  const result = await db.query<Product>(
    'SELECT * FROM products WHERE id = $1;',
    [id]
  );

  if (result.rows.length === 0) {
    return c.json({ error: 'Product not found' }, 404);
  }

  return c.json(result.rows[0]);
});

app.post('/products', async (c) => {
  // Fail randomly 50% of the time
  if (Math.random() < 0.5) {
    return c.json({ error: 'Internal server error' }, 500);
  }

  const body = await c.req.json();
  const result = await db.query<Product>(
    'INSERT INTO products (name, manufacture, output, price, width, height) VALUES ($1, $2, $3, $4, $5, $6) RETURNING *;',
    [
      body.name,
      body.manufacture,
      body.output,
      body.price,
      body.width,
      body.height,
    ]
  );
  return c.json(result.rows[0], 201);
});

app.put('/products/:id', async (c) => {
  const id = parseInt(c.req.param('id'));
  const body = await c.req.json();
  const result = await db.query<Product>(
    'UPDATE products SET name = $1, manufacture = $2, output = $3, price = $4, width = $5, height = $6 WHERE id = $7 RETURNING *;',
    [
      body.name,
      body.manufacture,
      body.output,
      body.price,
      body.width,
      body.height,
      id,
    ]
  );

  if (result.rows.length === 0) {
    return c.json({ error: 'Product not found' }, 404);
  }

  return c.json(result.rows[0]);
});

app.delete('/products/:id', async (c) => {
  const id = parseInt(c.req.param('id'));
  const result = await db.query(
    'DELETE FROM products WHERE id = $1 RETURNING *;',
    [id]
  );

  if (result.rows.length === 0) {
    return c.json({ error: 'Product not found' }, 404);
  }

  return c.json({ message: 'Product deleted' });
});

const port = 3001;

// Initialize database before starting server
await initDb();
console.log(`Server running on http://localhost:${port}`);

serve({
  fetch: app.fetch,
  port,
});
