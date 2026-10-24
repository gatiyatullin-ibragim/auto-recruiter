import { useMutation, useQuery, useQueryClient } from "@tanstack/react-query";
import { ChangeEvent, FormEvent, useEffect, useState } from "react";

type Candidate = {
  name: string;
  email: string;
  resume_text: string;
  skills?: string[] | null;
  criteria: Criteria;
};

type Criteria = {
  position: string;
  roles?: string[] | null;
  locations?: string[] | null;
  work_mode: "remote" | "hybrid" | "any" | "";
  employment_type: string;
  min_salary: number;
  currency: string;
  remote: boolean;
};

type Vacancy = {
  id: string;
  source: string;
  title: string;
  company: string;
  description: string;
  location: string;
  url: string;
  remote: boolean;
  skills: string[];
};

const emptyCandidate: Candidate = {
  name: "",
  email: "",
  resume_text: "",
  skills: [],
  criteria: {
    position: "",
    roles: [],
    locations: [],
    work_mode: "",
    employment_type: "",
    min_salary: 0,
    currency: "KZT",
    remote: false,
  },
};

function parseList(value: string): string[] {
  return value.split(",").map((item) => item.trim()).filter(Boolean);
}

function normalizeCandidate(candidate: Candidate): Candidate {
  return {
    ...candidate,
    skills: candidate.skills ?? [],
    criteria: {
      ...candidate.criteria,
      roles: candidate.criteria.roles ?? [],
      locations: candidate.criteria.locations ?? [],
    },
  };
}

async function request<T>(url: string, options?: RequestInit): Promise<T> {
  const response = await fetch(url, options);
  if (!response.ok) throw new Error(await response.text());
  return response.json();
}

