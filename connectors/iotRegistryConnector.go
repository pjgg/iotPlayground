package connectors

import (
	"fmt"
	"os"
	"sync"

	log "github.com/Sirupsen/logrus"
	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"
	cloudiot "google.golang.org/api/cloudiot/v1"
)

type IotRegistryConnector struct {
	Client    *cloudiot.Service
	projectID string
	region    string
}

type IotRegistryConnectorInterface interface {
	CreateRegistry(registryID string, config []*cloudiot.EventNotificationConfig) (*cloudiot.DeviceRegistry, error)
	DeleteRegistry(registryID string) (*cloudiot.Empty, error)
	GetRegistry(registryID string) (*cloudiot.DeviceRegistry, error)
	GenerateTopicName(topicName string) string
	ListRegistries() ([]*cloudiot.DeviceRegistry, error)
	SetRegistryIam(registryID string, member string, role string) (*cloudiot.Policy, error)
	GetRegistryIam(registryID string) (*cloudiot.Policy, error)
}

var onceRegistry sync.Once
var iotRegistryConnector IotRegistryConnector

func NewIotRegistryConnector(protocol Protocol, projectID string, region string) IotRegistryConnectorInterface {

	onceRegistry.Do(func() {
		ctx := context.Background()

		if protocol == HTTP {
			httpClient, err := google.DefaultClient(ctx, cloudiot.CloudPlatformScope)
			if err != nil {
				log.Fatalln(err.Error())
			}

			iotRegistryConnector.Client, err = cloudiot.New(httpClient)
			if err != nil {
				log.Fatalln(err.Error())
			}
			iotRegistryConnector.projectID = projectID
			iotRegistryConnector.region = region
		}

	})

	return &iotRegistryConnector
}

func (iotConnector *IotRegistryConnector) GenerateTopicName(topicName string) (fullTopicName string) {
	fullTopicName = fmt.Sprintf("projects/%s/topics/%s", iotConnector.projectID, topicName)
	return
}

func (iotConnector *IotRegistryConnector) CreateRegistry(registryID string, config []*cloudiot.EventNotificationConfig) (registry *cloudiot.DeviceRegistry, err error) {

	registryDef := cloudiot.DeviceRegistry{
		Id: registryID,
		EventNotificationConfigs: config,
	}

	parentPath := fmt.Sprintf("projects/%s/locations/%s", iotConnector.projectID, iotConnector.region)
	registry, err = iotConnector.Client.Projects.Locations.Registries.Create(parentPath, &registryDef).Do()
	if err != nil {
		log.Errorln(err.Error())
	} else {
		log.Debugln("Created registry:")
		log.Debugln("\tID: ", registry.Id)
		log.Debugln("\tHTTP: ", registry.HttpConfig.HttpEnabledState)
		log.Debugln("\tMQTT: ", registry.MqttConfig.MqttEnabledState)
		log.Debugln("\tName: ", registry.Name)
	}

	return
}

func (iotConnector *IotRegistryConnector) DeleteRegistry(registryID string) (empty *cloudiot.Empty, err error) {
	name := fmt.Sprintf("projects/%s/locations/%s/registries/%s", iotConnector.projectID, iotConnector.region, registryID)
	empty, err = iotConnector.Client.Projects.Locations.Registries.Delete(name).Do()
	if err != nil {
		log.Errorln(err.Error())
	}

	log.Debugln("Deleted registry")
	return
}

func (iotConnector *IotRegistryConnector) GetRegistry(registryID string) (registry *cloudiot.DeviceRegistry, err error) {
	parent := fmt.Sprintf("projects/%s/locations/%s/registries/%s", iotConnector.projectID, iotConnector.region, registryID)
	registry, err = iotConnector.Client.Projects.Locations.Registries.Get(parent).Do()
	if err != nil {
		log.Errorln(err.Error())
	}

	return
}

func (iotConnector *IotRegistryConnector) GetRegistryIam(registryID string) (policy *cloudiot.Policy, err error) {
	var req cloudiot.GetIamPolicyRequest
	path := fmt.Sprintf("projects/%s/locations/%s/registries/%s", iotConnector.projectID, iotConnector.region, registryID)
	policy, err = iotConnector.Client.Projects.Locations.Registries.GetIamPolicy(path, &req).Do()
	if err != nil {
		log.Errorln(err.Error())
	}

	log.Debugln("Policy:")
	for _, binding := range policy.Bindings {
		log.Debugln(os.Stdout, "Role: %s\n", binding.Role)
		for _, member := range binding.Members {
			log.Debugln(os.Stdout, "\tMember: %s\n", member)
		}
	}

	return
}

func (iotConnector *IotRegistryConnector) ListRegistries() (registries []*cloudiot.DeviceRegistry, err error) {
	parentPath := fmt.Sprintf("projects/%s/locations/%s", iotConnector.projectID, iotConnector.region)
	response, err := iotConnector.Client.Projects.Locations.Registries.List(parentPath).Do()
	if err != nil {
		log.Errorln(err.Error())
	}

	fmt.Println("Registries:")
	for _, registry := range response.DeviceRegistries {
		log.Debugln("\t", registry.Name)
	}

	return response.DeviceRegistries, err
}

func (iotConnector *IotRegistryConnector) SetRegistryIam(registryID string, member string, role string) (policy *cloudiot.Policy, err error) {
	req := cloudiot.SetIamPolicyRequest{
		Policy: &cloudiot.Policy{
			Bindings: []*cloudiot.Binding{
				{
					Members: []string{member},
					Role:    role,
				},
			},
		},
	}
	path := fmt.Sprintf("projects/%s/locations/%s/registries/%s", iotConnector.projectID, iotConnector.region, registryID)
	policy, err = iotConnector.Client.Projects.Locations.Registries.SetIamPolicy(path, &req).Do()
	if err != nil {
		log.Errorln(err.Error())
	}

	log.Debugln("Policy setted!")

	return
}
