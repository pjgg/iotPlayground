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

// HTTPIotRegistryConnector handler registry admin request.
type HTTPIotRegistryConnector struct {
	Client    *cloudiot.Service
	projectID string
	region    string
}

// HTTPIotRegistryConnectorInterface define registry IoT admin request behavior.
type HTTPIotRegistryConnectorInterface interface {
	CreateRegistry(registryID string, config []*cloudiot.EventNotificationConfig) (*cloudiot.DeviceRegistry, error)
	DeleteRegistry(registryID string) (*cloudiot.Empty, error)
	GetRegistry(registryID string) (*cloudiot.DeviceRegistry, error)
	GenerateTopicName(topicName string) string
	ListRegistries() ([]*cloudiot.DeviceRegistry, error)
	SetRegistryIam(registryID string, member string, role string) (*cloudiot.Policy, error)
	GetRegistryIam(registryID string) (*cloudiot.Policy, error)
}

var onceRegistry sync.Once
var iotRegistryConnector HTTPIotRegistryConnector

// NewHTTPIotRegistryConnector create a single instance of HTTPIotRegistryConnector
func NewHTTPIotRegistryConnector(protocol Protocol, projectID string, region string) HTTPIotRegistryConnectorInterface {

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

// GenerateTopicName create a topic name according google spec.
func (iotConnector *HTTPIotRegistryConnector) GenerateTopicName(topicName string) (fullTopicName string) {
	fullTopicName = fmt.Sprintf("projects/%s/topics/%s", iotConnector.projectID, topicName)
	return
}

// CreateRegistry create a device registry witha  given registryID.
func (iotConnector *HTTPIotRegistryConnector) CreateRegistry(registryID string, config []*cloudiot.EventNotificationConfig) (registry *cloudiot.DeviceRegistry, err error) {

	registryDef := cloudiot.DeviceRegistry{
		Id: registryID,
		EventNotificationConfigs: config,
	}

	parentPath := fmt.Sprintf("projects/%s/locations/%s", iotConnector.projectID, iotConnector.region)
	if registry, err = iotConnector.Client.Projects.Locations.Registries.Create(parentPath, &registryDef).Do(); err == nil {
		log.Debugln("Created registry:")
		log.Debugln("\tID: ", registry.Id)
		log.Debugln("\tHTTP: ", registry.HttpConfig.HttpEnabledState)
		log.Debugln("\tMQTT: ", registry.MqttConfig.MqttEnabledState)
		log.Debugln("\tName: ", registry.Name)
	}

	return
}

// DeleteRegistry remove an existing registry based in his registryID.
func (iotConnector *HTTPIotRegistryConnector) DeleteRegistry(registryID string) (empty *cloudiot.Empty, err error) {
	name := fmt.Sprintf("projects/%s/locations/%s/registries/%s", iotConnector.projectID, iotConnector.region, registryID)
	if iotConnector.Client.Projects.Locations.Registries.Delete(name).Do(); err == nil {
		log.Debugln("Deleted registry")
	}

	return
}

// GetRegistry retrieve a registry based in his registryID.
func (iotConnector *HTTPIotRegistryConnector) GetRegistry(registryID string) (registry *cloudiot.DeviceRegistry, err error) {
	parent := fmt.Sprintf("projects/%s/locations/%s/registries/%s", iotConnector.projectID, iotConnector.region, registryID)
	registry, err = iotConnector.Client.Projects.Locations.Registries.Get(parent).Do()

	return
}

// GetRegistryIam retrieve registry Iam based in his registryID.
func (iotConnector *HTTPIotRegistryConnector) GetRegistryIam(registryID string) (policy *cloudiot.Policy, err error) {
	var req cloudiot.GetIamPolicyRequest
	path := fmt.Sprintf("projects/%s/locations/%s/registries/%s", iotConnector.projectID, iotConnector.region, registryID)
	if policy, err = iotConnector.Client.Projects.Locations.Registries.GetIamPolicy(path, &req).Do(); err == nil {
		log.Debugln("Policy:")
		for _, binding := range policy.Bindings {
			log.Debugln(os.Stdout, "Role: %s\n", binding.Role)
			for _, member := range binding.Members {
				log.Debugln(os.Stdout, "\tMember: %s\n", member)
			}
		}
	}

	return
}

// SetRegistryIam update or create a registry Iam for a given registryID.
func (iotConnector *HTTPIotRegistryConnector) SetRegistryIam(registryID string, member string, role string) (policy *cloudiot.Policy, err error) {
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
	if policy, err = iotConnector.Client.Projects.Locations.Registries.SetIamPolicy(path, &req).Do(); err == nil {
		log.Debugln("Policy setted!")
	}

	return
}

// ListRegistries retrieve a list of registries of the current project.
func (iotConnector *HTTPIotRegistryConnector) ListRegistries() (registries []*cloudiot.DeviceRegistry, err error) {
	parentPath := fmt.Sprintf("projects/%s/locations/%s", iotConnector.projectID, iotConnector.region)
	if response, err := iotConnector.Client.Projects.Locations.Registries.List(parentPath).Do(); err == nil {
		fmt.Println("Registries:")
		for _, registry := range response.DeviceRegistries {
			log.Debugln("\t", registry.Name)
		}
		registries = response.DeviceRegistries
	}

	return registries, err
}
