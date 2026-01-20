package utils

import (
	"path/filepath"
	"runtime"
	"strings"

	ctestglobals "k8s.io/kubernetes/test/ctest/ctestglobals"
)

func GetCurrentFileName() string {
	// Get caller info (skip 1 to get the caller's caller)
	_, filename, _, ok := runtime.Caller(1)
	if !ok {
		return "unknown"
	}

	// Get just the filename without extension
	base := filepath.Base(filename)
	name := strings.TrimSuffix(base, filepath.Ext(base))
	return name
}

// GetItemByTestInfo searches for the first item in HardcodedConfig
// that contains the given testInfo string in its TestInfo slice
func GetItemByTestInfo(configs ctestglobals.HardcodedConfig, testInfo string) (ctestglobals.HardcodedConfigItem, bool) {
	for _, item := range configs {
		for _, info := range item.TestInfo {
			if strings.Contains(info, testInfo) {
				return item, true
			}
		}
	}
	return ctestglobals.HardcodedConfigItem{}, false
}

// More specific version - exact match
func GetItemByExactTestInfo(configs ctestglobals.HardcodedConfig, testInfo string) (ctestglobals.HardcodedConfigItem, bool) {
	for _, item := range configs {
		for _, info := range item.TestInfo {
			if info == testInfo {
				return item, true
			}
		}
	}
	return ctestglobals.HardcodedConfigItem{}, false
}

// More specific version - exact match
func GetItemByExactTestInfoAndField(configs ctestglobals.HardcodedConfig, testInfo string, field string) (ctestglobals.HardcodedConfigItem, bool) {
	for _, item := range configs {
		for _, info := range item.TestInfo {
			if info == testInfo {
				if item.Field == field {
					return item, true
				}
			}
		}
	}
	return ctestglobals.HardcodedConfigItem{}, false
}

// Case-insensitive search
func GetItemByTestInfoCI(configs ctestglobals.HardcodedConfig, testInfo string) (ctestglobals.HardcodedConfigItem, bool) {
	searchLower := strings.ToLower(testInfo)
	for _, item := range configs {
		for _, info := range item.TestInfo {
			if strings.Contains(strings.ToLower(info), searchLower) {
				return item, true
			}
		}
	}
	return ctestglobals.HardcodedConfigItem{}, false
}

// Version that returns multiple matches
func GetItemsByTestInfo(configs ctestglobals.HardcodedConfig, testInfo string) []ctestglobals.HardcodedConfigItem {
	var matches []ctestglobals.HardcodedConfigItem
	for _, item := range configs {
		for _, info := range item.TestInfo {
			if strings.Contains(info, testInfo) {
				matches = append(matches, item)
				break // Don't add same item multiple times
			}
		}
	}
	return matches
}
