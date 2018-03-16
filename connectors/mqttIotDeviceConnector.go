package connectors

import (
	"crypto/tls"
	"sync"

	"github.com/eclipse/paho.mqtt.golang"

	log "github.com/Sirupsen/logrus"
	paho "github.com/eclipse/paho.mqtt.golang"
	"github.com/pjgg/iotPlayground/configuration"
)

type MqttIotDeviceConnector struct {
	MQTT_Client    mqtt.Client
	publicKeyPath  string
	privateKeyPath string
	projectID      string
	region         string
	keyType        KeyType
	registryID     string
}

type MqttIotDeviceConnectorInterface interface {
	PublishMsg(toDeviceID, topicName, msg string) mqtt.Token
}

var onceMqttDevice sync.Once
var mqttIotDeviceConnector MqttIotDeviceConnector

func NewMqttIotConnector(registryID, MQTT_deviceID string) MqttIotDeviceConnectorInterface {

	onceMqttDevice.Do(func() {
		conf := configuration.New()
		//ctx := context.Background()

		mqttIotDeviceConnector.registryID = registryID
		mqttIotDeviceConnector.publicKeyPath = conf.DevicePublicKeyPath
		mqttIotDeviceConnector.privateKeyPath = conf.DevicePrivateKeyPath
		mqttIotDeviceConnector.keyType = RSA_PEM
		mqttIotDeviceConnector.projectID = conf.GcloudProjectID
		mqttIotDeviceConnector.region = conf.GcloudRegion
		jwt, _ := GenerateJWT(conf.GcloudProjectID, conf.DevicePrivateKeyPath, 60)
		opts := paho.NewClientOptions()

		opts.SetClientID("projects/" + conf.GcloudProjectID + "/locations/" + conf.GcloudRegion + "/registries/" + registryID + "/devices/" + MQTT_deviceID).
			AddBroker(conf.MqttEndpoint).
			SetUsername("unused").
			SetTLSConfig(&tls.Config{MinVersion: tls.VersionTLS12}).
			SetPassword(jwt).
			SetProtocolVersion(4) // Use MQTT 3.1.1

		opts.CleanSession = true
		opts.AutoReconnect = true

		cli := paho.NewClient(opts)
		if token := cli.Connect(); token.Wait() && token.Error() != nil {
			// Unable to connect to the MQTT broker.
			log.Errorln("MQTT Unable to connect:")
			log.Fatalln(token.Error())
		}
		mqttIotDeviceConnector.MQTT_Client = cli

	})

	return &mqttIotDeviceConnector
}

func (iotConnector *MqttIotDeviceConnector) PublishMsg(toDeviceID, topicName, msg string) (token mqtt.Token) {
	token = iotConnector.MQTT_Client.Publish("/devices/"+toDeviceID+"/"+topicName, 0, false, msg)
	token.Wait()
	return
}