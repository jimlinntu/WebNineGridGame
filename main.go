package main
import (
    "net/http"
    "github.com/gin-gonic/gin"
    "github.com/gin-contrib/cors"
    "fmt"
    // "io/ioutil"
)

type User struct {
    Account string `json:"account"`
    Password string `json:"password"`
}

func main(){
    router := gin.Default()
    router.Use(cors.Default())
    router.GET("/ping", func(c *gin.Context){
        c.JSON(200, gin.H{
            "message": "pong",
        })
    })
    router.POST("/auth", func(c *gin.Context){
        var user User
        if c.ShouldBindJSON(&user) == nil{
            fmt.Printf("account: %s, password: %s\n", user.Account, user.Password)
        }
        c.JSON(http.StatusOK, gin.H{
            "account": user.Account,
            "password": user.Password,
        })
    })
    router.Run("0.0.0.0:80")
}
