package types

import (
	"github.com/muneebhashone/go-fiber-api/enums"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Room struct {
	ID        primitive.ObjectID `bson:"_id" json:"id"`
	Type      enums.RoomType     `bson:"type" json:"type"`
	BasePrice float64            `bson:"base_price" json:"base_price"`
	Price     float64            `bson:"price" json:"price"`
	HotelID   primitive.ObjectID `bson:"hotel_id" json:"hotel_id"`
}
