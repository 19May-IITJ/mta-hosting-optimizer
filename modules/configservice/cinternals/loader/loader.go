package loader

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"mta2/modules/configservice/cinternals/constants"
	"mta2/modules/configservice/cpkg/ipconfig"
	"mta2/modules/natsmodule"
	"mta2/modules/utility"
	"strings"
	"time"

	"os"
	"path/filepath"
	"sort"
)

var (
	Ticker     *time.Ticker
	FLAGTOSAVE bool
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
		log.Printf("Error Abs %v\n", err)
		return err
	}
	if jsonFile, err = os.Open(absPath); err != nil {
		defer jsonFile.Close()
		log.Printf("Error Open %v\n", err)
		return err
	}
	defer jsonFile.Close()
	if byteValue, err = os.ReadFile(jsonFile.Name()); err != nil {
		log.Printf("Error ReadFile %v\n", err)
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
		log.Printf("Error Unmarshal %v\n", err)
		return err
	}

}

// wrapper over Unmarshal JSON
func decode(byteValue []byte, list *[]*ipconfig.IPConfigData) error {
	return json.NewDecoder(bytes.NewReader(byteValue)).Decode(list)
}

// Binary Search wrapper for slices
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

// A TTL to save the data to the file, after 30sec of serving latest request of Refresh Data
func TTLForFileSaving(ctx context.Context, ipl ipconfig.IPListI, nc natsmodule.NATSConnInterface) {
	log.Println("Started TTL handler to save in DB")
	path := os.Getenv(constants.DBPATH)
	if path == "" {
		path = constants.DEFAULTPATH
	}
	for {
		select {
		case <-ctx.Done():
			Ticker.Stop()
			log.Println("Got response to shutdown TTL")
			constants.Datamutex.Lock()

			// Marshal the entire updated data
			updatedData, err := json.MarshalIndent(ipl.GetIPValues(), "", "  ")
			if err != nil {
				fmt.Println("Error marshaling JSON:", err)
				return
			}

			// Write back the entire JSON file with the updated entry
			err = os.WriteFile(path, updatedData, 0644)
			if err != nil {
				fmt.Println("Error writing file:", err)
				return
			}
			FLAGTOSAVE = false
			serviceDown := "Roll Back"

			bye, _ := json.Marshal(serviceDown)
			if err := nc.Publish(constants.CONFIGSERVICE_PUB_SUBJECT, bye); err != nil {
				log.Println("Error in PUB ", err)
			}
			constants.Datamutex.Unlock()
			log.Println("JSON entry updated successfully!")
			log.Println("Send Signal to hosting for Roll Back")
			utility.TaskChan <- serviceDown
			return
		case <-Ticker.C:
			log.Println("ticker hit")
			if FLAGTOSAVE {
				constants.Datamutex.Lock()

				// Marshal the entire updated data
				updatedData, err := json.MarshalIndent(ipl.GetIPValues(), "", "  ")
				if err != nil {
					fmt.Println("Error marshaling JSON:", err)
					return
				}

				// Write back the entire JSON file with the updated entry
				err = os.WriteFile(path, updatedData, 0644)
				if err != nil {
					fmt.Println("Error writing file:", err)
					return
				}
				FLAGTOSAVE = false
				constants.Datamutex.Unlock()
				fmt.Println("JSON entry updated successfully!")
			}
		}
	}
}
