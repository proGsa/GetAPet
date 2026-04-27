import type { Pet, PetCreatePayload, PetUpdatePayload } from "../types/pet";
import type { VetPassportUpsertPayload } from "../types/vetPassport";

export const createEmptyPetCreatePayload = (): PetCreatePayload => ({
  vet_passport_id: "",
  pet_name: "",
  species: "",
  pet_age: 0,
  color: "",
  pet_gender: "Не указан",
  breed: "",
  pedigree: false,
  good_with_children: true,
  good_with_animals: true,
  pet_description: "",
  price: 0,
});

export const normalizePetGender = (value: string): string => {
  if (!value) {
    return "Не указан";
  }

  if (value === "Мальчик" || value === "Девочка" || value === "Не указан") {
    return value;
  }

  return "Не указан";
};

export const toPetUpdatePayload = (payload: PetCreatePayload, isActive: boolean): PetUpdatePayload => ({
  pet_name: payload.pet_name,
  species: payload.species,
  pet_age: payload.pet_age,
  color: payload.color,
  pet_gender: payload.pet_gender,
  breed: payload.breed,
  pedigree: payload.pedigree,
  good_with_children: payload.good_with_children,
  good_with_animals: payload.good_with_animals,
  pet_description: payload.pet_description,
  is_active: isActive,
  price: payload.price,
});

export const toPetUpdatePayloadFromPet = (pet: Pet, isActive: boolean): PetUpdatePayload => ({
  pet_name: pet.pet_name,
  species: pet.species,
  pet_age: pet.pet_age,
  color: pet.color,
  pet_gender: pet.pet_gender,
  breed: pet.breed,
  pedigree: pet.pedigree,
  good_with_children: pet.good_with_children,
  good_with_animals: pet.good_with_animals,
  pet_description: pet.pet_description,
  is_active: isActive,
  price: pet.price,
});

export const createEmptyVetPassportPayload = (): VetPassportUpsertPayload => ({
  chipping: false,
  sterilization: false,
  health_issues: "",
  vaccinations: "",
  parasite_treatments: "",
});
