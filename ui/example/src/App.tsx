import { useMemo, useState } from "react";
import type { ColumnDef } from "@tanstack/react-table";
import {
  AsyncSelectField,
  Button,
  createFormSchema,
  InputField,
  SelectField,
  useCustomForm,
  type SelectOption,
} from "@basesdk/ui";
import * as yup from "yup";

type DepartmentOption = SelectOption<{
  area: string;
}>;

const roleOptions: SelectOption[] = [
  { value: "designer", label: "Designer" },
  { value: "frontend", label: "Frontend" },
  { value: "backend", label: "Backend" },
  { value: "product", label: "Product" },
];

const stackOptions: SelectOption[] = [
  { value: "react", label: "React" },
  { value: "typescript", label: "TypeScript" },
  { value: "tailwind", label: "Tailwind" },
  { value: "tanstack-query", label: "TanStack Query" },
  { value: "react-hook-form", label: "React Hook Form" },
];

const planOptions: SelectOption[] = [
  { value: "starter", label: "Starter" },
  { value: "growth", label: "Growth" },
  { value: "scale", label: "Scale" },
];

const countryOptions: SelectOption[] = [
  { value: "pe", label: "Peru" },
  { value: "mx", label: "Mexico" },
  { value: "co", label: "Colombia" },
  { value: "cl", label: "Chile" },
];

const managerOptions: SelectOption[] = [
  { value: "ana-ramirez", label: "Ana Ramirez" },
  { value: "bruno-salazar", label: "Bruno Salazar" },
  { value: "camila-ruiz", label: "Camila Ruiz" },
  { value: "diego-castro", label: "Diego Castro" },
];

const departmentOptions: DepartmentOption[] = [
  { value: "eng-platform", label: "Platform", area: "Engineering" },
  { value: "eng-ui", label: "UI Systems", area: "Engineering" },
  { value: "ops-finance", label: "Finance Ops", area: "Operations" },
  { value: "ops-support", label: "Support", area: "Operations" },
  { value: "growth-sales", label: "Sales", area: "Growth" },
  { value: "growth-partners", label: "Partners", area: "Growth" },
];

const schema = createFormSchema((validators) => ({
  fullName: validators.requiredString(),
  email: validators.requiredEmail(),
  celular: validators.requiredCelular(),
  role: validators.requiredString(),
  manager: validators.requiredString(),
  stacks: validators.requiredStringArray(),
  plan: validators.requiredString(),
  country: validators.requiredString(),
  departments: validators.requiredStringArray(),
  active: validators.requiredBoolean(),
  notes: validators.nullableString(),
}));

type DemoFormValues = yup.InferType<typeof schema>;

const wait = (ms: number) => new Promise((resolve) => setTimeout(resolve, ms));

