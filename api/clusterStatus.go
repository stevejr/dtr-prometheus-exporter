package api

import (
  "time"
	"encoding/json"
)

// CSAPIEndpoint - the DTR API endpoint for Cluster Status
const CSAPIEndpoint = "api/v0/meta/cluster_status" 

// ClusterStatus - DTR Cluster Status API response
type ClusterStatus struct {
	RethinkSystemTables     RethinkSystemTables    `json:"rethink_system_tables"`
  ReplicaHealth           map[string]string      `json:"replica_health"`
	ReplicaTimestamp        map[string]time.Time   `json:"replica_timestamp"`
	ReplicaReadonly         map[string]bool        `json:"replica_readonly"`
	GcLockHolder            string                 `json:"gc_lock_holder"` 
}
// ClusterConfig - Cluster config
type ClusterConfig struct {
	HeartbeatTimeoutSecs    int                     `json:"heartbeat_timeout_secs"`
	ID                   string                     `json:"id"`
}
// DbConfig - DB Config
type DbConfig struct {
	ID                    string                    `json:"id"`
	Name                  string                    `json:"name"`
}
// CanonicalAddresses - Network address - 1 per replica
type CanonicalAddresses struct {
	Host                  string                    `json:"host"`
	Port                  int                       `json:"port"`
}
// ConnectedTo - Replicas connected to each other
type ConnectedTo struct {
  ReplicaID             map[string]bool         
}
// Network - The Server Network details
type Network struct {
	CanonicalAddresses    []CanonicalAddresses      `json:"canonical_addresses"`
	ClusterPort           int                       `json:"cluster_port"`
	ConnectedTo           map[string]bool           `json:"connected_to"`
	Hostname              string                    `json:"hostname"`
	HTTPAdminPort         string                    `json:"http_admin_port"`
	ReqlPort              int                       `json:"reql_port"`
	TimeConnected         time.Time                 `json:"time_connected"`
}
// Process - The Server Process details
type Process struct {
	Argv                  []string                  `json:"argv"`
	CacheSizeMb           int                       `json:"cache_size_mb"`
	Pid                   int                       `json:"pid"`
	TimeStarted           time.Time                 `json:"time_started"`
	Version               string                    `json:"version"`
}
// ServerStatus -  Server status details
type ServerStatus struct {
	ID                    string                    `json:"id"`
	Name                  string                    `json:"name"`
	Network               Network                   `json:"network"`
	Process               Process                   `json:"process"`
}
// Replicas - Replica details used in TableStatus
type Replicas struct {
	Server                string                    `json:"server"`
	State                 string                    `json:"state"`
}
// QueryEngine - Query Metrics used in Stats
type QueryEngine struct {
  ClientConnections     int                       `json:"client_connections,omitempty"`
	ClientsActive         int                       `json:"clients_active,omitempty"`
	QueriesPerSec         float64                   `json:"queries_per_sec,omitempty"`
	QueriesTotal          int                       `json:"queries_total,omitempty"`
	ReadDocsPerSec        float64                   `json:"read_docs_per_sec,omitempty"`
	ReadDocsTotal         int                       `json:"read_docs_total,omitempty"`
	WrittenDocsPerSec     float64                   `json:"written_docs_per_sec,omitempty"`
	WrittenDocsTotal      int                       `json:"written_docs_total,omitempty"`
}
// Cache - Used cache
type Cache struct {
	InUseBytes            int                       `json:"in_use_bytes,omitempty"`
}
// SpaceUsage - Disk space used by table
type SpaceUsage struct {
	DataBytes             int                       `json:"data_bytes,omitempty"`
	GarbageBytes          int                       `json:"garbage_bytes,omitempty"`
	MetadataBytes         int                       `json:"metadata_bytes,omitempty"`
	PreallocatedBytes     int                       `json:"preallocated_bytes,omitempty"`
}
// Disk - Disk related details by table
type Disk struct {
	ReadBytesPerSec       float64                   `json:"read_bytes_per_sec,omitempty"`
	ReadBytesTotal        int                       `json:"read_bytes_total,omitempty"`
	SpaceUsage            SpaceUsage                `json:"space_usage,omitempty"`
	WrittenBytesPerSec    float64                   `json:"written_bytes_per_sec,omitempty"`
	WrittenBytesTotal     int                       `json:"written_bytes_total,omitempty"`
}
// StorageEngine -  Storage details by table
type StorageEngine struct {
	Cache                 Cache                     `json:"cache,omitempty"`
	Disk                  Disk                      `json:"disk,omitempty"`
}
// Stats - Stats by table
type Stats struct {
  ID                    []string                  `json:"id"`
  QueryEngine           QueryEngine               `json:"query_engine,omitempty"`
  Server                string                    `json:"server,omitempty"`
  StorageEngine         StorageEngine             `json:"storage_engine,omitempty"`
  Table                 string                    `json:"table,omitempty"`   
}
// ShardsTS - Table sharding used in Table Status
type ShardsTS struct {
	PrimaryReplicas       []string                  `json:"primary_replicas"`
	Replicas              []Replicas                `json:"replicas"`
}                     
// ShardsTC - Table sharding used in Table Config
type ShardsTC struct {
	NonvotingReplicas     []interface{}             `json:"nonvoting_replicas"`
	PrimaryReplica        string                    `json:"primary_replica"`
	Replicas              []string                  `json:"replicas"`
}
// Status - Status of table used in Table Status
type Status struct {
	AllReplicasReady      bool                      `json:"all_replicas_ready"`
	ReadyForOutdatedReads bool                      `json:"ready_for_outdated_reads"`
	ReadyForReads         bool                      `json:"ready_for_reads"`
	ReadyForWrites        bool                      `json:"ready_for_writes"`
}
// TableStatus - Status of each table
type TableStatus struct {
	Db                    string                    `json:"db"`
	ID                    string                    `json:"id"`
	Name                  string                    `json:"name"`
	RaftLeader            string                    `json:"raft_leader"`
	Shards                []ShardsTS                `json:"shards"`
	Status                Status                    `json:"status"`
}
// TableConfig - Config of each table
type TableConfig struct {
	Db                    string                    `json:"db"`
	Durability            string                    `json:"durability"`
	ID                    string                    `json:"id"`
	Indexes               []string                  `json:"indexes"`
	Name                  string                    `json:"name"`
	PrimaryKey            string                    `json:"primary_key"`
	Shards                []ShardsTC                `json:"shards"`
	WriteAcks             string                    `json:"write_acks"`
}
// ReplicaHealth - Structure to hold the replica_health content
type ReplicaHealth struct {
	ReplicaID             map[string]string 
}
// ReplicaTimestamp - Structure to hold the replica_timestamp content
type ReplicaTimestamp struct {
	ReplicaID             map[string]time.Time
}
// ReplicaReadOnly - Structure to hold the replica_readonly content
type ReplicaReadOnly struct {
	ReplicaID             map[string]bool
}
// RethinkSystemTables - Structure to hold the rethink_system_tables content
type RethinkSystemTables struct {
	ClusterConfig         []ClusterConfig           `json:"cluster_config"`
	CurrentIssues         []interface{}             `json:"current_issues"`
	DbConfig              []DbConfig                `json:"db_config"`
	ServerStatus          []ServerStatus            `json:"server_status"`
	Stats                 []Stats                   `json:"stats"`
	TableConfig           []TableConfig             `json:"table_config"`
	TableStatus           []TableStatus             `json:"table_status"`
}

// CSReplicaHealth - main struct to hold the replica stats
type CSReplicaHealth struct {
  ID                  string
  HealthyCount        int
}

// GetCSReplicaHealthStats - Get the ReplicaHealth stats
func GetCSReplicaHealthStats(jsonData []byte) (*[]CSReplicaHealth, error) {
  var result []CSReplicaHealth
  var healthCount int

	var dtrClusterStatus ClusterStatus
		
	err := json.Unmarshal(jsonData, &dtrClusterStatus)
  if err != nil {
    return nil, err
  }
  
  for replica, health := range dtrClusterStatus.ReplicaHealth {
    if health == "OK" {
      healthCount = 1
    }
    result = append(result, CSReplicaHealth{
      ID: replica,
      HealthyCount: healthCount,
    }) 
  }

  return &result, nil
}