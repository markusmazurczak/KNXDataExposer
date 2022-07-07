package knx

import (
	"errors"
	"fmt"

	"gitea.realmottek.duckdns.org/nexus/KNXDataExposer/handler"
	"github.com/spf13/viper"
	"github.com/vapourismo/knx-go/knx"
	"github.com/vapourismo/knx-go/knx/dpt"
	"go.uber.org/zap"
)

//Connects to an KNX gateway which can be either a tunnel or a router connection.
//	- g: The connection string (IP:PORT)
//	- d: The calculated datapoint structure map
//	- dh: The DataHandler for communicating with the DB
func Start_knx_listener(logger *zap.Logger, g string, d map[string]string, dh *handler.DataHandler) error {
	if viper.GetString("knx.gatewayType") == "tunnel" {
		return start_tunnel_listener(logger, g, d, dh)
	} else if viper.GetString("knx.gatewayType") == "router" {
		return start_router_listener(logger, g, d, dh)
	} else {
		logger.Fatal("KNX", zap.String("gatewayType", viper.GetString("knx.gatewayType")), zap.String("msg", "Invalid gateway type"))
		return errors.New("Invalid KNX gateway type configured")
	}
}

//private function to connect using a tunnel
func start_tunnel_listener(logger *zap.Logger, g string, d map[string]string, dh *handler.DataHandler) error {
	client, err := knx.NewGroupTunnel(g, knx.DefaultTunnelConfig)
	if err != nil {
		logger.Fatal(err.Error())
	}
	defer client.Close()

	for msg := range client.Inbound() {
		if val, ok := d[msg.Destination.String()]; ok {
			temp, _ := dpt.Produce(val)
			err := temp.Unpack(msg.Data)
			if err != nil {
				continue
			}
			save_data(logger, msg.Destination.String(), val, temp, dh)
		} else {
			logger.Debug("Ignoring", zap.String("ga", msg.Destination.String()))
		}
	}
	return err
}

//private function to connect using a router
func start_router_listener(logger *zap.Logger, g string, d map[string]string, dh *handler.DataHandler) error {
	client, err := knx.NewGroupRouter(g, knx.DefaultRouterConfig)
	if err != nil {
		logger.Fatal(err.Error())
	}
	defer client.Close()

	for msg := range client.Inbound() {
		if val, ok := d[msg.Destination.String()]; ok {
			temp, _ := dpt.Produce(val)
			err := temp.Unpack(msg.Data)
			if err != nil {
				continue
			}
			save_data(logger, msg.Destination.String(), val, temp, dh)
		} else {
			logger.Debug("Ignoring", zap.String("ga", msg.Destination.String()))
		}
	}
	return err
}

//Saves the value of an KNX datagram with it's corresponding group address to the database
func save_data(logger *zap.Logger, dest_ga string, dpt string, d dpt.DatapointValue, dh *handler.DataHandler) error {
	v := fmt.Sprintf("%s", d)
	logger.Debug("Data", zap.String("ga", dest_ga), zap.String("dpt", dpt), zap.String("value", v))
	return handler.Insert_dataset(dest_ga, v, dh)
}
