package controller

import (
	"fmt"
	"net/http"

	"biodata-server/database"
	"biodata-server/models"

	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/gin-gonic/gin"
)
var collMhs = database.Db.Collection("mahasiswa")
var validate = validator.New()

func GetMhs(c *gin.Context){
	ctx := c.Request.Context()
	var storeMhs []models.MahasiswaPub
	cur, err := collMhs.Find(ctx, bson.M{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "gagal mendapatkan mahasiswa dari database"})
		return
	}
	if cur.Next(ctx) {
		var mahasiswa models.MahasiswaPub
		cur.Decode(&mahasiswa)
		storeMhs = append(storeMhs, mahasiswa)
	}
	if len(storeMhs) == 0 {
		c.JSON(http.StatusNoContent, gin.H{"message": "data belum ada"})
		return
	}
	c.IndentedJSON(http.StatusOK, storeMhs)
}

func PostMhs(c *gin.Context){
	ctx := c.Request.Context()
	var mahasiswa models.MahasiswaPriv
	err := c.BindJSON(&mahasiswa)
	if err != nil {
		c.JSON(http.StatusInternalServerError, "Gagal bind data")
		return
	}
	err = validate.Struct(mahasiswa)
	if err != nil {
		c.JSON(http.StatusBadRequest, "Data tidak sesuai format")
		return
	}

	fillCariMhs := bson.M{
		"$or": []bson.M{
			{"username": mahasiswa.Username},
			{"npm": mahasiswa.Npm},
		},
	}
	result := collMhs.FindOne(ctx, fillCariMhs)
	if result.Err() == nil {
		message := fmt.Sprintf("akun dengan npm %s sudah ada", mahasiswa.Npm)
		c.JSON(http.StatusConflict, gin.H{"message": message})
		return
	}
	_, err = collMhs.InsertOne(ctx, mahasiswa)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Gagal menambahkan data ke databse"})
		return
	}
	c.IndentedJSON(http.StatusCreated, gin.H{"message": "berhasil menambahkan data dengan npm " + mahasiswa.Npm})
}