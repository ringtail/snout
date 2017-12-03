package sysctl

import (
	log "github.com/Sirupsen/logrus"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

const (
	SYSTEM_PROC_DIR = "/proc/sys/"
)

// export api

func All() map[string]string {
	return walkAndCatFiles(SYSTEM_PROC_DIR)
}

func Find(names ...string) map[string]string {
	settings := make(map[string]string)
	for _, name := range names {
		settings_file_path := convert_settings_to_path(name)
		if _, err := os.Stat(settings_file_path); os.IsNotExist(err) {
			log.Warnf("Failed to query settings %v, Because of %v\n", name, err.Error())
			continue
		}
		bytes, err := ioutil.ReadFile(settings_file_path)
		if err != nil {
			log.Warnf("Failed to read settings %v from proc, Because of %v", name, err.Error())
			continue
		}
		settings[name] = string(bytes)
	}
	return settings
}

func Apply(name string, value string) error {
	settings_file_path := convert_settings_to_path(name)
	if _, err := os.Stat(settings_file_path); os.IsNotExist(err) {
		log.Warnf("Failed to apply settings change, Because of %v", err.Error())
		return err
	}
	err := ioutil.WriteFile(settings_file_path, []byte(value), 0644)
	if err != nil {
		log.Warnf("Failed to apply settings change, Because of %v", err.Error())
	}
	return err
}

// internal api

func walkAndCatFiles(rootPath string) map[string]string {
	settings := make(map[string]string)
	if !strings.Contains(rootPath, SYSTEM_PROC_DIR) {
		log.Warnf("The specific path is invalid %v", rootPath)
		return nil
	}
	fi, err := os.Stat(rootPath)
	if err != nil {
		log.Warnf("Failed to open files, Because of %v \n", rootPath)
		return nil
	}
	if fi.IsDir() == false {
		bytes, err := ioutil.ReadFile(rootPath)
		if err != nil {
			log.Warnf("Failed to read files, Because of %v", err.Error())
			return nil
		}
		settings_name := convert_path_to_settings(rootPath)
		settings[settings_name] = strings.Replace(strings.Trim(string(bytes), "\n"), "\t", " ", -1)
		return settings
	} else {
		files, err := ioutil.ReadDir(rootPath)
		if err != nil {
			log.Warnf("Failed to list dir, Because of %v", err.Error())
			return nil
		}
		for _, file := range files {
			base_path := filepath.Join(rootPath, file.Name())
			subfolder_settings := walkAndCatFiles(base_path)
			for key, value := range subfolder_settings {
				settings[key] = value
			}
		}
	}
	return settings
}

func convert_path_to_settings(path string) string {
	name := path[len(SYSTEM_PROC_DIR):]
	settings_arr := strings.Split(name, "/")
	return strings.Join(settings_arr, ".")
}

func convert_settings_to_path(name string) string {
	settings_arr := strings.Split(name, ".")
	return filepath.Join(SYSTEM_PROC_DIR, strings.Join(settings_arr, "/"))
}
