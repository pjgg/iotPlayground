package connectors

import (
	"io/ioutil"
	"time"

	log "github.com/Sirupsen/logrus"
	jwt "github.com/dgrijalva/jwt-go"
)

// Protocol must be HTTP or MQTT
type Protocol int

const (
	// HTTP ...
	HTTP Protocol = 1 + iota
	// MQTT ...
	MQTT
)

var protocolName = [...]string{
	"HTTP",
	"MQTT",
}

func (protocol Protocol) String() string {
	return protocolName[protocol-1]
}

// KeyType must be RSA_PEM or ES256_PEM
type KeyType int

const (
	// RSA_PEM ...
	rsaPem KeyType = 1 + iota
	// ES256_PEM ...
	es256Pem
)

var keyTypeName = [...]string{
	"RSA_PEM",
	"ES256_PEM",
}

func (keyType KeyType) String() string {
	return keyTypeName[keyType-1]
}

// GenerateJWT will generate a signed JWT token
func GenerateJWT(projectID, privateKeyFullPath string, expireTimeMin int) (string, error) {
	privateKeyBytes, err := ioutil.ReadFile(privateKeyFullPath)
	if err != nil {
		log.Errorln(err.Error())
		return "", err
	}
	privateKey, err := jwt.ParseRSAPrivateKeyFromPEM(privateKeyBytes)
	if err != nil {
		log.Errorln(err.Error())
		return "", err
	}

	t := time.Now()
	token := jwt.NewWithClaims(jwt.GetSigningMethod(jwt.SigningMethodRS256.Alg()), &jwt.StandardClaims{
		IssuedAt:  t.Unix(),
		ExpiresAt: t.Add(time.Minute * time.Duration(expireTimeMin)).Unix(),
		Audience:  projectID,
	})
	pass, err := token.SignedString(privateKey)

	if err != nil {
		log.Errorln(err.Error())
		return "", err
	}

	return pass, nil
}

// None type represent a empty type.
type None int

const (
	// NONE ...
	NONE None = 1 + iota
)

var noneName = [...]string{
	"NONE",
}

func (none None) String() string {
	return noneName[none-1]
}

// QoS type represent the message consumption type.
type QoS int

const (
	// AtMostOnce ...
	AtMostOnce QoS = 1 + iota
	// AtLeastOnce ...
	AtLeastOnce
)

var qosValue = [...]byte{
	0,
	1,
}

// Value returns consumption type value.
func (qos QoS) Value() byte {
	return qosValue[qos-1]
}
