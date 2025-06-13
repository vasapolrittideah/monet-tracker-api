package validator

import (
	"context"
	"strings"

	"github.com/charmbracelet/log"
	"github.com/go-playground/locales/en"
	"github.com/go-playground/locales/th"
	ut "github.com/go-playground/universal-translator"
	v "github.com/go-playground/validator/v10"
	en_translations "github.com/go-playground/validator/v10/translations/en"
	th_translations "github.com/go-playground/validator/v10/translations/th"
	"github.com/vasapolrittideah/money-tracker-api/shared/httperror"
	"google.golang.org/grpc/metadata"
)

var (
	val *v.Validate
	uni *ut.UniversalTranslator
)

func Init() {
	enLocale := en.New()
	thLocale := th.New()

	uni = ut.New(enLocale, enLocale, thLocale)
	val = v.New()

	enTrans, found := uni.GetTranslator("en")
	if !found {
		log.Fatalf("translator for 'en' not found")
	}
	if err := en_translations.RegisterDefaultTranslations(val, enTrans); err != nil {
		log.Fatalf("failed to register en translations: %v", err)
	}

	thTrans, found := uni.GetTranslator("th")
	if !found {
		log.Fatalf("translator for 'th' not found")
	}
	if err := th_translations.RegisterDefaultTranslations(val, thTrans); err != nil {
		log.Fatalf("failed to register th translations: %v", err)
	}
}

func getLang(ctx context.Context) string {
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		keys := []string{"accept-language", "grpcgateway-accept-language"}
		for _, key := range keys {
			if langs := md.Get(key); len(langs) > 0 {
				return strings.Split(langs[0], ",")[0]
			}
		}
	}
	return "en"
}

func ValidateInput(ctx context.Context, input any) *httperror.HTTPValidationError {
	err := val.Struct(input)
	if err == nil {
		return nil
	}

	lang := getLang(ctx)
	trans, _ := uni.GetTranslator(lang)

	var details []httperror.ValidationError
	for _, fe := range err.(v.ValidationErrors) {
		details = append(details, httperror.ValidationError{
			Field:   fe.Field(),
			Message: fe.Translate(trans),
		})
	}

	res := httperror.NewValidationError(details)

	return &res
}
