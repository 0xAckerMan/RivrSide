package data

type PackagePlan struct{
    CommonFields
    Name string `json:"name"`
    Price int `json:"price"`
}

