package main

import (
	"crypto/md5"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"reflect"
	"strings"

	log "github.com/sirupsen/logrus"
)

const (
	configFilesEnvVar = "SIDE_CAR_CONFIG_FILES"
)

// Server contains the configuration in memory
type Server struct {
	// hashes stores the hash of the file path [filePath]hash
	Hashes map[string]string
	Logger *log.Logger
}

func (s *Server) checkLiveness(w http.ResponseWriter, r *http.Request) {
	thashes, err := loadHashes()
	if err != nil {
		s.Logger.Error(err)
		w.WriteHeader(http.StatusInternalServerError)
		err = json.NewEncoder(w).Encode(Status{Message: "error checking config files", RestartRequired: true})
		if err != nil {
			s.Logger.Error(err)
		}
		return
	}
	equals := reflect.DeepEqual(thashes, s.Hashes)
	if !equals {
		s.Logger.Info("config files have changed")
		s.Logger.WithFields(log.Fields{"files": "cached"}).Info(s.Hashes)
		s.Logger.WithFields(log.Fields{"files": "current"}).Info(thashes)
		w.WriteHeader(http.StatusInternalServerError)
		err = json.NewEncoder(w).Encode(Status{Message: "config files have changed", RestartRequired: true})
		if err != nil {
			s.Logger.Error(err)
		}
		return
	}
	w.WriteHeader(http.StatusOK)
	// Ok, nothing to log here
	err = json.NewEncoder(w).Encode(Status{Message: "config files are ok", RestartRequired: false})
	if err != nil {
		s.Logger.Error(err)
	}
	return
}

// Status represent the status of the response
type Status struct {
	Message         string `json:"message"`
	RestartRequired bool   `json:"restartRequired"`
}

func main() {
	port := flag.String("p", "9090", "port for the service")
	flag.Parse()
	mux := http.NewServeMux()
	server := new(Server)
	Logger := log.New()
	hs, err := loadHashes()
	if err != nil {
		// this error is ignored and the flow continues to the listen and serve proces
		// the handler will response with the error
		Logger.Warn("something wrong loading hash config files", err.Error())
	}
	server.Hashes = hs
	server.Logger = Logger
	Logger.WithFields(log.Fields{"files": "cached"}).Info(server.Hashes)

	// routing
	mux.HandleFunc("/check/liveness", server.checkLiveness)

	Logger.Info("listen and serve port:", *port)
	err = http.ListenAndServe(":"+*port, mux)
	if err != nil {
		Logger.Error(err)
		return
	}
}

// loadHashes loads the hashes from the env var
func loadHashes() (map[string]string, error) {
	hashes := map[string]string{}
	stringPath := os.Getenv(configFilesEnvVar)
	if stringPath == "" {
		return hashes, errors.New("no paths to load in env var")
	}
	paths := strings.Split(stringPath, ",")

	if len(paths) > 0 {
		for _, path := range paths {
			hash, err := hash(path)
			if err != nil {
				return hashes, err
			}
			hashes[path] = hash
		}
	}
	return hashes, nil
}

// hash returns the hash of the content file otherwise returns an error
func hash(file string) (string, error) {
	f, err := os.Open(file)
	if err != nil {
		return "", err
	}
	defer f.Close()

	h := md5.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}
