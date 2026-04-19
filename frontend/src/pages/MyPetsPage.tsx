import { useEffect, useMemo, useState } from "react";
import { Link } from "react-router-dom";
import { petsApi } from "../api/pets";
import { vetPassportsApi } from "../api/vetPassports";
import {
  PetWithPassportForm,
  type PetWithPassportFormValues,
} from "../components/forms/PetWithPassportForm";
import { AlertMessage } from "../components/ui/AlertMessage";
import { EmptyState } from "../components/ui/EmptyState";
import { LoadingState } from "../components/ui/LoadingState";
import { useAuth } from "../hooks/useAuth";
import type { Pet, PetCreatePayload } from "../types/pet";
import type { VetPassport, VetPassportUpsertPayload } from "../types/vetPassport";
import { getErrorMessage } from "../utils/error";
import { formatPrice } from "../utils/format";
import {
  createEmptyPetCreatePayload,
  createEmptyVetPassportPayload,
  toPetUpdatePayload,
} from "../utils/pet";

const toPetCreatePayloadFromPet = (pet: Pet): PetCreatePayload => ({
  vet_passport_id: pet.vet_passport_id,
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
  price: pet.price,
});

const toInitialEditValue = (
  pet: Pet,
  passport: VetPassport | undefined,
): PetWithPassportFormValues => ({
  pet: toPetCreatePayloadFromPet(pet),
  passport: passport
    ? {
        chipping: passport.chipping,
        sterilization: passport.sterilization,
        health_issues: passport.health_issues,
        vaccinations: passport.vaccinations,
        parasite_treatments: passport.parasite_treatments,
      }
    : createEmptyVetPassportPayload(),
});

