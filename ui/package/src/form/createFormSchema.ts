import * as yup from "yup";

const requiredMessage = "Este campo es requerido";
const emailMessage = "Introduce una direccion de correo electronico valida";
const requiredEmailMessage = "El correo electronico es obligatorio";
const invalidMatchMessage = "El valor de este campo es invalido";
const phoneMessage = "Introduce un numero de celular valido";

const transformNullIfEmpty = (value: unknown) => {
  if (value == null) return null;
  if (typeof value !== "string") return value;
  if (value.trim() === "") return null;
  return value.trim();
};

const requiredString = (errorMessage: string = requiredMessage) => {
  return yup.string().transform(transformNullIfEmpty).required(errorMessage);
};

const nullableString = () => {
  return yup.string().transform(transformNullIfEmpty).default(null).nullable();
};

const requiredEmail = (errorMessage: string = requiredEmailMessage) => {
  return yup.string().transform(transformNullIfEmpty).email(emailMessage).required(errorMessage);
};

const nullableEmail = (errorMessage: string = emailMessage) => {
  return yup
    .string()
    .transform(transformNullIfEmpty)
    .test("emailTest", errorMessage, (value) => {
      if (value == null || value === "") return true;
      return yup.string().email().isValidSync(value);
    })
    .default(null)
    .nullable();
};

const requiredMatchString = (
  regex: RegExp,
  matchErrorMessage = invalidMatchMessage,
  requiredErrorMessage = requiredMessage,
) => {
  return yup
    .string()
    .transform(transformNullIfEmpty)
    .trim()
    .test("regexValidation", matchErrorMessage, (value) => {
      if (value == null || value === "") return true;
      return value.match(regex) != null;
    })
    .required(requiredErrorMessage);
};

const nullableOrMatchString = (regex: RegExp, errorMessage = invalidMatchMessage) => {
  return yup
    .string()
    .transform(transformNullIfEmpty)
    .test("regexValidation", errorMessage, (value) => {
      if (value == null || value === "") return true;
      return value.match(regex) != null;
    })
    .default(null)
    .nullable();
};

type StringNumbersOnlyOptions = {
  minLength?: number;
  maxLength?: number;
  errorMessage?: string;
};

const numbersOnly = (options?: StringNumbersOnlyOptions) => {
  const minLength = options?.minLength == null ? undefined : Math.max(options.minLength, 1);
  const maxLength = options?.maxLength == null ? undefined : Math.max(options.maxLength, minLength ?? 1);
  const regex = /^[0-9]+$/;

  let validator = yup.string().transform(transformNullIfEmpty).test("regexValidation", (value, context) => {
    if (value == null || value === "") return true;

    if (value.match(regex) == null) {
      return context.createError({ message: "El campo debe contener solo caracteres numericos" });
    }

    return true;
  });

  if (minLength != null) {
    validator = validator.min(minLength, `El campo debe tener al menos ${minLength} caracteres`);
  }

  if (maxLength != null) {
    validator = validator.max(maxLength, `La longitud maxima permitida es de ${maxLength} caracteres`);
  }

  return validator;
};

const requiredStringNumbersOnly = (options?: StringNumbersOnlyOptions) => {
  return numbersOnly(options).required(options?.errorMessage ?? requiredMessage);
};

const nullableStringNumbersOnly = (options?: StringNumbersOnlyOptions) => {
  return numbersOnly(options).default(null).nullable();
};

const requiredBoolean = (errorMessage: string = requiredMessage) => {
  return yup.boolean().typeError("El tipo de dato debe ser booleano").transform(transformNullIfEmpty).required(errorMessage);
};

const nullableBoolean = () => {
  return yup.boolean().typeError("El tipo de dato debe ser booleano").transform(transformNullIfEmpty).default(null).nullable();
};

const requiredNumber = (errorMessage: string = requiredMessage) => {
  return yup
    .number()
    .transform((value, originalValue) => {
      if (originalValue == null || originalValue === "") return undefined;
      return Number.isNaN(value) ? undefined : value;
    })
    .typeError("Debe ser numerico")
    .required(errorMessage);
};

const nullableNumber = () => {
  return yup
    .number()
    .transform((value, originalValue) => {
      if (originalValue == null || originalValue === "") return null;
      return Number.isNaN(value) ? null : value;
    })
    .typeError("Debe ser numerico")
    .default(null)
    .nullable();
};

const requiredPositiveNumber = (errorMessage: string = requiredMessage) => {
  return yup.number().typeError("Este campo debe ser un valor numero").min(0, "Este campo debe ser un valor numero no negativo").required(errorMessage);
};

const requiredPositiveInteger = (errorMessage: string = requiredMessage) => {
  return requiredPositiveNumber(errorMessage).integer("Este campo debe ser un numero entero");
};

const nullablePositiveNumber = () => {
  return yup
    .number()
    .transform((value, originalValue) => {
      if (originalValue == null || originalValue === "") return null;
      return Number.isNaN(value) ? null : value;
    })
    .typeError("Este campo debe ser un valor numero")
    .min(0, "Este campo debe ser un valor numero no negativo")
    .default(null)
    .nullable();
};

const nullablePositiveInteger = () => {
  return nullablePositiveNumber().integer("Este campo debe ser un numero entero");
};

const phone = (errorMessage: string) => {
  const regex = /^(?:\+51)\s?[0-9]{3}\s?[0-9]{3}\s?[0-9]{3}$/;

  return yup
    .string()
    .transform((value: unknown) => {
      if (value == null) return null;
      if (typeof value !== "string") return value;
      if (value.trim() === "") return null;
      if (value.trim().startsWith("+")) return value.trim();
      return `+51 ${value.trim()}`;
    })
    .test("regexValidation", errorMessage, (value) => {
      if (value == null || value === "") return true;
      return value.match(regex) != null;
    })
    .default(null);
};

const nullableCelular = (errorMessage: string = phoneMessage) => {
  return phone(errorMessage).nullable();
};

const requiredCelular = (errorMessage: string = phoneMessage) => {
  return phone(errorMessage).required(requiredMessage);
};

const requiredStringArray = (errorMessage: string = requiredMessage) => {
  return yup.array().of(requiredString()).min(1, "Debe contener al menos un elemento").required(errorMessage);
};

const nullableStringArray = () => {
  return yup.array().of(nullableString()).default(null).nullable();
};

export const validators = {
  transformNullIfEmpty,
  requiredString,
  nullableString,
  requiredEmail,
  nullableEmail,
  requiredMatchString,
  nullableOrMatchString,
  requiredStringNumbersOnly,
  nullableStringNumbersOnly,
  requiredBoolean,
  nullableBoolean,
  requiredNumber,
  nullableNumber,
  requiredPositiveNumber,
  requiredPositiveInteger,
  nullablePositiveNumber,
  nullablePositiveInteger,
  nullableCelular,
  requiredCelular,
  requiredStringArray,
  nullableStringArray,
};

type ShapeFactory<TShape extends yup.ObjectShape> = (_validators: typeof validators) => TShape;

export const createFormSchema = <TShape extends yup.ObjectShape>(buildShape: ShapeFactory<TShape>) => {
  return yup.object(buildShape(validators));
};