export default function App() {
  const [submittedValues, setSubmittedValues] = useState<DemoFormValues | null>(null);

  const form = useCustomForm(schema, {
    defaultValues: {
      fullName: "",
      email: "",
      celular: "",
      role: "frontend",
      manager: "camila-ruiz",
      stacks: ["react", "typescript"],
      plan: "growth",
      country: "pe",
      departments: ["eng-ui"],
      active: true,
      notes: null,
    },
  });

  const departmentColumns = useMemo<ColumnDef<DepartmentOption, unknown>[]>(
    () => [
      {
        accessorKey: "area",
        header: "Area",
        cell: ({ row }) => <span className="text-slate-500">{row.original.area}</span>,
      },
    ],
    [],
  );

  const loadCountries = async () => {
    await wait(300);
    return countryOptions;
  };

  const loadDepartments = async () => {
    await wait(450);
    return departmentOptions;
  };

  return (
    <main className="min-h-screen bg-[color:var(--ui-bg)] px-6 py-16 text-[color:var(--ui-text)] transition-colors">
      <div className="mx-auto flex max-w-6xl flex-col gap-8 rounded-3xl border border-[color:var(--ui-border)] bg-[color:var(--ui-panel)] p-8 shadow-2xl shadow-black/10 transition-colors">
        <div className="space-y-3">
          <span className="inline-flex rounded-full border border-[color:var(--ui-accent)]/30 bg-[color:var(--ui-accent)]/10 px-3 py-1 text-xs font-semibold uppercase tracking-[0.24em] text-[color:var(--ui-accent)]">
            UI package playground
          </span>
          <h1 className="text-4xl font-semibold tracking-tight text-[color:var(--ui-text)]">Formulario real con Yup + RHF</h1>
          <p className="max-w-3xl text-sm leading-6 text-[color:var(--ui-text-muted)]">
            Este ejemplo usa `createFormSchema`, `useCustomForm`, `InputField`, `SelectField` y
            `AsyncSelectField` desde el paquete local `@basesdk/ui`.
          </p>
        </div>

        <div className="grid gap-8 lg:grid-cols-[minmax(0,1.5fr)_minmax(320px,0.9fr)]">
          <form
            onSubmit={form.handleSubmit((values) => setSubmittedValues(values))}
            className="grid gap-6 rounded-2xl border border-[color:var(--ui-border)] bg-[color:var(--ui-panel-muted)] p-6"
          >
            <section className="grid gap-4 md:grid-cols-2">
              <InputField
                form={form}
                name="fullName"
                label="Nombre completo"
                placeholder="Ada Lovelace"
              />
              <InputField
                form={form}
                name="email"
                label="Correo electronico"
                type="email"
                placeholder="ada@basesdk.dev"
              />
              <InputField
                form={form}
                name="celular"
                label="Celular"
                placeholder="999 111 222"
              />
              <InputField
                form={form}
                name="active"
                variant="boolean"
                label="Estado"
                yesLabel="Activo"
                noLabel="Inactivo"
              />
            </section>

            <section className="grid gap-4 md:grid-cols-2">
              <SelectField
                form={form}
                name="role"
                variant="native"
                label="Rol principal"
                options={roleOptions}
              />
              <SelectField
                form={form}
                name="manager"
                variant="combobox"
                label="Manager"
                options={managerOptions}
                placeholder="Selecciona un manager"
                searchPlaceholder="Buscar manager"
              />
            </section>

            <section className="grid gap-4 md:grid-cols-2">
              <SelectField
                form={form}
                name="plan"
                variant="card"
                label="Plan"
                options={planOptions}
              />
            </section>

            <section className="grid gap-4">
              <SelectField
                form={form}
                name="stacks"
                variant="combobox"
                label="Stack del equipo"
                info="Multi-select visible sin ocultar seleccionados"
                options={stackOptions}
                multiple
                placeholder="Selecciona tecnologias"
                searchPlaceholder="Buscar tecnologia"
              />
              <AsyncSelectField
                form={form}
                name="country"
                variant="native"
                label="Pais"
                queryKey={["countries"]}
                loadOptions={loadCountries}
                placeholder="Selecciona un pais"
              />
              <AsyncSelectField
                form={form}
                name="departments"
                variant="table"
                label="Departamentos"
                queryKey={["departments"]}
                loadOptions={loadDepartments}
                columns={departmentColumns}
                searchKeys={["label", "area"]}
                multiple
                pageSize={4}
              />
            </section>

            <section className="grid gap-4">
              <InputField
                form={form}
                name="notes"
                label="Notas"
                placeholder="Notas internas opcionales"
              />
              <InputField
                form={form}
                name="active"
                variant="checkbox"
                label="Visible en dashboard"
                info="Mismo valor booleano con otra representacion visual"
              />
            </section>

            <div className="flex flex-wrap items-center gap-3">
              <Button type="submit">Guardar formulario</Button>
              <Button variant="secondary" type="button" onClick={() => form.reset()}>
                Reset
              </Button>
            </div>
          </form>

          <aside className="grid gap-6">
            <section className="rounded-2xl border border-[color:var(--ui-border)] bg-[color:var(--ui-panel-muted)] p-6">
              <h2 className="text-sm font-semibold uppercase tracking-[0.18em] text-[color:var(--ui-text-muted)]">Estado</h2>
              <dl className="mt-4 grid gap-3 text-sm text-[color:var(--ui-text-muted)]">
                <div className="flex items-center justify-between gap-4">
                  <dt className="text-[color:var(--ui-text-soft)]">Valido</dt>
                  <dd>{form.formState.isValid ? "Si" : "No"}</dd>
                </div>
                <div className="flex items-center justify-between gap-4">
                  <dt className="text-[color:var(--ui-text-soft)]">Dirty</dt>
                  <dd>{form.formState.isDirty ? "Si" : "No"}</dd>
                </div>
                <div className="flex items-center justify-between gap-4">
                  <dt className="text-[color:var(--ui-text-soft)]">Enviado</dt>
                  <dd>{form.formState.submitCount}</dd>
                </div>
              </dl>
            </section>

            <section className="rounded-2xl border border-[color:var(--ui-border)] bg-[color:var(--ui-panel-muted)] p-6">
              <h2 className="text-sm font-semibold uppercase tracking-[0.18em] text-[color:var(--ui-text-muted)]">Payload</h2>
              <pre className="mt-4 overflow-auto rounded-xl border border-[color:var(--ui-border)] bg-[color:var(--ui-bg)] p-4 text-xs leading-6 text-[color:var(--ui-accent)]">
                {JSON.stringify(submittedValues ?? form.getValues(), null, 2)}
              </pre>
            </section>
          </aside>
        </div>
      </div>
    </main>
  );
}
