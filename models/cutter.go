package models

type KeyCutterConfig struct {
	SecretSize int
}

type KeyCutter struct {
	Config KeyCutterConfig
}

func (c KeyCutter) CutKey(app *Application) (*Key, error) {

	key := Key{
		ApplicationID: app.ID,
	}

	_, err := key.SetRandomSecret(c.Config.SecretSize)

	if err != nil {
		return nil, err
	}

	return &key, nil
}