export default function App() {
  const queryClient = useQueryClient();
  const profile = useQuery({
    queryKey: ["profile"],
    queryFn: async () => normalizeCandidate(await request<Candidate>("/api/profile")),
  });
  const [form, setForm] = useState<Candidate>(emptyCandidate);
  const [skillsInput, setSkillsInput] = useState("");
  const [rolesInput, setRolesInput] = useState("");
  const [locationsInput, setLocationsInput] = useState("");
  const [resumeFile, setResumeFile] = useState<File | null>(null);
  const [formError, setFormError] = useState("");
  const [searchStarted, setSearchStarted] = useState(false);
  const [searchStatus, setSearchStatus] = useState("");
  const [vacancies, setVacancies] = useState<Vacancy[]>([]);
  const [sourceFilter, setSourceFilter] = useState("all");

  useEffect(() => {
    if (profile.data) {
      setForm(profile.data);
      setSkillsInput((profile.data.skills ?? []).join(", "));
      setRolesInput((profile.data.criteria.roles ?? []).join(", "));
      setLocationsInput((profile.data.criteria.locations ?? []).join(", "));
    }
  }, [profile.data]);

  const saveProfile = useMutation({
    mutationFn: (candidate: Candidate) =>
      request<Candidate>("/api/profile", {
        method: "PUT",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(candidate),
      }),
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ["profile"] }),
  });

  const uploadResume = useMutation({
    mutationFn: (file: File) => {
      const body = new FormData();
      body.append("file", file);
      return request<{ characters: number }>("/api/resume", { method: "POST", body });
    },
    onSuccess: () => queryClient.invalidateQueries({ queryKey: ["profile"] }),
  });

  const current = normalizeCandidate(form);
  const update = (patch: Partial<Candidate>) => setForm({ ...current, ...patch });

  function submit(event: FormEvent) {
    event.preventDefault();
    setFormError("");
    const candidate = buildCandidate();
    if (!validateSearch(candidate)) return;
    saveProfile.mutate(candidate);
  }

  function buildCandidate(): Candidate {
    return {
      ...current,
      skills: parseList(skillsInput),
      criteria: {
        ...current.criteria,
        roles: parseList(rolesInput),
        locations: parseList(locationsInput),
      },
    };
  }

  function validateSearch(candidate: Candidate): boolean {
    if (candidate.name.trim().length < 2) {
      setFormError("Введите имя длиной не менее 2 символов.");
      return false;
    }
    if (!candidate.criteria.position.trim() || candidate.criteria.position.trim().length < 2) {
      setFormError("Укажите желаемую позицию.");
      return false;
    }
    if (candidate.criteria.work_mode === "hybrid" && !candidate.criteria.locations?.length) {
      setFormError("Для офисного или гибридного формата укажите локацию.");
      return false;
    }
    return true;
  }

  function handleFile(event: ChangeEvent<HTMLInputElement>) {
    const file = event.target.files?.[0];
    if (file) {
      setResumeFile(file);
      uploadResume.mutate(file);
    }
  }

  function startSearch() {
    setFormError("");
    const candidate = buildCandidate();
    if (!validateSearch(candidate)) return;
    saveProfile.mutate(candidate, {
      onSuccess: async () => {
        setSearchStarted(true);
        setSearchStatus("Запускаем источники...");
        const result = await request<{ job_id: string }>("/api/search", {
          method: "POST",
          headers: { "Content-Type": "application/json" },
          body: JSON.stringify(candidate.criteria),
        });
        const events = new EventSource(`/api/search/${result.job_id}/events`);
        events.addEventListener("progress", (event) => {
          const progress = JSON.parse((event as MessageEvent).data) as { status: string; message?: string; found?: number };
          setSearchStatus(progress.message ?? `${progress.status}: найдено ${progress.found ?? 0}`);
          if (progress.status === "done") {
            events.close();
            request<Vacancy[]>("/vacancies").then(setVacancies).catch(() => setSearchStatus("Не удалось загрузить вакансии."));
          }
        });
        events.onerror = () => {
          events.close();
          setSearchStatus("Потеряно соединение с прогрессом поиска.");
        };
      },
    });
  }

  return (
    <main className="layout">
      <header>
        <p className="eyebrow">AI-NATIVE CAREER WORKSPACE</p>
        <h1>Auto Recruiter</h1>
        <p className="subtitle">Профиль и критерии поиска в одном месте.</p>
      </header>

      {profile.isLoading ? <p>Загрузка профиля...</p> : null}
      {profile.isError ? <p className="error">{profile.error.message}</p> : null}

      <form className="card" onSubmit={submit}>
        <h2>Профиль кандидата</h2>
        <label>
          Имя
          <input value={current.name} onChange={(e) => update({ name: e.target.value })} />
        </label>
        <label>
          Email
          <input type="email" value={current.email} onChange={(e) => update({ email: e.target.value })} />
        </label>
        <label>
          Желаемая позиция
          <input
            value={current.criteria.position}
            onChange={(e) => update({ criteria: { ...current.criteria, position: e.target.value } })}
            placeholder="Например, Go backend developer"
          />
        </label>
        <label>
          Локация
          <input
            value={locationsInput}
            onChange={(e) => setLocationsInput(e.target.value)}
            placeholder="Город или страна"
          />
        </label>
        <label>
          Навыки через запятую
          <input
            value={skillsInput}
            onChange={(e) => setSkillsInput(e.target.value)}
          />
        </label>
        <label>
          Текст резюме
          <textarea value={current.resume_text} onChange={(e) => update({ resume_text: e.target.value })} />
        </label>
        <div className="actions">
          <button type="submit" disabled={saveProfile.isPending}>Сохранить профиль</button>
          <label className="upload">
            {resumeFile ? `Резюме: ${resumeFile.name}` : "Загрузить резюме"}
            <input type="file" accept=".txt,.pdf,.docx" onChange={handleFile} />
          </label>
        </div>
        {uploadResume.isSuccess ? <p className="success">Резюме загружено.</p> : null}
        {saveProfile.isError ? <p className="error">{saveProfile.error.message}</p> : null}
        {formError ? <p className="error">{formError}</p> : null}
        <button type="button" className="primary-action" onClick={startSearch} disabled={saveProfile.isPending || !current.criteria.position.trim()}>
          Начать поиск
        </button>
        {searchStarted ? <p className="success">Профиль сохранён. Подключение источников поиска выполняется следующим этапом.</p> : null}
        {searchStatus ? <p className="search-status">{searchStatus}</p> : null}
      </form>

      <section className="card">
        <h2>Критерии поиска</h2>
        <label className="checkbox">
          <input
            type="checkbox"
            checked={current.criteria.remote}
            onChange={(e) => update({ criteria: { ...current.criteria, remote: e.target.checked } })}
          />
          Только удалённая работа
        </label>
        <label>
          Формат работы
          <select
            value={current.criteria.work_mode}
            onChange={(e) => update({ criteria: { ...current.criteria, work_mode: e.target.value as Criteria["work_mode"] } })}
          >
            <option value="">Не выбрано</option>
            <option value="remote">Удалённо</option>
            <option value="hybrid">Офис или гибрид</option>
            <option value="any">Не важно</option>
          </select>
        </label>
        <div className="grid">
          <label>
            Роли через запятую
            <input
              value={rolesInput}
              onChange={(e) => setRolesInput(e.target.value)}
            />
          </label>
          <label>
            Локации через запятую
            <input
              value={locationsInput}
              onChange={(e) => setLocationsInput(e.target.value)}
            />
          </label>
          <label>
            Минимальная зарплата
            <input type="number" min="0" value={current.criteria.min_salary} onChange={(e) => update({ criteria: { ...current.criteria, min_salary: Number(e.target.value) } })} />
          </label>
          <label>
            Валюта
            <input value={current.criteria.currency} onChange={(e) => update({ criteria: { ...current.criteria, currency: e.target.value.toUpperCase() } })} />
          </label>
        </div>
      </section>

      {vacancies.length > 0 ? (
        <section className="card">
          <div className="section-heading">
            <div>
              <h2>Вакансии</h2>
              <p className="subtitle">Найдено: {vacancies.length}</p>
            </div>
            <select value={sourceFilter} onChange={(e) => setSourceFilter(e.target.value)}>
              <option value="all">Все платформы</option>
              {[...new Set(vacancies.map((vacancy) => vacancy.source))].map((source) => (
                <option key={source} value={source}>{source}</option>
              ))}
            </select>
          </div>
          <div className="vacancies">
            {vacancies.filter((vacancy) => sourceFilter === "all" || vacancy.source === sourceFilter).map((vacancy) => (
              <article className="vacancy" key={vacancy.id}>
                <div>
                  <span className="badge">{vacancy.source}</span>
                  {vacancy.remote ? <span className="badge remote">Удалённо</span> : null}
                </div>
                <h3>{vacancy.title}</h3>
                <p>{vacancy.company} · {vacancy.location}</p>
                <p>{vacancy.description}</p>
                <a href={vacancy.url} target="_blank" rel="noreferrer">Открыть оригинал</a>
              </article>
            ))}
          </div>
        </section>
      ) : null}
    </main>
  );
}
