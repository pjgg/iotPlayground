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

type MqttIotDeviceConnectorTestSuite struct {
	suite.Suite
	configuration *configuration.Configuration
	registryID    string
	deviceIDOne   string
	deviceIDTwo   string
}

func (suite *MqttIotDeviceConnectorTestSuite) TestPublishMsg() {
	connectorDevices := connectors.NewMqttIotConnector(suite.registryID, suite.deviceIDOne)
	token := connectorDevices.PublishMsg(suite.deviceIDOne, "events", "test")

	if token.WaitTimeout(time.Minute*time.Duration(10)) && token.Error() != nil {
		assert.Failf(suite.T(), "error publish MQTT msg: ", token.Error().Error())
	}

}

func (suite *MqttIotDeviceConnectorTestSuite) SetupTest() {
	connector := connectors.NewIotRegistryConnector(connectors.HTTP, suite.configuration.GcloudProjectID, suite.configuration.GcloudRegion)

	eventNotificationConfigs := []*cloudiot.EventNotificationConfig{
		{
			PubsubTopicName: connector.GenerateTopicName("events"),
		},
	}

	_, err := connector.CreateRegistry(suite.registryID, eventNotificationConfigs)

	if err != nil {
		assert.Failf(suite.T(), "error when trying to create a register", err.Error())
	} else {

		connectorHttpDevices := connectors.NewHttpIotConnector(suite.registryID)

		connectorHttpDevices.CreateDevice(suite.deviceIDOne)
		connectorHttpDevices.CreateDevice(suite.deviceIDTwo)
	}
}

func (suite *MqttIotDeviceConnectorTestSuite) TearDownTest() {
	connector := connectors.NewIotRegistryConnector(connectors.HTTP, suite.configuration.GcloudProjectID, suite.configuration.GcloudRegion)
	connectorHttpDevices := connectors.NewHttpIotConnector(suite.registryID)

	deviceList, _ := connectorHttpDevices.ListDevices()
	for _, device := range deviceList {
		connectorHttpDevices.DeleteDevice(device.Id)
	}

	connector.DeleteRegistry(suite.registryID)
}

func TestMqttIotDeviceConnectorTestSuite(t *testing.T) {

	configInit()
	rand.Seed(time.Now().UnixNano())

	iotReg := new(MqttIotDeviceConnectorTestSuite)
	iotReg.configuration = configuration.New()
	iotReg.registryID = "test-registry-" + randStringRunes(4)
	iotReg.deviceIDOne = "test-device-" + randStringRunes(4)
	iotReg.deviceIDTwo = "test-device-" + randStringRunes(4)

	suite.Run(t, iotReg)
}
