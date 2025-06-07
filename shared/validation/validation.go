package validation

import (
	"context"
	"encoding/json"
	"net/http"
	"regexp"

	bufvalidate "buf.build/gen/go/bufbuild/protovalidate/protocolbuffers/go/buf/validate"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type validationErrorResponse struct {
	Code    int                     `json:"code"`
	Message string                  `json:"message"`
	Details []validationErrorDetail `json:"details"`
}

type validationErrorDetail struct {
	Field   string `json:"field"`
	Tag     string `json:"tag"`
	Message string `json:"message"`
}

func ErrorHandler(
	ctx context.Context,
	mux *runtime.ServeMux,
	marshaler runtime.Marshaler,
	w http.ResponseWriter,
	r *http.Request,
	err error,
) {
	st, ok := status.FromError(err)
	if !ok {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	if st.Code() == codes.InvalidArgument {
		for _, d := range st.Details() {
			if violations, ok := d.(*bufvalidate.Violations); ok {
				details := make([]validationErrorDetail, 0, len(violations.Violations))
				for _, v := range violations.Violations {
					field := ""
					if len(v.Field.Elements) > 0 && v.Field.Elements[0].FieldName != nil {
						field = *v.Field.Elements[0].FieldName
					}

					re := regexp.MustCompile(`^value`)
					message := re.ReplaceAllString(
						*v.Message,
						cases.Title(language.English, cases.Compact).String(field),
					)

					details = append(details, validationErrorDetail{
						Field:   field,
						Tag:     *v.RuleId,
						Message: message,
					})
				}

				errorResponse := validationErrorResponse{
					Code:    int(codes.InvalidArgument),
					Message: "validation failed",
					Details: details,
				}

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusBadRequest)
				if err := json.NewEncoder(w).Encode(errorResponse); err != nil {
					http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
				}

				return
			}
		}
	}

	runtime.DefaultHTTPErrorHandler(ctx, mux, marshaler, w, r, err)
}
