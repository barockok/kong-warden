package warden

type WardenQueryAPI interface {
	Match(map[string]interface{}) bool
}
