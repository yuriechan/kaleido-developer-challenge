package domain

type ItemState int64

const (
	ItemStateUnspecified ItemState = iota
	ItemStateListed
	ItemStateSold
)

type Item struct {
	ID    string    `json:"id"`
	Name  string    `json:"name"`
	State ItemState `json:"state"`
}
