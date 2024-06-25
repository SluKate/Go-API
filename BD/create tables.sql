CREATE TABLE measures (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255)
);

INSERT INTO measures (name)
VALUES ('кг'),
       ('литр'),
       ('штука');

CREATE TABLE products (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255),
    quantity FLOAT(24),
    unit_cost FLOAT(24),
    measure INTEGER REFERENCES measures (id) ON DELETE CASCADE
);

INSERT INTO products (name, quantity, unit_cost, measure)
VALUES ('Мандарины', 30, 20.0, 1),
       ('Вода', 20, 40.7, 2),
       ('Яблоки', 66, 23.7, 3);