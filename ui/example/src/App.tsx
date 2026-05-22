import React from "react";
import { Application } from "@basesdk/ui";
import { appRegistry } from "./app/registry";
import { appModules } from "./app/modules";

const getBaseUrl = () => import.meta.env.VITE_API_URL ?? window.location.origin;

export default function App() {
  return (
    <React.StrictMode>
      <Application
        getBaseUrl={getBaseUrl}
        modules={appModules}
        registry={appRegistry}
        appTitle="Base ERP"
        appSubtitle="Selecciona un modulo para entrar a su espacio de trabajo."
        loginTitle="Tenant Login"
        loginSubtitle="Accede al entorno principal del ERP."
      />
    </React.StrictMode>
  );
}
