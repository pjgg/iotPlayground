package device_test

import (
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/pjgg/iotPlayground/configuration"
	"github.com/pjgg/iotPlayground/connectors"
	"github.com/spf13/viper"
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
	msg := "test"
	connectorDevices := connectors.NewMQTTIotConnector(suite.registryID, suite.deviceIDOne)
	token := connectorDevices.PublishMsg(suite.deviceIDOne, suite.configuration.DeviceTelemetryTopic, msg, connectors.AtMostOnce)

	if token.WaitTimeout(time.Minute*time.Duration(10)) && token.Error() != nil {
		assert.NoError(suite.T(), token.Error(), "error publish MQTT")
	}

}

func (suite *MqttIotDeviceConnectorTestSuite) SetupTest() {
	connector := connectors.NewHTTPIotRegistryConnector(connectors.HTTP, suite.configuration.GcloudProjectID, suite.configuration.GcloudRegion)

	eventNotificationConfigs := []*cloudiot.EventNotificationConfig{
		{
			PubsubTopicName: connector.GenerateTopicName(suite.configuration.DeviceTelemetryTopic),
		},
	}

	_, err := connector.CreateRegistry(suite.registryID, eventNotificationConfigs)
	assert.NoError(suite.T(), err, "error publish MQTT")

	connectorHTTPDevices := connectors.NewHTTPIotConnector(suite.registryID)
	connectorHTTPDevices.SwapToRegistry(suite.registryID)

	connectorHTTPDevices.CreateDevice(suite.deviceIDOne)
	connectorHTTPDevices.CreateDevice(suite.deviceIDTwo)

}

func (suite *MqttIotDeviceConnectorTestSuite) TearDownTest() {
	connector := connectors.NewHTTPIotRegistryConnector(connectors.HTTP, suite.configuration.GcloudProjectID, suite.configuration.GcloudRegion)
	connectorHttpDevices := connectors.NewHTTPIotConnector(suite.registryID)

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

func randStringRunes(n int) string {
	var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func configInit() {
	viper.SetConfigName("config")
	configPath, exist := os.LookupEnv("CONFIG_PATH")
	if exist {
		viper.AddConfigPath(configPath)
	}
	viper.AddConfigPath("../../")
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}
}
