package cmd

import (
	"encoding/json"
	"log"
	"log/slog"

	"github.com/hashicorp/memberlist"
)

const PING = "ping"

type ClusterConfig struct {
	NodeName      string
	Host          string
	Port          int
	MasterAddress string //"ip:port"
}

type ClusterMessage struct {
	ActorTag string
	Msg      []byte
}

// called when a node joins
func (ah *ActorHub) NotifyJoin(node *memberlist.Node) {
	slog.Info("Node joined", "node", node.Name)
}

// called when a node leaves
func (ah *ActorHub) NotifyLeave(node *memberlist.Node) {
	slog.Info("Node left", "node", node.Name)
}

// called when a node is updated
func (ah *ActorHub) NotifyUpdate(node *memberlist.Node) {
	slog.Info("Node updated", "node", node.Name)
}

//

/*Message notification events*/
// NodeMeta part of the memberlist Delegate
func (ah *ActorHub) NodeMeta(limit int) []byte {
	return []byte(ah.nodeName)
}

// NotifyMsg handles incoming messages this is to satisfy the Delegate implementation
func (ah *ActorHub) NotifyMsg(msg []byte) {
	if string(msg) == PING {
		return
	}
	cm := ClusterMessage{}
	err := json.Unmarshal(msg, &cm)
	if err != nil {
		return
	}
	ah.eventChan <- Event{eventType: SendMessage, data: cm.Msg, tag: cm.ActorTag}
}

// GetBroadcasts is called when user data messages can be broadcast
func (ah *ActorHub) GetBroadcasts(overhead, limit int) [][]byte {
	// Example: Prepare broadcast message(s)
	return [][]byte{[]byte(PING)}
}

// LocalState is used for a TCP Push/Pull
func (ah *ActorHub) LocalState(join bool) []byte {
	// Example: Return local state
	return []byte("local-state")
}

// MergeRemoteState is invoked after a TCP Push/Pull
func (ah *ActorHub) MergeRemoteState(buf []byte, join bool) {
	// Example: Merge remote state

}

// NoOpLogger is a logger that discards all log messages
type NoOpLogger struct{}

// Write is a no-op function to discard log messages
func (l *NoOpLogger) Write(p []byte) (n int, err error) {
	return len(p), nil
}

// newNetActor ... create a NewActor and connect to master node if provided
func (ah *ActorHub) ConnectToCluster(clusterCfg ClusterConfig) error {

	cfg := memberlist.DefaultWANConfig()
	cfg.Name = clusterCfg.NodeName
	cfg.BindAddr = clusterCfg.Host // Bind address
	cfg.BindPort = clusterCfg.Port // Bind port
	cfg.Delegate = ah
	cfg.Events = ah
	cfg.Label = ActorsHub
	cfg.Logger = log.New(&NoOpLogger{}, "", 0)
	list, err := memberlist.Create(cfg)
	if err != nil {
		return nil
	}
	if len(clusterCfg.MasterAddress) > 0 {
		// Join an existing cluster (optional)
		_, err = list.Join([]string{clusterCfg.MasterAddress})
		if err != nil {
			return nil
		}
	}
	ah.nodeName = clusterCfg.NodeName
	ah.cluster = list
	return nil
}
