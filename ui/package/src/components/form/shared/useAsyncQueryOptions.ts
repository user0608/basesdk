import { useQuery } from "@tanstack/react-query";
import type { QueryKey } from "@tanstack/react-query";

type UseAsyncQueryOptionsArgs<TOption> = {
  queryKey: QueryKey;
  loadOptions: () => Promise<TOption[]>;
  loadOnMount?: boolean;
  loadOnInitialValue?: boolean;
  hasInitialValue?: boolean;
  errorMessage?: string;
};

export const useAsyncQueryOptions = <TOption>({
  queryKey,
  loadOptions,
  loadOnMount = true,
  loadOnInitialValue = false,
  hasInitialValue = false,
  errorMessage = "No se pudieron cargar las opciones",
}: UseAsyncQueryOptionsArgs<TOption>) => {
  const enabled = loadOnMount || (loadOnInitialValue && hasInitialValue);

  const query = useQuery<TOption[]>({
    queryKey,
    queryFn: loadOptions,
    enabled,
  });

  return {
    options: query.data ?? [],
    loaded: query.isFetched,
    error: query.isError ? errorMessage : undefined,
    isPending: query.isPending || query.isFetching,
    refresh: async () => {
      await query.refetch();
    },
  };
};
