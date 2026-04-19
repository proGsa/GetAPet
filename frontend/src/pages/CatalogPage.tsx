import { useEffect, useMemo, useState } from "react";
import { petsApi } from "../api/pets";
import { usersApi } from "../api/users";
import { AlertMessage } from "../components/ui/AlertMessage";
import { EmptyState } from "../components/ui/EmptyState";
import { LoadingState } from "../components/ui/LoadingState";
import { useAuth } from "../hooks/useAuth";
import type { Pet } from "../types/pet";
import type { User } from "../types/user";
import { getErrorMessage } from "../utils/error";
import { formatPrice } from "../utils/format";

export function CatalogPage() {
  const { user, token, mode } = useAuth();

  const [pets, setPets] = useState<Pet[]>([]);
  const [usersMap, setUsersMap] = useState<Record<string, User>>({});
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const [query, setQuery] = useState("");
  const [speciesFilter, setSpeciesFilter] = useState("all");
  const [activeOnly, setActiveOnly] = useState(true);
  const [myOnly, setMyOnly] = useState(false);

  const canFilterByMine = Boolean(user && mode === "seller");

  useEffect(() => {
    if (!canFilterByMine) {
      setMyOnly(false);
    }
  }, [canFilterByMine]);

  useEffect(() => {
    const loadData = async () => {
      setLoading(true);
      setError(null);

      try {
        const petsResponse = await petsApi.list();
        setPets(petsResponse);

        if (!token) {
          setUsersMap({});
          return;
        }

        try {
          const usersResponse = await usersApi.list(token);
          const nextUsersMap = usersResponse.reduce<Record<string, User>>((accumulator, currentUser) => {
            accumulator[currentUser.id] = currentUser;
            return accumulator;
          }, {});
          setUsersMap(nextUsersMap);
        } catch {
          setUsersMap({});
        }
      } catch (fetchError) {
        setError(getErrorMessage(fetchError, "Не удалось загрузить данные каталога"));
      } finally {
        setLoading(false);
      }
    };

    void loadData();
  }, [token]);

  const speciesOptions = useMemo(() => {
    const options = new Set<string>();
    pets.forEach((pet) => {
      if (pet.species) {
        options.add(pet.species);
      }
    });

    return ["all", ...Array.from(options).sort((left, right) => left.localeCompare(right))];
  }, [pets]);

  const filteredPets = useMemo(() => {
    const normalizedQuery = query.trim().toLowerCase();

    return pets.filter((pet) => {
      if (activeOnly && !pet.is_active) {
        return false;
      }

      if (myOnly && user && pet.seller_id !== user.id) {
        return false;
      }

      if (speciesFilter !== "all" && pet.species !== speciesFilter) {
        return false;
      }

      if (!normalizedQuery) {
        return true;
      }

      const searchable = [pet.pet_name, pet.species, pet.breed, pet.color, pet.pet_description]
        .join(" ")
        .toLowerCase();

      return searchable.includes(normalizedQuery);
    });
  }, [activeOnly, myOnly, pets, query, speciesFilter, user]);

  return (
    <section className="page-content">
      <div className="page-title-row">
        <h1>Каталог питомцев</h1>
        <p>Объявления от частных владельцев и приютов.</p>
      </div>

      <div className="filter-bar catalog-filter-bar">
        <label>
          Поиск
          <input
            value={query}
            placeholder="Кличка, вид, порода"
            onChange={(event) => setQuery(event.target.value)}
          />
        </label>

        <label>
          Вид
          <select value={speciesFilter} onChange={(event) => setSpeciesFilter(event.target.value)}>
            {speciesOptions.map((species) => (
              <option key={species} value={species}>
                {species === "all" ? "Все виды" : species}
              </option>
            ))}
          </select>
        </label>
        <div className="catalog-filter-toggles">
          <label className="toggle-field compact-toggle">
            <input
              type="checkbox"
              checked={activeOnly}
              onChange={(event) => setActiveOnly(event.target.checked)}
            />
            Только активные
          </label>

          {canFilterByMine ? (
            <label className="toggle-field compact-toggle">
              <input
                type="checkbox"
                checked={myOnly}
                onChange={(event) => setMyOnly(event.target.checked)}
              />
              Только мои объявления
            </label>
          ) : null}
        </div>
      </div>

      {loading ? <LoadingState label="Загрузка каталога..." /> : null}
      {error ? <AlertMessage variant="error">{error}</AlertMessage> : null}

      {!loading && !error && filteredPets.length === 0 ? (
        <EmptyState
          title="Питомцы не найдены"
          description="Измените фильтры или отключите режим показа только активных объявлений."
        />
      ) : null}

      {!loading && !error && filteredPets.length > 0 ? (
        <div className="card-grid">
          {filteredPets.map((pet) => {
            const seller = usersMap[pet.seller_id];
            return (
              <article key={pet.id} className="pet-card">
                <div className="pet-card-main">
                  <p className="pet-card-species">{pet.species || "Неизвестный вид"}</p>
                  <h2>{pet.pet_name}</h2>
                  <p className="pet-card-desc">{pet.pet_description || "Описание пока не добавлено"}</p>
                </div>

                <dl className="pet-metadata">
                  <div>
                    <dt>Возраст</dt>
                    <dd>{pet.pet_age}</dd>
                  </div>
                  <div>
                    <dt>Порода</dt>
                    <dd>{pet.breed || "Не указана"}</dd>
                  </div>
                  <div>
                    <dt>Продавец</dt>
                    <dd>{seller?.fio ?? "Продавец не указан"}</dd>
                  </div>
                  <div>
                    <dt>Цена</dt>
                    <dd>{formatPrice(pet.price)}</dd>
                  </div>
                </dl>

                <button type="button" className="primary-button inline-button" disabled>
                  Подробнее
                </button>
              </article>
            );
          })}
        </div>
      ) : null}
    </section>
  );
}
