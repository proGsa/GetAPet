CREATE TABLE IF NOT EXISTS vet_passport (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    chipping BOOLEAN DEFAULT FALSE,
    sterilization BOOLEAN DEFAULT FALSE,
    health_issues TEXT,
    vaccinations TEXT,
    parasite_treatments TEXT
);

CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    fio VARCHAR(255) NOT NULL,
    telephone_number VARCHAR(20) NOT NULL,
    city VARCHAR(50),
    user_login VARCHAR(50) UNIQUE,
    user_password VARCHAR(255),
    status VARCHAR(20),
    user_description TEXT
);

CREATE TABLE IF NOT EXISTS pet (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    vet_passport_id UUID REFERENCES vet_passport(id),
    seller_id UUID REFERENCES users(id),
    pet_name VARCHAR(255) NOT NULL,
    species VARCHAR(50) NOT NULL,
    pet_age INT NOT NULL,
    color VARCHAR(50),
    pet_gender VARCHAR(20),
    breed VARCHAR(255),
    pedigree BOOLEAN DEFAULT FALSE,
    good_with_children BOOLEAN DEFAULT TRUE,
    good_with_animals BOOLEAN DEFAULT TRUE,
    pet_description TEXT,
    is_active BOOLEAN DEFAULT TRUE,
    price DECIMAL(10,2)
);

CREATE TABLE IF NOT EXISTS purchase_request (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    pet_id UUID NOT NULL REFERENCES pet(id),
    seller_id UUID NOT NULL REFERENCES users(id),
    status VARCHAR(50) DEFAULT 'pending',
    request_date TIMESTAMP DEFAULT NOW()
);
