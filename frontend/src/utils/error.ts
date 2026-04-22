export const getErrorMessage = (error: unknown, fallback = "Что-то пошло не так"): string => {
  if (error instanceof Error) {
    return error.message;
  }

  return fallback;
};
