package main

import (
    "github.com/gin-gonic/gin"
    "net/http"
	"github.com/gin-contrib/cors"
)

type User struct {
    Nombre   string `json:"nombre"`
    Apellido string `json:"apellido"`
    Email    string `json:"email"`
}

var usuarios = []User{
    {Nombre: "Juan", Apellido: "Pérez", Email: "juan.perez@example.com"},
    {Nombre: "María", Apellido: "González", Email: "maria.gonzalez@example.com"},
    {Nombre: "Carlos", Apellido: "Ramírez", Email: "carlos.ramirez@example.com"},
    {Nombre: "Ana", Apellido: "Martínez", Email: "ana.martinez@example.com"},
}
func Getusuarios(c*gin.Context){
	c.IndentedJSON(http.StatusOK,usuarios)
}
func PostUsers(c*gin.Context){
	var newuser User
	if err := c.BindJSON(&newuser); err != nil{
		return
	}
	usuarios = append(usuarios, newuser)
	c.IndentedJSON(http.StatusOK,usuarios)
}

func main(){
	router := gin.Default()
	router.Use(cors.Default())
	router.POST("/nuevousuario",PostUsers)
	router.GET("/usuarios",Getusuarios)
	router.Run("localhost:8080") 

}