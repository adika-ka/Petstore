CREATE TABLE IF NOT EXISTS orders (
    id SERIAL PRIMARY KEY,
    pet_id INT NOT NULL,
    quantity INT NOT NULL CHECK (quantity > 0),
    ship_date TIMESTAMP NOT NULL,
    status VARCHAR(20) CHECK (status IN ('placed', 'approved', 'delivered')),
    complete BOOLEAN NOT NULL DEFAULT FALSE,
    FOREIGN KEY (pet_id) REFERENCES pets(id) ON DELETE CASCADE
);

CREATE INDEX idx_orders_pet_id ON orders(pet_id);