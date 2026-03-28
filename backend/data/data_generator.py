import csv
import uuid
import random
import bcrypt
from faker import Faker
from datetime import datetime, timedelta

fake = Faker("ru_RU")
NUM_VET_PASSPORTS = 20
NUM_USERS = 10
NUM_PETS = 25
NUM_PURCHASE_REQUESTS = 30

GO_BCRYPT_COST = 10

OUTPUT_FILES = {
    "vet_passport": "vet_passport.csv",
    "users": "users.csv",
    "pet": "pet.csv",
    "purchase_request": "purchase_request.csv",
    "password_mapping": "password_mapping.csv",
}
CITIES = [
    "Москва", "Санкт-Петербург", "Казань", "Екатеринбург", "Новосибирск",
    "Самара", "Нижний Новгород", "Ростов-на-Дону", "Уфа", "Краснодар"
]

USER_STATUSES = ["buyer", "seller", "buyer", "seller", "buyer"]

SPECIES_BREEDS = {
    "Кошка": [
        "Британская", "Сиамская", "Мейн-кун", "Сфинкс",
        "Шотландская вислоухая", "Персидская", "Бенгальская"
    ],
    "Собака": [
        "Лабрадор", "Шпиц", "Хаски", "Корги",
        "Немецкая овчарка", "Йоркширский терьер", "Пудель"
    ],
    "Попугай": [
        "Волнистый попугай", "Корелла", "Жако", "Неразлучник"
    ],
    "Кролик": [
        "Декоративный кролик", "Карликовый рекс", "Львиноголовый"
    ],
    "Хомяк": [
        "Сирийский", "Джунгарский", "Кэмпбелла"
    ]
}

PET_COLORS = [
    "Белый", "Черный", "Серый", "Рыжий", "Коричневый",
    "Бежевый", "Пятнистый", "Полосатый", "Золотистый"
]

PET_GENDERS = ["Мальчик", "Девочка"]

HEALTH_ISSUES_OPTIONS = [
    "Нет проблем",
    "Аллергия на корм",
    "Чувствительное пищеварение",
    "Проблемы с суставами",
    "Повышенное давление",
    "Перенесенная операция",
    "Нет хронических заболеваний"
]

USER_DESCRIPTIONS = [
    "Любит животных и заботится о них",
    "Ищет питомца для семьи",
    "Занимается разведением животных",
    "Продает питомцев с документами",
    "Ответственный владелец",
    "Любитель кошек",
    "Любитель собак",
    "Ищет домашнего любимца"
]

PET_DESCRIPTIONS = [
    "Очень ласковый и дружелюбный",
    "Активный и игривый",
    "Спокойный характер",
    "Хорошо идет на контакт",
    "Любит детей и внимание",
    "Приучен к лотку/пеленке",
    "Здоровый и ухоженный питомец",
    "Подойдет для семьи"
]

REQUEST_STATUSES = ["pending", "approved", "rejected"]
vet_passports = []
users = []
pets = []
purchase_requests = []
password_mapping = []
def generate_date_list(count_min=1, count_max=3, start_year=2022, end_year=2026):
    dates = []
    for _ in range(random.randint(count_min, count_max)):
        year = random.randint(start_year, end_year)
        month = random.randint(1, 12)
        day = random.randint(1, 28)
        dates.append(f"{year:04d}-{month:02d}-{day:02d}")
    return ",".join(sorted(set(dates)))


def generate_phone_number():
    return f"+79{random.randint(100000000, 999999999)}"


def generate_price(species):
    if species == "Кошка":
        return round(random.uniform(3000, 35000), 2)
    if species == "Собака":
        return round(random.uniform(5000, 80000), 2)
    if species == "Попугай":
        return round(random.uniform(1500, 25000), 2)
    if species == "Кролик":
        return round(random.uniform(1000, 12000), 2)
    return round(random.uniform(500, 7000), 2)


def save_csv(filename, fieldnames, rows):
    with open(filename, "w", newline="", encoding="utf-8-sig") as f:
        writer = csv.DictWriter(f, fieldnames=fieldnames, delimiter=",")
        writer.writeheader()
        writer.writerows(rows)


def hash_password_like_go(raw_password: str, cost: int = GO_BCRYPT_COST) -> str:
    password_bytes = raw_password.encode("utf-8")
    if len(password_bytes) > 72:
        raise ValueError("Пароль длиннее 72 байт, bcrypt в Go такой пароль не принимает.")

    hashed = bcrypt.hashpw(
        password_bytes,
        bcrypt.gensalt(rounds=cost)
    )
    return hashed.decode("utf-8")
def create_vet_passports():
    for _ in range(NUM_VET_PASSPORTS):
        vet_passports.append({
            "id": str(uuid.uuid4()),
            "chipping": random.choice([True, False]),
            "sterilization": random.choice([True, False]),
            "health_issues": random.choice(HEALTH_ISSUES_OPTIONS),
            "vaccinations": generate_date_list(1, 3, 2022, 2026),
            "parasite_treatments": generate_date_list(1, 3, 2022, 2026),
        })


