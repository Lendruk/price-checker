CREATE TABLE IF NOT EXISTS vendorEntries (
  id INTEGER PRIMARY KEY,
  fullName TEXT,
  price REAL,
  url TEXT,
  vendor INTEGER,
  sku TEXT,
  availability INTEGER,
  lastUpdated INTEGER,
  universalId INTEGER,
  FOREIGN KEY(universalId) REFERENCES products(id)
);

CREATE TABLE IF NOT EXISTS productHistory (
  id INTEGER PRIMARY KEY,
  vendorEntryId INTEGER,
  price REAL,
  availability INTEGER,
  updatedAt INTEGER,
  FOREIGN KEY(vendorEntryId) REFERENCES vendorEntries(id)
);

CREATE TABLE IF NOT EXISTS products (id INTEGER PRIMARY KEY, sku TEXT UNIQUE);

CREATE TABLE IF NOT EXISTS users (
  id INTEGER PRIMARY KEY
);

CREATE TABLE IF NOT EXISTS watchlists (
  user INTEGER,
  product INTEGER,
  FOREIGN KEY(user) REFERENCES users(id),
  FOREIGN KEY(product) REFERENCES products(id)
);

CREATE TABLE IF NOT EXISTS webhooks (
  id INTEGER PRIMARY KEY,
  hook TEXT UNIQUE
);

CREATE TABLE IF NOT EXISTS webhook_users (
  hook INTEGER,
  user INTEGER,
  FOREIGN KEY(hook) REFERENCES webhooks (id),
  FOREIGN KEY (user) REFERENCES users (id)
);