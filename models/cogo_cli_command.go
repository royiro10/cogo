package models

type CogoCLIInfo struct {
	Commad string
	Args   []string
}

type CogoCLICommand func(CogoCLIInfo) error
