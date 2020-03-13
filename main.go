package main
import (
    "os"
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
    QuestionOrder []int `json:"questionorder"`
    QuestionIndex int `json:"questionindex"`
    AnswerText string `json:"answertext"`
    AnswerBase64Str string `json:"answerbase64str"`
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

func createDatabaseCollection (client *mongo.Client, to_clear_collection bool) *mongo.Collection{
    collection := client.Database(databaseName).Collection("users")
    // Clear all users
    if to_clear_collection {
        _, err := collection.DeleteMany(context.TODO(), bson.D{{}})
        if err != nil {
            log.Fatal(err)
        }
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
            log.Printf("[*] Did not find any team match password: %s", user.Token)
            c.AbortWithStatus(http.StatusNotAcceptable)
            return
        }
        log.Printf("[*] Get account: %s", user.Account)
        c.Set("account", user.Account)
        c.Next()
        return
    }
}

func (user *User) saveGridNumbersAndIntialize(collection *mongo.Collection) bool{
    var question_order []int
    question_index := 0
    for i:= 0; i < 9; i++ {
        question_order = append(question_order, i)
    }
    // shuffle question order
    rand.Shuffle(len(question_order), func(i, j int){
        question_order[i], question_order[j] = question_order[j], question_order[i]
    })
    // check if user's grid numbers is valid
    log.Printf("User %s submitted grid numbers are: %v", user.GridNumbers)
    if len(user.GridNumbers) != 9 {
        log.Printf("Grid numbers' length should be 9!")
        return false
    }
    _, err := collection.UpdateOne(context.TODO(),
        bson.D{
            { "account", user.Account },
        }, bson.D{
            {"$set", bson.D{
                        {"gridnumbers", user.GridNumbers},
                        {"questionorder", question_order},
                        {"questionindex", question_index},
                    },
            },
        },
    )
    // Fail
    if err != nil {
        log.Printf("Something went wrong when user %s try to submit his grid numbers, his grid numbers are: %v", user.Account, user.GridNumbers)
        return false
    }
    // Success
    user.QuestionOrder = question_order
    user.QuestionIndex = question_index
    return true
}

func (user *User) getGridNumbersByAccount(collection *mongo.Collection) bool{
    log.Printf("[*] User %s try to get grid numbers", user.Account)
    err := collection.FindOne(context.TODO(),
            bson.D{{"account", user.Account}},
        ).Decode(user)

    if err != nil {
        // Team not found or database error
        return false
    }

    if user.GridNumbers == nil {
        log.Printf("[*] User %s's grid numbers have not been submitted yet", user.Account)
        return false
    }
    log.Printf("[*] User %s's grid numbers are: %v", user.Account, user.GridNumbers)
    return true
}

func (user *User) saveAnswer(){
    // Find this user by token(password) and update its answer
}

func main(){
    to_clear_collection := os.Args[1]
    //
    rand.Seed(115813)
    // 
    client := createMongoClient()
    collection := createDatabaseCollection(client, to_clear_collection == "true")
    initialize_users(collection, 10)

    router := gin.Default()
    router.Use(cors.Default())

    private_router := router.Group("/user")
    private_router.Use(check_authentication(collection))

    // private part (authentication required)
    {
        private_router.POST("/get_gridnumbers", func(c *gin.Context){
            var user User
            account, _ := c.MustGet("account").(string)
            user.Account = account
            found := user.getGridNumbersByAccount(collection)

            if !found {
                c.AbortWithStatus(http.StatusNotFound)
                return
            }

            // Return grid numbers to the front end client
            c.JSON(http.StatusOK, gin.H{
                "gridNumbers": user.GridNumbers,
                "questionOrder": user.QuestionOrder,
                "questionIndex": user.QuestionIndex,
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
                success := user.saveGridNumbersAndIntialize(collection)
                // if error
                if !success {
                    c.String(http.StatusNotFound, notfound)
                    return
                }

                c.JSON(http.StatusOK, gin.H{
                    "gridNumbers": user.GridNumbers,
                    "questionOrder": user.QuestionOrder,
                    "questionIndex": user.QuestionIndex,
                })
            }else{
                c.String(http.StatusNotFound, "Page Not Found")
            }
            return
        })

        private_router.POST("/push_answer", func (c *gin.Context){
            // account, _ := c.MustGet("account").(string)
            var user User
            if c.ShouldBindBodyWith(&user, binding.JSON) == nil {
                // TODO: save answer text and image base64 string to database
                c.String(http.StatusOK, "Answer Received")
            }else{
                c.String(http.StatusNotAcceptable, "Input is weird")
            }
        })

        // A user can check what he have answered for this question
        private_router.GET("/get_previous_answer", func(c *gin.Context){
            // TODO:
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
