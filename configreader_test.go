package main

import "testing"

func TestConfigReader(t *testing.T) {
	_, err := ReadConfig()

	if err != nil {
		t.Errorf("ConfigReader returned an error")
	}
}