package main

import (
	"fmt"
	"io"
	"os"

	"github.com/NatsuiroGinga/mydocker/container"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
)

func logContainer(containerId string) {
	logFileLocation := fmt.Sprintf(container.InfoLocFormat, containerId) + container.GetLogfile(containerId)
	file, err := os.Open(logFileLocation)
	if err != nil {
		logrus.Errorf("log container open file %s error %v", logFileLocation, err)
		return
	}
	content, err := io.ReadAll(file)
	if err != nil {
		logrus.Errorf("Log container read file %s error %v", logFileLocation, err)
		return
	}

	_, err = fmt.Fprint(os.Stdout, string(content))
	if err != nil {
		log.Errorf("Log container Fprint  error %v", err)
		return
	}
}
