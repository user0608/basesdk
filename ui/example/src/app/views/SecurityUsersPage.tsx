import { Button, useModal } from "@basesdk/ui";

export const SecurityUsersPage = () => {
  const dialog = useModal();

  return (
    <section className="flex justify-end">
      <Button
        type="button"
        onClick={() =>
          dialog.open({
            title: "Nuevo usuario",
            description: "Formulario de ejemplo renderizado dentro del manejador global de dialogos.",
            content: ({ close }) => (
              <form className="grid gap-3">
                <label className="grid gap-1.5 text-sm font-medium text-ui-text">
                  Nombre
                  <input
                    className="rounded-xl border border-ui-border bg-ui-surface px-3 py-2 text-sm outline-none transition-colors focus:border-ui-primary focus:ring-2 focus:ring-ui-focus"
                    placeholder="Nombre del usuario"
                  />
                </label>
                <label className="grid gap-1.5 text-sm font-medium text-ui-text">
                  Correo
                  <input
                    className="rounded-xl border border-ui-border bg-ui-surface px-3 py-2 text-sm outline-none transition-colors focus:border-ui-primary focus:ring-2 focus:ring-ui-focus"
                    placeholder="usuario@empresa.com"
                  />
                </label>

                <div className="flex justify-end gap-2 pt-1">
                  <Button type="button" variant="secondary" onClick={close}>
                    Cancelar
                  </Button>
                  <Button type="button" onClick={close}>
                    Guardar
                  </Button>
                </div>
              </form>
            ),
          })
        }
      >
        Nuevo usuario
      </Button>
    </section>
  );
};
