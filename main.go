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
    "github.com/gin-gonic/gin"
    "github.com/gin-gonic/gin/binding"
    "github.com/gin-contrib/cors"
    "github.com/gin-contrib/static"
    "fmt"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"
    "go.mongodb.org/mongo-driver/mongo/readpref"
    "github.com/googollee/go-socket.io"
    "github.com/googollee/go-engine.io"
    "net/http"
    "github.com/googollee/go-engine.io/transport/polling"
    "github.com/googollee/go-engine.io/transport/websocket"
    "github.com/googollee/go-engine.io/transport"
)

var databaseName = "WebNineGameDB"
var notfound = "Page Not Found"
var DATA_FOLDER = "./data"
var ADMIN = "admin"

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

func createDatabaseCollection (client *mongo.Client, to_clear_collection bool) (*mongo.Collection){
    collection := client.Database(databaseName).Collection("users")
    // Clear all users
    if to_clear_collection {
        log.Printf("[*] Clearing collections")
        _, err := collection.DeleteMany(context.TODO(), bson.D{{}})
        if err != nil {
            log.Fatal(err)
        }
    }
    return collection
}

func initialize_users(collection *mongo.Collection, max_team int){
    team_prefix := "team"
    admin := ADMIN
    // Write to the file
    f, err := os.Create("./account.txt")
    if err != nil{
        log.Fatal(err)
    }
    defer f.Close()

    // Set username
    for i := 1; i <= max_team; i++ {
        password := team_prefix + strconv.Itoa(i) + strconv.Itoa(rand.Intn(10000))
        account := team_prefix + strconv.Itoa(i)
        user := User{
            Account: account,
            Password: password,
            QuestionIndex : -1,
        }
        _, err := f.WriteString(account + " " + password + "\n")
        if err != nil{
            log.Fatal(err)
        }
        // Insert it into the collection
        _, err = collection.InsertOne(context.TODO(), user)
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
    _, err = f.WriteString(admin + " " + password + "\n")
    if err != nil{
        log.Fatal(err)
    }
    _, err = collection.InsertOne(context.TODO(), user)
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
            { "gridnumbers", nil}, // only update gridnumbers that have not initiailized!
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

func (user *User) saveAnswer(collection *mongo.Collection) bool{
    // Find this user by account and update its answer if its questionIndex != -
    filter := bson.D{
            {"account", user.Account },
            {"questionindex", bson.D{{"$ne", -1}}}, // questionIndex should not equal to -1
        }

    // Register this team's answer text and image
    update := bson.D{
            {"$set", bson.D{{"answertext", user.AnswerText}, {"answerbase64str", user.AnswerBase64Str}}},
        }

    err := collection.FindOneAndUpdate(context.TODO(), filter, update).Err()
    if err != nil {
        log.Printf("[*] User %s try to saveAnswer but failed", user.Account)
        return false
    }
    log.Printf("[*] User %s's answer text: %s successfully saved!", user.Account, user.AnswerText)
    return true
}

func (user *User) getAnswer(collection *mongo.Collection) bool {
    filter := bson.D{
            {"account", user.Account},
            {"questionindex", bson.D{{"$ne", -1}}},
        }
    err := collection.FindOne(context.TODO(), filter).Decode(user)
    if err != nil{
        log.Printf("[*] User %s getAnswer attempt failed!", user.Account)
        return false
    }
    return true
}

func (user *User) updateQuestionIndex(collection *mongo.Collection) bool{
    filter := bson.D{
            {"account", user.Account},
            {"questionindex", -1}, // update questionindex only when this questionindex == -1
        }
    if user.QuestionIndex < 0 || user.QuestionIndex > 8{
        log.Printf("[!] user.QuestionIndex should be in [0, 8], but the user give %d", user.QuestionIndex)
        return false
    }
    log.Printf("[*] User %s selected question index is %d", user.Account, user.QuestionIndex)
    update := bson.D{
            {"$set", bson.D{{"questionindex", user.QuestionIndex}}},
        }
    var before_updated_user User
    err := collection.FindOneAndUpdate(context.TODO(), filter, update).Decode(&before_updated_user)
    if err != nil {
        log.Printf("[!] User %s try to update question index to %d but failed! (It may due to current database qidx != -1)", user.Account, user.QuestionIndex)
        return false
    }
    user.GridNumbers = before_updated_user.GridNumbers
    return true
}

// Admin gets all information
func getAll(collection *mongo.Collection, questions []Question) ([]*User, []*Question, bool){
    var users []*User
    var chosen_questions []*Question
    filter := bson.D{
        {
            "account", bson.D{{"$ne", ADMIN}}, // account should not be equal to ADMIN
        },
    }
    cursor, err := collection.Find(context.TODO(), filter)
    if err != nil {
        log.Printf("Somethings went wrong in getAll()")
        return nil, nil, false
    }
    for cursor.Next(context.TODO()){
        var user User
        cursor.Decode(&user)
        if user.GridNumbers == nil {
            user.GridNumbers = make([]int, 0) // avoid gin.H empty slice bug
        }
        // append that user's memory address
        users = append(users, &user)
        if user.QuestionIndex != -1 && user.GridNumbers != nil{
            chosen_questions = append(chosen_questions, &questions[user.GridNumbers[user.QuestionIndex]-1]) // return that question's memory address
        }else{
            chosen_questions = append(chosen_questions, nil) // append nil
        }
    }
    // Get each user's corresponding question
    return users, chosen_questions, true
}

func resetAll(collection *mongo.Collection) bool{
    // reset GridNumbers
    filter := bson.D{
        {"account", bson.D{{"$ne", ADMIN}}},
    }
    update := bson.D{
        {"$set", bson.D{
                {"gridnumbers", nil},
                {"questionindex", -1},
                {"questionfinishedmask", nil},
                {"answertext", ""},
                {"answerbase64str", ""},
            },
        },
    }
    if result, err := collection.UpdateMany(context.TODO(), filter, update); err != nil || result.MatchedCount == 0{
        return false
    }
    return true
}

// https://github.com/simagix/mongo-go-examples/blob/master/examples/transaction_test.go
func approve_answer(user_account string, collection *mongo.Collection) bool{
    // Find questionindex != -1 and update it to 1 (act as a lock)
    var user User
    filter := bson.D{
        {"account", user_account},
        {"questionindex", bson.D{{"$ne", -1}}}, // to be approved answer cannot be -1!
    }
    update := bson.D{
        {"$set", bson.D{
            {"questionindex", -1},
        }}, // reset questionindex to -1
    }
    // NOTE: Even if two processes race into this line, it will only be one process succeeded!
    // (Because one process must first change questionindex to -1!)
    err := collection.FindOneAndUpdate(context.TODO(), filter, update).Decode(&user)
    // If this operation failed, it means it is a race condition,
    // we give up the below updating question_finished_mask operation
    if err != nil{
        log.Printf("[*] FindOneAndUpdate in approve_answer failed!")
        return false
    }
    if user.QuestionIndex == -1 {
        log.Printf("[*] It is weird.... this should not happen in approve_answer function")
        return false
    }
    // Update question_finished_mask
    filter = bson.D{
        {"account", user_account},
    }
    update = bson.D{
        {"$set", bson.D{
            {"questionfinishedmask." + strconv.Itoa(user.QuestionIndex), true}, // mark that question as solved!
            {"answertext", ""},
            {"answerbase64str", ""},
        },
        },
    }
    if err = collection.FindOneAndUpdate(context.TODO(), filter, update).Err(); err != nil{
        return false
    }
    return true
}

func skip_answer(user_account string, collection *mongo.Collection) bool{
    // Set this account's questionindex to -1
    filter := bson.D{
        {"account", user_account},
        {"questionindex", bson.D{
            {"$ne", -1}, // qidx should not be -1!
        }},
    }
    update := bson.D{
        {"$set", bson.D{
            {"questionindex", -1},
            {"answertext", ""},
            {"answerbase64str", ""},
        }},
    }
    if err := collection.FindOneAndUpdate(context.TODO(), filter, update).Err(); err != nil{
        log.Printf("[*] Skip %s's question failed!", user_account)
        return false
    }
    return true
}

// https://github.com/googollee/go-socket.io/issues/194#issuecomment-481234861
func generateEngineOptions() *engineio.Options{
    poll_transport := polling.Default
    websocket_transport := websocket.Default
    websocket_transport.CheckOrigin = func(req *http.Request) bool{
        return true
    }
    options := engineio.Options{
        Transports: []transport.Transport{
            poll_transport,
            websocket_transport,
        },
    }
    return &options
}

func initialize_socket_io() (*socketio.Server){
    options := generateEngineOptions()
    server, err := socketio.NewServer(options)
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
    user_collection := createDatabaseCollection(client, to_clear_collection == "true")
    if to_clear_collection == "true" {
        initialize_users(user_collection, 10)
    }
    // Load questions from the csv file
    questions := load_question_from_csv("questions.csv")
    //
    server := initialize_socket_io()
    defer server.Close()

    router := gin.Default()
    config := cors.DefaultConfig()
    config.AllowAllOrigins = true
    if false{
        config.AllowOrigins = []string{"http://localhost:8080"}
    }
    config.AllowCredentials = true

    router.Use(cors.New(config))
    // Serve static file: https://github.com/gin-gonic/gin/issues/75#issuecomment-223592440
    router.Use(static.Serve("/", static.LocalFile("./dist", true)))

    socket_router := router.Group("/socket_api")
    {
        // Socket initialization
        socket_router.GET("/socket/", gin.WrapH(server))
        socket_router.POST("/socket/", gin.WrapH(server))
    }

    private_router := router.Group("/user")
    private_router.Use(check_authentication(user_collection))

    // private part (authentication required)
    {
        // Get all information
        private_router.POST("/get_gridnumbers", func(c *gin.Context){
            var user User
            account, _ := c.MustGet("account").(string)
            user.Account = account
            found := user.getGridNumbersByAccount(user_collection)

            if !found {
                c.AbortWithStatus(http.StatusNotFound)
                return
            }

            var question Question
            // If question index == -1, indicate that this team's questionIndex have not been determined
            if user.QuestionIndex != -1{
                question = questions[user.GridNumbers[user.QuestionIndex]-1]
            }

            // Return grid numbers to the front end client
            c.JSON(http.StatusOK, gin.H{
                "gridNumbers": user.GridNumbers,
                "questionIndex": user.QuestionIndex,
                "question_finished_mask": user.QuestionFinishedMask,
                "question": gin.H{
                    "description": question.Description,
                    "image": question.Base64Image,
                },
                "answertext": user.AnswerText,
                "answerbase64str": user.AnswerBase64Str,
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
            account, _ := c.MustGet("account").(string)
            var user User
            user.Account = account

            if c.ShouldBindBodyWith(&user, binding.JSON) == nil {
                log.Printf("[*] Received User %s's answer %s", account, user.AnswerText)
                success := user.saveAnswer(user_collection)
                if !success {
                    log.Printf("[*] User %s try to push_answer but failed!", account)
                    c.String(http.StatusNotAcceptable, "push_answer failed! (It may due to questoinIndex == -1)")
                    return
                }
                // TODO: save answer text and image base64 string to database
                c.String(http.StatusOK, "Answer Received")
                return
            }else{
                c.String(http.StatusNotAcceptable, "Input is weird")
                return
            }
        })

        // A user can check what he have answered for this question
        private_router.POST("/get_answer", func(c *gin.Context){
            account, _ := c.MustGet("account").(string)
            var user User
            user.Account = account
            log.Printf("account: %s \n", account)

            if c.ShouldBindBodyWith(&user, binding.JSON) == nil {
                log.Printf("[*] User %s try to get_answer", account)
                success := user.getAnswer(user_collection)
                if !success {
                    c.String(http.StatusNotAcceptable, "get Answer failed")
                    return
                }
                c.JSON(http.StatusOK, gin.H{
                    "answertext": user.AnswerText,
                    "answerbase64str": user.AnswerBase64Str,
                })
                return
            }else{
                c.String(http.StatusNotAcceptable, "/get_answer's input body is weird!")
                return
            }
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
                    // Return questionIndex and question
                    c.JSON(http.StatusOK, gin.H{
                        "questionIndex": user.QuestionIndex,
                        "question": gin.H{
                            "description": questions[user.GridNumbers[user.QuestionIndex]-1].Description,
                            "image": questions[user.GridNumbers[user.QuestionIndex]-1].Base64Image,
                        },
                    })
                }else{
                    c.String(http.StatusNotAcceptable, "")
                }
            }else{
                c.String(http.StatusNotAcceptable, "Select question fail!")
            }
        })

        // (Admin) get all information
        private_router.POST("/get_all", func(c *gin.Context){
            account, _ := c.MustGet("account").(string)
            if account != ADMIN{
                c.String(http.StatusNotAcceptable, "You are not admin! Get out!")
                return
            }
            users, chosen_questions, success := getAll(user_collection, questions)
            if !success {
                c.String(http.StatusNotAcceptable, "getAll() failed")
                return
            }

            c.JSON(http.StatusOK, gin.H{
                "users": users,
                "questions":chosen_questions,
            })
            return
        })
        private_router.POST("/reset_all", func(c *gin.Context){
            account, _ := c.MustGet("account").(string)
            if account != ADMIN{
                c.String(http.StatusNotAcceptable, "You are not admin! Get out!")
                return
            }
            if success := resetAll(user_collection); !success{
                c.String(http.StatusNotAcceptable, "resetAll() failed")
            }
            c.String(http.StatusOK, "reset_all route succeeded!")
            return
        })
        // (Admin) approve this answer
        private_router.POST("/approve_answer", func(c *gin.Context){
            account, _ := c.MustGet("account").(string)
            if account != ADMIN {
                c.String(http.StatusNotAcceptable, "You are not admin! Get out!")
                return
            }
            var target_user User
            if c.ShouldBindBodyWith(&target_user, binding.JSON) != nil{
                c.String(http.StatusNotAcceptable, "[!] c.ShouldBindBodyWith failed!")
                return
            }
            fmt.Printf("target_user account: %s\n", target_user.Account)
            if success := approve_answer(target_user.Account, user_collection); !success{
                c.String(http.StatusNotAcceptable, "[!] approve_answer failed!")
                return
            }
            c.String(http.StatusOK, "[*] Update question finished mask succeeded!")
            return
        })
        // (Admin) skip this team's question
        private_router.POST("/skip_answer", func(c *gin.Context){
            account, _ := c.MustGet("account").(string)
            if account != ADMIN {
                c.String(http.StatusNotAcceptable, "You are not admin! Get out!")
                return
            }
            var target_user User
            if c.ShouldBindBodyWith(&target_user, binding.JSON) != nil{
                c.String(http.StatusNotAcceptable, "[!] c.ShouldBindBodyWith failed!")
                return
            }
            fmt.Printf("target_user account: %s\n", target_user.Account)
            if success := skip_answer(target_user.Account, user_collection); !success{
                c.String(http.StatusNotAcceptable, "[!] skip_answer failed!")
                return
            }
            c.String(http.StatusOK, "[*] Skipping that question succeeded!")
            return
        })
    }

    // Perform authentication
    api_router := router.Group("/api")
    {
        api_router.POST("/auth", func(c *gin.Context){
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
    }
    router.Run("0.0.0.0:80")
    destroyMongoClient(client)
}
