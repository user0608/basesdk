import { useMutation, useQueryClient } from "@tanstack/react-query";
import type { MutationFunction, QueryKey, UseMutationOptions } from "@tanstack/react-query";
import { useConfirmDialog } from "../components/dialog/useConfirmDialog";
import { useToast } from "../components/toast/useToast";

export type UseMutateOptions<TData, TError, TVariables, TContext> = {
  mutationFn: MutationFunction<TData, TVariables>;
  mutationOptions?: Omit<
    UseMutationOptions<TData, TError, TVariables, TContext>,
    "mutationFn" | "onSuccess" | "onError" | "onSettled"
  >;
  onSuccess?: (data: TData, variables: TVariables, context: TContext | undefined) => Promise<unknown> | unknown;
  onError?: (error: TError, variables: TVariables, context: TContext | undefined) => Promise<unknown> | unknown;
  onSettled?: (
    data: TData | undefined,
    error: TError | null,
    variables: TVariables,
    context: TContext | undefined,
  ) => Promise<unknown> | unknown;
  invalidateQueryKey?: QueryKey | QueryKey[];
  requireConfirm?: boolean;
  confirmTitle?: string;
  confirmMessage?: string;
  confirmLabel?: string;
  cancelLabel?: string;
  loadingMessage?: string;
  successMessage?: string;
  getErrorMessage?: (error: TError) => string;
};

const toQueryKeys = (queryKey: QueryKey | QueryKey[] | undefined): QueryKey[] => {
  if (!queryKey) return [];
  if (Array.isArray(queryKey) && queryKey.every(Array.isArray)) return queryKey as QueryKey[];
  return [queryKey as QueryKey];
};

export const useMutate = <TData = unknown, TError = Error, TVariables = void, TContext = unknown>({
  mutationFn,
  mutationOptions,
  onSuccess,
  onError,
  onSettled,
  invalidateQueryKey,
  requireConfirm = false,
  confirmTitle = "Confirmar",
  confirmMessage = "Esta seguro de que desea proceder con esta operacion?",
  confirmLabel = "Confirmar",
  cancelLabel = "Cancelar",
  loadingMessage = "Procesando la solicitud...",
  successMessage = "La operacion se completo con exito.",
  getErrorMessage,
}: UseMutateOptions<TData, TError, TVariables, TContext>) => {
  const confirm = useConfirmDialog();
  const toast = useToast();
  const queryClient = useQueryClient();
  const mutation = useMutation<TData, TError, TVariables, TContext>({
    ...mutationOptions,
    mutationFn,
  });

  const mutate = async (variables: TVariables) => {
    if (requireConfirm) {
      const confirmed = await confirm({
        title: confirmTitle,
        content: confirmMessage,
        confirmLabel,
        cancelLabel,
      });

      if (!confirmed) return undefined;
    }

    let data: TData | undefined;
    let error: TError | null = null;
    let context: TContext | undefined;
    const toastId = toast.loading(loadingMessage);

    try {
      data = await mutation.mutateAsync(variables, {
        onSuccess: (_data, _variables, _context) => {
          context = _context;
        },
      });

      for (const queryKey of toQueryKeys(invalidateQueryKey)) {
        await queryClient.invalidateQueries({ queryKey });
      }

      await onSuccess?.(data, variables, context);
      toast.success(successMessage, { id: toastId });
      return data;
    } catch (caught) {
      error = caught as TError;
      await onError?.(error, variables, context);
      toast.error(getErrorMessage?.(error) ?? (error instanceof Error ? error.message : "No se pudo completar la operacion."), {
        id: toastId,
      });
      throw caught;
    } finally {
      await onSettled?.(data, error, variables, context);
    }
  };

  return {
    mutate,
    mutation,
    data: mutation.data,
    error: mutation.error,
    isError: mutation.isError,
    isPending: mutation.isPending,
    isSuccess: mutation.isSuccess,
    reset: mutation.reset,
  };
};
