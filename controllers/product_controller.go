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

var productCollection *mongo.Collection = configs.GetCollection(configs.DB, "product")
var validateProduct = validator.New()

// Create
func CreateProduct() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

		var product models.Product
		defer cancel()
		//validate request body
		if err := c.BindJSON(&product); err != nil {
			c.JSON(http.StatusBadRequest, responses.ProductResponses{Status: http.StatusBadRequest, Message: "error",
				Data: map[string]interface{}{"data": err.Error()}})
			return
		}
		//validator untuk required field
		if ValidatorErr := validate.Struct(&product); ValidatorErr != nil {
			c.JSON(http.StatusBadRequest, responses.ProductResponses{Status: http.StatusBadRequest,
				Message: "Error", Data: map[string]interface{}{"data": ValidatorErr.Error()}})
		}

		//insert data
		newProduct := models.Product{
			Id:           primitive.NewObjectID(),
			Nama_product: product.Nama_product,
			Jumlah:       product.Jumlah,
			Dibuat:       product.Dibuat,
			Ketersediaan: product.Ketersediaan,
		}
		result, err := productCollection.InsertOne(ctx, newProduct)
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.ProductResponses{
				Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{
					"data": err.Error(),
				},
			})
			return
		}
		c.JSON(http.StatusCreated, responses.ProductResponses{
			Status: http.StatusCreated, Message: "sukses",
			Data: map[string]interface{}{"data": result},
		})

	}
}

// get all
func GetAllProducts() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var products []models.Product

		defer cancel()

		result, err := productCollection.Find(ctx, bson.M{})
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.ProductResponses{
				Status:  http.StatusInternalServerError,
				Message: "error",
				Data: map[string]interface{}{
					"data": err.Error(),
				},
			})
			return
		}
		//baca data produk
		defer result.Close(ctx)
		for result.Next(ctx) {
			var singleProduct models.Product
			if err = result.Decode(&singleProduct); err != nil {
				c.JSON(http.StatusInternalServerError, responses.ProductResponses{
					Status:  http.StatusInternalServerError,
					Message: "error",
					Data: map[string]interface{}{
						"data": err.Error(),
					},
				})
			}
			products = append(products, singleProduct)
		}
		c.JSON(http.StatusOK, responses.ProductResponses{
			Status:  http.StatusOK,
			Message: "sukses",
			Data: map[string]interface{}{
				"data": products,
			},
		})
	}
}

// get product by id
func GetProducts() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		productId := c.Param("productId")
		var product models.Product
		defer cancel()

		objId, _ := primitive.ObjectIDFromHex(productId)
		err := productCollection.FindOne(ctx, bson.M{"id": objId}).Decode(&product)
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.ProductResponses{
				Status:  http.StatusInternalServerError,
				Message: "error",
				Data: map[string]interface{}{
					"data": err.Error(),
				},
			})
			return
		}
		c.JSON(http.StatusOK, responses.ProductResponses{
			Status:  http.StatusOK,
			Message: "berhasil",
			Data: map[string]interface{}{
				"data": product,
			},
		})
	}
}

// Edit Product
func EditProduct() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		productId := c.Param("productId")
		var product models.Product
		defer cancel()
		objId, _ := primitive.ObjectIDFromHex(productId)
		//validasi request body
		if err := c.BindJSON(&product); err != nil {
			c.JSON(http.StatusInternalServerError, responses.ProductResponses{
				Status:  http.StatusInternalServerError,
				Message: "error",
				Data: map[string]interface{}{
					"data": err.Error(),
				},
			})
			return
		}
		update := bson.M{
			"Nama_product": product.Nama_product,
			"Jumlah":       product.Jumlah,
			"Dibuat":       product.Dibuat,
			"Ketersediaan": product.Ketersediaan,
		}
		result, err := productCollection.UpdateOne(ctx, bson.M{
			"id": objId,
		}, bson.M{"$set": update})

		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.ProductResponses{
				Status:  http.StatusInternalServerError,
				Message: "error",
				Data: map[string]interface{}{
					"data": err.Error(),
				},
			})
			return
		}
		//GET UPDATE PRODUCT DETAILS
		var updateProduct models.Product
		if result.MatchedCount == 1 {
			err := productCollection.FindOne(ctx, bson.M{
				"id": objId,
			}).Decode(&updateProduct)
			if err != nil {
				c.JSON(http.StatusInternalServerError, responses.ProductResponses{
					Status:  http.StatusInternalServerError,
					Message: "error",
					Data: map[string]interface{}{
						"data": err.Error(),
					},
				})
				return
			}
		}
		c.JSON(http.StatusOK, responses.ProductResponses{
			Status:  http.StatusOK,
			Message: "sukses",
			Data: map[string]interface{}{
				"data": updateProduct,
			},
		})

	}
}

// DELETE PRODUCT
func DeleteProduct() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		productId := c.Param("productId")
		defer cancel()

		objId, _ := primitive.ObjectIDFromHex(productId)
		result, err := productCollection.DeleteOne(ctx, bson.M{"id": objId})
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.ProductResponses{
				Status:  http.StatusInternalServerError,
				Message: "error",
				Data: map[string]interface{}{
					"data": err.Error(),
				},
			})
			return
		}

		if result.DeletedCount < 1 {
			c.JSON(http.StatusNotFound, responses.ProductResponses{
				Status:  http.StatusNotFound,
				Message: "error",
				Data:    map[string]interface{}{"data": "Not found"},
			})
			return
		}

		c.JSON(http.StatusOK, responses.ProductResponses{
			Status:  http.StatusOK,
			Message: "berhasil",
			Data: map[string]interface{}{
				"data": "sukses",
			},
		})
	}
}
