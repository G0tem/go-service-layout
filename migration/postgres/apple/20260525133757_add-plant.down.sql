BEGIN;
-- Откатываем СТРОГО в обратном порядке, учитывая зависимости БД:
DROP INDEX IF EXISTS idx_apple_plant_id;
DROP INDEX IF EXISTS idx_plant_name;

-- Сначала убираем связь, потом колонку (или одной командой, но явно надёжнее)
ALTER TABLE apple DROP CONSTRAINT IF EXISTS fk_apple_plant;
ALTER TABLE apple DROP COLUMN IF EXISTS plant_id;

-- В конце удаляем родительскую таблицу
DROP TABLE IF EXISTS plant;
COMMIT;