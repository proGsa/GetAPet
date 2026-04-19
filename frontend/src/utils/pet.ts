import type { PetCreatePayload, PetUpdatePayload } from "../types/pet";
import type { VetPassportUpsertPayload } from "../types/vetPassport";

export const createEmptyPetCreatePayload = (): PetCreatePayload => ({
  vet_passport_id: "",
  pet_name: "",
  species: "",
  pet_age: 0,
  color: "",
  pet_gender: "male",
  breed: "",
  pedigree: false,
  good_with_children: true,
  good_with_animals: true,
  pet_description: "",
  price: 0,
});

export const toPetUpdatePayload = (payload: PetCreatePayload): PetUpdatePayload => ({
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
  price: payload.price,
});

export const createEmptyVetPassportPayload = (): VetPassportUpsertPayload => ({
  chipping: false,
  sterilization: false,
  health_issues: "",
  vaccinations: "",
  parasite_treatments: "",
});
