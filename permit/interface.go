package permit

type WardenQueryAPI interface {
	Match(map[string]interface{}) bool
}
