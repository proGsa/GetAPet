CREATE TABLE IF NOT EXISTS vet_passport (
    id SERIAL PRIMARY KEY,                  -- Уникальный идентификатор записи
    chipping BOOLEAN DEFAULT FALSE,         -- Чипирование (TRUE/FALSE)
    sterilization BOOLEAN DEFAULT FALSE,    -- Стерилизация (TRUE/FALSE)
    health_issues TEXT,                     -- Проблемы со здоровьем
    vaccinations TEXT,                      -- Прививки (хранить список дат)
    parasite_treatments TEXT                -- Обработка от паразитов (хранить список дат)
);

CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,                  -- Уникальный идентификатор записи
    fio VARCHAR(255) NOT NULL,              -- ФИО пользователя
    telephone_number VARCHAR(20) NOT NULL,  -- Номер телефона
    city VARCHAR(50),                        -- Город пользователя
    user_login VARCHAR(50) UNIQUE,          -- Логин пользователя
    user_password VARCHAR(255),             -- Пароль 
    status VARCHAR(20),                      -- Статус пользователя
    user_description TEXT                    -- Дополнительная информация о пользователе
);

CREATE TABLE IF NOT EXISTS pet (
    id SERIAL PRIMARY KEY,                    -- Уникальный идентификатор питомца
    vet_passport_id INT REFERENCES vet_passport(id),  -- Ссылка на ветпаспорт
    seller_id INT REFERENCES users(id),      -- Ссылка на владельца / продавца
    pet_name VARCHAR(255) NOT NULL,          -- Кличка питомца
    species VARCHAR(50) NOT NULL,            -- Вид животного (например: кошка, собака)
    pet_age INT NOT NULL,                     -- Возраст питомца (в годах)
    color VARCHAR(50),                        -- Окрас
    pet_gender VARCHAR(20),                   -- Пол питомца
    breed VARCHAR(255),                        -- Порода
    pedigree BOOLEAN DEFAULT FALSE,           -- Наличие родословной
    good_with_children BOOLEAN DEFAULT TRUE,  -- Ладит с детьми
    good_with_animals BOOLEAN DEFAULT TRUE,   -- Ладит с другими животными
    pet_description TEXT,                     -- Дополнительная информация о питомце
    is_active BOOLEAN DEFAULT TRUE,           -- Актуальность объявления
    price DECIMAL(10,2)                       -- Цена питомца
);

CREATE TABLE IF NOT EXISTS purchase_request (
    id SERIAL PRIMARY KEY,                   -- Уникальный идентификатор заявки
    pet_id INT NOT NULL REFERENCES pet(id),  -- Ссылка на питомца
    seller_id INT NOT NULL REFERENCES users(id),  -- Ссылка на продавца
    status VARCHAR(50) DEFAULT 'pending',    -- Статус заявки (например: pending, approved, rejected)
    request_date TIMESTAMP DEFAULT NOW()     -- Дата и время создания заявки
);