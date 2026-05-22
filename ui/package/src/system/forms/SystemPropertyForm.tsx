import { Button } from "../../components/actions/Button";
import { InputField } from "../../components/form/InputField";
import { SelectField } from "../../components/form/SelectField";
import { createFormSchema } from "../../form/createFormSchema";
import { useCustomForm } from "../../form/useCustomForm";
import { useMutate } from "../../query/useMutate";
import type { PropertyDataType, PropertyResponse, TenantPropertyResponse } from "../types";
import { useSystemService } from "../useSystemService";

const dataTypeOptions: Array<{ value: PropertyDataType; label: string }> = [
  { value: "string", label: "Texto" },
  { value: "int", label: "Entero" },
  { value: "float", label: "Decimal" },
  { value: "bool", label: "Booleano" },
  { value: "json", label: "JSON" },
];

const schema = createFormSchema((validators) => ({
  key: validators.requiredString(),
  value: validators.requiredString(),
  dataType: validators.requiredString(),
  description: validators.nullableString(),
}));

export type SystemPropertyFormProps = {
  property?: PropertyResponse | TenantPropertyResponse;
  tenantCodigo?: string;
  close?: () => void;
};

type PropertyFormValues = {
  key: string;
  value: string;
  dataType: PropertyDataType;
  description: string | null;
};

const sensitiveKeys = new Set(["jwt_token_ttl"]);

const validatePropertyValue = (dataType: PropertyDataType, value: string) => {
  const trimmed = value.trim();
  if (dataType === "string") return null;
  if (dataType === "int" && !/^-?\d+$/.test(trimmed)) return "El valor debe ser un entero valido";
  if (dataType === "float" && (trimmed === "" || Number.isNaN(Number(trimmed)))) return "El valor debe ser un numero valido";
  if (dataType === "bool" && !["1", "0", "true", "false", "yes", "no", "on", "off", "enabled", "disabled"].includes(trimmed.toLowerCase())) return "El valor debe ser booleano";
  if (dataType === "json") {
    try {
      JSON.parse(value);
    } catch {
      return "El valor debe ser JSON valido";
    }
  }
  return null;
};

export const SystemPropertyForm = ({ property, tenantCodigo, close }: SystemPropertyFormProps) => {
  const system = useSystemService();
  const service = tenantCodigo ? system.tenantProperties(tenantCodigo) : system.properties;
  const queryKey = tenantCodigo ? ["system", "tenants", tenantCodigo, "properties"] : ["system", "properties"];
  const isEdit = Boolean(property);
  const form = useCustomForm(schema, {
    defaultValues: {
      key: property?.key ?? "",
      value: property?.value ?? "",
      dataType: property?.dataType ?? "string",
      description: property?.description ?? null,
    },
  });
  const dataType = form.watch("dataType") as PropertyDataType;
  const key = form.watch("key");
  const save = useMutate({
    mutationFn: (values: PropertyFormValues) => {
      if (isEdit && property) return service.update(property.key, values);
      return service.create(values);
    },
    invalidateQueryKey: queryKey,
    successMessage: isEdit ? "Propiedad actualizada." : "Propiedad creada.",
  });

  return (
    <form
      className="grid gap-3"
      onSubmit={form.handleSubmit(async (values) => {
        const errorMessage = validatePropertyValue(values.dataType as PropertyDataType, values.value);
        if (errorMessage) {
          form.setError("value", { message: errorMessage });
          return;
        }

        await save.mutate(values as PropertyFormValues);
        close?.();
      })}
    >
      <InputField form={form} name="key" label="Key" readOnly={isEdit} />
      <SelectField form={form} name="dataType" label="Tipo" options={dataTypeOptions} />
      {key && sensitiveKeys.has(String(key)) && (
        <p className="rounded-xl bg-ui-danger/10 px-3 py-2 text-sm leading-5 text-ui-danger">
          Esta propiedad es sensible. Un valor incorrecto puede cerrar sesiones o impedir el acceso.
        </p>
      )}
      <InputField
        form={form}
        name="value"
        label="Valor"
        multiline={dataType === "json"}
        rows={dataType === "json" ? 8 : 4}
        placeholder={dataType === "json" ? "{\n  \"enabled\": true\n}" : undefined}
      />
      <InputField form={form} name="description" label="Descripcion" />

      <div className="flex justify-end gap-2">
        <Button type="button" variant="secondary" onClick={close}>
          Cancelar
        </Button>
        <Button type="submit" loading={save.isPending}>
          Guardar
        </Button>
      </div>
    </form>
  );
};
