package config

import (
	"testing"
)

func TestGenerateImagehubTemplate(t *testing.T) {
	_, err := GenerateImagehubTemplate("file://../../templates/imagehub-filter.yaml")
	if err != nil {
		t.Error(err)
	}
}
