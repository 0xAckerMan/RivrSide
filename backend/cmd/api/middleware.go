package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/0xAckerMan/cmd/utils"
	"github.com/0xAckerMan/internal/data"
)

func (app *Application) UserMiddleware (next http.Handler) http.Handler{
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        var token string

        cookie, err := r.Cookie("token")
        AuthorizationHeader := r.Header.Get("Authorization")
        fields := strings.Fields(AuthorizationHeader)

        if len(fields) != 0 && fields[0] == "Bearer"{
            token = fields[1]
        } else{
            if err != nil{
                token = cookie.Value
            }
        }

        sub, err := utils.ValidateToken(token, os.Getenv("JWT_SECRET"))
        if err != nil{
            app.UnAuthorizedResponse(w,r)
            return
        }
        fmt.Println(sub)

        var user data.User
        res := app.DB.First(&user, "id = ?", sub)
        if res.Error != nil{
            app.UnAuthorizedResponse(w,r)
            return
        }

        ctx := r.Context()
        ctx = context.WithValue(ctx, "user", user)
        next.ServeHTTP(w,r.WithContext(ctx))
    })
} 
