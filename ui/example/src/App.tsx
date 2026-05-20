import { TenantApplication } from "@basesdk/ui";
import { componentRegistry } from "./app/registry";
import { appModules } from "./app/modules";

export default function App() {
  return (
    <TenantApplication
      modules={appModules}
      registry={componentRegistry}
      appTitle="Base ERP"
      appSubtitle="Selecciona un modulo para entrar a su espacio de trabajo."
      loginTitle="Tenant Login"
      loginSubtitle="Accede al entorno principal del ERP."
      defaultTenantCodigo="tenant_default"
    />
  );
}
