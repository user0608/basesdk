import { useServices } from "../services/ServiceProvider";

export const useSystemService = () => useServices().system;
