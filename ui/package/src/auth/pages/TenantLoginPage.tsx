import { useState } from "react";
import { useLocation, useNavigate } from "react-router-dom";
import { Button } from "../../components/actions/Button";
import { InputField } from "../../components/form/InputField";
import { createFormSchema } from "../../form/createFormSchema";
import { useCustomForm } from "../../form/useCustomForm";
import { useAuth } from "../useAuth";

const schema = createFormSchema((validators) => ({
  tenantCodigo: validators.requiredString(),
  username: validators.requiredString(),
  password: validators.requiredString(),
}));

type TenantLoginPageProps = {
  title?: string;
  subtitle?: string;
  redirectTo?: string;
  defaultTenantCodigo?: string;
};

export const TenantLoginPage = ({
  title = "Tenant Login",
  subtitle = "Accede al entorno principal del ERP.",
  redirectTo = "/app",
  defaultTenantCodigo,
}: TenantLoginPageProps) => {
  const navigate = useNavigate();
  const location = useLocation();
  const { loginTenant } = useAuth();
  const [submitError, setSubmitError] = useState<string | null>(null);

  const form = useCustomForm(schema, {
    defaultValues: {
      tenantCodigo: defaultTenantCodigo ?? "",
      username: "",
      password: "",
    },
  });

  return (
    <div className="mx-auto w-full max-w-md rounded-3xl bg-ui-panel p-8 shadow-xl shadow-black/5 ring-1 ring-inset ring-ui-border/40">
      <div className="space-y-2">
        <h1 className="text-3xl font-semibold text-ui-text">{title}</h1>
        <p className="text-sm leading-6 text-ui-text-muted">{subtitle}</p>
      </div>

      <form
        className="mt-8 grid gap-4"
        onSubmit={form.handleSubmit(async (values) => {
          setSubmitError(null);

          try {
            await loginTenant(values);
            const nextPath = (location.state as { from?: { pathname?: string } } | null)?.from?.pathname ?? redirectTo;
            navigate(nextPath, { replace: true });
          } catch (error) {
            setSubmitError(error instanceof Error ? error.message : "No se pudo iniciar sesion");
          }
        })}
      >
        <InputField form={form} name="tenantCodigo" label="Tenant" placeholder="tenant_default" />
        <InputField form={form} name="username" label="Usuario" placeholder="kevin" />
        <InputField form={form} name="password" label="Password" type="password" placeholder="••••••••" />

        {submitError && <p className="text-sm text-ui-danger">{submitError}</p>}

        <Button type="submit" loading={form.formState.isSubmitting}>
          Iniciar sesion tenant
        </Button>
      </form>
    </div>
  );
};
