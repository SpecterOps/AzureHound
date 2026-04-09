package bloodhound

type bloodhoundResponse[T any] struct {
	Data T `json:"data"`
}
