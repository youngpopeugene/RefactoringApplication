package api

import (
	"app/internal/gateway"
	"app/internal/infrastructure/httputil"
	"github.com/gin-gonic/gin"
	jwt "github.com/kyfk/gin-jwt"
	"net/http"
)

func (s *server) Router() {
	auth, err := s.NewAuth()
	if err != nil {
		s.logger.Fatal(err)
	}
	s.group.Use(jwt.ErrorHandler)

	// Auth
	s.group.POST("/login", auth.Authenticate)
	s.group.POST("/register", s.Register)

	// Substations
	s.group.GET("/substations/:pk", s.GetSubstationByPK)

	// Locations
	s.group.GET("/locations", s.GetAllLocations)

	// Factories
	s.group.GET("/factories/:pk", s.GetFactoryByPK)

	// High Voltage Equipment
	s.group.GET("/equipment/high-voltage/:pk", s.GetRangeOfHighVoltageEquipmentByPK)

	// Cable Lines
	s.group.GET("/cable-lines/:pk", s.GetCableLineByPK)

	// Cells
	s.group.GET("/cells/kvl/:pk", s.GetCellKVLByPK)
	s.group.GET("/cells/tn/:pk", s.GetCellTNByPK)
	s.group.GET("/cells/tsn/:pk", s.GetCellTSNByPK)

	// NSS
	s.group.GET("/nss/:pk", s.GetNSSByPK)

	// Standard Voltages
	s.group.GET("/standard-voltages/:pk", s.GetRangeOfStandardVoltageByPK)

	// Fuses
	s.group.GET("/fuses/:pk", s.GetFuseByPK)

	// Transformers
	s.group.GET("/transformers/id/:pk", s.GetTransformerByPK)
	s.group.GET("/transformers/location/:location", s.GetTransformersByLocation)
	s.group.GET("/transformers/types/:pk", s.GetTypeOfTransformerByPK)

	// Requests
	s.group.GET("/requests", Dispatcher(auth), s.GetAllRequests)
	s.group.GET("/requests/:pk", s.GetRequestByPK)
	s.group.GET("/workers/:username/requests", s.GetRequestsByWorkerUsername)
	s.group.POST("/requests", Dispatcher(auth), s.CreateRequest)
	s.group.PUT("/requests", Worker(auth), s.UpdateRequest)

	// Users workers
	s.group.GET("users/workers", Dispatcher(auth), s.GetAllUsersWorkers)
}

func (s *server) GetSubstationByPK(c *gin.Context) {
	gw := gateway.NewGateway(s.db)
	obj, err := gw.GetSubstationByPK(c.Param("pk"))
	if err != nil {
		s.logger.Error(err.Error())
		httputil.NewResponse(c, http.StatusNotFound, "error", err.Error(), nil)
		return
	}
	httputil.NewResponse(c, http.StatusOK, "success", "", obj)
}

func (s *server) GetFactoryByPK(c *gin.Context) {
	gw := gateway.NewGateway(s.db)
	obj, err := gw.GetFactoryByPK(c.Param("pk"))
	if err != nil {
		s.logger.Error(err.Error())
		httputil.NewResponse(c, http.StatusNotFound, "error", err.Error(), nil)
		return
	}
	httputil.NewResponse(c, http.StatusOK, "success", "", obj)
}

func (s *server) GetRangeOfHighVoltageEquipmentByPK(c *gin.Context) {
	gw := gateway.NewGateway(s.db)
	obj, err := gw.GetRangeOfHighVoltageEquipmentByPK(c.Param("pk"))
	if err != nil {
		s.logger.Error(err.Error())
		httputil.NewResponse(c, http.StatusNotFound, "error", err.Error(), nil)
		return
	}
	httputil.NewResponse(c, http.StatusOK, "success", "", obj)
}

