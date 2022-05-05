package flink

type TotalResource struct {
	CpuCores          float64     `json:"cpuCores"`
	TaskHeapMemory    int         `json:"taskHeapMemory"`
	TaskOffHeapMemory int         `json:"taskOffHeapMemory"`
	ManagedMemory     int         `json:"managedMemory"`
	NetworkMemory     int         `json:"networkMemory"`
	ExtendedResources interface{} `json:"extendedResources"`
}

type FreeResource struct {
	CpuCores          float64     `json:"cpuCores"`
	TaskHeapMemory    int         `json:"taskHeapMemory"`
	TaskOffHeapMemory int         `json:"taskOffHeapMemory"`
	ManagedMemory     int         `json:"managedMemory"`
	NetworkMemory     int         `json:"networkMemory"`
	ExtendedResources interface{} `json:"extendedResources"`
}

type Hardware struct {
	CpuCores       int   `json:"cpuCores"`
	PhysicalMemory int64 `json:"physicalMemory"`
	FreeMemory     int   `json:"freeMemory"`
	ManagedMemory  int   `json:"managedMemory"`
}

type MemoryConfiguration struct {
	FrameworkHeap      int         `json:"frameworkHeap"`
	TaskHeap           int         `json:"taskHeap"`
	FrameworkOffHeap   int         `json:"frameworkOffHeap"`
	TaskOffHeap        int         `json:"taskOffHeap"`
	NetworkMemory      int         `json:"networkMemory"`
	ManagedMemory      int         `json:"managedMemory"`
	JvmMetaspace       int         `json:"jvmMetaspace"`
	JvmOverhead        int         `json:"jvmOverhead"`
	TotalFlinkMemory   interface{} `json:"totalFlinkMemory"`
	TotalProcessMemory int         `json:"totalProcessMemory"`
}

// TaskManager represents the task manager info.
type TaskManager struct {
	Id                     string               `json:"id"`
	Path                   string               `json:"path"`
	DataPort               int                  `json:"dataPort"`
	JmxPort                int                  `json:"jmxPort"`
	TimeSinceLastHeartbeat int64                `json:"timeSinceLastHeartbeat"`
	SlotsNumber            int                  `json:"slotsNumber"`
	FreeSlots              int                  `json:"freeSlots"`
	TotalResource          *TotalResource       `json:"totalResource"`
	FreeResource           *FreeResource        `json:"freeResource"`
	Hardware               *Hardware            `json:"hardware"`
	MemoryConfiguration    *MemoryConfiguration `json:"memoryConfiguration"`
}

// TaskManagers represents the response of `/taskmanagers`
type TaskManagers struct {
	TaskManagers []*TaskManager `json:"taskmanagers"`
}
