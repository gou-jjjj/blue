package commands

type Cmd struct {
	Name      string   `json:"name"`
	Summary   string   `json:"summary"`
	Group     string   `json:"group"`
	Arity     int      `json:"arity"`
	Key       string   `json:"key"`
	Value     string   `json:"value"`
	Arguments []string `json:"arguments"`
}
