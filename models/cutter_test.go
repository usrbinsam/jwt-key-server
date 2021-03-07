package models

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestApp(t *testing.T) {

	var (
		app    Application
		cutter KeyCutter
	)

	app = Application{
		Name:           "Test Application",
		SupportMessage: "Contact Us at 555-555-5555",
	}
	app.ID = 1

	cutter = KeyCutter{Config: KeyCutterConfig{32}}

	t.Run("TestCutKey", func(t *testing.T) {
		key, err := cutter.CutKey(&app)

		assert.Nil(t, err, "CutKey failed")
		assert.NotNil(t, key, "key failed to cut")
		assert.Equal(t, uint(1), key.ApplicationID, "Key didn't set for the correct application")
	})

}
