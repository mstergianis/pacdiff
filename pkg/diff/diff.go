package diff

import "encoding/json"

type Diff map[string]any

func ParseDiff(raw []byte) (Diff, error) {
	d := Diff{}
	err := json.Unmarshal(raw, &d)
	return d, err
}
