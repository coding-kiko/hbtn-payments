CREATE TABLE IF NOT EXISTS payments (
    id uuid default gen_random_uuid() PRIMARY KEY UNIQUE,
    created_at VARCHAR(255) default NOW(),
    month VARCHAR(255) UNIQUE,
    amount INTEGER,
    receipt_url VARCHAR(255),
    company VARCHAR(255) default 'N/A'
);