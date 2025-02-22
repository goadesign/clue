// Code generated by goa v3.20.0, DO NOT EDIT.
//
// front HTTP client CLI support package
//
// Command:
// $ goa gen goa.design/clue/example/weather/services/front/design -o
// services/front

package client

import (
	"encoding/json"
	"fmt"

	front "goa.design/clue/example/weather/services/front/gen/front"
)

// BuildTestAllPayload builds the payload for the front test_all endpoint from
// CLI flags.
func BuildTestAllPayload(frontTestAllBody string) (*front.TestAllPayload, error) {
	var err error
	var body TestAllRequestBody
	{
		err = json.Unmarshal([]byte(frontTestAllBody), &body)
		if err != nil {
			return nil, fmt.Errorf("invalid JSON for body, \nerror: %s, \nexample of valid JSON:\n%s", err, "'{\n      \"exclude\": [\n         \"Veritatis vel inventore voluptatem ab nulla.\",\n         \"In optio sed id.\",\n         \"Tenetur repellendus commodi asperiores.\",\n         \"Quaerat omnis vel quia ab dolorem qui.\"\n      ],\n      \"include\": [\n         \"Consectetur eos nihil accusamus reiciendis eligendi.\",\n         \"Aut autem non.\",\n         \"Aperiam possimus assumenda commodi aut facilis provident.\",\n         \"Aspernatur voluptatem et placeat deserunt.\"\n      ]\n   }'")
		}
	}
	v := &front.TestAllPayload{}
	if body.Include != nil {
		v.Include = make([]string, len(body.Include))
		for i, val := range body.Include {
			v.Include[i] = val
		}
	}
	if body.Exclude != nil {
		v.Exclude = make([]string, len(body.Exclude))
		for i, val := range body.Exclude {
			v.Exclude[i] = val
		}
	}

	return v, nil
}
