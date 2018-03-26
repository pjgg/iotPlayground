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

// MQTTIotDeviceConnector handler devices telemetry communication.
type MQTTIotDeviceConnector struct {
	MQTTClient     mqtt.Client
	publicKeyPath  string
	privateKeyPath string
	projectID      string
	region         string
	keyType        KeyType
	registryID     string
}

// MQTTIotDeviceConnectorInterface define device telemetry behavior.
type MQTTIotDeviceConnectorInterface interface {
	PublishMsg(toDeviceID, topicName, msg string, delivery QoS) mqtt.Token
}

var onceMqttDevice sync.Once
var mqttIotDeviceConnector MQTTIotDeviceConnector

const mqttRetries = 5
const mqttDelaySecond = 5

// NewMQTTIotConnector create a single MQTTIotDeviceConnector instance.
func NewMQTTIotConnector(registryID, MQTTdeviceID string) MQTTIotDeviceConnectorInterface {

	onceMqttDevice.Do(func() {
		conf := configuration.New()
		mqttIotDeviceConnector.registryID = registryID
		mqttIotDeviceConnector.publicKeyPath = conf.DevicePublicKeyPath
		mqttIotDeviceConnector.privateKeyPath = conf.DevicePrivateKeyPath
		mqttIotDeviceConnector.keyType = rsaPem
		mqttIotDeviceConnector.projectID = conf.GcloudProjectID
		mqttIotDeviceConnector.region = conf.GcloudRegion
		jwt, _ := GenerateJWT(conf.GcloudProjectID, conf.DevicePrivateKeyPath, conf.DeviceJwtExpirationInMin)
		opts := paho.NewClientOptions()

		opts.SetClientID("projects/" + conf.GcloudProjectID + "/locations/" + conf.GcloudRegion + "/registries/" + registryID + "/devices/" + MQTTdeviceID).
			AddBroker(conf.MqttEndpoint).
			SetUsername("unused").
			SetTLSConfig(&tls.Config{MinVersion: tls.VersionTLS12}).
			SetPassword(jwt).
			SetProtocolVersion(4) // Use MQTT 3.1.1

		opts.CleanSession = true
		opts.AutoReconnect = true

		log.Info("ClientID: " + opts.ClientID)
		mqttIotDeviceConnector.MQTTClient = paho.NewClient(opts)
		mqttIotDeviceConnector.mqttConnect(mqttRetries, mqttDelaySecond)

	})

	return &mqttIotDeviceConnector
}

// PublishMsg push a mqtt message to google mqtt broker. Thids message will be propagated to a pub/sub topic.
func (iotConnector *MQTTIotDeviceConnector) PublishMsg(toDeviceID, topicName, msg string, delivery QoS) (token mqtt.Token) {

	finalTopicName := fmt.Sprintf("/devices/%s/%s", toDeviceID, topicName)
	log.Info("Publish Msg to topic " + finalTopicName)
	if !iotConnector.MQTTClient.IsConnected() {
		log.Info("Client Not Connected. Reconnecting... ")
		mqttIotDeviceConnector.mqttConnect(mqttRetries, mqttDelaySecond)
	}

	token = iotConnector.MQTTClient.Publish(finalTopicName, delivery.Value(), false, msg)
	if token.Wait() && token.Error() != nil {
		log.Errorln("MQTT Publish telemetric fail:")
		panic(token.Error())
	}

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Disconecting Mqtt client ...", r)
			iotConnector.MQTTClient.Disconnect(1)
		}
	}()

	return
}

func (iotConnector *MQTTIotDeviceConnector) mqttConnect(retriesAmount, elapsed int) (success bool) {
	success = true
	if token := iotConnector.MQTTClient.Connect(); token.Wait() && token.Error() != nil {
		success = false
		log.Errorln("MQTT Unable to connect:")
		panic(token.Error())
	}

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Disconecting Mqtt client ...", r)
			iotConnector.MQTTClient.Disconnect(1)
			if retriesAmount > 0 {
				fmt.Println("Retrying Mqtt connection ...")
				iotConnector.connectionRetry(retriesAmount, elapsed)
			}
		}
	}()
	return
}

func (iotConnector *MQTTIotDeviceConnector) connectionRetry(retriesAmount, elapsed int) {
	retriesAmount--
	for range time.Tick(time.Duration(elapsed) * time.Second) {
		iotConnector.mqttConnect(retriesAmount, elapsed)
	}
}
