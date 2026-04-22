import { useEffect, useMemo, useState } from "react";
import { Link, useParams } from "react-router-dom";
import { petsApi } from "../api/pets";
import { purchaseRequestsApi } from "../api/purchaseRequests";
import { usersApi } from "../api/users";
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
import type { User } from "../types/user";
import type { VetPassport, VetPassportUpsertPayload } from "../types/vetPassport";
import { getErrorMessage } from "../utils/error";
import { formatPrice, shortId } from "../utils/format";
import {
  createEmptyVetPassportPayload,
  normalizePetGender,
  toPetUpdatePayload,
} from "../utils/pet";

const toPetCreatePayloadFromPet = (pet: Pet): PetCreatePayload => ({
  vet_passport_id: pet.vet_passport_id,
  pet_name: pet.pet_name,
  species: pet.species,
  pet_age: pet.pet_age,
  color: pet.color,
  pet_gender: normalizePetGender(pet.pet_gender),
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

export function PetPage() {
  const { id } = useParams();
  const { user, token, mode } = useAuth();

  const [pet, setPet] = useState<Pet | null>(null);
  const [seller, setSeller] = useState<User | null>(null);
  const [passport, setPassport] = useState<VetPassport | null>(null);

  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const [isSubmitting, setIsSubmitting] = useState(false);
  const [actionMessage, setActionMessage] = useState<string | null>(null);
  const [isEditing, setIsEditing] = useState(false);
  const [hasMyRequest, setHasMyRequest] = useState(false);

  useEffect(() => {
    if (!id) {
      setError("Некорректный идентификатор объявления");
      setLoading(false);
      return;
    }

    const loadData = async () => {
      setLoading(true);
      setError(null);

      try {
        const loadedPet = await petsApi.getById(id);
        setPet(loadedPet);

        if (token) {
          try {
            const loadedSeller = await usersApi.getById(loadedPet.seller_id, token);
            setSeller(loadedSeller);
          } catch {
            setSeller(null);
          }
        } else {
          setSeller(null);
        }

        if (loadedPet.vet_passport_id) {
          try {
            const loadedPassport = await vetPassportsApi.getById(loadedPet.vet_passport_id);
            setPassport(loadedPassport);
          } catch {
            setPassport(null);
          }
        } else {
          setPassport(null);
        }
      } catch (loadError) {
        setError(getErrorMessage(loadError, "Не удалось загрузить карточку питомца"));
      } finally {
        setLoading(false);
      }
    };

    void loadData();
  }, [id, token]);

  useEffect(() => {
    if (!token || !user || !pet || mode !== "buyer" || user.id === pet.seller_id || !pet.is_active) {
      setHasMyRequest(false);
      return;
    }

    const loadMyRequestState = async () => {
      try {
        const requests = await purchaseRequestsApi.listByBuyer(user.id, token);
        setHasMyRequest(requests.some((request) => request.pet_id === pet.id));
      } catch {
        setHasMyRequest(false);
      }
    };

    void loadMyRequestState();
  }, [mode, pet, token, user]);

  const canManagePet = useMemo(
    () => Boolean(user && token && pet && mode === "seller" && user.id === pet.seller_id),
    [mode, pet, token, user],
  );

  const canCreateRequest = useMemo(
    () =>
      Boolean(
        user &&
          token &&
          pet &&
          mode === "buyer" &&
          user.id !== pet.seller_id &&
          pet.is_active &&
          !hasMyRequest,
      ),
    [hasMyRequest, mode, pet, token, user],
  );

  const canViewRequestStatus = useMemo(
    () =>
      Boolean(
        user &&
          token &&
          pet &&
          mode === "buyer" &&
          user.id !== pet.seller_id &&
          pet.is_active &&
          hasMyRequest,
      ),
    [hasMyRequest, mode, pet, token, user],
  );

  const shouldShowInactiveHint = useMemo(
    () => Boolean(user && pet && mode === "buyer" && user.id !== pet.seller_id && !pet.is_active),
    [mode, pet, user],
  );

  const hasActionButtons = useMemo(
    () => canCreateRequest || canViewRequestStatus || canManagePet,
    [canCreateRequest, canManagePet, canViewRequestStatus],
  );

  const noActionsHint = useMemo(() => {
    if (!token || !pet || !user || hasActionButtons || shouldShowInactiveHint) {
      return null;
    }

    if (mode === "seller" && user.id !== pet.seller_id) {
      return "Вы находитесь в режиме продавца, переключите в режим покупателя, чтобы оставить заявку.";
    }

    if (mode === "buyer" && user.id === pet.seller_id) {
      return "Это Ваше объявление. Переключитесь в режим продавца для управления.";
    }

    return "Для этого объявления в текущем режиме нет доступных действий.";
  }, [hasActionButtons, mode, pet, shouldShowInactiveHint, token, user]);

  const editValue = useMemo(() => {
    if (!pet) {
      return null;
    }

    return toInitialEditValue(pet, passport ?? undefined);
  }, [passport, pet]);

  const handleCreateRequest = async () => {
    if (!pet || !user || !token) {
      return;
    }

    setActionMessage(null);
    setIsSubmitting(true);

    try {
      await purchaseRequestsApi.create(
        {
          pet_id: pet.id,
          buyer_id: user.id,
        },
        token,
      );
      setHasMyRequest(true);
      setActionMessage("Заявка успешно отправлена.");
    } catch (createError) {
      setActionMessage(getErrorMessage(createError, "Не удалось отправить заявку"));
    } finally {
      setIsSubmitting(false);
    }
  };

  const handleUpdate = async (payload: {
    pet: PetCreatePayload;
    passport: VetPassportUpsertPayload;
  }) => {
    if (!pet || !token) {
      return;
    }

    setActionMessage(null);
    setIsSubmitting(true);

    try {
      if (pet.vet_passport_id) {
        const updatedPassport = await vetPassportsApi.update(pet.vet_passport_id, payload.passport);
        setPassport(updatedPassport);
      }

      const nextPetPayload: PetCreatePayload = {
        ...payload.pet,
        vet_passport_id: pet.vet_passport_id,
      };
      const updatedPet = await petsApi.update(
        pet.id,
        toPetUpdatePayload(nextPetPayload),
        token,
      );
      setPet((current) => {
        if (!current) {
          return current;
        }

        return {
          ...current,
          ...updatedPet,
          seller_id: current.seller_id,
          vet_passport_id: current.vet_passport_id,
          is_active: current.is_active,
        };
      });

      setIsEditing(false);
      setActionMessage("Объявление обновлено.");
    } catch (updateError) {
      setActionMessage(getErrorMessage(updateError, "Не удалось обновить объявление"));
    } finally {
      setIsSubmitting(false);
    }
  };

  const handleToggleActive = async () => {
    if (!pet || !token) {
      return;
    }

    setActionMessage(null);
    setIsSubmitting(true);

    try {
      if (pet.is_active) {
        const deactivated = await petsApi.update(
          pet.id,
          toPetUpdatePayload(toPetCreatePayloadFromPet(pet)),
          token,
        );
        setPet((current) => {
          if (!current) {
            return current;
          }

          return {
            ...current,
            ...deactivated,
            seller_id: current.seller_id,
            vet_passport_id: current.vet_passport_id,
            is_active: false,
          };
        });
        setActionMessage("Объявление переведено в неактивные.");
      } else {
        const recreated = await petsApi.create(toPetCreatePayloadFromPet(pet), token);
        await petsApi.remove(pet.id, token);
        setPet((current) => {
          if (!current) {
            return current;
          }

          return {
            ...current,
            ...recreated,
            seller_id: recreated.seller_id || current.seller_id,
            vet_passport_id: recreated.vet_passport_id || current.vet_passport_id,
            is_active: true,
          };
        });

        setActionMessage("Объявление снова активно.");
      }
    } catch (toggleError) {
      setActionMessage(getErrorMessage(toggleError, "Не удалось изменить активность объявления"));
    } finally {
      setIsSubmitting(false);
    }
  };

  return (
    <section className="page-content">
      {loading ? <LoadingState label="Загрузка карточки питомца..." /> : null}
      {error ? <AlertMessage variant="error">{error}</AlertMessage> : null}

      {!loading && !error && !pet ? <EmptyState title="Питомец не найден" /> : null}

      {!loading && !error && pet ? (
        <>
          {actionMessage ? <AlertMessage variant="info">{actionMessage}</AlertMessage> : null}

          {!token ? (
            <AlertMessage>
              Чтобы выполнить действие, <Link to="/login">войдите</Link> в аккаунт.
            </AlertMessage>
          ) : null}

          {shouldShowInactiveHint ? (
            <AlertMessage variant="info">На неактивное объявление нельзя оставить заявку.</AlertMessage>
          ) : null}

          {noActionsHint ? <AlertMessage variant="info">{noActionsHint}</AlertMessage> : null}

          <div className="panel-header-row">
            <div className="page-title-row">
              <h1>{pet.pet_name}</h1>
              <p>
                {pet.species || "Вид не указан"}, {pet.breed || "Порода не указана"}
              </p>
              <p>
                <span className={`status-badge ${pet.is_active ? "status-active" : "status-inactive"}`}>
                  {pet.is_active ? "Объявление активно" : "Объявление неактивно"}
                </span>
              </p>
            </div>

            <div className="pet-actions pet-details-actions">
              {canCreateRequest ? (
                <button
                  type="button"
                  className="primary-button inline-button"
                  disabled={isSubmitting}
                  onClick={() => {
                    void handleCreateRequest();
                  }}
                >
                  {isSubmitting ? "Отправка..." : "Оставить заявку"}
                </button>
              ) : null}

              {canViewRequestStatus ? (
                <button type="button" className="secondary-button inline-button">
                  Посмотреть статус заявки
                </button>
              ) : null}

              {canManagePet ? (
                <>
                  <button
                    type="button"
                    className="secondary-button inline-button"
                    disabled={isSubmitting}
                    onClick={() => {
                      void handleToggleActive();
                    }}
                  >
                    {pet.is_active ? "Сделать неактивным" : "Сделать активным"}
                  </button>
                  <button
                    type="button"
                    className="secondary-button inline-button"
                    onClick={() => setIsEditing((current) => !current)}
                  >
                    {isEditing ? "Закрыть редактирование" : "Редактировать"}
                  </button>
                </>
              ) : null}
            </div>
          </div>

          {canManagePet && isEditing && editValue ? (
            <article className="panel">
              <h2>Редактирование объявления</h2>
              <PetWithPassportForm
                initialValue={editValue}
                submitLabel="Сохранить объявление"
                isSubmitting={isSubmitting}
                onSubmit={handleUpdate}
              />
            </article>
          ) : (
            <article className="panel">
              <h2>Информация о питомце</h2>
              <dl className="details-grid">
                <div>
                  <dt>Возраст</dt>
                  <dd>{pet.pet_age}</dd>
                </div>
                <div>
                  <dt>Окрас</dt>
                  <dd>{pet.color || "Не указан"}</dd>
                </div>
                <div>
                  <dt>Пол</dt>
                  <dd>{pet.pet_gender || "Не указан"}</dd>
                </div>
                <div>
                  <dt>Цена</dt>
                  <dd>{formatPrice(pet.price)}</dd>
                </div>
                <div>
                  <dt>Родословная</dt>
                  <dd>{pet.pedigree ? "Да" : "Нет"}</dd>
                </div>
                <div>
                  <dt>Ладит с детьми</dt>
                  <dd>{pet.good_with_children ? "Да" : "Нет"}</dd>
                </div>
                <div>
                  <dt>Ладит с животными</dt>
                  <dd>{pet.good_with_animals ? "Да" : "Нет"}</dd>
                </div>
              </dl>

              <p className="text-block">{pet.pet_description || "Описание отсутствует"}</p>

              {passport ? (
                <div className="sub-panel">
                  <h3>Ветпаспорт</h3>
                  <p>Чипирование: {passport.chipping ? "да" : "нет"}</p>
                  <p>Стерилизация: {passport.sterilization ? "да" : "нет"}</p>
                  <p>Проблемы со здоровьем: {passport.health_issues || "нет"}</p>
                  <p>Вакцинации: {passport.vaccinations || "не указаны"}</p>
                  <p>Обработки от паразитов: {passport.parasite_treatments || "не указаны"}</p>
                </div>
              ) : null}

              <div className="sub-panel">
                <h3>Продавец</h3>
                {seller ? (
                  <>
                    <p>{seller.fio}</p>
                    <p>{seller.telephone_number || "Телефон не указан"}</p>
                    <p>{seller.city || "Город не указан"}</p>
                  </>
                ) : (
                  <p>{`ID продавца: ${shortId(pet.seller_id)}`}</p>
                )}
              </div>
            </article>
          )}
        </>
      ) : null}
    </section>
  );
}