def create_users():
    used_logins = set()

    for _ in range(NUM_USERS):
        user_id = str(uuid.uuid4())

        login = fake.user_name()
        while login in used_logins:
            login = fake.user_name()
        used_logins.add(login)
        raw_password = fake.password(
            length=12,
            special_chars=True,
            digits=True,
            upper_case=True,
            lower_case=True
        )

        hashed_password = hash_password_like_go(raw_password, GO_BCRYPT_COST)

        user = {
            "id": user_id,
            "fio": fake.name(),
            "telephone_number": generate_phone_number(),
            "city": random.choice(CITIES),
            "user_login": login,
            "user_password": hashed_password,
            "status": random.choice(USER_STATUSES),
            "user_description": random.choice(USER_DESCRIPTIONS),
        }
        users.append(user)

        password_mapping.append({
            "user_id": user_id,
            "user_login": login,
            "raw_password": raw_password,
            "hashed_password": hashed_password,
            "bcrypt_cost": GO_BCRYPT_COST,
        })


def create_pets():
    seller_ids = [user["id"] for user in users if user["status"] == "seller"]
    if not seller_ids:
        raise ValueError("Нет пользователей со статусом seller для генерации pet.")

    used_vet_passport_ids = set()

    for _ in range(NUM_PETS):
        available_passports = [
            vp["id"] for vp in vet_passports
            if vp["id"] not in used_vet_passport_ids
        ]
        vet_passport_id = (
            random.choice(available_passports)
            if available_passports
            else random.choice(vet_passports)["id"]
        )
        used_vet_passport_ids.add(vet_passport_id)

        species = random.choice(list(SPECIES_BREEDS.keys()))
        breed = random.choice(SPECIES_BREEDS[species])

        if species in ["Кошка", "Собака", "Кролик"]:
            good_with_children = random.choice([True, True, True, False])
            good_with_animals = random.choice([True, True, False])
        else:
            good_with_children = random.choice([True, True, False])
            good_with_animals = random.choice([True, False])

        pets.append({
            "id": str(uuid.uuid4()),
            "vet_passport_id": vet_passport_id,
            "seller_id": random.choice(seller_ids),
            "pet_name": fake.first_name(),
            "species": species,
            "pet_age": random.randint(1, 15),
            "color": random.choice(PET_COLORS),
            "pet_gender": random.choice(PET_GENDERS),
            "breed": breed,
            "pedigree": random.choice([True, False]),
            "good_with_children": good_with_children,
            "good_with_animals": good_with_animals,
            "pet_description": random.choice(PET_DESCRIPTIONS),
            "is_active": random.choice([True, True, True, False]),
            "price": generate_price(species),
        })


def create_purchase_requests():
    if not pets:
        raise ValueError("Сначала нужно создать pets.")

    for _ in range(NUM_PURCHASE_REQUESTS):
        chosen_pet = random.choice(pets)

        request_date = datetime.now() - timedelta(
            days=random.randint(0, 60),
            hours=random.randint(0, 23),
            minutes=random.randint(0, 59)
        )

        purchase_requests.append({
            "id": str(uuid.uuid4()),
            "pet_id": chosen_pet["id"],
            "seller_id": chosen_pet["seller_id"],
            "status": random.choices(
                REQUEST_STATUSES,
                weights=[60, 25, 15],
                k=1
            )[0],
            "request_date": request_date.strftime("%Y-%m-%d %H:%M:%S"),
        })
def save_all_to_csv():
    save_csv(
        OUTPUT_FILES["vet_passport"],
        ["id", "chipping", "sterilization", "health_issues", "vaccinations", "parasite_treatments"],
        vet_passports
    )

    save_csv(
        OUTPUT_FILES["users"],
        ["id", "fio", "telephone_number", "city", "user_login", "user_password", "status", "user_description"],
        users
    )

    save_csv(
        OUTPUT_FILES["pet"],
        [
            "id", "vet_passport_id", "seller_id", "pet_name", "species", "pet_age",
            "color", "pet_gender", "breed", "pedigree", "good_with_children",
            "good_with_animals", "pet_description", "is_active", "price"
        ],
        pets
    )

    save_csv(
        OUTPUT_FILES["purchase_request"],
        ["id", "pet_id", "seller_id", "status", "request_date"],
        purchase_requests
    )

    save_csv(
        OUTPUT_FILES["password_mapping"],
        ["user_id", "user_login", "raw_password", "hashed_password", "bcrypt_cost"],
        password_mapping
    )


def main():
    create_vet_passports()
    create_users()
    create_pets()
    create_purchase_requests()
    save_all_to_csv()

    print("CSV-файлы успешно созданы:")
    for _, filename in OUTPUT_FILES.items():
        print(f"- {filename}")


if __name__ == "__main__":
    main()

