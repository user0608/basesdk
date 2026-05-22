import { useState } from "react";
import { useNavigate } from "react-router-dom";
import { Button } from "../../components/actions/Button";
import { InputField } from "../../components/form/InputField";
import { createFormSchema } from "../../form/createFormSchema";
import { useCustomForm } from "../../form/useCustomForm";
import { useAuth } from "../useAuth";

const schema = createFormSchema((validators) => ({
  password: validators.requiredString(),
}));

export const TenantPasswordChangePage = () => {
  const navigate = useNavigate();
  const { completeTenantPasswordChange, httpApi, logoutTenant, tenantSession } = useAuth();
  const [submitError, setSubmitError] = useState<string | null>(null);
  const form = useCustomForm(schema, { defaultValues: { password: "" } });

  return (
    <main className="grid min-h-screen place-items-center bg-ui-bg px-4 py-6 text-ui-text sm:px-6">
      <div className="w-full max-w-sm rounded-2xl bg-ui-panel p-5 shadow-xl shadow-ui-text/5 ring-1 ring-inset ring-ui-border/45 sm:max-w-md sm:p-6">
        <div className="space-y-1.5 border-b border-ui-border/25 pb-4">
          <div className="text-xs font-semibold uppercase tracking-[0.18em] text-ui-accent">Password requerido</div>
          <h1 className="text-2xl font-semibold tracking-tight text-ui-text">Cambia tu password</h1>
          <p className="text-sm leading-5 text-ui-text-muted">
            La cuenta {tenantSession?.username ?? "actual"} debe actualizar su password antes de continuar.
          </p>
        </div>

        <form
          className="mt-5 grid gap-3.5"
          onSubmit={form.handleSubmit(async (values) => {
            setSubmitError(null);
            try {
              await httpApi.patch<void>({ path: "/api/v1/me/password", auth: "tenant", data: { password: values.password } });
              completeTenantPasswordChange();
              navigate("/app", { replace: true });
            } catch (error) {
              setSubmitError(error instanceof Error ? error.message : "No se pudo actualizar el password");
            }
          })}
        >
          <InputField form={form} name="password" label="Nuevo password" type="password" />

          {submitError && <p className="rounded-xl bg-ui-danger/10 px-3 py-2 text-sm text-ui-danger">{submitError}</p>}

          <Button type="submit" className="w-full" loading={form.formState.isSubmitting}>
            Actualizar password
          </Button>
          <Button type="button" variant="secondary" className="w-full" onClick={logoutTenant}>
            Cerrar sesion
          </Button>
        </form>
      </div>
    </main>
  );
};
