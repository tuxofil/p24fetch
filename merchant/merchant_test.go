package merchant

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMD5Hex(t *testing.T) {
	assert.Equal(t, "d41d8cd98f00b204e9800998ecf8427e", md5hex(""))
	assert.Equal(t, "0cc175b9c0f1b6a831c399e269772661", md5hex("a"))
}

func TestSHA1Hex(t *testing.T) {
	assert.Equal(t, "da39a3ee5e6b4b0d3255bfef95601890afd80709", sha1hex(""))
	assert.Equal(t, "86f7e437faa5a7fce15d1ddcb9eaeaea377667b8", sha1hex("a"))
}
