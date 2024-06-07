package duk

import "encoding/json"

type HandlerFunc func(*Ctx)

type Ctx struct {
	Payload string
}

func NewCtx(payload string) *Ctx {
	return &Ctx{
		Payload: payload,
	}
}

func (c *Ctx) Parse(out any) error {
	err := json.Unmarshal([]byte(c.Payload), out)
	if err != nil {
		return err
	}

	return nil
}
