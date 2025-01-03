package tunnel

type Action string

const (
	ActionCreate Action = "create"
	ActionDelete Action = "delete"
	ActionUpdate Action = "update"
	ActionList   Action = "list"
	ActionFind   Action = "find"
	ActionGet    Action = "get"
)
