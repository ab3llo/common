package discovery

import (
	"context"
	"os"
	"strconv"

	servicedirectory "cloud.google.com/go/servicedirectory/apiv1"
	"cloud.google.com/go/servicedirectory/apiv1/servicedirectorypb"
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

func RegisterServiceWithConsul(serviceName, grpcAddr, consulHost, consulPort string) {
	consulConfig := consulapi.DefaultConfig()
	consulConfig.Address = consulHost + ":" + consulPort
	client, err := consulapi.NewClient(consulConfig)
	if err != nil {
		log.Error("Failed to create Consul client: " + err.Error())
		os.Exit(1)
	}

	// Extract port from grpcAddr (e.g., ":5001" -> "5001")
	port := grpcAddr
	if port[0] == ':' {
		port = port[1:]
	}

	registration := &consulapi.AgentServiceRegistration{
		ID:      serviceName,
		Name:    serviceName,
		Address: serviceName, // Use service name as address for Docker networking
		Port:    parsePort(port),
	}

	if err := client.Agent().ServiceRegister(registration); err != nil {
		log.Error("Failed to register service with Consul: " + err.Error())
		os.Exit(1) // Exit if we can't register the service
	}
	log.Info("Service registered with Consul: " + registration.Name)
}

func RegisterService(serviceName, grpcAddr, gcpProject, gcpRegion, namespace string) {
	ctx := context.Background()

	client, err := servicedirectory.NewRegistrationClient(ctx)
	if err != nil {
		log.Error("Failed to create service directory registration client: " + err.Error())
		os.Exit(1)
	}
	defer client.Close()

	parent := "projects/" + gcpProject + "/locations/" + gcpRegion + "/namespaces/" + namespace

	serviceReq := &servicedirectorypb.CreateServiceRequest{
		Parent:    parent,
		ServiceId: serviceName,
		Service: &servicedirectorypb.Service{
			Annotations: map[string]string{
				"key1": "value1",
				"key2": "value2",
			},
		},
	}

	srv, err := client.CreateService(ctx, serviceReq)

	if err != nil {
		log.Error("Failed to create service: " + err.Error())
		// Proceed even if service creation fails, it might already exist
	}

	_, err = client.CreateEndpoint(
		ctx,
		&servicedirectorypb.CreateEndpointRequest{
			Parent:     parent,
			EndpointId: serviceName + "-endpoint",
			Endpoint: &servicedirectorypb.Endpoint{
				Address: grpcAddr,
				Port:    int32(parsePort(grpcAddr)),
				Annotations: map[string]string{
					"key1": "value1",
					"key2": "value2",
				},
			},
		},
	)
	if err != nil {
		log.Error("Failed to create endpoint: " + err.Error())
		// Proceed even if service creation fails, it might already exist
	}

	log.Info("Service registered with Service Directory: " + serviceName)
}

func GetInstance(serviceName, gcpProject, gcpRegion, namespace string) (*servicedirectorypb.Endpoint, error) {
	ctx := context.Background()

	lookupClient, err := servicedirectory.NewLookupClient(ctx)
	if err != nil {
		log.Error("Failed to create service directory lookup client: " + err.Error())
		return nil, err
	}
	defer lookupClient.Close()

	req := &servicedirectorypb.ResolveServiceRequest{
		Name: "projects/" + gcpProject + "/locations/" + gcpRegion + "/namespaces/" + namespace + "/services/" + serviceName,
	}

	resp, err := lookupClient.ResolveService(ctx, req)
	if err != nil {
		log.Error("Failed to resolve service: " + err.Error())
		return nil, err
	}

	if len(resp.GetService().GetEndpoints()) == 0 {
		return nil, os.ErrNotExist
	}

	return resp.GetService().GetEndpoints()[0], nil
}

func GetInstanceWithConsul(serviceName, consulHost, consulPort string) *consulapi.ServiceEntry {
	consulConfig := consulapi.DefaultConfig()
	consulConfig.Address = consulHost + ":" + consulPort
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

	return instances[0]
}