func (s *server) GetCableLineByPK(c *gin.Context) {
	gw := gateway.NewGateway(s.db)
	obj, err := gw.GetCableLineByPK(c.Param("pk"))
	if err != nil {
		s.logger.Error(err.Error())
		httputil.NewResponse(c, http.StatusNotFound, "error", err.Error(), nil)
		return
	}
	httputil.NewResponse(c, http.StatusOK, "success", "", obj)
}

func (s *server) GetCellKVLByPK(c *gin.Context) {
	gw := gateway.NewGateway(s.db)
	obj, err := gw.GetCellKVLByPK(c.Param("pk"))
	if err != nil {
		s.logger.Error(err.Error())
		httputil.NewResponse(c, http.StatusNotFound, "error", err.Error(), nil)
		return
	}
	httputil.NewResponse(c, http.StatusOK, "success", "", obj)
}

func (s *server) GetFuseByPK(c *gin.Context) {
	gw := gateway.NewGateway(s.db)
	obj, err := gw.GetFuseByPK(c.Param("pk"))
	if err != nil {
		s.logger.Error(err.Error())
		httputil.NewResponse(c, http.StatusNotFound, "error", err.Error(), nil)
		return
	}
	httputil.NewResponse(c, http.StatusOK, "success", "", obj)
}

func (s *server) GetCellTNByPK(c *gin.Context) {
	gw := gateway.NewGateway(s.db)
	obj, err := gw.GetCellTNByPK(c.Param("pk"))
	if err != nil {
		s.logger.Error(err.Error())
		httputil.NewResponse(c, http.StatusNotFound, "error", err.Error(), nil)
		return
	}
	httputil.NewResponse(c, http.StatusOK, "success", "", obj)
}

func (s *server) GetCellTSNByPK(c *gin.Context) {
	gw := gateway.NewGateway(s.db)
	obj, err := gw.GetCellTSNByPK(c.Param("pk"))
	if err != nil {
		s.logger.Error(err.Error())
		httputil.NewResponse(c, http.StatusNotFound, "error", err.Error(), nil)
		return
	}
	httputil.NewResponse(c, http.StatusOK, "success", "", obj)
}

func (s *server) GetNSSByPK(c *gin.Context) {
	gw := gateway.NewGateway(s.db)
	obj, err := gw.GetNSSByPK(c.Param("pk"))
	if err != nil {
		s.logger.Error(err.Error())
		httputil.NewResponse(c, http.StatusNotFound, "error", err.Error(), nil)
		return
	}
	httputil.NewResponse(c, http.StatusOK, "success", "", obj)
}

func (s *server) GetRangeOfStandardVoltageByPK(c *gin.Context) {
	gw := gateway.NewGateway(s.db)
	obj, err := gw.GetRangeOfStandardVoltageByPK(c.Param("pk"))
	if err != nil {
		s.logger.Error(err.Error())
		httputil.NewResponse(c, http.StatusNotFound, "error", err.Error(), nil)
		return
	}
	httputil.NewResponse(c, http.StatusOK, "success", "", obj)
}

func (s *server) GetTypeOfTransformerByPK(c *gin.Context) {
	gw := gateway.NewGateway(s.db)
	obj, err := gw.GetTypeOfTransformerByPK(c.Param("pk"))
	if err != nil {
		s.logger.Error(err.Error())
		httputil.NewResponse(c, http.StatusNotFound, "error", err.Error(), nil)
		return
	}
	httputil.NewResponse(c, http.StatusOK, "success", "", obj)
}

func (s *server) GetTransformerByPK(c *gin.Context) {
	gw := gateway.NewGateway(s.db)
	obj, err := gw.GetTransformerByPK(c.Param("pk"))
	if err != nil {
		s.logger.Error(err.Error())
		httputil.NewResponse(c, http.StatusNotFound, "error", err.Error(), nil)
		return
	}
	httputil.NewResponse(c, http.StatusOK, "success", "", obj)
}

func (s *server) GetRequestByPK(c *gin.Context) {
	gw := gateway.NewGateway(s.db)
	obj, err := gw.GetRequestByPK(c.Param("pk"))
	if err != nil {
		s.logger.Error(err.Error())
		httputil.NewResponse(c, http.StatusNotFound, "error", err.Error(), nil)
		return
	}
	httputil.NewResponse(c, http.StatusOK, "success", "", obj)
}

