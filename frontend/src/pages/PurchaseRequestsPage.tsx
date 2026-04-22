import { useCallback, useEffect, useMemo, useState } from "react";
import { Link, useSearchParams } from "react-router-dom";
import { petsApi } from "../api/pets";
import { purchaseRequestsApi } from "../api/purchaseRequests";
import { usersApi } from "../api/users";
import { AlertMessage } from "../components/ui/AlertMessage";
import { EmptyState } from "../components/ui/EmptyState";
import { LoadingState } from "../components/ui/LoadingState";
import { useAuth } from "../hooks/useAuth";
import type { Pet } from "../types/pet";
import type { PurchaseRequest } from "../types/purchaseRequest";
import type { User } from "../types/user";
import { getErrorMessage } from "../utils/error";
import { formatDateTime, shortId } from "../utils/format";

const DEFAULT_FILTER = "all";
const MODERATED_STATUSES = new Set(["pending", "approved", "rejected"]);

interface RequestViewModel {
  request: PurchaseRequest;
  pet: Pet | null;
  buyer: User | null;
  seller: User | null;
}

const toPetsMap = (pets: Pet[]): Record<string, Pet> =>
  pets.reduce<Record<string, Pet>>((accumulator, currentPet) => {
    accumulator[currentPet.id] = currentPet;
    return accumulator;
  }, {});

const toUsersMap = (users: User[]): Record<string, User> =>
  users.reduce<Record<string, User>>((accumulator, currentUser) => {
    accumulator[currentUser.id] = currentUser;
    return accumulator;
  }, {});

const normalizeStatus = (status: string): string => status.trim().toLowerCase();

const isKnownStatus = (status: string): boolean =>
  status === "pending" || status === "approved" || status === "rejected";

const statusBadgeClassName = (status: string): string => {
  const normalized = normalizeStatus(status);
  return isKnownStatus(normalized)
    ? `status-badge status-${normalized}`
    : "status-badge";
};

const requestStatusText = (status: string): string => {
  const normalized = normalizeStatus(status);
  if (normalized === "pending") {
    return "На рассмотрении";
  }

  if (normalized === "approved") {
    return "Одобрено";
  }

  if (normalized === "rejected") {
    return "Отказано";
  }

  return status;
};

const petLabel = (pet: Pet | null, petId: string): string => {
  const name = pet?.pet_name?.trim() || `Питомец ${shortId(petId)}`;
  const species = pet?.species?.trim() || "вид не указан";
  return `${name} (${species})`;
};

const userLabel = (user: User | null, userId: string, fallbackTitle: string): string => {
  const fio = user?.fio?.trim();
  if (fio) {
    return fio;
  }

  return `${fallbackTitle} ${shortId(userId)}`;
};

const userPhoneLabel = (user: User | null): string => {
  const phone = user?.telephone_number?.trim();
  return phone ? phone : "Телефон не указан";
};

