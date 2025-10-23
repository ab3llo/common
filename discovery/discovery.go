package discovery

import (
	"math/rand"
	"os"
	"strconv"

	consulapi "github.com/hashicorp/consul/api"
	"github.com/hmlylab/common/logger"
)

var (
	log = logger.NewLogger()
)

func parsePort(portStr string) int {
	port, err := strconv.Atoi(portStr)
	if err != nil {
		log.Error("Could not parse port: " + err.Error())
		return 0
	}
	return port
}

func RegisterServiceWithConsul(serviceName, addr, consulHost, consulPort string) {
	consulConfig := consulapi.DefaultConfig()
	consulConfig.Address = resolveConsulAddress(consulPort, consulHost)
	client, err := consulapi.NewClient(consulConfig)
	if err != nil {
		log.Error("Failed to create Consul client: " + err.Error())
		os.Exit(1)
	}

	// Extract port from grpcAddr (e.g., ":5001" -> "5001")
	port := addr
	if port[0] == ':' {
		port = port[1:]
	}

	registration := &consulapi.AgentServiceRegistration{
		ID:      serviceName,
		Name:    serviceName,
		Address: serviceName + ":" + port, // Use service name as address for Docker networking
		Port:    parsePort(port),
	}

	if err := client.Agent().ServiceRegister(registration); err != nil {
		log.Error("Failed to register service with Consul: " + err.Error())
		os.Exit(1) // Exit if we can't register the service
	}
	log.Info("Service registered with Consul: " + registration.Name)
}

func GetInstanceWithConsul(serviceName, consulHost, consulPort string) *consulapi.ServiceEntry {
	consulConfig := consulapi.DefaultConfig()
	consulConfig.Address = resolveConsulAddress(consulPort, consulHost)
	client, err := consulapi.NewClient(consulConfig)
	if err != nil {
		log.Error("Failed to create Consul client: " + err.Error())
		os.Exit(1)
	}

	instances, _, err := client.Health().Service(serviceName, "", true, nil)
	if err != nil {
		log.Error("Failed to get service from Consul: " + err.Error())
		os.Exit(1)
	}

	if len(instances) == 0 {
		log.Error("No service found with name: " + serviceName)
		os.Exit(1)
	}

	n := rand.Intn(len(instances))
	return instances[n]
}

func resolveConsulAddress(consulPort, consulHost string) string {
	consulAddress := ""
	if consulPort == "" {
		consulAddress = consulHost
	} else {
		consulAddress = consulHost + ":" + consulPort
	}
	return consulAddress
}
