package device

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"os"
	"sync"

	log "github.com/Sirupsen/logrus"
	"github.com/pjgg/iotPlayground/configuration"
	"github.com/pjgg/iotPlayground/connectors"
	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"
	cloudiot "google.golang.org/api/cloudiot/v1"
)

// HTTPIotDeviceConnector handler devices admin request.
type HTTPIotDeviceConnector struct {
	HTTPClient     *cloudiot.Service
	publicKeyPath  string
	privateKeyPath string
	projectID      string
	region         string
	keyType        connectors.KeyType
	registryID     string
}

// HTTPIotDeviceConnectorInterface define device IoT admin request behavior.
type HTTPIotDeviceConnectorInterface interface {
	SwapToRegistry(registryID string)
	CreateDevice(deviceID string) (*cloudiot.Device, error)
	DeleteDevice(deviceID string) (*cloudiot.Empty, error)
	GetDevice(deviceID string) (*cloudiot.Device, error)
	SetDeviceConfig(deviceID string, configData string) (*cloudiot.DeviceConfig, error)
	GetDeviceConfigs(deviceID string) ([]*cloudiot.DeviceConfig, error)
	GetDeviceStates(deviceID string) ([]*cloudiot.DeviceState, error)
	ListDevices() ([]*cloudiot.Device, error)
	PatchDevice(deviceID string, newDevice *cloudiot.Device, field string) (*cloudiot.Device, error)
}

var onceHTTPDevice sync.Once
var httpIotDeviceConnector HTTPIotDeviceConnector

// NewDeviceHTTPIotConnector create a single instance of HTTPIotDeviceConnector
func NewDeviceHTTPIotConnector(registryID string) HTTPIotDeviceConnectorInterface {

	onceHTTPDevice.Do(func() {
		conf := configuration.New()
		ctx := context.Background()

		httpClient, err := google.DefaultClient(ctx, cloudiot.CloudPlatformScope)
		if err != nil {
			log.Fatalln(err.Error())
		}

		httpIotDeviceConnector.HTTPClient, err = cloudiot.New(httpClient)
		if err != nil {
			log.Fatalln(err.Error())
		}

		httpIotDeviceConnector.registryID = registryID
		httpIotDeviceConnector.publicKeyPath = conf.DevicePublicKeyPath
		httpIotDeviceConnector.privateKeyPath = conf.DevicePrivateKeyPath
		httpIotDeviceConnector.keyType = connectors.RsaPem
		httpIotDeviceConnector.projectID = conf.GcloudProjectID
		httpIotDeviceConnector.region = conf.GcloudRegion

	})

	return &httpIotDeviceConnector
}

// SwapToRegistry overwrite HTTP otDeviceConnector local device registry ID, so all device request will be thrown against this registryID
func (iotConnector *HTTPIotDeviceConnector) SwapToRegistry(registryID string) {
	iotConnector.registryID = registryID
}

// CreateDevice will create a device over a previous given registryID.
func (iotConnector *HTTPIotDeviceConnector) CreateDevice(deviceID string) (device *cloudiot.Device, err error) {

	keyBytes, err := ioutil.ReadFile(iotConnector.publicKeyPath)
	if err != nil {
		log.Error(err.Error())
	}

	deviceDef := cloudiot.Device{
		Id: deviceID,
		Credentials: []*cloudiot.DeviceCredential{
			{
				PublicKey: &cloudiot.PublicKeyCredential{
					Format: iotConnector.keyType.String(),
					Key:    string(keyBytes),
				},
			},
		},
	}

	parent := fmt.Sprintf("projects/%s/locations/%s/registries/%s", iotConnector.projectID, iotConnector.region, iotConnector.registryID)
	if device, err = iotConnector.HTTPClient.Projects.Locations.Registries.Devices.Create(parent, &deviceDef).Do(); err == nil {
		log.Debugln("Successfully created device.")
		log.Debugln("\tID: ", device.Id)
		log.Debugln("\tName: ", device.Name)
	}

	return
}

// DeleteDevice will delete a device over a previous given registryID.
func (iotConnector *HTTPIotDeviceConnector) DeleteDevice(deviceID string) (response *cloudiot.Empty, err error) {
	path := fmt.Sprintf("projects/%s/locations/%s/registries/%s/devices/%s", iotConnector.projectID, iotConnector.region, iotConnector.registryID, deviceID)
	if response, err = iotConnector.HTTPClient.Projects.Locations.Registries.Devices.Delete(path).Do(); err == nil {
		log.Debugln("Deleted device!")
	}

	return
}

