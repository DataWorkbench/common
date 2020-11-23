package qerror

var (
	Internal = &Error{
		code:   "InternalError",
		desc:   "We encountered an internal error. Please try again.",
		status: 500,
	}

	NotAcceptable = &Error{
		code:   "NotAcceptable",
		desc:   "The request is not acceptable.",
		status: 406,
	}

	MethodNotAllowed = &Error{
		code:   "MethodNotAllowed",
		desc:   "The specified method is not allowed against this resource.",
		status: 405,
	}

	// InvalidParams render message if parameter error in request URL-Query.
	InvalidParams = &Error{
		code:   "InvalidParams",
		desc:   "The parameter you provided is invalid.",
		status: 400,
	}
	// InvalidHeaders render message if parameter error in request Headers.
	InvalidHeaders = &Error{
		code:   "InvalidHeaders",
		desc:   "The header you provided is invalid.",
		status: 400,
	}
	// InvalidHeaders render message if parameter error in request Body.
	InvalidJSON = &Error{
		code:   "InvalidJSON",
		desc:   "The request body with json format is invalid.",
		status: 400,
	}
	// InvalidHeaders render message if any parameter error in request without Query|Headers|Body
	InvalidRequest = &Error{
		code:   "InvalidRequest",
		desc:   "The requests you provided is invalid.",
		status: 400,
	}

	// PermissionDenied render message when use have no permission to accessing resource
	PermissionDenied = &Error{
		code:   "PermissionDenied",
		desc:   "You don't have enough permission to accomplish this request",
		status: 403,
	}

	// NotActive render message if accessing resource is inactive and operation not allows.
	NotActive = &Error{
		code:   "NotActive",
		desc:   "The resource you are accessing does not active.",
		status: 403,
	}

	// NotExists render message if accessing resource not exists.
	NotExists = &Error{
		code:   "NotExists",
		desc:   "The resource you are accessing does not exist.",
		status: 404,
	}

	// AlreadyExists render message if creating resource already exists.
	AlreadyExists = &Error{
		code:   "AlreadyExists",
		desc:   "The resource you are creating already exists.",
		status: 409,
	}
)
