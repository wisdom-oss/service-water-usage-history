package routes

import (
	"fmt"

	wisdomType "github.com/wisdom-oss/commonTypes/v2"
)

var ErrPageTooLarge = wisdomType.WISdoMError{
	Type:   "https://www.rfc-editor.org/rfc/rfc9110#section-15.5.14",
	Status: 413,
	Title:  "Page Size Too Large",
	Detail: fmt.Sprintf("Due to limitations on the system side, the selected page size is too large too handle. Please select a value smaller than %d", MaxPageSize),
}
