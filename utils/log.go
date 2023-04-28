package utils

import (
	"os"

	"github.com/sirupsen/logrus"
)

func CreateLogFile() (*os.File, error) {
<<<<<<< HEAD
	logFile, err := os.OpenFile("C:/Users/LENOVO/Documents/Final_project/inc-final-project/api.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
=======
	logFile, err := os.OpenFile(DotEnv("LOG_LOCATION"), os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
>>>>>>> 6d181691126b8bc5f2a3876fdf834723f8441d05
	if err != nil {
		return nil, err
	}

	logrus.SetOutput(logFile)

	return logFile, nil
}
