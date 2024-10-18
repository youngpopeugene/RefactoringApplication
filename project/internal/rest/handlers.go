package rest

import (
	"app/internal/models"
	"database/sql"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"strconv"
	"time"
)

type handler struct {
	db *gorm.DB
}

func NewHandler(db *gorm.DB) *handler {
	return &handler{db: db}
}

//Get by primary key

func (h *handler) GetSubstationByPK(c *gin.Context) {
	pk := c.Param("pk")
	var obj = models.Substation{
		Name: pk,
	}
	h.db.First(&obj)
	c.JSON(200, obj)
}

func (h *handler) GetFactoryByPK(c *gin.Context) {
	pk := c.Param("pk")
	var obj = models.Factory{
		Name: pk,
	}
	h.db.First(&obj)
	c.JSON(200, obj)
}

func (h *handler) GetRangeOfHighVoltageEquipmentByPK(c *gin.Context) {
	pk := c.Param("pk")
	pkInt, _ := strconv.Atoi(pk)
	var obj = models.RangeOfHighVoltageEquipment{
		ID: pkInt,
	}
	h.db.First(&obj)
	c.JSON(200, obj)
}

func (h *handler) GetCableLineByPK(c *gin.Context) {
	pk := c.Param("pk")
	var obj = models.CableLine{
		Mark: pk,
	}
	h.db.First(&obj)
	c.JSON(200, obj)
}

func (h *handler) GetTireSectionByPK(c *gin.Context) {
	pk := c.Param("pk")
	var obj = models.TireSection{
		Name: pk,
	}
	h.db.First(&obj)
	c.JSON(200, obj)
}

func (h *handler) GetCellKVLByPK(c *gin.Context) {
	pk := c.Param("pk")
	var obj = models.CellKVL{
		DispatchName: pk,
	}
	h.db.First(&obj)
	c.JSON(200, obj)
}

func (h *handler) GetFuseByPK(c *gin.Context) {
	pk := c.Param("pk")
	var obj = models.Fuse{
		Mark: pk,
	}
	h.db.First(&obj)
	c.JSON(200, obj)
}

func (h *handler) GetCellTNByPK(c *gin.Context) {
	pk := c.Param("pk")
	var obj = models.CellTN{
		DispatchName: pk,
	}
	h.db.First(&obj)
	c.JSON(200, obj)
}

func (h *handler) GetCellTSNByPK(c *gin.Context) {
	pk := c.Param("pk")
	var obj = models.CellTSN{
		DispatchName: pk,
	}
	h.db.First(&obj)
	c.JSON(200, obj)
}

func (h *handler) GetNSSByPK(c *gin.Context) {
	pk := c.Param("pk")
	pkInt, _ := strconv.Atoi(pk)
	var obj = models.NSS{
		ID: pkInt,
	}
	h.db.First(&obj)
	c.JSON(200, obj)
}

func (h *handler) GetRangeOfStandardVoltageByPK(c *gin.Context) {
	pk := c.Param("pk")
	pkInt, _ := strconv.Atoi(pk)
	var obj = models.RangeOfStandardVoltage{
		ID: pkInt,
	}
	h.db.First(&obj)
	c.JSON(200, obj)
}

func (h *handler) GetTypeOfTransformerByPK(c *gin.Context) {
	pk := c.Param("pk")
	var obj = models.TypeOfTransformer{
		Type: pk,
	}
	h.db.First(&obj)
	c.JSON(200, obj)
}

func (h *handler) GetTransformerByPK(c *gin.Context) {
	pk := c.Param("pk")
	pkInt, _ := strconv.Atoi(pk)
	var obj = models.Transformer{
		FactoryNumber: pkInt,
	}
	h.db.First(&obj)
	c.JSON(200, obj)
}

