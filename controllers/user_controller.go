package controllers

import (
	"context"
	"crudgolang/configs"
	"crudgolang/models"
	"crudgolang/responses"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var userCollection *mongo.Collection = configs.GetCollection(configs.DB, "users")
var validate = validator.New()

// CREATE
func CreateUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var user models.User
		defer cancel()

		//validate the request body
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, responses.UserResponse{Status: http.StatusBadRequest, Message: "error",
				Data: map[string]interface{}{"data": err.Error()}})
			return
		}
		//use validator untuk validate required field
		if validationErr := validate.Struct(&user); validationErr != nil {
			c.JSON(http.StatusBadRequest, responses.UserResponse{Status: http.StatusBadRequest,
				Message: "error", Data: map[string]interface{}{"data": validationErr.Error()}})
			return
		}

		newUser := models.User{
			Id:       primitive.NewObjectID(),
			Name:     user.Name,
			Location: user.Location,
			Title:    user.Title,
		}

		result, err := userCollection.InsertOne(ctx, newUser)
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.UserResponse{Status: http.StatusInternalServerError,
				Message: "error", Data: map[string]interface{}{"data": err.Error()}})
		}
		c.JSON(http.StatusCreated, responses.UserResponse{Status: http.StatusCreated, Message: "success",
			Data: map[string]interface{}{"data": result}})
	}
}

// READ
func GetAllUsers() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var users []models.User
		defer cancel()

		results, err := userCollection.Find(ctx, bson.M{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.UserResponse{
				Status:  http.StatusInternalServerError,
				Message: "error",
				Data:    map[string]interface{}{"data": err.Error()},
			})
			return
		}

		//baca dari db dalam cara yang optimal
		defer results.Close(ctx)
		for results.Next(ctx) {
			var singleUser models.User
			if err = results.Decode(&singleUser); err != nil {
				c.JSON(http.StatusInternalServerError,
					responses.UserResponse{Status: http.StatusInternalServerError,
						Message: "error",
						Data:    map[string]interface{}{"data": err.Error()}})
			}
			users = append(users, singleUser)
		}

		c.JSON(http.StatusOK,
			responses.UserResponse{
				Status:  http.StatusOK,
				Message: "sukses",
				Data:    map[string]interface{}{"data": users},
			})

	}
}

// get user
func GetUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		userId := c.Param("userId")
		var user models.User
		defer cancel()

		objId, _ := primitive.ObjectIDFromHex(userId)
		err := userCollection.FindOne(ctx, bson.M{"id": objId}).Decode(&user)
		if err != nil {
			c.JSON(http.StatusInternalServerError,
				responses.UserResponse{
					Status:  http.StatusInternalServerError,
					Message: "error",
					Data:    map[string]interface{}{"data": err.Error()},
				})
			return
		}
		c.JSON(http.StatusOK, responses.UserResponse{
			Status:  http.StatusOK,
			Message: "sukses",
			Data: map[string]interface{}{
				"data": user,
			},
		})
		return
	}
}

// edit data
func EditUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		userId := c.Param("userId")
		var user models.User
		defer cancel()

		objId, _ := primitive.ObjectIDFromHex(userId)
		//validasi request body
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, responses.UserResponse{
				Status:  http.StatusBadRequest,
				Message: "error",
				Data: map[string]interface{}{
					"data": err.Error(),
				},
			})
			return
		}
		update := bson.M{"name": user.Name,
			"location": user.Location,
			"title":    user.Title}
		result, err := userCollection.UpdateOne(ctx,
			bson.M{
				"id": objId,
			}, bson.M{"$set": update})

		if err != nil {
			c.JSON(http.StatusInternalServerError,
				responses.UserResponse{
					Status:  http.StatusInternalServerError,
					Message: "error",
					Data: map[string]interface{}{
						"data": err.Error()},
				})
			return
		}

		//get update user details
		var updateUser models.User
		if result.MatchedCount == 1 {
			err := userCollection.FindOne(ctx,
				bson.M{"id": objId}).Decode(&updateUser)
			if err != nil {
				c.JSON(http.StatusInternalServerError, responses.UserResponse{
					Status:  http.StatusInternalServerError,
					Message: "error",
					Data: map[string]interface{}{
						"data": err.Error(),
					},
				})
				return
			}
		}
		c.JSON(http.StatusOK, responses.UserResponse{
			Status:  http.StatusOK,
			Message: "berhasil",
			Data: map[string]interface{}{
				"data": updateUser,
			},
		})
	}

}

// delete user
func DeleteUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		userId := c.Param("userId")
		defer cancel()

		objId, _ := primitive.ObjectIDFromHex(userId)

		result, err := userCollection.DeleteOne(ctx, bson.M{"id": objId})

		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.UserResponse{
				Status:  http.StatusInternalServerError,
				Message: "error",
				Data: map[string]interface{}{
					"data": err.Error(),
				},
			})
			return
		}

		if result.DeletedCount < 1 {
			c.JSON(http.StatusNotFound,
				responses.UserResponse{Status: http.StatusNotFound,
					Message: "error",
					Data:    map[string]interface{}{"data": "data not found"}})
			return
		}

		c.JSON(http.StatusOK,
			responses.UserResponse{Status: http.StatusOK,
				Message: "sukses", Data: map[string]interface{}{"data": "data berhasil dihapus"}})

	}

}
