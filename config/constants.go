package config

import "time"

const (
	DatabaseQueryTimeLayout string = `'YYYY-MM-DD"T"HH24:MI:SS"."MS"Z"TZ'`
	DatabaseTimeLayout      string = time.RFC3339
	ErrTheSameId                   = "cannot use the same uuid for 'id' and 'parent_id' fields"
	ErrRpcNodFoundAndNoRows        = "rpc error: code = NotFound desc = no rows in result set"
	ErrNoRows                      = "no rows in result set"
	ErrObjectType                  = "object type error: code =  NodFound"
	ErrEnvNodFound                 = "No .env file found"
	AgentSmsText                   = " Код подтверждения для PaynetBNPL"
	ExpireTimeout                  = 2
)
