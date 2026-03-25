INSERT INTO vet_passport (chipping, sterilization, health_issues, vaccinations, parasite_treatments)
VALUES
(TRUE, TRUE, 'Нет проблем', '2023-01-10,2023-07-10', '2023-01-15,2023-07-15'),
(FALSE, TRUE, 'Аллергия на корм', '2023-02-20', '2023-02-25'),
(TRUE, FALSE, 'Повышенное давление', '2023-03-05,2023-09-05', '2023-03-10');

INSERT INTO users (fio, telephone_number, city, user_login, user_password, status, user_description)
VALUES
('Иванов Иван Иванович', '+79161234567', 'Москва', 'ivan123', 'password_hash_1', 'buyer', 'Любитель собак'),
('Петрова Мария Геннадьевна', '+79261234568', 'Санкт-Петербург', 'maria_p', 'password_hash_2', 'buyer', 'Ищет кошку'),
('Сидоров Сергей Владимирович', '+79361234569', 'Казань', 'sergey_s', 'password_hash_3', 'seller', 'Продает щенков');

INSERT INTO pet (vet_passport_id, seller_id, pet_name, species, pet_age, color, pet_gender, breed, pedigree, good_with_children, good_with_animals, pet_description, is_active, price)
VALUES
(1, 1, 'Барсик', 'Кошка', 2, 'Серый', 'Мальчик', 'Британская', TRUE, TRUE, TRUE, 'Очень ласковый', TRUE, 5000.00),
(2, 2, 'Мурка', 'Кошка', 1, 'Черный', 'Девочка', 'Сиамская', FALSE, TRUE, TRUE, 'Игривый характер', TRUE, 4000.00),
(3, 3, 'Шарик', 'Собака', 3, 'Коричневый', 'Мальчик', 'Лабрадор', TRUE, TRUE, TRUE, 'Активный щенок', TRUE, 15000.00);

INSERT INTO purchase_request (pet_id, seller_id, status, request_date)
VALUES
(1, 1, 'pending', '2026-03-08 10:00:00'),
(2, 2, 'approved', '2026-03-07 14:30:00'),
(3, 3, 'rejected', '2026-03-06 09:15:00');
