package envvar

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestHealthy(t *testing.T) {

	// setup
	_ = os.Setenv("EV1", "value1")
	_ = os.Setenv("EV2", "value2")
	_ = os.Setenv("EV3", "")

	// verify
	cl := NewClient()
	{
		val, err := cl.LoadSecret("EV1")
		assert.NoError(t, err)
		assert.Equal(t, "value1", val)
	}
	{
		val, err := cl.LoadSecret("EV2")
		assert.NoError(t, err)
		assert.Equal(t, "value2", val)
	}
	{
		val, err := cl.LoadSecret("EV3")
		assert.NoError(t, err)
		assert.Equal(t, "", val)
	}

}

func TestInjectFaulty(t *testing.T) {

	// setup
	_ = os.Setenv("EV1", "value1")
	_ = os.Setenv("EV3", "")
	_ = os.Unsetenv("EV4")

	// verify
	cl := NewClient()
	{
		val, err := cl.LoadSecret("EV1")
		assert.NoError(t, err)
		assert.Equal(t, "value1", val)
	}
	{
		val, err := cl.LoadSecret("EV3")
		assert.NoError(t, err)
		assert.Equal(t, "", val)
	}
	{
		_, err := cl.LoadSecret("EV4")
		assert.Error(t, err)
	}

}
