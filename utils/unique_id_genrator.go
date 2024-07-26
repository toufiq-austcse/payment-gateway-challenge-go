package utils

import (
	"github.com/google/uuid"
	"strings"
)

func GenerateUUID() (string, error) {
	uuidWithHyphen, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}
	newUuId := strings.Replace(uuidWithHyphen.String(), "-", "", -1)
	return newUuId, nil
}
