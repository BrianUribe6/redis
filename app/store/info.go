package store

import (
	"fmt"
)

type ReplicationInfo struct {
	role                       string
	connectedSlaves            uint
	masterFailoverState        string
	masterReplId               string
	masterReplOffset           int
	secondReplOffset           int
	replBacklogActive          int
	replBacklogSize            int
	replBacklogFirstByteOffset int
	replBacklogHistLen         int
}

const (
	MASTER_ROLE = "master"
	SLAVE_ROLE  = "slave"
)

var Info ReplicationInfo = ReplicationInfo{
	role: MASTER_ROLE,
}

func (r *ReplicationInfo) SetRole(role string) {
	if role != MASTER_ROLE && role != SLAVE_ROLE {
		panic("Invalid role")
	}
	r.role = role
}

func (r *ReplicationInfo) String() string {
	return fmt.Sprint(
		"role:", r.role, "\n",
		"connected_slaves:", r.connectedSlaves, "\n",
		"master_failover_state:", r.masterFailoverState, "\n",
		"master_replid:", r.masterReplId, "\n",
		"master_repl_offset:", r.masterReplOffset, "\n",
		"second_repl_offset:", r.secondReplOffset, "\n",
		"repl_backlog_active:", r.replBacklogActive, "\n",
		"repl_backlog_size:", r.replBacklogSize, "\n",
		"repl_backlog_first_byte_offset:", r.replBacklogFirstByteOffset, "\n",
		"repl_backlog_histlen:", r.replBacklogHistLen, "\n",
	)

}
