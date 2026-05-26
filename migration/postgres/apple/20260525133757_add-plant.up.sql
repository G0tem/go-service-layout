BEGIN;
-- Создаём родительскую таблицу
CREATE TABLE plant (
    id UUID PRIMARY KEY,
    name TEXT,
    status TEXT,
    created_at TIMESTAMPTZ DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMPTZ,
    deleted_at TIMESTAMPTZ
);
CREATE INDEX idx_plant_name ON plant (name);

-- Добавляем колонку и FK к существующей таблице apple
ALTER TABLE apple ADD COLUMN plant_id UUID;
ALTER TABLE apple ADD CONSTRAINT fk_apple_plant 
    FOREIGN KEY (plant_id) REFERENCES plant(id) ON DELETE CASCADE;

-- Индекс на новую связь
CREATE INDEX idx_apple_plant_id ON apple (plant_id);
COMMIT;