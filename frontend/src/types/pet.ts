export interface Pet {
  id: string;
  vet_passport_id: string;
  seller_id: string;
  pet_name: string;
  species: string;
  pet_age: number;
  color: string;
  pet_gender: string;
  breed: string;
  pedigree: boolean;
  good_with_children: boolean;
  good_with_animals: boolean;
  pet_description: string;
  is_active: boolean;
  price: number;
}

export interface PetCreatePayload {
  vet_passport_id: string;
  pet_name: string;
  species: string;
  pet_age: number;
  color: string;
  pet_gender: string;
  breed: string;
  pedigree: boolean;
  good_with_children: boolean;
  good_with_animals: boolean;
  pet_description: string;
  price: number;
}

export interface PetUpdatePayload {
  pet_name: string;
  species: string;
  pet_age: number;
  color: string;
  pet_gender: string;
  breed: string;
  pedigree: boolean;
  good_with_children: boolean;
  good_with_animals: boolean;
  pet_description: string;
  price: number;
}
