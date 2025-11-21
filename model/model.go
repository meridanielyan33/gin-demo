package model

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username  string `gorm:"unique" json:"username"`
	Email     string `gorm:"unique" json:"email"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Age       int    `json:"age"`
	Password  string `json:"-" form:"password"`
}

type UserLoginRequest struct {
	Email    string `json:"email`
	Password string `json:"password`
}

type UserLoginResponse struct {
	Message  string `json:"message"`
	JWTToken string
}

type UserLogoutRequest struct {
	Email string `json:"email"`
}

type UserLogoutResponse struct {
	Message string `json:"message"`
}

type Movie struct {
	ID          primitive.ObjectID   `bson:"_id,omitempty" json:"id"`
	Title       string               `bson:"title" json:"title"`
	ReleaseYear int                  `bson:"release_year" json:"release_year"`
	DirectorID  primitive.ObjectID   `bson:"director_id" json:"director_id"`
	Actors      []primitive.ObjectID `bson:"actors" json:"actors"`

	Director      *Director `bson:"-" json:"director"`
	ActorsDetails []Actor   `bson:"-" json:"actors_details"`
}

type MovieResponse struct {
	Title         string    `bson:"title" json:"title"`
	ReleaseYear   int       `bson:"release_year" json:"release_year"`
	Director      *Director `bson:"-" json:"director"`
	ActorsDetails []Actor   `bson:"-" json:"actors_details"`
}

type Actor struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	FirstName string             `json:"first_name" bson:"first_name"`
	LastName  string             `json:"last_name" bson:"last_name"`
	BirthDate time.Time          `bson:"birth_date" json:"birth_date"`
}

type Director struct {
	ID        primitive.ObjectID `bson:"_id,omitempty" json:"id"`
	FirstName string             `json:"first_name" bson:"first_name"`
	LastName  string             `json:"last_name" bson:"last_name"`
	BirthDate time.Time          `bson:"birth_date" json:"birth_date"`
}
