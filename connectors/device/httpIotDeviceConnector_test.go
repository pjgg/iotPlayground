package device_test

import (
	"encoding/base64"
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/pjgg/iotPlayground/configuration"
	"github.com/pjgg/iotPlayground/connectors"
	"github.com/pjgg/iotPlayground/connectors/device"
	"github.com/pjgg/iotPlayground/connectors/registry"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	cloudiot "google.golang.org/api/cloudiot/v1"
)

func (suite *IotDeviceConnectorTestSuite) TestCreateDevice() {
	deviceID := "my-test-device" + randStringRunes(4)
	connector := device.NewDeviceHTTPIotConnector(suite.registryID)

	device, err := connector.CreateDevice(deviceID)

	assert.NoError(suite.T(), err, "UnexpectedError")
	assert.NotNil(suite.T(), device)
	assert.EqualValues(suite.T(), device.Id, deviceID)
	assert.EqualValues(suite.T(), device.Blocked, false)

}

func (suite *IotDeviceConnectorTestSuite) TestGetDevice() {
	deviceID := "my-test-device" + randStringRunes(4)
	connectorDevices := device.NewDeviceHTTPIotConnector(suite.registryID)
	connectorDevices.SwapToRegistry(suite.registryID)

	connectorDevices.CreateDevice(deviceID)
	device, err := connectorDevices.GetDevice(deviceID)

	assert.NoError(suite.T(), err, "UnexpectedError")
	assert.NotNil(suite.T(), device)
	assert.EqualValues(suite.T(), device.Id, deviceID)
	assert.EqualValues(suite.T(), device.Blocked, false)

}

func (suite *IotDeviceConnectorTestSuite) TestSetDeviceConfig() {
	deviceID := "my-test-device" + randStringRunes(4)
	connectorDevices := device.NewDeviceHTTPIotConnector(suite.registryID)
	connectorDevices.SwapToRegistry(suite.registryID)

	connectorDevices.CreateDevice(deviceID)
	config, err := connectorDevices.SetDeviceConfig(deviceID, "{networkID:'myNetworkID'}")

	assert.NoError(suite.T(), err, "UnexpectedError")
	assert.NotNil(suite.T(), config)
	decodedData, _ := base64.StdEncoding.DecodeString(config.BinaryData)
	assert.EqualValues(suite.T(), decodedData, "{networkID:'myNetworkID'}")

}

func (suite *IotDeviceConnectorTestSuite) TestGetDeviceConfigs() {
	deviceID := "my-test-device" + randStringRunes(4)
	connectorDevices := device.NewDeviceHTTPIotConnector(suite.registryID)
	connectorDevices.CreateDevice(deviceID)
	connectorDevices.SetDeviceConfig(deviceID, "{networkID:'myNetworkID'}")
	config, err := connectorDevices.GetDeviceConfigs(deviceID)

	assert.NoError(suite.T(), err, "UnexpectedError")
	assert.EqualValues(suite.T(), len(config) > 0, true)

}

func (suite *IotDeviceConnectorTestSuite) TestGetDeviceStates() {
	deviceID := "my-test-device" + randStringRunes(4)
	connectorDevices := device.NewDeviceHTTPIotConnector(suite.registryID)
	connectorDevices.SwapToRegistry(suite.registryID)

	connectorDevices.CreateDevice(deviceID)

	deviceStates, err := connectorDevices.GetDeviceStates(deviceID)
	assert.NoError(suite.T(), err, "UnexpectedError")
	assert.EqualValues(suite.T(), len(deviceStates), 0)

}

func (suite *IotDeviceConnectorTestSuite) TestListDevices() {
	deviceID := "my-test-device" + randStringRunes(4)
	connectorDevices := device.NewDeviceHTTPIotConnector(suite.registryID)
	connectorDevices.SwapToRegistry(suite.registryID)

	connectorDevices.CreateDevice(deviceID)

	devices, err := connectorDevices.ListDevices()
	assert.NoError(suite.T(), err, "UnexpectedError")
	assert.EqualValues(suite.T(), len(devices), 1)

}

func (suite *IotDeviceConnectorTestSuite) TestPatchDevice() {
	deviceID := "my-test-device" + randStringRunes(4)
	connectorDevices := device.NewDeviceHTTPIotConnector(suite.registryID)
	connectorDevices.SwapToRegistry(suite.registryID)

	connectorDevices.CreateDevice(deviceID)
	deviceToUpdate := &cloudiot.Device{
		Metadata: map[string]string{"key": "value"},
	}

	deviceUpdated, err := connectorDevices.PatchDevice(deviceID, deviceToUpdate, "Metadata")
	assert.NoError(suite.T(), err, "UnexpectedError")
	assert.EqualValues(suite.T(), deviceUpdated.Metadata["key"], "value")

}

type IotDeviceConnectorTestSuite struct {
	suite.Suite
	configuration *configuration.Configuration
	registryID    string
}

func (suite *IotDeviceConnectorTestSuite) SetupTest() {
	connector := registry.NewHTTPIotRegistryConnector(connectors.HTTP, suite.configuration.GcloudProjectID, suite.configuration.GcloudRegion)

	eventNotificationConfigs := []*cloudiot.EventNotificationConfig{
		{
			PubsubTopicName: connector.GenerateTopicName(suite.configuration.DeviceTelemetryTopic),
		},
	}

	_, err := connector.CreateRegistry(suite.registryID, eventNotificationConfigs)
	assert.NoError(suite.T(), err, "UnexpectedError")

}

func (suite *IotDeviceConnectorTestSuite) TearDownTest() {
	connector := registry.NewHTTPIotRegistryConnector(connectors.HTTP, suite.configuration.GcloudProjectID, suite.configuration.GcloudRegion)
	connectorDevices := device.NewDeviceHTTPIotConnector(suite.registryID)
	deviceList, _ := connectorDevices.ListDevices()
	for _, device := range deviceList {
		connectorDevices.DeleteDevice(device.Id)
	}

	connector.DeleteRegistry(suite.registryID)
}

func TestIotDeviceConnectorTestSuite(t *testing.T) {

	configInit()
	rand.Seed(time.Now().UnixNano())

	iotReg := new(IotDeviceConnectorTestSuite)
	iotReg.configuration = configuration.New()
	iotReg.registryID = "test-registry-" + randStringRunes(4)
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
