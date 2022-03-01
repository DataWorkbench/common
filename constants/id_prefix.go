package constants

// Declares id prefix
const (
	// IdPrefixWorkspace represents the id prefix of `workspace`.
	IdPrefixWorkspace = "wks-"

	// IdPrefixRoleSystem represents the id prefix of `role system`.
	IdPrefixRoleSystem = "ros-"

	// IdPrefixRoleCustom represents the id prefix of `role custom`.
	IdPrefixRoleCustom = "roc-"

	// IdPrefixStreamJob represents the id prefix of `stream job`.
	IdPrefixStreamJob = "stj-"

	// IdPrefixSyncJob represents the id prefix of `stream job`.
	IdPrefixSyncJob = "syj-"

	// IdPrefixStreamInstance represents the id prefix of `stream instance`.
	IdPrefixStreamInstance = "sti-"

	// IdPrefixSyncInstance represents the id prefix of `stream instance`.
	IdPrefixSyncInstance = "syi-"

	// IdPrefixFlinkCluster represents the id prefix of `cluster flink`.
	IdPrefixFlinkCluster = "cfi-"

	// IdPrefixNetwork represents the id prefix of `network`.
	IdPrefixNetwork = "net-"

	// IdPrefixDatasource represents the id prefix of `datasource meta`.
	IdPrefixDatasource = "som-"

	// IdPrefixResourceFile represents the id prefix of `resource file`.
	IdPrefixResourceFile = "res-"

	// IdPrefixProjectModule represents the id prefix of `project module`
	IdPrefixProjectModule = "pmo-"

	// IdPrefixMonitorRule represents the id prefix of `monitor rule`.
	IdPrefixMonitorRule = "mor-"

	// IdPrefixUDF represents the id prefix of `UDF`.
	IdPrefixUDF = "udf-"
)

// Defines the id for IdGenerator. To prevent ID conflicts.
// The newly added instance must be in the last.
const (
	IdInstanceWorkspace int64 = iota + 1
	IdInstanceStreamJob
	IdInstanceStreamInstance
	IdInstanceSyncJob
	IdInstanceSyncInstance
	IdInstanceFlinkCluster
	IdInstanceNetwork
	IdInstanceDataSource
	IdInstanceResourceFile
	IdInstanceUDF
)

// Defines the id for IdGenerator. To prevent version conflicts.
// The newly added instance must be in the last.
const (
	VerInstanceStreamJob int64 = iota + 201
	VerInstanceSyncJob
	VerInstanceResourceFile
)
