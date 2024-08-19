package analysis

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/fs"
	"math"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strconv"
	"time"

	"github.com/ag-hn/speed-analysis/filesystem"
)

type AddrData struct {
	//** Ip Address of RSU */
	Ip string `json:"ip"`
	//** Longitude */
	Lng float64 `json:"lon"`
	//** Latitude */
	Lat float64 `json:"lat"`
	//** Mac address */
	Addr  string `json:"addr"`
	Flags int    `json:"flags"`
	//** Strength of signal */
	Rssi int `json:"rssi"`
	//** Distance from RSU */
	Seen int `json:"seen"`
	//** Timestamp of entry */
	Time string `json:"time"`
	//** Capture of entry (e.g., 'wifi', 'bluetooth', etc.) */
	CaptureType string `json:"type"`
}

type AddrResponse []AddrData

const ADDRESS_DATA_DIRECTORY = "./__input/addr-data"
const OUTPUT_DATA_FILE = "./__output/"

// const DURATION_MINUTE_DELAY_THRESHOLD = 10
const DURATION_MINUTE_DELAY_THRESHOLD = 10.0
// In minutes
const MINIMUM_DURATION_LOGGING_THRESHOLD = 0.1
// In miles
const MINIMUM_SPEED_LOGGING_THRESHOLD = 0.0001
// Minimum number on records which consider a valid item
const MINIMUM_PROCESSING_COUNT_LOGGING_THRESHOLD = 2

func ListProcessFilePaths() (paths []fs.DirEntry, err error) {
	var files []fs.DirEntry
	var directoryPath string

	if !filepath.IsAbs(ADDRESS_DATA_DIRECTORY) {
		directoryPath, err = filepath.Abs(ADDRESS_DATA_DIRECTORY)
		if err != nil {
			return nil, err
		}
	} else {
		directoryPath = ADDRESS_DATA_DIRECTORY
	}

	directoryInfo, err := os.Stat(directoryPath)
	if err != nil {
		return nil, err
	}

	if !directoryInfo.IsDir() {
		return nil, errors.New("ListProcessFilePaths+Given file is not a directory")
	}

	files, err = filesystem.GetDirectoryListingByType(directoryPath, filesystem.FilesListingType, true)
	if err != nil {
		return nil, err
	}

	return files, nil
}

func ProcessFilePath(file fs.DirEntry) (item []ProcessedItem, err error) {
	fileInfo, err := file.Info()
	if err != nil {
		return []ProcessedItem{}, err
	}

	isSymlink := fileInfo.Mode()&os.ModeSymlink != 0

	if isSymlink || file.IsDir() {
		return []ProcessedItem{}, errors.New("ProcessFilePath+Can only process files. Given Symlink or Directories")
	}

	filePath := filepath.Join(ADDRESS_DATA_DIRECTORY, file.Name())
	potentialAddressData, err := os.ReadFile(filePath)
	if err != nil {
		return []ProcessedItem{}, err
	}

	var data AddrResponse
	err = json.Unmarshal(potentialAddressData, &data)
	if err != nil {
		return []ProcessedItem{}, err
	}

	return processResponse(data), nil
}

func processResponse(res AddrResponse) []ProcessedItem {
	sorted := res[:]
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Time < sorted[j].Time
	})
	processByIp := []ProcessedItem{}
	output := ""

	workingAddresses := sorted

	var previous AddrData
	var last AddrData
	distance := 0.0
    count := 0

	var previousTime time.Time
	var minTime time.Time
	var maxTime time.Time
	for _, addr := range workingAddresses {
		if previous == (AddrData{}) {
			previous = addr
            last = previous
			t, err := strconv.Atoi(previous.Time)
			if err != nil {
				panic("ProcessFilePath+Cannot parse time property: " + previous.Addr)
			}
			tt := time.UnixMilli(int64(t))
			output += fmt.Sprintf("\nStart new processing - %s,%s: rssi %d | Time %s | Lat %f | Lng %f\n", addr.Ip, addr.Addr, addr.Rssi, tt.String(), addr.Lat, addr.Lng)
			minTime = tt
			maxTime = tt

			distance = 0
            count = 0

			continue

		}

		previousTimeInt, err := strconv.Atoi(previous.Time)
		if err != nil {
			panic("ProcessFilePath+Cannot parse time property: " + previous.Addr)
		}
		previousTime = time.UnixMilli(int64(previousTimeInt))

		timeAsInt, err := strconv.Atoi(addr.Time)
		if err != nil {
			panic("ProcessFilePath+Cannot parse time property: " + previous.Addr)
		}

		currentTime := time.UnixMilli(int64(timeAsInt))

		if previousTime.Sub(currentTime).Abs().Minutes() > DURATION_MINUTE_DELAY_THRESHOLD {
            last = previous
			previous = AddrData{}

            appendProcessedItem(distance, maxTime, minTime, addr.Ip, addr.Addr, addr.Rssi, count, &output, &processByIp)

			continue
		}

		if currentTime.Unix() > maxTime.Unix() {
			maxTime = currentTime
		}
		if currentTime.Unix() < minTime.Unix() {
			minTime = currentTime
		}

		cX := addr.Lat
		cY := addr.Lng
		distance += getDistance(previous.Lat, previous.Lng, addr.Lat, addr.Lng)

        count += 1
		previous = addr
		output += fmt.Sprintf("\n%s,%s: Distance %f | rssi %d | prev %s | curr %s  | Lat %f | Lng % f\n", addr.Ip, addr.Addr, distance, addr.Rssi, currentTime.String(), previousTime.String(), cX, cY)
	}

    appendProcessedItem(distance, maxTime, minTime, last.Ip, last.Addr, last.Rssi, count, &output, &processByIp)

	newFile := path.Join(OUTPUT_DATA_FILE, last.Addr)
	filesystem.WriteToFile(newFile, output)

	return processByIp
}

// Distance approximation using the Haversine formula (https://en.wikipedia.org/wiki/Haversine_formula)
// lat1, lng1, lat2, lng2 are expected as `deg`
// Return value is given in miles.
func getDistance(lat1 float64, lng1 float64, lat2 float64, lng2 float64) float64 {
	R := 6371000.0 // Earth Radius in meters
	degToRad := math.Pi / 180
	distInMeters := R * degToRad * math.Sqrt(math.Pow(math.Cos(lat1*degToRad)*(lng1-lng2), 2)+math.Pow(lat1-lat2, 2))
	distInMiles := distInMeters / 1609.34
	return distInMiles
}

// *__Mutates `output` and `list`
func appendProcessedItem(distance float64, maxTime time.Time, minTime time.Time, ip string, addr string, rssi int, count int, output *string, list *[]ProcessedItem) bool {
	duration := maxTime.Sub(minTime)
    formattedDuration := duration.Hours()
	speed := distance / formattedDuration

    if count >= MINIMUM_PROCESSING_COUNT_LOGGING_THRESHOLD && speed >= MINIMUM_SPEED_LOGGING_THRESHOLD && duration.Minutes() >= MINIMUM_DURATION_LOGGING_THRESHOLD {
        *output += fmt.Sprintf("\n-----\nduration: %fh or %fmin\nrssi: %d\n(min,max): (%s,%s)\ndistance: %f\nspeed: %f\n\n", formattedDuration, duration.Minutes(), rssi, minTime.String(), maxTime.String(), distance, speed)
        item := ProcessedItem{
            Name:  ip,
            Addr:  addr,
            Speed: fmt.Sprintf("%f", speed),
            Lat:   "",
            Lng:   "",
        }
        *list = append(*list, item)
    }

    return true;
}
