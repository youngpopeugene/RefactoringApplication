package rest

import (
	"github.com/gin-gonic/gin"
	jwt "github.com/kyfk/gin-jwt"
	"gorm.io/gorm"
	"log"
)

func SetupRouter(db *gorm.DB) *gin.Engine {
	router := gin.Default()
	routerGroup := router.Group("/api/v1")
	InitRoutes(routerGroup, db)
	return router
}

func InitRoutes(routerGroup *gin.RouterGroup, db *gorm.DB) {
	h := NewHandler(db)

	//Auth
	auth, err := h.NewAuth()
	if err != nil {
		log.Fatal(err)
	}
	routerGroup.Use(jwt.ErrorHandler)
	routerGroup.POST("/login", auth.Authenticate)
	routerGroup.POST("/register", h.Register)

	//Get by primary key
	routerGroup.GET("/substation/:pk", h.GetSubstationByPK)
	routerGroup.GET("/factory/:pk", h.GetFactoryByPK)
	routerGroup.GET("/rangeOfHighVoltageEquipment/:pk",
		h.GetRangeOfHighVoltageEquipmentByPK)
	routerGroup.GET("/cableLine/:pk", h.GetCableLineByPK)
	routerGroup.GET("/cellKVL/:pk", h.GetCellKVLByPK)
	routerGroup.GET("/fuse/:pk", h.GetFuseByPK)
	routerGroup.GET("/cellTN/:pk", h.GetCellTNByPK)
	routerGroup.GET("/cellTSN/:pk", h.GetCellTSNByPK)
	routerGroup.GET("/nss/:pk", h.GetNSSByPK)
	routerGroup.GET("/rangeOfStandardVoltage/:pk",
		h.GetRangeOfStandardVoltageByPK)
	routerGroup.GET("/typeOfTransformer/:pk",
		h.GetTypeOfTransformerByPK)
	routerGroup.GET("/transformer/:pk", h.GetTransformerByPK)
	routerGroup.GET("/request/:pk", h.GetRequestByPK)

	//Get all users workers
	routerGroup.GET("/workers", Dispatcher(auth), h.GetAllUsersWorkers)

	//Get all locations
	routerGroup.GET("/locations", h.GetAllLocations)

	//Get transformers by location of substation
	routerGroup.GET("/transformers/:location", h.GetTransformersByLocation)

	//Get all requests
	routerGroup.GET("/requests", Dispatcher(auth), h.GetAllRequests)

	//Get requests by worker_username
	routerGroup.GET("/requests/:workerUsername", h.GetRequestsByWorkerUsername)

	//Create request
	routerGroup.POST("/createRequest", Dispatcher(auth), h.CreateRequest)

	//Update request (close it)
	routerGroup.PUT("/updateRequest", Worker(auth), h.UpdateRequest)
}
