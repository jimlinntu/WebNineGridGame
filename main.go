package main
import (
    "log"
    "time"
    "context"
    "net/http"
    "github.com/gin-gonic/gin"
    "github.com/gin-contrib/cors"
    "fmt"
    // "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
    "go.mongodb.org/mongo-driver/mongo/readpref"
    // "io/ioutil"
)

var databaseName = "WebNineGameDB"

// A team is a user, admin is also a user
type User struct {
    Account string `json:"account"`
    Password string `json:"password"`
    GridNumbers []int `json:"gridnumbers"`
}


func createMongoClient() *mongo.Client{
    client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://172.17.0.1:17990"))
    ctx, _ := context.WithTimeout(context.Background(), 100*time.Second)
    err = client.Connect(ctx)
    ctx, _ = context.WithTimeout(context.Background(), 10*time.Second)
    err = client.Ping(ctx, readpref.Primary())
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("[*] Client Creation Success!")
    return client
}

func destroyMongoClient(client *mongo.Client) {
    err := client.Disconnect(context.TODO())
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("[*] Client Destruction Success!")
}

func createDatabaseCollection (client *mongo.Client){
    _ = client.Database(databaseName).Collection("users")
}


func main(){
    // 
    client := createMongoClient()
    createDatabaseCollection(client)
    destroyMongoClient(client)

    router := gin.Default()
    router.Use(cors.Default())
    router.GET("/ping", func(c *gin.Context){
        c.JSON(200, gin.H{
            "message": "pong",
        })
    })

    router.POST("/push_gridnumbers", func(c *gin.Context){
        var user User
        if c.ShouldBindJSON(&user) == nil{
            fmt.Printf("%v\n", user.GridNumbers)
        }
        c.JSON(http.StatusOK, gin.H{
            "gridNumbers": user.GridNumbers,
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
            "token": "TODOTOKEN!",
            "gridNumbers": []int{}, // TODO
        })
    })
    router.Run("0.0.0.0:80")
}
