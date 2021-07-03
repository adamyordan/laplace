
package config

import (
	"testing"
)

func TestConfigInit(t *testing.T) {
	_,err := ConfigInit()
	if err != nil {
		t.Error(err)
	}
}

func TestSetDefaults(t *testing.T) {
	err := SetDefaults()
	if err != nil {
		t.Error(err)
	}
}