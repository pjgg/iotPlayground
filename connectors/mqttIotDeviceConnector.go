package connectors

import (
	"crypto/tls"
	"fmt"
	"sync"
	"time"

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
	PublishMsg(toDeviceID, topicName, msg string, delivery QoS) mqtt.Token
}

var onceMqttDevice sync.Once
var mqttIotDeviceConnector MqttIotDeviceConnector

const MQTT_RETRIES = 5
const MQTT_DELAY_SECOND = 5

func NewMqttIotConnector(registryID, MQTT_deviceID string) MqttIotDeviceConnectorInterface {

	onceMqttDevice.Do(func() {
		conf := configuration.New()
		mqttIotDeviceConnector.registryID = registryID
		mqttIotDeviceConnector.publicKeyPath = conf.DevicePublicKeyPath
		mqttIotDeviceConnector.privateKeyPath = conf.DevicePrivateKeyPath
		mqttIotDeviceConnector.keyType = RSA_PEM
		mqttIotDeviceConnector.projectID = conf.GcloudProjectID
		mqttIotDeviceConnector.region = conf.GcloudRegion
		jwt, _ := GenerateJWT(conf.GcloudProjectID, conf.DevicePrivateKeyPath, conf.DeviceJwtExpirationInMin)
		opts := paho.NewClientOptions()

		opts.SetClientID("projects/" + conf.GcloudProjectID + "/locations/" + conf.GcloudRegion + "/registries/" + registryID + "/devices/" + MQTT_deviceID).
			AddBroker(conf.MqttEndpoint).
			SetUsername("unused").
			SetTLSConfig(&tls.Config{MinVersion: tls.VersionTLS12}).
			SetPassword(jwt).
			SetProtocolVersion(4) // Use MQTT 3.1.1

		opts.CleanSession = true
		opts.AutoReconnect = true

		log.Info("ClientID: " + opts.ClientID)
		mqttIotDeviceConnector.MQTT_Client = paho.NewClient(opts)
		mqttIotDeviceConnector.mqttConnect(MQTT_RETRIES, MQTT_DELAY_SECOND)

	})

	return &mqttIotDeviceConnector
}

func (iotConnector *MqttIotDeviceConnector) PublishMsg(toDeviceID, topicName, msg string, delivery QoS) (token mqtt.Token) {

	finalTopicName := fmt.Sprintf("/devices/%s/%s", toDeviceID, topicName)
	log.Info("Publish Msg to topic " + finalTopicName)
	if !iotConnector.MQTT_Client.IsConnected() {
		log.Info("Client Not Connected. Reconnecting... ")
		mqttIotDeviceConnector.mqttConnect(MQTT_RETRIES, MQTT_DELAY_SECOND)
	}

	token = iotConnector.MQTT_Client.Publish(finalTopicName, delivery.Value(), false, msg)
	if token.Wait() && token.Error() != nil {
		log.Errorln("MQTT Publish telemetric fail:")
		panic(token.Error())
	}

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Disconecting Mqtt client ...", r)
			iotConnector.MQTT_Client.Disconnect(1)
		}
	}()

	return
}

func (iotConnector *MqttIotDeviceConnector) mqttConnect(retriesAmount, elapsed int) (success bool) {
	success = true
	if token := iotConnector.MQTT_Client.Connect(); token.Wait() && token.Error() != nil {
		success = false
		log.Errorln("MQTT Unable to connect:")
		panic(token.Error())
	}

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Disconecting Mqtt client ...", r)
			iotConnector.MQTT_Client.Disconnect(1)
			if retriesAmount > 0 {
				fmt.Println("Retrying Mqtt connection ...")
				iotConnector.connectionRetry(retriesAmount, elapsed)
			}
		}
	}()
	return
}

func (iotConnector *MqttIotDeviceConnector) connectionRetry(retriesAmount, elapsed int) {
	retriesAmount--
	for range time.Tick(time.Duration(elapsed) * time.Second) {
		iotConnector.mqttConnect(retriesAmount, elapsed)
	}
}
