package models

type CogoCLIInfo struct {
	Commad string
	Args   []string
	Flags  *CogoCLIFlags
}

type CogoCLICommand func(CogoCLIInfo) error