export function MyPetsPage() {
  const { user, token, mode } = useAuth();

  const [pets, setPets] = useState<Pet[]>([]);
  const [passportsMap, setPassportsMap] = useState<Record<string, VetPassport>>({});
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [message, setMessage] = useState<string | null>(null);
  const [isSubmitting, setIsSubmitting] = useState(false);
  const [editingPetId, setEditingPetId] = useState<string | null>(null);

  const isSellerMode = Boolean(user && mode === "seller");

  const myPets = useMemo(() => {
    if (!user) {
      return [];
    }

    return pets
      .filter((pet) => pet.seller_id === user.id)
      .sort((left, right) => Number(right.is_active) - Number(left.is_active));
  }, [pets, user]);

  const createInitialValue = useMemo<PetWithPassportFormValues>(
    () => ({
      pet: createEmptyPetCreatePayload(),
      passport: createEmptyVetPassportPayload(),
    }),
    [],
  );

  useEffect(() => {
    if (!user || !token || mode !== "seller") {
      setLoading(false);
      return;
    }

    const loadData = async () => {
      setLoading(true);
      setError(null);

      try {
        const [petsResponse, passportsResponse] = await Promise.all([
          petsApi.list(),
          vetPassportsApi.list(),
        ]);

        setPets(petsResponse);

        const nextPassportsMap = passportsResponse.reduce<Record<string, VetPassport>>(
          (accumulator, passport) => {
            accumulator[passport.id] = passport;
            return accumulator;
          },
          {},
        );
        setPassportsMap(nextPassportsMap);
      } catch (loadError) {
        setError(getErrorMessage(loadError, "Не удалось загрузить объявления продавца"));
      } finally {
        setLoading(false);
      }
    };

    void loadData();
  }, [mode, token, user]);

  if (!user || !token) {
    return (
      <section className="page-content narrow-page">
        <article className="panel">
          <h1>Мои объявления</h1>
          <p>
            Для доступа к этому разделу необходимо <Link to="/login">войти</Link>.
          </p>
        </article>
      </section>
    );
  }

  if (!isSellerMode) {
    return (
      <section className="page-content narrow-page">
        <article className="panel">
          <h1>Мои объявления</h1>
          <p>Этот раздел доступен в режиме продавца. Переключите режим в верхнем меню.</p>
        </article>
      </section>
    );
  }

  const handleCreate = async (payload: {
    pet: PetCreatePayload;
    passport: VetPassportUpsertPayload;
  }) => {
    if (!token) {
      return;
    }

    setIsSubmitting(true);
    setMessage(null);

    let createdPassportId: string | null = null;

    try {
      const createdPassport = await vetPassportsApi.create(payload.passport);
      createdPassportId = createdPassport.id;

      const createdPet = await petsApi.create(
        {
          ...payload.pet,
          vet_passport_id: createdPassport.id,
        },
        token,
      );

      setPets((current) => [createdPet, ...current]);
      setPassportsMap((current) => ({ ...current, [createdPassport.id]: createdPassport }));
      setMessage("Объявление создано.");
    } catch (createError) {
      if (createdPassportId) {
        try {
          await vetPassportsApi.remove(createdPassportId);
        } catch {
          // ignore rollback cleanup errors
        }
      }
      setMessage(getErrorMessage(createError, "Не удалось создать объявление"));
    } finally {
      setIsSubmitting(false);
    }
  };

  const handleUpdate = async (
    pet: Pet,
    payload: {
      pet: PetCreatePayload;
      passport: VetPassportUpsertPayload;
    },
  ) => {
    if (!token) {
      return;
    }

    setIsSubmitting(true);
    setMessage(null);

    try {
      const updatedPassport = await vetPassportsApi.update(pet.vet_passport_id, payload.passport);
      const updatedPet = await petsApi.update(pet.id, toPetUpdatePayload(payload.pet), token);

      setPassportsMap((current) => ({
        ...current,
        [updatedPassport.id]: updatedPassport,
      }));
      setPets((current) =>
        current.map((item) => {
          if (item.id !== pet.id) {
            return item;
          }

          return {
            ...item,
            ...updatedPet,
            seller_id: item.seller_id,
            vet_passport_id: item.vet_passport_id,
          };
        }),
      );
      setEditingPetId(null);
      setMessage("Объявление обновлено.");
    } catch (updateError) {
      setMessage(getErrorMessage(updateError, "Не удалось обновить объявление"));
    } finally {
      setIsSubmitting(false);
    }
  };

  const handleDelete = async (petId: string) => {
    if (!token) {
      return;
    }

    const shouldDelete = window.confirm("Удалить это объявление?");
    if (!shouldDelete) {
      return;
    }

    setIsSubmitting(true);
    setMessage(null);

    try {
      await petsApi.remove(petId, token);
      setPets((current) => current.filter((pet) => pet.id !== petId));
      setMessage("Объявление удалено.");
    } catch (deleteError) {
      setMessage(getErrorMessage(deleteError, "Не удалось удалить объявление"));
    } finally {
      setIsSubmitting(false);
    }
  };

  const handleToggleActive = async (pet: Pet) => {
    if (!token) {
      return;
    }

    setIsSubmitting(true);
    setMessage(null);

    try {
      if (pet.is_active) {
        const deactivated = await petsApi.update(pet.id, toPetUpdatePayload(toPetCreatePayloadFromPet(pet)), token);
        setPets((current) =>
          current.map((item) => {
            if (item.id !== pet.id) {
              return item;
            }

            return {
              ...item,
              ...deactivated,
              seller_id: item.seller_id,
              vet_passport_id: item.vet_passport_id,
              is_active: false,
            };
          }),
        );
        setMessage("Объявление переведено в неактивные.");
      } else {
        // Backend currently does not expose direct "activate" in update payload.
        // Recreate ad as active, then remove archived one.
        const recreated = await petsApi.create(toPetCreatePayloadFromPet(pet), token);
        await petsApi.remove(pet.id, token);
        setPets((current) =>
          current.map((item) => {
            if (item.id !== pet.id) {
              return item;
            }

            return {
              ...item,
              ...recreated,
              seller_id: recreated.seller_id || item.seller_id,
              vet_passport_id: recreated.vet_passport_id || item.vet_passport_id,
              is_active: true,
            };
          }),
        );
        setEditingPetId((current) => (current === pet.id ? null : current));
        setMessage("Объявление снова активно.");
      }
    } catch (toggleError) {
      setMessage(getErrorMessage(toggleError, "Не удалось изменить активность объявления"));
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <section className="page-content">
      <div className="page-title-row">
        <h1>Мои объявления</h1>
      </div>

      {loading ? <LoadingState label="Загрузка объявлений..." /> : null}
      {error ? <AlertMessage variant="error">{error}</AlertMessage> : null}
      {message ? <AlertMessage>{message}</AlertMessage> : null}

      <article className="panel">
        <h2>Создать объявление</h2>
        <PetWithPassportForm
          initialValue={createInitialValue}
          submitLabel="Создать объявление"
          isSubmitting={isSubmitting}
          onSubmit={handleCreate}
        />
      </article>

      {!loading && !error && myPets.length === 0 ? (
        <EmptyState
          title="Пока нет объявлений"
          description="Создайте первое объявление через форму выше."
        />
      ) : null}

      <div className="stack-list">
        {myPets.map((pet) => {
          const isEditing = editingPetId === pet.id;
          const editValue = toInitialEditValue(pet, passportsMap[pet.vet_passport_id]);
          return (
            <article key={pet.id} className="panel compact-panel">
              <div className="panel-header-row">
                <div>
                  <h2>{pet.pet_name}</h2>
                  <p>
                    {pet.species} • {formatPrice(pet.price)}
                  </p>
                </div>

                <div className="button-row">
                  <button
                    type="button"
                    className="secondary-button"
                    onClick={() => setEditingPetId((current) => (current === pet.id ? null : pet.id))}
                  >
                    {isEditing ? "Закрыть" : "Редактировать"}
                  </button>
                  <button
                    type="button"
                    className="secondary-button"
                    disabled={isSubmitting}
                    onClick={() => {
                      void handleToggleActive(pet);
                    }}
                  >
                    {pet.is_active ? "Сделать неактивным" : "Сделать активным"}
                  </button>
                  <button
                    type="button"
                    className="danger-button"
                    disabled={isSubmitting}
                    onClick={() => {
                      void handleDelete(pet.id);
                    }}
                  >
                    Удалить
                  </button>
                </div>
              </div>

              <p className="text-block">{pet.pet_description || "Описание отсутствует"}</p>

              {isEditing ? (
                <PetWithPassportForm
                  initialValue={editValue}
                  submitLabel="Сохранить объявление"
                  isSubmitting={isSubmitting}
                  onSubmit={(payload) => handleUpdate(pet, payload)}
                />
              ) : null}
            </article>
          );
        })}
      </div>
    </section>
  );
}