export function PurchaseRequestsPage() {
  const { user, token, mode } = useAuth();
  const [searchParams] = useSearchParams();

  const [requests, setRequests] = useState<PurchaseRequest[]>([]);
  const [petsMap, setPetsMap] = useState<Record<string, Pet>>({});
  const [usersMap, setUsersMap] = useState<Record<string, User>>({});

  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [message, setMessage] = useState<string | null>(null);
  const [activeRequestId, setActiveRequestId] = useState<string | null>(null);

  const initialStatusFilter = useMemo(() => {
    const value = normalizeStatus(searchParams.get("status") ?? DEFAULT_FILTER);
    return value === DEFAULT_FILTER || isKnownStatus(value) ? value : DEFAULT_FILTER;
  }, [searchParams]);

  const initialPetFilter = useMemo(() => searchParams.get("pet") ?? DEFAULT_FILTER, [searchParams]);
  const initialSellerQuery = useMemo(() => searchParams.get("seller") ?? "", [searchParams]);

  const [statusFilter, setStatusFilter] = useState(initialStatusFilter);
  const [petFilter, setPetFilter] = useState(initialPetFilter);
  const [sellerQuery, setSellerQuery] = useState(initialSellerQuery);

  const loadData = useCallback(async () => {
    if (!user || !token) {
      setRequests([]);
      setPetsMap({});
      setUsersMap({});
      setLoading(false);
      return;
    }

    setLoading(true);
    setError(null);

    try {
      const requestPromise =
        mode === "seller"
          ? purchaseRequestsApi.listBySeller(user.id, token)
          : purchaseRequestsApi.listByBuyer(user.id, token);

      const [loadedRequests, loadedPets, loadedUsers] = await Promise.all([
        requestPromise,
        petsApi.list(),
        usersApi.list(token).catch(() => [] as User[]),
      ]);

      setRequests(loadedRequests);
      setPetsMap(toPetsMap(loadedPets));
      setUsersMap(toUsersMap(loadedUsers));
    } catch (loadError) {
      setError(getErrorMessage(loadError, "Не удалось загрузить заявки"));
    } finally {
      setLoading(false);
    }
  }, [mode, token, user]);

  useEffect(() => {
    void loadData();
  }, [loadData]);

  useEffect(() => {
    setMessage(null);
  }, [mode]);

  const viewModels = useMemo<RequestViewModel[]>(
    () =>
      requests.map((request) => {
        const pet = petsMap[request.pet_id] ?? null;
        const buyer =
          usersMap[request.buyer_id] ??
          (user && request.buyer_id === user.id ? user : null);
        const seller =
          usersMap[request.seller_id] ??
          (user && request.seller_id === user.id ? user : null);

        return {
          request,
          pet,
          buyer,
          seller,
        };
      }),
    [petsMap, requests, user, usersMap],
  );

  const petOptions = useMemo(() => {
    const options = new Map<string, string>();

    viewModels.forEach((item) => {
      if (!options.has(item.request.pet_id)) {
        options.set(item.request.pet_id, petLabel(item.pet, item.request.pet_id));
      }
    });

    return Array.from(options.entries()).sort((left, right) => left[1].localeCompare(right[1], "ru"));
  }, [viewModels]);

  const normalizedSellerQuery = sellerQuery.trim().toLowerCase();

  const filteredRequests = useMemo(
    () =>
      viewModels.filter((item) => {
        const requestStatus = normalizeStatus(item.request.status);
        if (statusFilter !== DEFAULT_FILTER && requestStatus !== statusFilter) {
          return false;
        }

        if (petFilter !== DEFAULT_FILTER && item.request.pet_id !== petFilter) {
          return false;
        }

        if (mode !== "buyer" || !normalizedSellerQuery) {
          return true;
        }

        const sellerFio = item.seller?.fio?.toLowerCase() ?? "";
        return sellerFio.includes(normalizedSellerQuery);
      }),
    [mode, normalizedSellerQuery, petFilter, statusFilter, viewModels],
  );

  const handleDelete = async (requestId: string) => {
    if (!token) {
      return;
    }

    if (!window.confirm("Удалить эту заявку?")) {
      return;
    }

    setMessage(null);
    setActiveRequestId(requestId);

    try {
      await purchaseRequestsApi.remove(requestId, token);
      setRequests((current) => current.filter((item) => item.id !== requestId));
      setMessage("Заявка удалена.");
    } catch (deleteError) {
      setMessage(getErrorMessage(deleteError, "Не удалось удалить заявку"));
    } finally {
      setActiveRequestId(null);
    }
  };

  const handleStatusChange = async (requestId: string, status: "approved" | "rejected") => {
    if (!token) {
      return;
    }

    setMessage(null);
    setActiveRequestId(requestId);

    try {
      await purchaseRequestsApi.updateStatus(requestId, { status }, token);
      await loadData();
      setMessage(
        status === "approved"
          ? "Заявка одобрена. Остальные заявки по этому объявлению автоматически отклонены."
          : "Заявка отклонена.",
      );
    } catch (updateError) {
      setMessage(getErrorMessage(updateError, "Не удалось обновить статус заявки"));
    } finally {
      setActiveRequestId(null);
    }
  };

  if (!user || !token) {
    return (
      <section className="page-content narrow-page">
        <article className="panel">
          <h1>Заявки</h1>
          <p>
            Для доступа к этому разделу необходимо <Link to="/login">войти</Link>.
          </p>
        </article>
      </section>
    );
  }

  const pageTitle = mode === "seller" ? "Заявки на мои объявления" : "Мои заявки";
  const pageDescription =
    mode === "seller"
      ? "Управляйте откликами покупателей на Ваши объявления."
      : "Здесь отображаются Ваши отклики на объявления продавцов.";

  return (
    <section className="page-content">
      <div className="page-title-row">
        <h1>{pageTitle}</h1>
        <p>{pageDescription}</p>
      </div>

      <div className="filter-bar requests-filter-bar">
        <label>
          Статус
          <select value={statusFilter} onChange={(event) => setStatusFilter(event.target.value)}>
            <option value={DEFAULT_FILTER}>Все статусы</option>
            <option value="pending">На рассмотрении</option>
            <option value="approved">Одобрено</option>
            <option value="rejected">Отказано</option>
          </select>
        </label>

        <label>
          Объявление
          <select value={petFilter} onChange={(event) => setPetFilter(event.target.value)}>
            <option value={DEFAULT_FILTER}>Все объявления</option>
            {petOptions.map(([petId, label]) => (
              <option key={petId} value={petId}>
                {label}
              </option>
            ))}
          </select>
        </label>

        {mode === "buyer" ? (
          <label>
            ФИО продавца
            <input
              value={sellerQuery}
              placeholder="Введите ФИО продавца"
              onChange={(event) => setSellerQuery(event.target.value)}
            />
          </label>
        ) : null}
      </div>

      {loading ? <LoadingState label="Загрузка заявок..." /> : null}
      {error ? <AlertMessage variant="error">{error}</AlertMessage> : null}
      {message ? <AlertMessage variant="info">{message}</AlertMessage> : null}

      {!loading && !error && filteredRequests.length === 0 ? (
        <EmptyState
          title="Заявки не найдены"
          description="Измените фильтры или проверьте позже."
        />
      ) : null}

      {!loading && !error && filteredRequests.length > 0 ? (
        <div className="stack-list">
          {filteredRequests.map((item) => {
            const normalizedStatus = normalizeStatus(item.request.status);
            const currentPetLabel = petLabel(item.pet, item.request.pet_id);
            const buyerName = userLabel(item.buyer, item.request.buyer_id, "Покупатель");
            const sellerName = userLabel(item.seller, item.request.seller_id, "Продавец");
            const buyerPhone = userPhoneLabel(item.buyer);
            const sellerPhone = userPhoneLabel(item.seller);
            const canApprove = MODERATED_STATUSES.has(normalizedStatus) && normalizedStatus !== "approved";
            const canReject = MODERATED_STATUSES.has(normalizedStatus) && normalizedStatus !== "rejected";

            return (
              <article key={item.request.id} className="panel compact-panel request-card">
                <div className="panel-header-row">
                  <div className="request-card-title">
                    <div className="pet-card-tags">
                      <h2>{item.pet?.pet_name || `Питомец ${shortId(item.request.pet_id)}`}</h2>
                      <span className={statusBadgeClassName(item.request.status)}>
                        {requestStatusText(item.request.status)}
                      </span>
                    </div>
                    <p>{item.pet?.species || "Вид не указан"}</p>
                    <p className="hint-text">Создана {formatDateTime(item.request.request_date)}</p>
                  </div>

                  <div className="button-row">
                    <Link to={`/pets/${item.request.pet_id}`} className="secondary-button inline-button">
                      Открыть объявление
                    </Link>

                    {mode === "buyer" ? (
                      <button
                        type="button"
                        className="danger-button"
                        disabled={activeRequestId === item.request.id}
                        onClick={() => {
                          void handleDelete(item.request.id);
                        }}
                      >
                        Удалить заявку
                      </button>
                    ) : (
                      <>
                        <button
                          type="button"
                          className="primary-button"
                          disabled={!canApprove || activeRequestId === item.request.id}
                          onClick={() => {
                            void handleStatusChange(item.request.id, "approved");
                          }}
                        >
                          Одобрить
                        </button>
                        <button
                          type="button"
                          className="danger-button"
                          disabled={!canReject || activeRequestId === item.request.id}
                          onClick={() => {
                            void handleStatusChange(item.request.id, "rejected");
                          }}
                        >
                          Отклонить
                        </button>
                      </>
                    )}
                  </div>
                </div>

                <dl className="pet-metadata request-metadata">
                  <div>
                    <dt>Покупатель</dt>
                    <dd>
                      <div>{buyerName}</div>
                      <div className="hint-text request-person-phone">{buyerPhone}</div>
                    </dd>
                  </div>
                  <div>
                    <dt>Продавец</dt>
                    <dd>
                      <div>{sellerName}</div>
                      <div className="hint-text request-person-phone">{sellerPhone}</div>
                    </dd>
                  </div>
                  <div>
                    <dt>Питомец</dt>
                    <dd>{currentPetLabel}</dd>
                  </div>
                </dl>
              </article>
            );
          })}
        </div>
      ) : null}
    </section>
  );
}
