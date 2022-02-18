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

	// IdPrefixMonitorRule represents the id prefix of monitor rule.
	IdPrefixMonitorRule = "mor-"

	// IdPrefixUDF represents the id prefix of UDF.
	IdPrefixUDF = "udf-"

	// IdPrefixResourceFile represents the id prefix of resource file.
	IdPrefixResourceFile = "res-"

	// FIXME: removed follow.
	//SourceTablesIDPrefix = "sot-"
	//JobIDPrefix          = "job-"
)

// Defines the id for IdGenerator. To prevent ID conflicts.
const (
	IdInstanceWorkspace int64 = iota + 1
	IdInstanceStreamJob
	IdInstanceSyncJob
	IdInstanceStreamInstance
	IdInstanceFlinkCluster
	IdInstanceNetwork
	IdInstanceDataSource
	IdInstanceUDF
	IdInstanceResourceFile
)

// Defines the id for IdGenerator. To prevent version conflicts.
const (
	VerInstanceStreamJob int64 = iota + 201
	VerInstanceSyncJob   int64 = iota + 201
	VerInstanceResourceFile
)
