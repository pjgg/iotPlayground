package configuration

import (
	"os"
	"sync"

	log "github.com/Sirupsen/logrus"
	"github.com/spf13/viper"
)

type Configuration struct {
	GcloudProjectID      string
	GcloudRegion         string
	DevicePublicKeyPath  string
	DevicePrivateKeyPath string
	MqttEndpoint         string
}

var onceConfiguration sync.Once
var ConfigurationInstance *Configuration

func New() *Configuration {
	onceConfiguration.Do(func() {
		ConfigurationInstance = &Configuration{}

		if len(os.Getenv("GCLOUD_PROJECT")) == 0 {
			log.Fatalln("missing required ENV GCLOUD_PROJECT")
		}

		if len(os.Getenv("GOOGLE_APPLICATION_CREDENTIALS")) == 0 {
			log.Fatalln("missing required ENV GOOGLE_APPLICATION_CREDENTIALS")
		}

		ConfigurationInstance.GcloudProjectID = viper.GetString("gcloud.projectID")
		ConfigurationInstance.GcloudRegion = viper.GetString("gcloud.region")
		ConfigurationInstance.DevicePublicKeyPath = viper.GetString("device.publicKeyPath")
		ConfigurationInstance.DevicePrivateKeyPath = viper.GetString("device.privateKeyPath")
		ConfigurationInstance.MqttEndpoint = viper.GetString("gcloud.mqtt")

		log.WithFields(log.Fields{
			"GcloudProjectID":      ConfigurationInstance.GcloudProjectID,
			"GcloudRegion":         ConfigurationInstance.GcloudRegion,
			"DevicePublicKeyPath":  ConfigurationInstance.DevicePublicKeyPath,
			"DevicePrivateKeyPath": ConfigurationInstance.DevicePrivateKeyPath,
			"MqttEndpoint":         ConfigurationInstance.MqttEndpoint,
		}).Info("configuration loaded")
	})

	return ConfigurationInstance
}
