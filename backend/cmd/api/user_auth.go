package main

import (
	"net/http"
	"os"
	"time"

	"github.com/0xAckerMan/cmd/utils"
	"github.com/0xAckerMan/internal/data"
	"golang.org/x/crypto/bcrypt"
)

func (app *Application) UserLogin(w http.ResponseWriter, r *http.Request){
    var login data.Login

    err := app.readJSON(w,r,&login)
    if err != nil{
        app.serverErrorResponse(w,r,err)
        return
    }

    var user *data.User

    result := app.DB.Where("email = ?", login.Email).First(&user)
    if result.Error != nil{
        app.serverErrorResponse(w,r,result.Error)
        return
    }

    if result.RowsAffected == 0{
        message := "email or password is wrong, check and try again"
        app.errorResponse(w,r,http.StatusNotFound,envelope{"error": message})
    }
    
    err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(login.Password))
    if err != nil{
        message := "check your password and try again"
        app.errorResponse(w,r,http.StatusUnprocessableEntity,envelope{"error":message})
        return
    }

    var ttl = time.Hour * 24 * 14
    token, err := utils.GenerateToken(ttl, user.ID, os.Getenv("JWT_SECRET"))
    if err != nil{
        app.serverErrorResponse(w,r,err)
        return
    }

    http.SetCookie(w,&http.Cookie{
        Name: "token",
        Value: token,
        Expires: time.Now().Add(ttl),
        HttpOnly: true,
        MaxAge: int(ttl.Seconds()),
        Domain: "localhost",
        Path: "/",
        SameSite: http.SameSiteNoneMode,
        Secure: true,
    })

    err = app.writeJSON(w,http.StatusOK,envelope{"token": token},nil)
    if err != nil{
        app.serverErrorResponse(w,r,err)
        return
    }

    http.Redirect(w,r,"/",http.StatusOK)
}

func (app *Application) SignOut(w http.ResponseWriter, r *http.Request){
    http.SetCookie(w,&http.Cookie{
        Name: "token",
        Value: "",
        Expires: time.Now(),
        HttpOnly: true,
        MaxAge: -1,
        Domain: "localhost",
        Path: "/",
        SameSite: http.SameSiteNoneMode,
        Secure: true,
    })

    http.Redirect(w,r,"/",http.StatusSeeOther)
}
