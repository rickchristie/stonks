package setup

import _ "embed"

//go:embed schema/note.init.sql
var note string

type sqlMap struct {
	Key   string
	Value string
}

// Init is ordered. Tests and local setup execute schema files in this order.
var Init = []sqlMap{
	{Key: "note", Value: note},
}

func InitMap() map[string]string {
	ret := make(map[string]string, len(Init))
	for _, value := range Init {
		ret[value.Key] = value.Value
	}
	return ret
}
