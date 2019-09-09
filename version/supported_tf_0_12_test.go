// +build !tf_0_11

package version

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSupported(t *testing.T) {
	want := TF12
	if !Supported(want) {
		t.Errorf("Supported(%s) should return true", want)
	}
}

func TestLeastSupportedVersion(t *testing.T) {
	want := TF12
	got := LeastSupportedVersion()
	if got != want {
		t.Errorf("LeastSupportedVersion() = %s; want %s", got, want)
	}
}

func TestAllSupportedVersions(t *testing.T) {
	want := []string{TF12}
	got := AllSupportedVersions()
	assert.Equal(t, want, got)
}
