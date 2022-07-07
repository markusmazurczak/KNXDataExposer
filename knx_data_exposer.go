package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"

	"gitea.realmottek.duckdns.org/nexus/KNXDataExposer/db"
	"gitea.realmottek.duckdns.org/nexus/KNXDataExposer/handler"
	my_knx "gitea.realmottek.duckdns.org/nexus/KNXDataExposer/knx"
	"gitea.realmottek.duckdns.org/nexus/KNXDataExposer/util"
)

var g errgroup.Group

func main() {
	logger := util.Init_logger()
	defer logger.Sync()

	_, err := util.Init_config(logger)
	if err != nil {
		logger.Sugar().Fatalf("Cannot parse config: %s", err)
		panic(err)
	}

	server_string := fmt.Sprintf("%s:%d", viper.GetString("server.bindIP"), viper.GetInt("server.bindPort"))
	gateway_string := fmt.Sprintf("%s:%d", viper.GetString("knx.gatewayIP"), viper.GetInt("knx.gatewayPort"))

	dpts := create_datapoint_structure(logger, viper.GetStringMapStringSlice("datapoints"))

	dataHandler := handler.DataHandler{
		DB: db.Init(logger),
	}

	logger.Info("REST-Server", zap.String("url", server_string))
	gin_server := &http.Server{
		Addr:         server_string,
		Handler:      router(logger, &dataHandler),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	g.Go(func() error {
		return gin_server.ListenAndServe()
	})

	logger.Info("KNX router", zap.String("url", gateway_string))
	g.Go(func() error {
		//return start_knx_listener(logger, gateway_string, dpts, &dataHandler)
		return my_knx.Start_knx_listener(logger, gateway_string, dpts, &dataHandler)
	})

	if err := g.Wait(); err != nil {
		logger.Sugar().Fatalf("%v", err)
	}
}

//Changes the datapoint configuration structure in a way that group addresses are keys
//and the DPT's are values
func create_datapoint_structure(logger *zap.Logger, s map[string][]string) map[string]string {
	r := make(map[string]string)

	for k, a := range s {
		//Convert into valid KNX dpt format
		name := k[4:len(k)-3] + "." + k[len(k)-3:]
		for _, v := range a {
			logger.Debug("Transform datapoint", zap.String("ga", v), zap.String("dpt", name))
			r[v] = name
		}

	}
	return r
}

//GIN middleware function to lookup a value for a given knx group address and return it
func router(logger *zap.Logger, dh *handler.DataHandler) http.Handler {
	router := gin.New()
	router.SetTrustedProxies(nil)
	router.Use(gin.Recovery())

	router.GET("/dataset", func(c *gin.Context) {
		logger.Info("Got call", zap.String("remote_ip", c.ClientIP()), zap.String("group_address", c.Query("ga")))
		d, err := handler.Get_dataset(c.Query("ga"), dh)
		if err != nil {
			logger.Error("Lookup error", zap.String("remote_ip", c.ClientIP()), zap.String("group_address", c.Query("ga")), zap.String("msg", err.Error()))
		}
		c.JSON(200, d)
	})
	return router
}
