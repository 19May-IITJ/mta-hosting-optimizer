package config

import (
	"encoding/json"
	"log"
	"mta2/internal/constants"
	"mta2/pkg/ipconfig"
	"os"
	"path/filepath"
	"strconv"
)

func LoadConfigThreshold() int {
	defaultThreshold := 1

	threshold := os.Getenv(constants.MTA_THRESHOLD)
	x, err := strconv.Atoi(threshold)
	if err != nil || x <= 0 {
		x = defaultThreshold
	}

	return x
}
func LoadConfigIPConfiguration(c ipconfig.Configuration, ips ipconfig.IPListI) (err error) {
	var (
		jsonFile  *os.File
		byteValue []byte
		absPath   string
	)

	path := os.Getenv(constants.DBPATH)
	if path == "" {
		path = constants.DEFAULTPATH
	}
	if absPath, err = filepath.Abs(path); err != nil {
		log.Printf("Error %v\n", err)
		return err
	}
	if jsonFile, err = os.Open(absPath); err != nil {
		defer jsonFile.Close()
		log.Printf("Error %v\n", err)
		return err
	}
	defer jsonFile.Close()
	if byteValue, err = os.ReadFile(jsonFile.Name()); err != nil {
		log.Printf("Error %v\n", err)
		return err
	}
	list := ips.GetIPList().IpConfigList
	if err = decode(byteValue, &list); err == nil {
		ips.SetIPList(list)
		for _, ip := range ips.GetIPList().IpConfigList {
			if c.Contains(ip.Hostname) {
				if ip.Status {
					val, _ := c.GetValue(ip.Hostname)
					val++
					c.Put(ip.Hostname, val)
				}
			} else {
				if ip.Status {
					c.Put(ip.Hostname, 1)
				}
			}

		}
		return nil
	} else {
		log.Printf("Error %v\n", err)
		return err
	}

}
func decode(byteValue []byte, list *[]*ipconfig.IPConfig) error {
	return json.Unmarshal(byteValue, &list)
}
