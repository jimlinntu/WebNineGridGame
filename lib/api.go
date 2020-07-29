package lib

import (
    "log"
    "context"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo"
    "strconv"
)

func RejectAnswer(user_account string, collection *mongo.Collection) bool{
    filter := bson.D{
        {"account", user_account},
        {"questionindex", bson.D{
            {"$ne", -1},
        }},
    }

    update := bson.D{
        {"$set", bson.D{
            {"isrejected", true}, // set this user current answer is rejected
        }},
    }

    if err := collection.FindOneAndUpdate(context.TODO(), filter, update).Err(); err != nil{
        log.Printf("[*] Reject %s's question failed!", user_account)
        return false
    }
    return true
}

func PetitionSkipAnswer(user_account string, collection *mongo.Collection) bool{
    filter := bson.D{
        {"account", user_account},
        {"questionindex", bson.D{
            {"$ne", -1},
        }},
    }

    update := bson.D{
        {"$set", bson.D{
            {"haspetition", true}, // set this user intent to petition for skipping
        }},
    }

    if err := collection.FindOneAndUpdate(context.TODO(), filter, update).Err(); err != nil{
        log.Printf("[*] The user %s petition for skipping failed!", user_account)
        return false
    }
    return true
}

func DeleteFinished(user_account string, question_index int, num_grid int, collection *mongo.Collection) bool{
    if question_index < 0 || question_index >= num_grid{
        log.Printf("[*] Try to delete account: %s's question_index: %d failed", user_account, question_index)
        return false
    }
    filter := bson.D{
        {"account", user_account},
    }
    update := bson.D{
        {"$set", bson.D{
            {"questionfinishedmask." + strconv.Itoa(question_index), false}, // mark that question as solved!
        }},
    }
    if err := collection.FindOneAndUpdate(context.TODO(), filter, update).Err(); err != nil{
        log.Printf("[*] Try to delete account: %s's question_index: %d failed when updating the datebase", user_account, question_index)
    }
    return true
}
