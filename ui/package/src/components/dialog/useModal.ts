import { useContext } from "react";
import { DialogContext } from "./DialogProvider";

export const useModal = () => {
  const context = useContext(DialogContext);

  if (!context) {
    throw new Error("useModal debe usarse dentro de DialogProvider");
  }

  return context;
};
