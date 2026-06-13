export const ApiErrNone = '';
export const ApiErrValidation = 'Validation';
export const ApiErrNotFound = 'NotFound';
export const ApiErrInternalError = 'InternalError';

export type ApiErr =
	| typeof ApiErrNone
	| typeof ApiErrValidation
	| typeof ApiErrNotFound
	| typeof ApiErrInternalError;

export type ApiResponse = {
	error: ApiErr | string;
};
