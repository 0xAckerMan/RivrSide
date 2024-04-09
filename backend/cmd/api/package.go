package main

import (
	"net/http"

	"github.com/0xAckerMan/internal/data"
)

func (app *Application) HandleGetAllPackages(w http.ResponseWriter, r *http.Request) {
	var packages []data.PackagePlan

	package_plan := app.DB.Find(&packages)

	if package_plan.RowsAffected == 0 {
		app.writeJSON(w, http.StatusOK, envelope{"package": "No package plan found at the moment"}, nil)
		return
	}

	if package_plan.Error != nil {
		app.serverErrorResponse(w, r, ErrRecordNotFound)
		return
	}

	err := app.writeJSON(w, http.StatusOK, envelope{"packages": &packages}, nil)
	if err != nil {
		app.serverErrorResponse(w,r,err)
		return
	}
}

func (app *Application) HandleGetSinglePackage(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDparam(w, r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	var plan data.PackagePlan

    result := app.DB.First(&plan, id)

    if result.Error != nil{
        app.noRecordFoundResponse(w,r)
        return
    }

	err = app.writeJSON(w, 200, envelope{"plan": plan}, nil)
	if err != nil {
		app.serverErrorResponse(w,r,err)
		return
	}
}

func (app *Application) HandleCreatePackagePlan(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name  string `json:"name"`
		Price int    `json:"price"`
	}

    err := app.readJSON(w,r,&input)
    if err != nil{
        app.serverErrorResponse(w,r,err)
        return
    }

    package_plan := &data.PackagePlan{
        Name: input.Name,
        Price: input.Price,
    }

    err = app.DB.Create(&package_plan).Error
    if err != nil{
        app.serverErrorResponse(w,r,err)
    }

    err = app.writeJSON(w,http.StatusCreated,envelope{"package_plan": package_plan}, nil)
    if err != nil{
        app.serverErrorResponse(w,r,err)
        return
    }
}

func (app *Application) HandleUpdatePackagePlan(w http.ResponseWriter, r *http.Request) {
    id, err := app.readIDparam(w,r)
    if err != nil{
        app.notFoundResponse(w,r)
        return
    }

    var package_plan data.PackagePlan

    if app.DB.First(&package_plan, id).Error != nil{
        app.notFoundResponse(w,r)
        return
    }

	var input struct{
        Name *string `json:"name"`
        Price *int `json:"price"`
    }

    err = app.readJSON(w,r,&input)
    if err != nil {
        app.serverErrorResponse(w,r,err)
        return
    }

    if input.Name != nil{
        package_plan.Name = *input.Name
    }

    if input.Price != nil{
        package_plan.Price = *input.Price
    }

    err = app.DB.Save(&package_plan).Error
    if err != nil{
        app.serverErrorResponse(w,r,err)
        return
    }

    err = app.writeJSON(w, http.StatusOK, envelope{"package_plan": &package_plan},nil)
    if err !=nil{
        app.serverErrorResponse(w,r,err)
        return
    }
}

func (app *Application) HandleDeletePackagePlan(w http.ResponseWriter, r *http.Request) {
    id, err := app.readIDparam(w,r)
    if err != nil{
        app.notFoundResponse(w,r)
        return
    }

    var package_plan data.PackagePlan

    result := app.DB.Delete(&package_plan, id)
    if result.RowsAffected == 0{
        app.noRecordFoundResponse(w,r)
        return
    }

    if result.Error != nil{
        app.serverErrorResponse(w,r,err)
        return
    }

    err = app.writeJSON(w,http.StatusOK,envelope{"success": "record deleted successfuly"},nil)
    if err != nil{
        app.serverErrorResponse(w,r, err)
        return
    }
}