func (s *server) GetAllUsersWorkers(c *gin.Context) {
	gw := gateway.NewGateway(s.db)
	users := gw.GetAllUsersWorkers()
	if len(users) == 0 {
		s.logger.Warn("no users found")
		httputil.NewResponse(c, http.StatusNotFound, "error", "no users found", nil)
		return
	}
	httputil.NewResponse(c, http.StatusOK, "success", "", users)
}

func (s *server) GetAllLocations(c *gin.Context) {
	gw := gateway.NewGateway(s.db)
	locations := gw.GetAllLocations()
	if len(locations) == 0 {
		s.logger.Warn("no locations found")
		httputil.NewResponse(c, http.StatusNotFound, "error", "no locations found", nil)
		return
	}
	httputil.NewResponse(c, http.StatusOK, "success", "", locations)
}

func (s *server) GetTransformersByLocation(c *gin.Context) {
	gw := gateway.NewGateway(s.db)
	transformers := gw.GetTransformersByLocation(c.Param("location"))
	if len(transformers) == 0 {
		s.logger.Warn("no transformers found")
		httputil.NewResponse(c, http.StatusNotFound, "error", "no transformers found", nil)
		return
	}
	httputil.NewResponse(c, http.StatusOK, "success", "", transformers)
}

func (s *server) GetAllRequests(c *gin.Context) {
	gw := gateway.NewGateway(s.db)
	requests := gw.GetAllRequests()
	if len(requests) == 0 {
		s.logger.Warn("no requests found")
		httputil.NewResponse(c, http.StatusNotFound, "error", "no requests found", nil)
		return
	}
	httputil.NewResponse(c, http.StatusOK, "success", "", requests)
}

func (s *server) GetRequestsByWorkerUsername(c *gin.Context) {
	gw := gateway.NewGateway(s.db)
	requests := gw.GetRequestsByWorkerUsername(c.Param("username"))
	if len(requests) == 0 {
		s.logger.Warn("no requests found for worker")
		httputil.NewResponse(c, http.StatusNotFound, "error", "no requests found for worker", nil)
		return
	}
	httputil.NewResponse(c, http.StatusOK, "success", "", requests)
}

func (s *server) CreateRequest(c *gin.Context) {
	var obj struct {
		WorkerUsername           string `json:"worker_username"`
		TransformerFactoryNumber int    `json:"transformer_factory_number"`
	}
	if err := c.ShouldBind(&obj); err != nil {
		s.logger.Error("invalid input")
		httputil.NewResponse(c, http.StatusBadRequest, "error", err.Error(), nil)
		return
	}

	gw := gateway.NewGateway(s.db)
	request, err := gw.CreateRequest(obj.WorkerUsername, obj.TransformerFactoryNumber)
	if err != nil {
		s.logger.Error("failed to create request")
		httputil.NewResponse(c, http.StatusBadRequest, "error", err.Error(), nil)
		return
	}
	httputil.NewResponse(c, http.StatusOK, "success", "", request)
}

func (s *server) UpdateRequest(c *gin.Context) {
	var obj struct {
		WorkerUsername           string `json:"worker_username"`
		TransformerFactoryNumber int    `json:"transformer_factory_number"`
	}
	if err := c.ShouldBind(&obj); err != nil {
		s.logger.Error("invalid input")
		httputil.NewResponse(c, http.StatusBadRequest, "error", err.Error(), nil)
		return
	}

	gw := gateway.NewGateway(s.db)
	request, err := gw.UpdateRequest(obj.WorkerUsername, obj.TransformerFactoryNumber)
	if err != nil {
		s.logger.Error("failed to update request")
		httputil.NewResponse(c, http.StatusBadRequest, "error", err.Error(), nil)
		return
	}
	httputil.NewResponse(c, http.StatusOK, "success", "", request)
}
