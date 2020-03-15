package main
import (
    "encoding/csv"
    "encoding/base64"
    "io"
    "io/ioutil"
    "bufio"
    "path"
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
    "github.com/googollee/go-socket.io"
    // "io/ioutil"
)

var databaseName = "WebNineGameDB"
var notfound = "Page Not Found"
var DATA_FOLDER = "./data"

// A team is a user, admin is also a user
type User struct {
    Account string `json:"account"`
    Password string `json:"password"`
    Token string `json:"token"`
    GridNumbers []int `json:"gridnumbers"`
    QuestionIndex int `json:"questionindex"`
    QuestionFinishedMask []bool `json:"questionfinishedmask"`
    AnswerText string `json:"answertext"`
    AnswerBase64Str string `json:"answerbase64str"`
}

type Question struct {
    Description string `json:"description"`
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

func createDatabaseCollection (client *mongo.Client, to_clear_collection bool) (*mongo.Collection, *mongo.Collection){
    collection := client.Database(databaseName).Collection("users")
    question_collection := client.Database(databaseName).Collection("questions")
    // Clear all users
    if to_clear_collection {
        log.Printf("[*] Clearing collections")
        _, err := collection.DeleteMany(context.TODO(), bson.D{{}})
        if err != nil {
            log.Fatal(err)
        }
    }
    return collection, question_collection
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

func (user * User)findUser(collection *mongo.Collection) bool{
    var result User
    filter := bson.D{
            {"account", user.Account},
            {"password", user.Password},
        }

    err := collection.FindOne(context.TODO(), filter).Decode(&result)
    if err != nil{
        // this account is not found
        log.Printf("[!] Cannot find %s", user.Account)
        return false
    }
    // Copy result field values to this user
    log.Printf("[!] User %s is logined!", user.Account)
    *user = result // default copy
    if user.GridNumbers == nil{
        user.GridNumbers = make([]int, 0) // avoid gin.H bug
    }
    return true
}

func (user *User) saveGridNumbersAndIntialize(collection *mongo.Collection) bool{
    question_index := -1 // have not chosen a question
    // Initiailize question mask
    question_finished_mask := make([]bool, 9)
    // check if user's grid numbers is valid
    log.Printf("User %s submitted grid numbers are: %v", user.Account, user.GridNumbers)
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
                        {"questionindex", question_index},
                        {"questionfinishedmask", question_finished_mask},
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
    user.QuestionFinishedMask = question_finished_mask
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
        user.GridNumbers = make([]int, 0) // avoid gin.H bug
        log.Printf("[*] User %s's grid numbers have not been submitted yet", user.Account)
        return false
    }
    log.Printf("[*] User %s's grid numbers are: %v, question finished mask are %v", user.Account, user.GridNumbers, user.QuestionFinishedMask)
    return true
}

func (user *User) saveAnswer(){
    // Find this user by token(password) and update its answer
}

func (user *User) updateQuestionIndex(collection *mongo.Collection) bool{
    filter := bson.D{
            {"account", user.Account},
            {"questionindex", -1}, // update questionindex only when this questionindex == -1
        }
    update := bson.D{
            {"$set", bson.D{{"questionindex", user.QuestionIndex}}},
        }
    err := collection.FindOneAndUpdate(context.TODO(), filter, update).Err()
    if err != nil {
        log.Printf("[!] User %s try to update question index to %d but failed! (It may due to qidx != -1)", user.Account, user.QuestionIndex)
        return false
    }
    return true
}

func initialize_socket_io() (*socketio.Server){
    server, err := socketio.NewServer(nil)
    if err != nil {
        log.Fatal(err)
    }
    //
    server.OnConnect("/", func(s socketio.Conn) error{
        s.SetContext("")
        log.Printf("[s] socket **id == %s** is connected", s.ID())
        return nil
    })
    server.OnDisconnect("/", func(s socketio.Conn, reason string){
        log.Printf("[s] socket **id == %s**  connection closed, reason: %s", s.ID(), reason)
    })
    go server.Serve()
    return server
}

func load_image_as_base64_string(image_path string) string{
    file, err := os.Open(image_path)
    defer file.Close()
    if err != nil {
        log.Fatal(err)
    }
    //
    reader := bufio.NewReader(file)
    content, err := ioutil.ReadAll(reader)
    if err != nil {
        log.Fatal(err)
    }
    encoded := base64.StdEncoding.EncodeToString(content)
    return encoded
}


func load_question_from_csv(filename string) []Question{
    file, err := os.Open(path.Join(DATA_FOLDER, filename))
    if err != nil {
        log.Fatal(err)
    }
    defer file.Close()
    var questions []Question
    r := csv.NewReader(file)
    for {
        record, err := r.Read()
        if err == io.EOF {
            break
        }
        if err != nil {
            log.Fatal(err)
        }

        description, image_path := record[0], record[1]

        questions = append(questions, Question{Description: description,
                                               Base64Image: load_image_as_base64_string(path.Join(DATA_FOLDER, image_path)),
                                           })
    }
    return questions
}

func main(){
    to_clear_collection := os.Args[1]
    //
    rand.Seed(115813)
    // 
    client := createMongoClient()
    user_collection, _ := createDatabaseCollection(client, to_clear_collection == "true")
    if to_clear_collection == "true" {
        initialize_users(user_collection, 10)
    }
    // TODO: use questions
    _ = load_question_from_csv("questions.csv")
    server := initialize_socket_io()
    defer server.Close()

    router := gin.Default()
    config := cors.DefaultConfig()
    config.AllowOrigins = []string{"http://localhost:8080"}
    config.AllowCredentials = true

    router.Use(cors.New(config))
    // Socket initialization
    router.GET("/socket/", gin.WrapH(server))
    router.POST("/socket/", gin.WrapH(server))

    private_router := router.Group("/user")
    private_router.Use(check_authentication(user_collection))

    // private part (authentication required)
    {
        // socket io
        private_router.POST("/get_gridnumbers", func(c *gin.Context){
            var user User
            account, _ := c.MustGet("account").(string)
            user.Account = account
            found := user.getGridNumbersByAccount(user_collection)

            if !found {
                c.AbortWithStatus(http.StatusNotFound)
                return
            }

            // Return grid numbers to the front end client
            c.JSON(http.StatusOK, gin.H{
                "gridNumbers": user.GridNumbers,
                "questionIndex": user.QuestionIndex,
                "question_finished_mask": user.QuestionFinishedMask,
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
                success := user.saveGridNumbersAndIntialize(user_collection)
                // if error
                if !success {
                    c.String(http.StatusNotFound, notfound)
                    return
                }

                c.JSON(http.StatusOK, gin.H{
                    "gridNumbers": user.GridNumbers,
                    "questionIndex": user.QuestionIndex,
                    "question_finished_mask": user.QuestionFinishedMask,
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
        private_router.GET("/get_current_answer", func(c *gin.Context){
            // TODO:
        })

        // select question
        private_router.POST("/select_question", func(c *gin.Context){
            // Update question Index
            account, _ := c.MustGet("account").(string)
            var user User
            if c.ShouldBindBodyWith(&user, binding.JSON) == nil{
                user.Account = account // assign its account
                success := user.updateQuestionIndex(user_collection)
                if success {
                    c.JSON(http.StatusOK, gin.H{
                        "questionIndex": user.QuestionIndex,
                    })
                }else{
                    c.String(http.StatusNotAcceptable, "")
                }
            }else{
                c.String(http.StatusNotAcceptable, "Select question fail!")
            }
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
            log.Printf("[*] Verifying account: %s, password: %s\n", user.Account, user.Password)
            success := user.findUser(user_collection)
            if !success {
                c.String(http.StatusNotAcceptable, "This account is not found")
            }
            c.JSON(http.StatusOK, gin.H{
                "account": user.Account,
                "password": user.Password,
                "token": user.Password,
            })
        }else{
            c.String(http.StatusNotFound, "Page Not Found")
        }
        return
    })
    router.Run("0.0.0.0:80")
    destroyMongoClient(client)
}
