package models

type CogoCLIInfo struct {
	Command string
	Args    []string
	Flags   *CogoCLIFlags
}

type CogoCLICommand func(CogoCLIInfo) error
