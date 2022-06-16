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

	// IdPrefixAlertPolicy represents the id prefix of `alert policy`.
	IdPrefixAlertPolicy = "alt-"

	// IdPrefixNotifier represents the id prefix of `Notifier`
	IdPrefixNotifier = "nof-"

	// IdPrefixDataServiceCluster represents the id prefix of `dataservice cluster`
	IdPrefixDataServiceCluster = "dsc-"

	// IdPrefixApiGroup represents the id prefix of `dataservice api group`
	IdPrefixApiGroup = "dsg-"

	// IdPrefixCustomerApi represents the id prefix of `dataservice customer api`
	IdPrefixCustomerApi = "dsa-"

	// IdPrefixApiRequestParam represents the id prefix of `customer api request param`
	IdPrefixApiRequestParam = "dsq-"

	// IdPrefixApiResponseParam represents the id prefix of `customer api response param`
	IdPrefixApiResponseParam = "dsp-"

	// IdPrefixCustomerApiVersion represents the id prefix of `dataservice customer api version`
	IdPrefixCustomerApiVersion = "dsv-"

	//// IdPrefixUDF represents the id prefix of `UDF`.
	//IdPrefixUDF = "udf-"
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
	IdInstanceAlertPolicy
	IdInstanceDataServiceCluster
	IdInstanceApiGroup
	IdInstanceCustomerApi
	IdInstanceApiRequestParam
	IdInstanceApiResponseParam
	IdInstanceCustomerApiVersion
)

// Defines the id for IdGenerator. To prevent version conflicts.
// The newly added instance must be in the last.
const (
	VerInstanceStreamJob int64 = iota + 201
	VerInstanceSyncJob
	VerInstanceResourceFile
)
