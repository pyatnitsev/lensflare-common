package consul

import (
	"fmt"
	"os"
	"strconv"

	consulapi "github.com/hashicorp/consul/api"
)

func RegisterService() error {
	consulAddr := os.Getenv("CONSUL_HTTP_ADDR")
	if consulAddr == "" {
		return fmt.Errorf("CONSUL_HTTP_ADDR not set")
	}

	portStr := os.Getenv("HTTP_PORT")
	if portStr == "" {
		return fmt.Errorf("HTTP_PORT not set")
	}
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return fmt.Errorf("Invalid HTTP_PORT: %v", err)
	}

	serviceID := os.Getenv("CONSUL_SERVICE_ID")
	if serviceID == "" {
		return fmt.Errorf("CONSUL_SERVICE_ID not set")
	}
	serviceName := os.Getenv("CONSUL_SERVICE_NAME")
	if serviceName == "" {
		return fmt.Errorf("CONSUL_SERVICE_NAME not set")
	}

	address := os.Getenv("HOST_ADDRESS")
	if address == "" {
		return fmt.Errorf("HOST_ADDRESS not set")
	}

	config := consulapi.DefaultConfig()
	config.Address = consulAddr
	client, err := consulapi.NewClient(config)
	if err != nil {
		return fmt.Errorf("NewClient error: %v", err)
	}

	check := &consulapi.AgentServiceCheck{
		HTTP:     fmt.Sprintf("http://%s:%d/health", address, port),
		Interval: "10s",
		Timeout:  "3s",
	}

	reg := &consulapi.AgentServiceRegistration{
		ID:      serviceID,
		Name:    serviceName,
		Address: address,
		Port:    port,
		Check:   check,
	}

	if err := client.Agent().ServiceRegister(reg); err != nil {
		return fmt.Errorf("ServiceRegister error: %v", err)
	}
	return nil
}