// GetDevice will retrieve a device, if is a member of a registryID.
func (iotConnector *HTTPIotDeviceConnector) GetDevice(deviceID string) (device *cloudiot.Device, err error) {
	path := fmt.Sprintf("projects/%s/locations/%s/registries/%s/devices/%s", iotConnector.projectID, iotConnector.region, iotConnector.registryID, deviceID)
	if device, err = iotConnector.HTTPClient.Projects.Locations.Registries.Devices.Get(path).Do(); err == nil {
		log.Debugln("\tId: ", device.Id)
		for _, credential := range device.Credentials {
			log.Debugln("\t\tCredential Expire: ", credential.ExpirationTime)
			log.Debugln("\t\tCredential Type: ", credential.PublicKey.Format)
			log.Debugln("\t\t--------")
		}
		log.Debugln("\tLast Config Ack: ", device.LastConfigAckTime)
		log.Debugln("\tLast Config Send: ", device.LastConfigSendTime)
		log.Debugln("\tLast Event Time: ", device.LastEventTime)
		log.Debugln("\tLast Heartbeat Time: ", device.LastHeartbeatTime)
		log.Debugln("\tLast State Time: ", device.LastStateTime)
		log.Debugln("\tNumId: ", device.NumId)

	}

	return
}

// GetDeviceConfigs will retrieve a device configuration, if is a member of a registryID.
func (iotConnector *HTTPIotDeviceConnector) GetDeviceConfigs(deviceID string) (configs []*cloudiot.DeviceConfig, err error) {
	path := fmt.Sprintf("projects/%s/locations/%s/registries/%s/devices/%s", iotConnector.projectID, iotConnector.region, iotConnector.registryID, deviceID)
	if response, err := iotConnector.HTTPClient.Projects.Locations.Registries.Devices.ConfigVersions.List(path).Do(); err == nil {
		log.Debugln("Successfully retrieved device config!")
		configs = response.DeviceConfigs
		for _, config := range response.DeviceConfigs {
			log.Debugln(config.Version, " : ", config.BinaryData)
		}
	}

	return
}

// GetDeviceStates will retrieve a device states, if is a member of a registryID.
func (iotConnector *HTTPIotDeviceConnector) GetDeviceStates(deviceID string) (states []*cloudiot.DeviceState, err error) {
	path := fmt.Sprintf("projects/%s/locations/%s/registries/%s/devices/%s", iotConnector.projectID, iotConnector.region, iotConnector.registryID, deviceID)
	if response, err := iotConnector.HTTPClient.Projects.Locations.Registries.Devices.States.List(path).Do(); err == nil {
		log.Debugln("Successfully retrieved device states!")
		states = response.DeviceStates
		for _, state := range response.DeviceStates {
			log.Debugln(state.UpdateTime, " : ", state.BinaryData)
		}
	}

	return
}

// ListDevices will retrieve a list of devices that are member of a registryID.
func (iotConnector *HTTPIotDeviceConnector) ListDevices() (devices []*cloudiot.Device, err error) {
	parent := fmt.Sprintf("projects/%s/locations/%s/registries/%s", iotConnector.projectID, iotConnector.region, iotConnector.registryID)
	if response, err := iotConnector.HTTPClient.Projects.Locations.Registries.Devices.List(parent).Do(); err == nil {
		log.Debugln("Successfully retrieved devices!")
		devices = response.Devices
		log.Debugln("Devices:")
		for _, device := range response.Devices {
			log.Debugln("\t", device.Id)
		}

	}
	return
}

// PatchDevice make a partial update over a device, if is a member of a registryID.
func (iotConnector *HTTPIotDeviceConnector) PatchDevice(deviceID string, newDevice *cloudiot.Device, field string) (device *cloudiot.Device, err error) {

	parent := fmt.Sprintf("projects/%s/locations/%s/registries/%s/devices/%s", iotConnector.projectID, iotConnector.region, iotConnector.registryID, deviceID)
	if device, err = iotConnector.HTTPClient.Projects.Locations.Registries.Devices.Patch(parent, newDevice).UpdateMask(field).Do(); err == nil {
		log.Debugln("Successfully patched device.")
	}

	return
}

// SetDeviceConfig will update and push to server a device configuration.
func (iotConnector *HTTPIotDeviceConnector) SetDeviceConfig(deviceID string, configData string) (deviceConfig *cloudiot.DeviceConfig, err error) {
	req := cloudiot.ModifyCloudToDeviceConfigRequest{
		BinaryData: base64.StdEncoding.EncodeToString([]byte(configData)),
	}

	path := fmt.Sprintf("projects/%s/locations/%s/registries/%s/devices/%s", iotConnector.projectID, iotConnector.region, iotConnector.registryID, deviceID)
	if deviceConfig, err = iotConnector.HTTPClient.Projects.Locations.Registries.Devices.ModifyCloudToDeviceConfig(path, &req).Do(); err == nil {
		fmt.Fprintf(os.Stdout, "Config set!\nVersion now: %d", deviceConfig.Version)
	}

	return
}
