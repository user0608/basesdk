import { uid } from "uid";

export const uniqueCode = (codigo?: string, len: number = 20) => {
  if (codigo == null || codigo.trim() === "") {
    codigo = uid(len);
  }

  return codigo;
};
