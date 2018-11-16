package server

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/KyberNetwork/node-monitor/collector"
	"github.com/KyberNetwork/node-monitor/storage"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Server struct {
	r         *gin.Engine
	collector *collector.Collector
	storage   Storage
}

func NewServer() *Server {
	r := gin.Default()

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	corsConfig.MaxAge = 5 * time.Minute

	storage := storage.NewStorage()

	r.Use(cors.New(corsConfig))
	return &Server{
		r:         r,
		collector: collector.NewCollector(storage),
		storage:   storage,
	}
}

func (self *Server) GetData(c *gin.Context) {
	data, err := self.storage.GetData()
	if err != nil {
		c.JSON(
			http.StatusOK,
			gin.H{
				"success": false,
				"reason":  err.Error(),
			},
		)
		return
	}
	c.JSON(
		http.StatusOK,
		gin.H{
			"success": true,
			"data":    data,
		},
	)
}

func (self *Server) GetDataNormal(c *gin.Context) {
	fromTime := c.Query("fromTime")
	toTime := c.Query("toTime")
	fromTimeUint, err := strconv.ParseUint(fromTime, 10, 64)
	if err != nil {
		log.Println(err)
		c.JSON(
			http.StatusOK,
			gin.H{
				"success": false,
				"reason":  "time is not a number",
			},
		)
		return
	}
	toTimeUint, err := strconv.ParseUint(toTime, 10, 64)
	if err != nil {
		toTimeUint = uint64(time.Now().Unix())
	}
	data, err := self.storage.GetDataNormal(fromTimeUint, toTimeUint)
	if err != nil {
		c.JSON(
			http.StatusOK,
			gin.H{
				"success": false,
				"reason":  err.Error(),
			},
		)
		return
	}
	c.JSON(
		http.StatusOK,
		gin.H{
			"success": true,
			"data":    data,
		},
	)
}

func (self *Server) GetStatNormal(c *gin.Context) {
	data := self.collector.GetStatNormalData()
	c.JSON(
		http.StatusOK,
		gin.H{
			"success": true,
			"data":    data,
		},
	)
}

func (self *Server) GetTotalRequest(c *gin.Context) {
	data := self.collector.GetTotalRequest()
	c.JSON(
		http.StatusOK,
		gin.H{
			"success": true,
			"data":    data,
		},
	)
}

func (self *Server) Run() {
	go self.collector.CollectData()
	// self.r.Use(static.Serve("/", static.LocalFile("/go/src/github.com/Gin/node-tracker/app", true)))

	api := self.r.Group("/api")
	api.GET("/stress-data", self.GetData)
	api.GET("/normal-data", self.GetDataNormal)
	api.GET("/stat-normal-data", self.GetStatNormal)
	api.GET("/request", self.GetTotalRequest)

	self.r.Run(":8000")
}
