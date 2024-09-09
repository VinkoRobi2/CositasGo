package main

import (
	"database/sql"
	"log"
	"net/http"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql" 
	"github.com/gorilla/sessions"
)

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

var db *sql.DB
var store = sessions.NewCookieStore([]byte("Super-secret"))


func init() {
	store.Options.HttpOnly = true
	store.Options.Secure = true
}

func LoginPost(c *gin.Context) {
	var user User
	if err:= c.BindJSON(&user);err != nil{
     c.JSON(http.StatusBadRequest,gin.H{
		"message":"Datos enviados invalidamente",
		"success":false,
	 })
	 return
	}
	var storedpass string
	err := db.QueryRow("SELECT password FROM usuarios WHERE username = ?", user.Username).Scan(&storedpass)
	if err != nil{
		if err == sql.ErrNoRows{
			c.JSON(http.StatusBadRequest,gin.H{
				"message":"Usuario no encontrado",
				"success":false,
			})
			return
		}
		c.JSON(http.StatusInternalServerError,gin.H{
			"message":"Error interno del servidor",
			"success":false,
		})
		return
	}
	if storedpass != user.Password{
		c.JSON(http.StatusUnauthorized,gin.H{
			"message":"Usuario o contra incorrecta",
			"success":false,
		})
		return
	}

	session, _ := store.Get(c.Request, "session")
	session.Values["authenticated"] = true
	session.Values["username"] = user.Username
	session.Save(c.Request, c.Writer)

	c.JSON(http.StatusAccepted, gin.H{
		"message":"inicio de sesion exitoso",
		"success":true,
	})
}


func Authmiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		session, err := store.Get(c.Request, "session")
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": "No se pudo recuperar la sesión",
			})
			c.Abort()
			return
		}
		auth, ok := session.Values["authenticated"].(bool)
		if !ok || !auth {
			c.JSON(http.StatusUnauthorized, gin.H{
				"message": "no estás autorizado",
			})
			c.Abort()
			return
		}
		c.Next()
	}
}

func Profile(c *gin.Context) {
	session, _ := store.Get(c.Request, "session")
	username := session.Values["username"].(string)
	c.JSON(http.StatusOK, gin.H{
		"message":  "Perfil de usuario",
		"username": username,
	})
}

func Register(c*gin.Context){
	var newuser User
	if err:= c.BindJSON(&newuser);err !=nil{
		c.JSON(http.StatusBadRequest,gin.H{
			"message":"Datos mal enviados",
			"success": false,
		})
	    c.Abort()
		return
	}
	if len(newuser.Username) < 3 || len(newuser.Password) < 6 {
        c.JSON(http.StatusBadRequest, gin.H{
            "message": "Nombre de usuario o contraseña demasiado cortos",
			"success": false,
        })
        return
    }
	var exist bool
	err := db.QueryRow("SELECT COUNT(*) > 0 FROM usuarios WHERE username = ?", newuser.Username).Scan(&exist)
	if err != nil{
		c.JSON(http.StatusInternalServerError,gin.H{
			"message":"Ocurrio un error en el servidor",
			"success": false,
		})
		return
	}
	if exist {
		c.JSON(http.StatusConflict,gin.H{
			"message":"El usuario ya existe",
			"success": false,
			"exist":true,
		})
		return
	}
	_, err = db.Exec("INSERT INTO usuarios (username, password) VALUES (?, ?)", newuser.Username, newuser.Password)
	if err!=nil{
		c.JSON(http.StatusBadRequest,gin.H{
			"message":"hubo un error el registar el usuario",
			"success": false,
		})
	}
	c.JSON(http.StatusAccepted,gin.H{
		"message":"Usuario registrado exitosamente",
		"success": true,
	})


}

func main() {
	var err error
	db, err = sql.Open("mysql", "root:admin@tcp(127.0.0.1:3306)/GO")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	router := gin.Default()

	router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:5173"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true, 
	}))

	router.POST("/login", LoginPost)
	router.POST("/register",Register)
	auth := router.Group("/auth")
	auth.Use(Authmiddleware())
	{
		auth.GET("/profile", Profile)
	}

	if err := router.Run(":8081"); err != nil {
		log.Fatal("Error al iniciar el servidor: ", err)
	}
}
