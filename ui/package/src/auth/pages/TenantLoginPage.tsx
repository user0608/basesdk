import { useEffect, useState } from "react";
import { Link, useLocation, useNavigate, useSearchParams } from "react-router-dom";
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
  const [searchParams] = useSearchParams();
  const { loginTenant } = useAuth();
  const [submitError, setSubmitError] = useState<string | null>(null);
  const tenantFromUrl = searchParams.get("tenant") ?? "";

  const form = useCustomForm(schema, {
    defaultValues: {
      tenantCodigo: tenantFromUrl || defaultTenantCodigo || "",
      username: "",
      password: "",
    },
  });
  const { setValue } = form;

  useEffect(() => {
    if (tenantFromUrl) {
      setValue("tenantCodigo", tenantFromUrl);
    }
  }, [setValue, tenantFromUrl]);

  return (
    <main className="grid min-h-screen place-items-center bg-ui-bg px-4 py-6 text-ui-text sm:px-6">
      <div className="w-full max-w-sm rounded-2xl bg-ui-panel p-5 shadow-xl shadow-ui-text/5 ring-1 ring-inset ring-ui-border/45 sm:max-w-md sm:p-6">
        <div className="space-y-1.5 border-b border-ui-border/25 pb-4">
          <div className="text-xs font-semibold uppercase tracking-[0.18em] text-ui-accent">Base ERP</div>
          <h1 className="text-2xl font-semibold tracking-tight text-ui-text sm:text-3xl">{title}</h1>
          <p className="text-sm leading-5 text-ui-text-muted">{subtitle}</p>
        </div>

        <form
          className="mt-5 grid gap-3.5"
          onSubmit={form.handleSubmit(async (values) => {
            setSubmitError(null);

            try {
              const session = await loginTenant(values);
              const nextPath = session.mustChangePassword
                ? "/change-password"
                : (location.state as { from?: { pathname?: string } } | null)?.from?.pathname ?? redirectTo;
              navigate(nextPath, { replace: true });
            } catch (error) {
              setSubmitError(error instanceof Error ? error.message : "No se pudo iniciar sesion");
            }
          })}
        >
          <InputField form={form} name="tenantCodigo" label="Tenant" placeholder="tenant_default" />
          <InputField form={form} name="username" label="Usuario" placeholder="kevin" />
          <InputField form={form} name="password" label="Password" type="password" placeholder="Password" />

          {submitError && <p className="rounded-xl bg-ui-danger/10 px-3 py-2 text-sm text-ui-danger">{submitError}</p>}

          <Button type="submit" className="mt-1 w-full" loading={form.formState.isSubmitting}>
            Iniciar sesion
          </Button>
        </form>

        <div className="mt-4 border-t border-ui-border/20 pt-3 text-center">
          <Link to="/system/login" className="text-xs font-medium text-ui-text-soft transition-colors hover:text-ui-primary">
            Entrar como administrador del sistema
          </Link>
        </div>
      </div>
    </main>
  );
};
