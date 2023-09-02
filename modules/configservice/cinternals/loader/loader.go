package loader

import (
	"encoding/json"
	"log"
	"mta2/modules/configservice/cinternals/constants"
	"mta2/modules/configservice/cpkg/ipconfig"
	"strings"

	"os"
	"path/filepath"
	"sort"
)

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
			if !c.Contains(ip.Hostname) {
				c.Put(ip.Hostname, &ipconfig.HostData{
					ActiveIP: 0,
				})
			}
			hd, _ := c.GetValue(ip.Hostname)
			if ip.Status {

				hd.HostedIP = append(hd.HostedIP,
					strings.Join([]string{ip.IPAddresses,
						constants.Active}, constants.Sep))
			} else {

				hd.HostedIP = append(hd.HostedIP,
					strings.Join([]string{ip.IPAddresses,
						constants.Inactive}, constants.Sep))
			}

		}
		for _, ip := range ips.GetIPValues() {

			if ip.Status {
				hostdata, _ := c.GetValue(ip.Hostname)
				hostdata.ActiveIP++
				// c.Put(ip.Hostname, hostdata)
			}
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
			return mid
		} else if midIP < targetIP {
			left = mid + 1
		} else {
			right = mid - 1
		}
	}
	return -1
}