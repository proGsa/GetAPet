import { useEffect, useState } from "react";
import type { PetCreatePayload } from "../../types/pet";
import type { VetPassportUpsertPayload } from "../../types/vetPassport";

export interface PetWithPassportFormValues {
  pet: PetCreatePayload;
  passport: VetPassportUpsertPayload;
}

interface PetWithPassportFormProps {
  initialValue: PetWithPassportFormValues;
  submitLabel: string;
  isSubmitting?: boolean;
  onSubmit: (payload: { pet: PetCreatePayload; passport: VetPassportUpsertPayload }) => Promise<void> | void;
}

const parseNumber = (value: string): number => {
  const parsed = Number(value);
  return Number.isFinite(parsed) ? parsed : 0;
};

export function PetWithPassportForm({
  initialValue,
  submitLabel,
  isSubmitting = false,
  onSubmit,
}: PetWithPassportFormProps) {
  const [form, setForm] = useState<PetWithPassportFormValues>(initialValue);

  useEffect(() => {
    setForm(initialValue);
  }, [initialValue]);

  const setPetString = (key: keyof PetCreatePayload, value: string) => {
    setForm((current) => ({
      ...current,
      pet: { ...current.pet, [key]: value },
    }));
  };

  const setPetNumber = (key: keyof PetCreatePayload, value: string) => {
    setForm((current) => ({
      ...current,
      pet: { ...current.pet, [key]: parseNumber(value) },
    }));
  };

  const setPetBoolean = (key: keyof PetCreatePayload, value: boolean) => {
    setForm((current) => ({
      ...current,
      pet: { ...current.pet, [key]: value },
    }));
  };

  const setPassportString = (key: keyof VetPassportUpsertPayload, value: string) => {
    setForm((current) => ({
      ...current,
      passport: { ...current.passport, [key]: value },
    }));
  };

  const setPassportBoolean = (key: keyof VetPassportUpsertPayload, value: boolean) => {
    setForm((current) => ({
      ...current,
      passport: { ...current.passport, [key]: value },
    }));
  };

  const handleSubmit = async (event: React.FormEvent<HTMLFormElement>) => {
    event.preventDefault();

    const normalizedPet: PetCreatePayload = {
      ...form.pet,
      pet_age: Number(form.pet.pet_age),
      price: Number(form.pet.price),
    };

    await onSubmit({
      pet: normalizedPet,
      passport: form.passport,
    });
  };

  return (
    <form className="form-grid" onSubmit={handleSubmit}>
      <label>
        Кличка
        <input
          value={form.pet.pet_name}
          onChange={(event) => setPetString("pet_name", event.target.value)}
          required
        />
      </label>

      <label>
        Вид
        <input
          value={form.pet.species}
          onChange={(event) => setPetString("species", event.target.value)}
          required
        />
      </label>

      <label>
        Возраст
        <input
          type="number"
          min={0}
          max={50}
          value={form.pet.pet_age}
          onChange={(event) => setPetNumber("pet_age", event.target.value)}
          required
        />
      </label>

      <label>
        Цена (RUB)
        <input
          type="number"
          min={0}
          step="100"
          value={form.pet.price}
          onChange={(event) => setPetNumber("price", event.target.value)}
          required
        />
      </label>

      <label>
        Окрас
        <input value={form.pet.color} onChange={(event) => setPetString("color", event.target.value)} required />
      </label>

      <label>
        Пол
        <select value={form.pet.pet_gender} onChange={(event) => setPetString("pet_gender", event.target.value)} required>
          <option value="male">Самец</option>
          <option value="female">Самка</option>
        </select>
      </label>

      <label>
        Порода
        <input value={form.pet.breed} onChange={(event) => setPetString("breed", event.target.value)} />
      </label>

      <label className="wide-label">
        Описание
        <textarea
          rows={4}
          value={form.pet.pet_description}
          onChange={(event) => setPetString("pet_description", event.target.value)}
        />
      </label>

      <label className="toggle-field">
        <input
          type="checkbox"
          checked={form.pet.pedigree}
          onChange={(event) => setPetBoolean("pedigree", event.target.checked)}
        />
        Родословная
      </label>

      <label className="toggle-field">
        <input
          type="checkbox"
          checked={form.pet.good_with_children}
          onChange={(event) => setPetBoolean("good_with_children", event.target.checked)}
        />
        Ладит с детьми
      </label>

      <label className="toggle-field">
        <input
          type="checkbox"
          checked={form.pet.good_with_animals}
          onChange={(event) => setPetBoolean("good_with_animals", event.target.checked)}
        />
        Ладит с животными
      </label>

      <section className="passport-block wide-label">
        <h3>Ветпаспорт</h3>

        <div className="form-grid">
          <label className="toggle-field">
            <input
              type="checkbox"
              checked={form.passport.chipping}
              onChange={(event) => setPassportBoolean("chipping", event.target.checked)}
            />
            Чипирование
          </label>

          <label className="toggle-field">
            <input
              type="checkbox"
              checked={form.passport.sterilization}
              onChange={(event) => setPassportBoolean("sterilization", event.target.checked)}
            />
            Стерилизация
          </label>

          <label className="wide-label">
            Проблемы со здоровьем
            <textarea
              rows={3}
              value={form.passport.health_issues}
              onChange={(event) => setPassportString("health_issues", event.target.value)}
            />
          </label>

          <label className="wide-label">
            Вакцинации
            <textarea
              rows={3}
              value={form.passport.vaccinations}
              onChange={(event) => setPassportString("vaccinations", event.target.value)}
            />
          </label>

          <label className="wide-label">
            Обработки от паразитов
            <textarea
              rows={3}
              value={form.passport.parasite_treatments}
              onChange={(event) => setPassportString("parasite_treatments", event.target.value)}
            />
          </label>
        </div>
      </section>

      <button type="submit" className="primary-button" disabled={isSubmitting}>
        {isSubmitting ? "Сохранение..." : submitLabel}
      </button>
    </form>
  );
}
