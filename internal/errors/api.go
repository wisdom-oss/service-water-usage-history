package apiErrors

import "github.com/wisdom-oss/common-go/v3/types"

var ErrInvalidPageSettings = types.ServiceError{
	Type:   "https://www.rfc-editor.org/rfc/rfc9110#section-15.5.1",
	Status: 400,
	Title:  "Missing Consumer ID",
	Detail: "The API could not detect any consumer id in the request",
}

var ErrEmptyConsumerID = types.ServiceError{
	Type:   "https://www.rfc-editor.org/rfc/rfc9110#section-15.5.1",
	Status: 400,
	Title:  "Missing Consumer ID",
	Detail: "The API could not detect any consumer id in the request",
}

var ErrInvalidConsumerID = types.ServiceError{
	Type:   "https://www.rfc-editor.org/rfc/rfc9110#section-15.5.1",
	Status: 400,
	Title:  "Invalid Consumer ID",
	Detail: "The consumer id is not in a valid format",
}

var ErrUnknownConsumer = types.ServiceError{
	Type:   "https://www.rfc-editor.org/rfc/rfc9110#section-15.5.5",
	Status: 404,
	Title:  "Unknown Consumer",
	Detail: "No consumer with the supplied id exists",
}
