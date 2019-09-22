package middleware

import (
	"database/sql"
	"errors"
	"log"
	"net/http"
	"strings"

	"github.com/jacky-htg/inventory/libraries/api"
	"github.com/jacky-htg/inventory/libraries/array"
	"github.com/jacky-htg/inventory/packages/auth/models"
)

func Auths(db *sql.DB, log *log.Logger, allow []string) api.Middleware {
	fn := func(before api.Handler) api.Handler {
		h := func(w http.ResponseWriter, r *http.Request) {
			var isAuth bool
			var str array.ArrString
			isAuth = true

			ctx := r.Context()
			curRoute := ctx.Value("url").(string)

			inArray, _ := str.InArray(curRoute, allow)
			if !inArray {
				var access models.Access
				var err error
				url := r.URL.String()
				controller := strings.Split(url, "/")[1]

				isAuth, err = access.IsAuth(
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

			before(w, r)
		}

		return h
	}

	return fn
}
