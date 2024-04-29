package store

import (
	"fmt"
	"strings"

	"github.com/google/uuid"
)

type ReplicationInfo struct {
	role                       string
	ConnectedSlaves            uint
	MasterFailoverState        string
	MasterReplId               string
	MasterReplOffset           int
	SecondReplOffset           int
	ReplBacklogActive          int
	ReplBacklogSize            int
	ReplBacklogFirstByteOffset int
	ReplBacklogHistLen         int
}

const (
	MASTER_ROLE = "master"
	SLAVE_ROLE  = "slave"
)

var Info ReplicationInfo = ReplicationInfo{
	role:             MASTER_ROLE,
	MasterReplId:     strings.Replace(uuid.NewString(), "-", "", -1),
	MasterReplOffset: 0,
}

func (r *ReplicationInfo) SetRole(role string) {
	if role != MASTER_ROLE && role != SLAVE_ROLE {
		panic("Invalid role")
	}
	r.role = role
}

func (r *ReplicationInfo) Role() string {
	return r.role
}

func (r *ReplicationInfo) String() string {
	return fmt.Sprint(
		"role:", r.role, "\n",
		"connected_slaves:", r.ConnectedSlaves, "\n",
		"master_failover_state:", r.MasterFailoverState, "\n",
		"master_replid:", r.MasterReplId, "\n",
		"master_repl_offset:", r.MasterReplOffset, "\n",
		"second_repl_offset:", r.SecondReplOffset, "\n",
		"repl_backlog_active:", r.ReplBacklogActive, "\n",
		"repl_backlog_size:", r.ReplBacklogSize, "\n",
		"repl_backlog_first_byte_offset:", r.ReplBacklogFirstByteOffset, "\n",
		"repl_backlog_histlen:", r.ReplBacklogHistLen, "\n",
	)

}
