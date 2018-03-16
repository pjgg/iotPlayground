package connectors

import (
	"io/ioutil"
	"time"

	log "github.com/Sirupsen/logrus"
	jwt "github.com/dgrijalva/jwt-go"
)

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

type KeyType int

const (
	// RSA_X509_PEM ...
	RSA_PEM KeyType = 1 + iota
	// ES256_PEM ...
	ES256_PEM
)

var keyTypeName = [...]string{
	"RSA_PEM",
	"ES256_PEM",
}

func (keyType KeyType) String() string {
	return keyTypeName[keyType-1]
}

func GenerateJWT(projectId, privateKeyFullPath string, expireTimeMin int) (string, error) {
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
		Audience:  projectId,
	})
	pass, err := token.SignedString(privateKey)

	/*jwtToken := jwt.New(jwt.SigningMethodRS256)
	claims := jwtToken.Claims.(jwt.MapClaims)
	if err != nil {
		log.Errorln(err.Error())
		return nil, err
	}
	claims["iat"] = time.Now().Unix()
	claims["exp"] = time.Now().Add(time.Minute * time.Duration(expireTimeMin)).Unix()
	claims["aud"] = projectId
	jwtTokenString, err := jwtToken.SignedString(privateKey)*/
	if err != nil {
		log.Errorln(err.Error())
		return "", err
	}

	return pass, nil
}

type None int

const (
	NONE None = 1 + iota
)

var noneName = [...]string{
	"NONE",
}

func (none None) String() string {
	return noneName[none-1]
}
