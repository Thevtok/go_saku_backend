package utils

import (
	"os"

	"github.com/sirupsen/logrus"
)

func CreateLogFile() (*os.File, error) {
    logFile, err := os.OpenFile("/Developments/Golang/src/final-project-inc/log/api.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
    if err != nil {
        return nil, err
    }

    logrus.SetOutput(logFile)

    return logFile, nil
}