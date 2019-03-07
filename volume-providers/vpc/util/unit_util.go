/*******************************************************************************
 * IBM Confidential
 * OCO Source Materials
 * IBM Cloud Container Service, 5737-D43
 * (C) Copyright IBM Corp. 2018 All Rights Reserved.
 * The source code for this program is not  published or otherwise divested of
 * its trade secrets, irrespective of what has been deposited with
 * the U.S. Copyright Office.
 ******************************************************************************/

package util

import (
	"github.com/IBM/ibmcloud-storage-volume-lib/lib/utils/reasoncode"
	"regexp"
	"strconv"
	"strings"
)

// Constants
const (
	DiskUnitsMB string = "MB"
	DiskUnitsGB string = "GB"
	DiskUnitsTB string = "TB"

	DiskTypeRAID1  string = "RAID1"
	DiskTypeRAID10 string = "RAID10"
	DiskTypeSAN    string = "SAN"
	DiskTypeSATA   string = "SATA"
	DiskTypeSSD    string = "SSD"
)

// DiskConfig ...
type DiskConfig struct {
	Count    int
	Size     int // Size is always in GB
	DiskType string
}

// GetMemorySizeInGB processes a string, and extracts the memory size. If the provided
// unit is in MB, this is converted to GB
func GetMemorySizeInGB(value string) (memoryGB int, err error) {

	memory, units, err := splitQuantityAndUnits(value, []string{"MB", "GB"}, false)

	if err != nil {
		return
	} else if units == "MB" {
		memoryGB = memory / 1024
	} else {
		memoryGB = memory
	}

	return

}

// GetNetworkPortSpeedInMbps processes a string, and extracts the network speed. If the provided
// speed is in Gbpg, this is converted to Mbps
func GetNetworkPortSpeedInMbps(value string) (networkSpeedMbps int, err error) {

	networkSpeed, units, err := splitQuantityAndUnits(value, []string{"MBPS", "GBPS"}, false)

	if err != nil {
		return
	} else if units == "GBPS" {
		networkSpeedMbps = networkSpeed * 1000
	} else {
		networkSpeedMbps = networkSpeed
	}

	return

}

// ProcessDiskConfig turns a StorageConfiguration string value into an array of DiskConfig objects
func ProcessDiskConfig(value string, permittedDiskTypes []string) (dcs []DiskConfig, err error) {

	// Split storage configuration into individual disks
	dcStrs := strings.Split(value, ";")
	for _, dcStr := range dcStrs {
		// Trim any excess whitespace from disk config
		dcStr = strings.TrimSpace(dcStr)

		var dc DiskConfig
		dc, err = splitStorageConfiguration(dcStr, permittedDiskTypes)
		if err != nil {
			return
		}

		// Perform some validation on RAID arrays
		if dc.DiskType == DiskTypeRAID10 {
			if dc.Count < 4 || dc.Count%2 != 0 {
				err = NewError(reasoncode.ErrorUnreconciledMachineConfig, "RAID 10 specified, but count not divisible by 2, or fewer than 4 drives")
				return
			}
		} else if dc.DiskType == DiskTypeRAID1 {
			if dc.Count%2 != 0 {
				err = NewError(reasoncode.ErrorUnreconciledMachineConfig, "RAID 1 specified, but count not divisible by 2")
				return
			}
		}

		dcs = append(dcs, dc)
	}

	return
}

// ProcessGPUConfig Takes a GPU configuration spec of the form:
//   <count1>x_<gpu_type1>;<count2>x_<gpu_type2>...
//
// It returns an array of SL keys with one entry for each instance of
// a GPU represented by the spec.
//
// e.g. 1xK80;2xV100 will produce an array with 3 entries:
// [0] = K80
// [1] = V100
// [2] = V100
func ProcessGPUConfig(configuration string, permittedGPUTypes []string) (gpuKeys []string, err error) {

	pattern := "^(([0-9]+)x_)(?i)("                 // match digits, and allow case insensitivity for suffixes
	pattern += strings.Join(permittedGPUTypes, "|") // add the suffixes in the form aaa|bbb|ccc
	pattern += ")"
	pattern += "$" // ensure no characters after the suffix
	re := regexp.MustCompile(pattern)

	// Get an array of our entries
	gpus := strings.Split(configuration, ";")
	for _, gpu := range gpus {
		split := re.FindStringSubmatch(gpu)

		if len(split) == 0 {
			err = NewError(reasoncode.ErrorUnreconciledMachineConfig, "Invalid format for GPU configuration: "+configuration)
			return
		}

		// We can't error here. We have pre-validated using regex
		count, _ := strconv.Atoi(split[2])

		// Add as many entries as was requested for this type
		for i := 0; i < count; i++ {
			gpuKeys = append(gpuKeys, split[3])
		}
	}

	return
}

// splitQuantityAndUnits validates that a field value contains both quantity and units, and is in the correct format
func splitQuantityAndUnits(value string, permittedSuffixes []string, suffixOptional bool) (quantity int, units string, err error) {

	pattern := "^([0-9]+)(?i)("                     // match digits, and allow case insensitivity for suffixes
	pattern += strings.Join(permittedSuffixes, "|") // add the suffixes in the form aaa|bbb|ccc
	pattern += ")"
	if suffixOptional {
		pattern += "{0,1}" // make the suffix optional
	}
	pattern += "$" // ensure no characters after the suffix

	re := regexp.MustCompile(pattern)
	split := re.FindStringSubmatch(value)

	if len(split) == 0 {
		err = NewError(reasoncode.ErrorUnreconciledMachineConfig, "Invalid format for config value: "+value)
	} else {
		// we can't error here. We have pre-validated using regex
		quantity, _ = strconv.Atoi(split[1])
		units = split[2]
	}

	return
}

func splitStorageConfiguration(value string, permittedDiskTypes []string) (dc DiskConfig, err error) {

	pattern := "^(([0-9]+)x_){0,1}(.+)_(?i)("        // match digits, and allow case insensitivity for suffixes
	pattern += strings.Join(permittedDiskTypes, "|") // add the suffixes in the form aaa|bbb|ccc
	pattern += ")"
	pattern += "$" // ensure no characters after the suffix

	re := regexp.MustCompile(pattern)
	split := re.FindStringSubmatch(value)

	if len(split) == 0 {
		err = NewError(reasoncode.ErrorUnreconciledMachineConfig, "Invalid format for storage value: "+value)
		return
	}

	size, units, err := splitQuantityAndUnits(split[3], []string{DiskUnitsMB, DiskUnitsGB, DiskUnitsTB}, false)
	if err != nil {
		return
	}

	// Convert all sizes to GB
	switch units {
	case DiskUnitsMB:
		{
			size /= 1000
		}
	case DiskUnitsTB:
		{
			size *= 1000
		}
	}

	var count int

	if split[2] != "" {
		// we can't error here. We have pre-validated using regex
		count, _ = strconv.Atoi(split[2])
	} else {
		count = 1
	}

	diskType := split[4]

	dc = DiskConfig{
		Count:    count,
		Size:     size,
		DiskType: diskType,
	}

	return
}
