package src

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var db = make(map[string]string)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "hello_world")
	})

	r.GET("/proofs/:image_id", GetEntry)

	r.POST("/proofs", CreateEntry)

	return r
}

type Entry struct {
	gorm.Model
	CID     string `json:"cid"`
	ImageID string `gorm:"primaryKey" json:"image_id"`
}

func ConnectToDatabase() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}

func CreateEntry(c *gin.Context) {
	//TODO: send http request to ipfs kube to add new entry
	image_id := c.Query("image_id")

	if image_id == "" {

		c.JSON(400, gin.H{"error": "imageId query parameter is required"})
		return
	}

	db, err := ConnectToDatabase()

	if err != nil {
		c.JSON(400, gin.H{"error": "Database Error" + err.Error()})
		return
	}

	fileheader, err := c.FormFile("file")

	if err != nil {
		c.JSON(400, gin.H{"error": "File Error" + err.Error()})
		return
	}

	src, err := fileheader.Open()

	if err != nil {
		c.JSON(400, gin.H{"error": "Couldnt open Fileheader" + err.Error()})
		return
	}
	defer src.Close()

	var buf bytes.Buffer

	writer := multipart.NewWriter(&buf)

	part, err := writer.CreateFormFile("file", "myimage.jpg")
	if err != nil {
		c.JSON(500, gin.H{"error": "Formfile error" + err.Error()})
		return
	}

	_, err = io.Copy(part, src)

	if err != nil {
		c.JSON(500, gin.H{"error": "Couldnt copy file" + err.Error()})
		return
	}

	// Send POST request to Kubo API
	resp, err := http.Post(
		"http://ipfs-service:5001/api/v0/add?pin=true", // <-- Add this
		writer.FormDataContentType(),
		&buf,
	)
	if err != nil {
		c.JSON(500, gin.H{"error": "Request to ipfs failed" + err.Error()})
		return
	}
	writer.Close()
	var ipfsResp struct {
		Hash string `json:"Hash"`
	}

	// Parse IPFS response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		c.JSON(500, gin.H{"error": "Request to ipfs failed" + err.Error()})
		return
	}

	if err := json.Unmarshal(body, &ipfsResp); err != nil {
		c.JSON(500, gin.H{"error": "Request to ipfs failed" + err.Error()})
		return
	}

	defer resp.Body.Close()

	//TODO: Afterwards save CID, and Image Id in database (is to be replaced at the end by smart contracts on blockchaiis to be replaced at the end by smart contracts on blockchainn)

	entry := Entry{
		CID:     ipfsResp.Hash,
		ImageID: image_id,
	}

	if err := db.Create(&entry).Error; err != nil {
		c.JSON(500, gin.H{"error": "failed to save to ipfs"})
		return
	}

	c.JSON(200, gin.H{"cid": ipfsResp.Hash, "image_id": image_id})

}
func GetEntry(c *gin.Context) {
	//TODO: Get Entry by first looking up image id from database, and getting the CID for kube IPFS:
	imageID := c.Param("image_id")

	db, err := ConnectToDatabase()
	if err != nil {
		c.JSON(500, gin.H{"error": "Could not connect to ipfs server" + err.Error()})
		return
	}

	var entry Entry
	err = db.First(&entry, "image_id = ?", imageID).Error
	if err != nil {
		c.JSON(404, gin.H{"error": "Entry not found" + err.Error()})
		return
	}

	resp, err := http.Get("http://ipfs-service:8080/ipfs/" + entry.CID)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to fetch from IPFS. Error:" + string(err.Error())})
		return
	}
	defer resp.Body.Close()

	c.Header("Content-Type", "application/zip")
	c.Header("Content-Disposition", `attachment; filename="proof.zip"`)
	io.Copy(c.Writer, resp.Body)

	//return "http://localhost:8080/ipfs/" + entry.CID or send file out directly
}
