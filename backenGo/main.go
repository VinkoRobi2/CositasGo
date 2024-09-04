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
	ID string `json:"id"`
}

var usuarios = []User{
    {Nombre: "Juan", Apellido: "Pérez", Email: "juan.perez@example.com",ID: "1"},
    {Nombre: "María", Apellido: "González", Email: "maria.gonzalez@example.com",ID: "2"},
    {Nombre: "Carlos", Apellido: "Ramírez", Email: "carlos.ramirez@example.com",ID: "3"},
    {Nombre: "Ana", Apellido: "Martínez", Email: "ana.martinez@example.com",ID: "4"},
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

func getuserbyid(c*gin.Context){
	id := c.Param("id")
	for _, a:= range usuarios{
		if a.ID == id{
			c.IndentedJSON(http.StatusOK,id)
		}
	}
	c.IndentedJSON(http.StatusNotFound, gin.H{"messag":"Usuario no encontrado"})
}

func main(){
	router := gin.Default()
	router.Use(cors.Default())
	router.GET("/usuariobyid/:id",getuserbyid)
	router.POST("/nuevousuario",PostUsers)
	router.GET("/usuarios",Getusuarios)
	router.Run("localhost:8080") 

}