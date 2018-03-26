package connectors_test

import (
	"math/rand"
	"testing"
	"time"

	"github.com/pjgg/iotPlayground/configuration"
	"github.com/pjgg/iotPlayground/connectors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	cloudiot "google.golang.org/api/cloudiot/v1"
)

func (suite *IotRegistryConnectorTestSuite) TestGenerateTopicName() {
	connector := connectors.NewHTTPIotRegistryConnector(connectors.HTTP, suite.configuration.GcloudProjectID, suite.configuration.GcloudRegion)

	expectedTopicName := "projects/" + suite.configuration.GcloudProjectID + "/topics/" + suite.configuration.DeviceTelemetryTopic
	result := connector.GenerateTopicName(suite.configuration.DeviceTelemetryTopic)

	assert.EqualValues(suite.T(), result, expectedTopicName)

}

func (suite *IotRegistryConnectorTestSuite) TestGetRegistry() {
	connector := connectors.NewHTTPIotRegistryConnector(connectors.HTTP, suite.configuration.GcloudProjectID, suite.configuration.GcloudRegion)

	registry, err := connector.GetRegistry(suite.registryID)
	assert.NoError(suite.T(), err, "UnexpectedError")
	assert.NotNil(suite.T(), registry)
	assert.EqualValues(suite.T(), registry.Id, suite.registryID)
	assert.EqualValues(suite.T(), registry.HttpConfig.HttpEnabledState, "HTTP_ENABLED")
}

func (suite *IotRegistryConnectorTestSuite) listRegistries() {
	connector := connectors.NewHTTPIotRegistryConnector(connectors.HTTP, suite.configuration.GcloudProjectID, suite.configuration.GcloudRegion)

	registries, _ := connector.ListRegistries()
	assert.EqualValues(suite.T(), 1, len(registries))
}

func (suite *IotRegistryConnectorTestSuite) setRegistryIamTest() {
	connector := connectors.NewHTTPIotRegistryConnector(connectors.HTTP, suite.configuration.GcloudProjectID, suite.configuration.GcloudRegion)

	policy, _ := connector.SetRegistryIam(suite.registryID, "pablosDevice@bq.com", "admin")
	assert.NotNil(suite.T(), policy)
	assert.EqualValues(suite.T(), policy.Version, 1)
}

func (suite *IotRegistryConnectorTestSuite) getRegistryIamTest() {
	connector := connectors.NewHTTPIotRegistryConnector(connectors.HTTP, suite.configuration.GcloudProjectID, suite.configuration.GcloudRegion)
	policy, _ := connector.GetRegistryIam(suite.registryID)
	assert.NotNil(suite.T(), policy)
	assert.EqualValues(suite.T(), policy.Version, 1)
}

type IotRegistryConnectorTestSuite struct {
	suite.Suite
	configuration *configuration.Configuration
	registryID    string
}

func (suite *IotRegistryConnectorTestSuite) SetupTest() {
	connector := connectors.NewHTTPIotRegistryConnector(connectors.HTTP, suite.configuration.GcloudProjectID, suite.configuration.GcloudRegion)

	eventNotificationConfigs := []*cloudiot.EventNotificationConfig{
		{
			PubsubTopicName: connector.GenerateTopicName(suite.configuration.DeviceTelemetryTopic),
		},
	}

	_, err := connector.CreateRegistry(suite.registryID, eventNotificationConfigs)
	assert.NoError(suite.T(), err, "UnexpectedError")

}

func (suite *IotRegistryConnectorTestSuite) TearDownTest() {
	connector := connectors.NewHTTPIotRegistryConnector(connectors.HTTP, suite.configuration.GcloudProjectID, suite.configuration.GcloudRegion)
	connector.DeleteRegistry(suite.registryID)
}

func TestIotRegistryConnectorTestSuite(t *testing.T) {

	configInit()
	rand.Seed(time.Now().UnixNano())

	iotReg := new(IotRegistryConnectorTestSuite)
	iotReg.configuration = configuration.New()
	iotReg.registryID = "test-registry-" + randStringRunes(4)
	suite.Run(t, iotReg)
}
