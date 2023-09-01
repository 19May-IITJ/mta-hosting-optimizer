package loader

import (
	"encoding/json"
	"log"
	"mta2/internal/constants"
	"mta2/pkg/ipconfig"
	"os"
	"path/filepath"
	"sort"
	"strconv"
)

// Load Config Threshold loads the MTA_THRESHOLD env variable default:1
func LoadConfigThreshold() int {
	defaultThreshold := 1

	threshold := os.Getenv(constants.MTA_THRESHOLD)
	x, err := strconv.Atoi(threshold)
	if err != nil || x <= 0 {
		x = defaultThreshold
	}

	return x
}

// Load Config IPConfiguration loads the IP configuration Data based on DBPATH env variable
func LoadConfigIPConfiguration(c ipconfig.ConfigServiceIPMap, ips ipconfig.IPListI) (err error) {
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
	list := ips.GetIPValues()
	if err = decode(byteValue, &list); err == nil {
		sort.Slice(list, func(i, j int) bool {
			return list[i].IPAddresses < list[j].IPAddresses
		})
		ips.SetIPList(list)
		for _, ip := range ips.GetIPValues() {
			c.Put(ip.IPAddresses, &ipconfig.IPState{
				State:    ip.Status,
				Hostname: ip.Hostname,
			})
		}
		return nil
	} else {
		log.Printf("Error %v\n", err)
		return err
	}

}

// wrapper over Unmarshal JSON
func decode(byteValue []byte, list *[]*ipconfig.IPConfigData) error {
	return json.Unmarshal(byteValue, &list)
}

func Search(s []*ipconfig.IPConfigData, targetIP string) int {
	left, right := 0, len(s)-1
	for left <= right {
		mid := left + (right-left)/2
		midIP := s[mid].IPAddresses

		if midIP == targetIP {
			log.Printf("Found %s at index %d\n", targetIP, mid)
			return mid
		} else if midIP < targetIP {
			left = mid + 1
		} else {
			right = mid - 1
		}
	}
	return -1
}
