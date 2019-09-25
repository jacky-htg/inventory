package middleware

import (
	"context"
	"database/sql"
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/jacky-htg/inventory/libraries/api"
	"github.com/jacky-htg/inventory/libraries/array"
	"github.com/jacky-htg/inventory/models"
)

// Auths middleware to auhtorization checking
func Auths(db *sql.DB, log *log.Logger, allow []string) api.Middleware {
	fn := func(before api.Handler) api.Handler {
		h := func(w http.ResponseWriter, r *http.Request) {
			var isAuth bool
			var str array.ArrString
			var user models.User
			isAuth = true

			ctx := r.Context()
			curRoute := ctx.Value(api.Ctx("url")).(string)

			inArray, _ := str.InArray(curRoute, allow)
			if !inArray {
				var access models.Access
				var err error
				url := r.URL.String()
				controller := strings.Split(url, "/")[1]

				isAuth, user, err = access.IsAuth(
					ctx,
					db,
					r.Header.Get("Token"),
					controller,
					strings.ToUpper(r.Method)+" "+curRoute,
				)

				if err == sql.ErrNoRows {
					log.Printf("ERROR : %+v", err)
					api.ResponseError(w, api.ErrForbidden(errors.New("Forbidden"), ""))
					return
				}

				if err != nil {
					log.Printf("ERROR : %+v", err)
					api.ResponseError(w, err)
					return
				}
			}

			if !isAuth {
				log.Print("ERROR : Forbidden")
				api.ResponseError(w, api.ErrForbidden(errors.New("Forbidden"), ""))
				return
			}

			ctx = context.WithValue(ctx, api.Ctx("auth"), user)
			before(w, r.WithContext(ctx))
		}

		return h
	}

	return fn
}
