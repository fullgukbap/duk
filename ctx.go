package duk

import "encoding/json"

type HandlerFunc func(*Ctx)

type Ctx struct {
	payload any
}

func NewCtx(payload any) *Ctx {
	return &Ctx{
		payload: payload,
	}
}

func (c *Ctx) Parse(out any) error {
	jsonBytes, err := json.Marshal(c.payload)
	if err != nil {
		return err
	}

	err = json.Unmarshal(jsonBytes, out)
	if err != nil {
		return err
	}

	return nil
}
