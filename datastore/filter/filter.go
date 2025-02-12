package filter

import "go.mongodb.org/mongo-driver/v2/bson"

func BeginsWith(a string) bson.M {
	return bson.M{"$regex": "^" + a, "$options": "i"}
}
