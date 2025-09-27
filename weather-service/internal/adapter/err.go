package adapter

import (
	"github.com/marcelofabianov/fault"
)

var (
	ErrExternalAPICall                   = fault.New("external api call failed", fault.WithCode(fault.Internal))
	ErrExternalAPIUnacceptableStatusCode = fault.New("external api returned an unacceptable status code", fault.WithCode(fault.Internal))
	ErrExternalAPIParse                  = fault.New("failed to parse external api response", fault.WithCode(fault.Internal))
	ErrLocationNotFound                  = fault.New("location not found in external api", fault.WithCode(fault.NotFound))
)
