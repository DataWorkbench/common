package constants

// Node statue.
const (
	NodeStatusEnabled  int8 = iota + 1 // => "enabled"
	NodeStatusDisabled                 // => "disabled"
)

// Node priority.
const (
	NodePriorityHighest int8 = iota + 1 // => "highest"
	NodePriorityHigh                    // => "high"
	NodePriorityMedium                  // => "medium"
	NodePriorityLow                     // => "low"
	NodePriorityLowest                  // => "lowest"
)

// Defines the supported node type.
const (
	NodeTypeVirtual int = iota + 1 // => "virtual"
	NodeTypeShell                  // => "shell"
	NodeTypeFlink
)
