package config

import "testing"

func TestDefaultMetadataConfig(t *testing.T) {
	c := DefaultMetadataConfig()
	if !c.Enabled {
		t.Error("expected Enabled to be true")
	}
	if c.MaxKeysPerPort != 16 {
		t.Errorf("expected MaxKeysPerPort=16, got %d", c.MaxKeysPerPort)
	}
}

func TestMetadataConfigValidateOK(t *testing.T) {
	c := DefaultMetadataConfig()
	if err := c.Validate(); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

func TestMetadataConfigValidateDisabled(t *testing.T) {
	c := MetadataConfig{Enabled: false, MaxKeysPerPort: 0}
	if err := c.Validate(); err != nil {
		t.Errorf("disabled config should always be valid, got: %v", err)
	}
}

func TestMetadataConfigValidateZeroMaxKeys(t *testing.T) {
	c := MetadataConfig{Enabled: true, MaxKeysPerPort: 0}
	if err := c.Validate(); err == nil {
		t.Error("expected error for MaxKeysPerPort=0")
	}
}

func TestMetadataConfigValidateNegativeMaxKeys(t *testing.T) {
	c := MetadataConfig{Enabled: true, MaxKeysPerPort: -1}
	if err := c.Validate(); err == nil {
		t.Error("expected error for negative MaxKeysPerPort")
	}
}

func TestMetadataConfigValidateExceedsLimit(t *testing.T) {
	c := MetadataConfig{Enabled: true, MaxKeysPerPort: 512}
	if err := c.Validate(); err == nil {
		t.Error("expected error for MaxKeysPerPort > 256")
	}
}

func TestMetadataConfigValidateBoundary(t *testing.T) {
	c := MetadataConfig{Enabled: true, MaxKeysPerPort: 256}
	if err := c.Validate(); err != nil {
		t.Errorf("expected 256 to be valid, got: %v", err)
	}
}
