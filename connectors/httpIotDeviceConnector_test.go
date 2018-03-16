package connectors_test

/*
import (
	"encoding/base64"
	"math/rand"
	"testing"
	"time"

	"github.com/pjgg/iotPlayground/configuration"
	"github.com/pjgg/iotPlayground/connectors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	cloudiot "google.golang.org/api/cloudiot/v1"
)

func (suite *IotDeviceConnectorTestSuite) TestCreateDevice() {
	deviceID := "my-test-device" + randStringRunes(4)
	connector := connectors.NewHttpIotConnector(suite.registryID)
	device, err := connector.CreateDevice(deviceID)

	if err != nil {
		assert.Failf(suite.T(), "error when trying to create a device", err.Error())
	} else {
		assert.EqualValues(suite.T(), device.Id, deviceID)
		assert.EqualValues(suite.T(), device.Blocked, false)
	}

}

func (suite *IotDeviceConnectorTestSuite) TestGetDevice() {
	deviceID := "my-test-device" + randStringRunes(4)
	connectorDevices := connectors.NewHttpIotConnector(suite.registryID)
	connectorDevices.CreateDevice(deviceID)
	device, err := connectorDevices.GetDevice(deviceID)

	if err != nil {
		assert.Failf(suite.T(), "error when trying to create a device", err.Error())
	} else {
		assert.EqualValues(suite.T(), device.Id, deviceID)
		assert.EqualValues(suite.T(), device.Blocked, false)
	}

}

func (suite *IotDeviceConnectorTestSuite) TestSetDeviceConfig() {
	deviceID := "my-test-device" + randStringRunes(4)
	connectorDevices := connectors.NewHttpIotConnector(suite.registryID)
	connectorDevices.CreateDevice(deviceID)
	config, err := connectorDevices.SetDeviceConfig(deviceID, "{networkID:'myNetworkID'}")

	if err != nil {
		assert.Failf(suite.T(), "set config error", err.Error())
	} else {
		decodedData, _ := base64.StdEncoding.DecodeString(config.BinaryData)
		assert.EqualValues(suite.T(), decodedData, "{networkID:'myNetworkID'}")
	}

}

func (suite *IotDeviceConnectorTestSuite) TestGetDeviceConfigs() {
	deviceID := "my-test-device" + randStringRunes(4)
	connectorDevices := connectors.NewHttpIotConnector(suite.registryID)
	connectorDevices.CreateDevice(deviceID)
	connectorDevices.SetDeviceConfig(deviceID, "{networkID:'myNetworkID'}")
	config, err := connectorDevices.GetDeviceConfigs(deviceID)

	if err != nil {
		assert.Failf(suite.T(), "set config error", err.Error())
	} else {
		assert.EqualValues(suite.T(), len(config) > 0, true)
	}

}

func (suite *IotDeviceConnectorTestSuite) TestGetDeviceStates() {
	deviceID := "my-test-device" + randStringRunes(4)
	connectorDevices := connectors.NewHttpIotConnector(suite.registryID)
	connectorDevices.CreateDevice(deviceID)

	deviceStates, err := connectorDevices.GetDeviceStates(deviceID)
	if err != nil {
		assert.Failf(suite.T(), "error retriving states", err.Error())
	} else {
		assert.EqualValues(suite.T(), len(deviceStates), 0)
	}

}

func (suite *IotDeviceConnectorTestSuite) TestListDevices() {
	deviceID := "my-test-device" + randStringRunes(4)
	connectorDevices := connectors.NewHttpIotConnector(suite.registryID)
	connectorDevices.CreateDevice(deviceID)

	devices, err := connectorDevices.ListDevices()
	if err != nil {
		assert.Failf(suite.T(), "error list devices", err.Error())
	} else {
		assert.EqualValues(suite.T(), len(devices), 1)
	}

}

func (suite *IotDeviceConnectorTestSuite) TestPatchDevice() {
	deviceID := "my-test-device" + randStringRunes(4)
	connectorDevices := connectors.NewHttpIotConnector(suite.registryID)
	connectorDevices.CreateDevice(deviceID)
	deviceToUpdate := &cloudiot.Device{
		Metadata: map[string]string{"key": "value"},
	}

	deviceUpdated, err := connectorDevices.PatchDevice(deviceID, deviceToUpdate, "Metadata")
	if err != nil {
		assert.Failf(suite.T(), "error path device", err.Error())
	} else {
		assert.EqualValues(suite.T(), deviceUpdated.Metadata["key"], "value")
	}

}

type IotDeviceConnectorTestSuite struct {
	suite.Suite
	configuration *configuration.Configuration
	registryID    string
}

func (suite *IotDeviceConnectorTestSuite) SetupTest() {
	connector := connectors.NewIotRegistryConnector(connectors.HTTP, suite.configuration.GcloudProjectID, suite.configuration.GcloudRegion)

	eventNotificationConfigs := []*cloudiot.EventNotificationConfig{
		{
			PubsubTopicName: connector.GenerateTopicName("pablo-test"),
		},
	}

	_, err := connector.CreateRegistry(suite.registryID, eventNotificationConfigs)

	if err != nil {
		assert.Failf(suite.T(), "error when trying to create a register", err.Error())
	}
}

func (suite *IotDeviceConnectorTestSuite) TearDownTest() {
	connector := connectors.NewIotRegistryConnector(connectors.HTTP, suite.configuration.GcloudProjectID, suite.configuration.GcloudRegion)
	connectorDevices := connectors.NewHttpIotConnector(suite.registryID)
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
*/
