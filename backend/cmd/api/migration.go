package main

import "github.com/0xAckerMan/internal/data"

func (app *Application) migrations() {
	app.DB.AutoMigrate(
		&data.Role{},
		&data.Room{},
		&data.Payment{},
		&data.PackagePlan{},
		&data.User{},
        &data.Subscription{},
	)
}
