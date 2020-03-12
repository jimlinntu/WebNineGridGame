package main
import (
    "math/rand"
    "strconv"
    "log"
    "time"
    "context"
    "net/http"
    "github.com/gin-gonic/gin"
    "github.com/gin-gonic/gin/binding"
    "github.com/gin-contrib/cors"
    "fmt"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
    "go.mongodb.org/mongo-driver/mongo/readpref"
    // "io/ioutil"
)

var databaseName = "WebNineGameDB"
var notfound = "Page Not Found"

// A team is a user, admin is also a user
type User struct {
    Account string `json:"account"`
    Password string `json:"password"`
    Token string `json:"token"`
    GridNumbers []int `json:"gridnumbers"`
}

type Question struct {
    Title string `json:"title"`
    Base64Image string `json:"base64image"`
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

func createDatabaseCollection (client *mongo.Client) *mongo.Collection{
    collection := client.Database(databaseName).Collection("users")
    // Clear all users
    _, err := collection.DeleteMany(context.TODO(), bson.D{{}})
    if err != nil {
        log.Fatal(err)
    }
    return collection
}

func initialize_users(collection *mongo.Collection, max_team int){
    team_prefix := "team"
    admin := "admin"

    // Set username
    for i := 1; i <= max_team; i++ {
        password := team_prefix + strconv.Itoa(i) + strconv.Itoa(rand.Intn(10000))
        user := User{
            Account: team_prefix + strconv.Itoa(i),
            Password: password,
        }
        // Insert it into the collection
        _, err := collection.InsertOne(context.TODO(), user)
        if err != nil {
            log.Fatal(err)
        }else{
            log.Printf("Team %d's password is: %s\n", i, password)
        }
    }
    // Insert Admin
    password := admin + strconv.Itoa(rand.Intn(1000000))
    user := User{
            Account: admin,
            Password: password,
    }
    _, err := collection.InsertOne(context.TODO(), user)
    if err != nil {
        log.Fatal(err)
    }else{
        log.Printf("Admin's password is: %s\n", password)
    }
}

func check_authentication(collection *mongo.Collection) gin.HandlerFunc{
    return func(c *gin.Context){
        var user User
        // if error, abort it
        if c.ShouldBindBodyWith(&user, binding.JSON) != nil {
            c.AbortWithStatus(http.StatusNotAcceptable)
            return
        }
        // if token is "", abort it
        if user.Token == "" {
            c.AbortWithStatus(http.StatusNotAcceptable)
            return
        }
        // Search (user.Token == password) user
        filter := bson.D{{"password", user.Token}}
        err := collection.FindOne(context.TODO(), filter).Decode(&user)
        if err != nil {
            log.Fatal(err)
        }
        log.Printf("[*] Get account: %s", user.Account)
        c.Set("account", user.Account)
        c.Next()
        return
    }
}

func (user *User) saveGridNumbers(collection *mongo.Collection) error{
    fmt.Println("haha")
    _, err := collection.UpdateOne(context.TODO(),
        bson.D{
            { "account", user.Account },
        }, bson.D{
            {"$set", bson.D{
                        {"gridnumbers", user.GridNumbers},
                    },
            },
        },
    )
    return err
}

func (user *User) getGridNumbersByAccount(collection *mongo.Collection){
}

func main(){
    //
    rand.Seed(115813)
    // 
    client := createMongoClient()
    collection := createDatabaseCollection(client)
    initialize_users(collection, 10)

    router := gin.Default()
    router.Use(cors.Default())

    private_router := router.Group("/user")
    private_router.Use(check_authentication(collection))

    // private part (authentication required)
    {
        private_router.POST("/get_gridnumbers", func(c *gin.Context){
            // account, _ := c.MustGet("account").(string)

            c.JSON(http.StatusOK, gin.H{
                "gridNumbers": []int{1,2,3,4,5,6,7,8,9},
            })
            return
        })
        private_router.POST("/push_gridnumbers", func(c *gin.Context){
            account, _ := c.MustGet("account").(string)
            var user User
            user.Account = account
            // Plug gridNumbers to the database
            if c.ShouldBindBodyWith(&user, binding.JSON) == nil{
                log.Printf("Get %s's gridNumbers: %v\n", account, user.GridNumbers)
                // Save grid numbers to database
                err := user.saveGridNumbers(collection)
                // if error
                if err != nil {
                    c.String(http.StatusNotFound, notfound)
                    return
                }

                c.JSON(http.StatusOK, gin.H{
                    "gridNumbers": user.GridNumbers,
                })
            }else{
                c.String(http.StatusNotFound, "Page Not Found")
            }
            return
        })
    }

    router.GET("/ping", func(c *gin.Context){
        c.JSON(200, gin.H{
            "message": "pong",
        })
    })

    // Perform authentication
    // TODO
    router.POST("/auth", func(c *gin.Context){
        var user User
        if c.ShouldBindJSON(&user) == nil{
            fmt.Printf("account: %s, password: %s\n", user.Account, user.Password)
            c.JSON(http.StatusOK, gin.H{
                "account": user.Account,
                "password": user.Password,
                "token": user.Password,
                "gridNumbers": []int{}, // TODO
            })
        }else{
            c.String(http.StatusNotFound, "Page Not Found")
        }
        return
    })
    router.Run("0.0.0.0:80")
    destroyMongoClient(client)
}
