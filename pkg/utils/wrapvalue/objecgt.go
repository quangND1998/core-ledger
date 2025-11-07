package wv

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
)

func ToObject(src, des any) error {
	bytes, err := json.Marshal(src)
	if err != nil {
		return err
	}
	err = json.Unmarshal(bytes, des)
	if err != nil {
		logrus.Info("parse object failed")
		return err
	}
	return nil
}
