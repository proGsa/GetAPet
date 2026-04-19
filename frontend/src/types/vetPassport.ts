export interface VetPassport {
  id: string;
  chipping: boolean;
  sterilization: boolean;
  health_issues: string;
  vaccinations: string;
  parasite_treatments: string;
}

export interface VetPassportUpsertPayload {
  chipping: boolean;
  sterilization: boolean;
  health_issues: string;
  vaccinations: string;
  parasite_treatments: string;
}
