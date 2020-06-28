package lib

import (
    "log"
    "context"
    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/mongo"
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

