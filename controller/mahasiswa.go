package controller

import (
	"fmt"
	"net/http"
	"strings"
	"sync"

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
	defer cur.Close(ctx)
	
	for cur.Next(ctx) {
		var mahasiswa models.MahasiswaPriv
		err := cur.Decode(&mahasiswa)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"message": "Gagal decode data mahasiswa"})
      return
		}
		storeMhs = append(storeMhs, models.MahasiswaPub{
			Nama: mahasiswa.Nama,
			Npm: mahasiswa.Npm,
			Angkatan: mahasiswa.Angkatan,
			Alamat: mahasiswa.Alamat,
			NoHp: mahasiswa.NoHp,
			NoHpKeluarga: mahasiswa.NoHpKeluarga,
			AsalSekolah: mahasiswa.AsalSekolah,
		})
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
		c.JSON(http.StatusInternalServerError, gin.H{"message": "Gagal bind data"})
		return
	}

	err = validate.Struct(mahasiswa)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Data tidak sesuai format"})
		return
	}

	var konflik []string
	var wg sync.WaitGroup
	var mut sync.Mutex
	wg.Add(3)
	go func() {
		defer wg.Done()
		err := collMhs.FindOne(ctx, bson.M{"mahasiswapub.npm": mahasiswa.Npm}).Err()
		if err == nil {
			mut.Lock()
			konflik = append(konflik, "npm")
			mut.Unlock()
		}
	}()
	go func() {
		defer wg.Done()
		err := collMhs.FindOne(ctx, bson.M{"username": mahasiswa.Username}).Err()
		if err == nil {
			mut.Lock()
			konflik = append(konflik, "username")
			mut.Unlock()
		}
	}()
	go func() {
		defer wg.Done()
		err := collMhs.FindOne(ctx, bson.M{"mahasiswapub.nama": mahasiswa.Nama}).Err()
		if err == nil {
			mut.Lock()
			konflik = append(konflik, "nama")
			mut.Unlock()
	}
	}()
	wg.Wait()

	if len(konflik) > 0 {
		message := fmt.Sprintf("akun dengan %s tersebur sudah terdaftar", strings.Join(konflik, ", "))
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