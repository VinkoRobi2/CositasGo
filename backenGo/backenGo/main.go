package main

import (
	"net/http"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type User struct {
	ID       string `json:"id"`
	Nombre   string `json:"nombre"`
	Apellido string `json:"apellido"`
	Email    string `json:"email"`
}

var users = []User{
	{
		ID:       "1",
		Nombre:   "Mario",
		Apellido: "Perez",
		Email:    "juan.z@example.com",
	},
	{
		ID:       "2",
		Nombre:   "Juan",
		Apellido: "Perez",
		Email:    "juan.perez@example.com",
	},
	{
		ID:       "3",
		Nombre:   "Alberto",
		Apellido: "Perez",
		Email:    "jez@example.com",
	},
}

func Authmiddelware() gin.HandlerFunc{
 return func(c *gin.Context) {
	api_key := "a8023085fc7952e30f599599086d062d"
    headers := c.GetHeader("api_key")
	if headers == ""{
		c.JSON(http.StatusUnauthorized,gin.H{
			"message":"Api_key missing",
		})
		c.Abort()
	}else if headers != api_key{
      c.JSON(http.StatusUnauthorized,gin.H{
		"message":"no estas autorizado",
	  })
	  c.Abort()
	}
	c.Next()
	}}



func GetUsers(c *gin.Context) {
	id := c.Query("id")
	if id != "" {
		for _, u := range users {
			if u.ID == id {
				c.IndentedJSON(http.StatusOK, u)
				return
			}
		}
		c.JSON(http.StatusNotFound, gin.H{
			"message": "No existe un usuario con ese ID",
		})
		return
	}
	c.IndentedJSON(http.StatusOK, users)
}

func FindByApellido(c *gin.Context) {
	apellido := strings.ToLower(c.Param("apellido"))
	for _, a := range users {
		if strings.ToLower(a.Apellido) == apellido {
			c.IndentedJSON(http.StatusOK, a)
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{
		"message": "No se encontró un item con este apellido",
	})
}

func PostUsers(c *gin.Context) {
	var newUser User
	if err := c.BindJSON(&newUser); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Error al procesar la solicitud",
		})
		return
	}
	users = append(users, newUser)
	c.IndentedJSON(http.StatusCreated, newUser)
}

func GetHello(c *gin.Context) {
	c.String(http.StatusOK, "Hello World")
}

func GetElementByID(c *gin.Context) {
	id := c.Param("id")
	for _, u := range users {
		if u.ID == id {
			c.IndentedJSON(http.StatusOK, u)
			return
		}
	}
	c.JSON(http.StatusNotFound, gin.H{
		"message": "No existe un elemento con ese ID",
	})
}

func main() {
	router := gin.Default()


	router.Use(cors.Default())


	authorized := router.Group("/")
	authorized.Use(Authmiddelware()) 

	// Endpoints protegidos
	authorized.GET("/users", GetUsers)
	authorized.GET("/getuser/:apellido", FindByApellido)
	authorized.POST("/users", PostUsers)
	authorized.GET("/user/:id", GetElementByID)

	router.GET("/hello", GetHello)

	router.Run("localhost:8080")

}
