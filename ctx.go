package duk

import "encoding/json"

type Ctx struct {
	Payload string
}

func (c *Ctx) Parse(out any) error {
	err := json.Unmarshal([]byte(c.Payload), out)
	if err != nil {
		return err
	}

	return nil
}
