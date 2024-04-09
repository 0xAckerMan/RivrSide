package data

type PackagePlan struct{
    CommonFields
    Name string `json:"name"`
    Price int `json:"price"`
}

func NewPackage() *PackagePlan{
    return &PackagePlan{
        Name: "Room only",
        Price: 4700,
    }
}