func (h *handler) GetRequestByPK(c *gin.Context) {
	pk := c.Param("pk")
	pkInt, _ := strconv.Atoi(pk)
	var obj = models.Request{
		ID: pkInt,
	}
	h.db.First(&obj)

	c.JSON(200, obj)
}

//Get all users workers

func (h *handler) GetAllUsersWorkers(c *gin.Context) {
	role := models.RoleWorker
	var users []models.User
	h.db.Where("role = ?", role).Find(&users)
	c.JSON(200, users)
}

//Get all unique location

func (h *handler) GetAllLocations(c *gin.Context) {
	var locations []string
	h.db.Model(&models.Substation{}).Select("DISTINCT location").Find(&locations)
	c.JSON(200, locations)
}

//Get transformers by location of substation

func (h *handler) GetTransformersByLocation(c *gin.Context) {
	location := c.Param("location")
	var substations []models.Substation
	result := h.db.Raw("SELECT * "+
		"FROM transformers "+
		"JOIN substations on transformers.substation = substations.name "+
		"WHERE substations.location = @location",
		sql.Named("location", location)).Find(&substations)
	if result.RowsAffected == 0 {
		c.JSON(400, "record not found")
		return
	}
	c.JSON(200, substations)
}

//Get all requests

func (h *handler) GetAllRequests(c *gin.Context) {
	var requests []models.Request
	h.db.Find(&requests)
	c.JSON(200, requests)
}

//Get requests by worker_username

func (h *handler) GetRequestsByWorkerUsername(c *gin.Context) {
	var requests []models.Request
	workerUsername := c.Param("workerUsername")
	h.db.Where("worker_username = ?", workerUsername).Find(&requests)
	c.JSON(200, requests)
}

//Create request

func (h *handler) CreateRequest(c *gin.Context) {
	var obj struct {
		WorkerUsername           string `json:"worker_username"`
		TransformerFactoryNumber int    `json:"transformer_factory_number"`
	}
	if err := c.ShouldBind(&obj); err != nil {
		c.JSON(400, err)
		return
	}

	var transformer = models.Transformer{
		FactoryNumber: obj.TransformerFactoryNumber,
	}
	h.db.First(&transformer)

	var user = models.User{
		Username: obj.WorkerUsername,
	}
	h.db.First(&user)
	if user.Role != models.RoleWorker {
		c.JSON(400, "Only worker's username could be passed")
		return
	}

	request := models.Request{
		TransformerFactoryNumber: transformer.FactoryNumber,
		WorkerUsername:           user.Username,
		IsCompleted:              false,
		DateOpened:               time.Now(),
	}

	h.db.Model(&transformer).Association("TransformerFactoryNumber").Append(&request)
	h.db.Model(&user).Association("WorkerUsername").Append(&request)
	h.db.Create(&request)
	c.JSON(200, request)
}

//Update request (close it)

func (h *handler) UpdateRequest(c *gin.Context) {
	var obj struct {
		WorkerUsername           string `json:"worker_username"`
		TransformerFactoryNumber int    `json:"transformer_factory_number"`
	}
	if err := c.ShouldBind(&obj); err != nil {
		c.JSON(400, err)
		return
	}

	var transformer = models.Transformer{
		FactoryNumber: obj.TransformerFactoryNumber,
	}
	result := h.db.First(&transformer)

	var user = models.User{
		Username: obj.WorkerUsername,
	}
	result = h.db.First(&user)

	if user.Role != models.RoleWorker {
		c.JSON(400, "Only worker's username could be passed")
		return
	}

	request := models.Request{
		TransformerFactoryNumber: transformer.FactoryNumber,
		WorkerUsername:           user.Username,
	}
	h.db.First(&request)
	request.IsCompleted = true
	request.DateClosed = time.Now()

	h.db.Model(&transformer).Association("TransformerFactoryNumber").Append(&request)
	h.db.Model(&user).Association("WorkerUsername").Append(&request)
	h.db.Save(&request)
	c.JSON(200, request)
}
