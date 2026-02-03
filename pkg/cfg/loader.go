package cfg

import (
	"encoding/json"
	"os"
)

var C config

func Load(file string) error {
	b, err := os.ReadFile(file)
	if err != nil {
		return err
	}

	c := config{}
	if err := json.Unmarshal(b, &c); err != nil {
		return err
	} else {
		C = c
	}

	return nil
}
