package constants

const (
	UnseenServerName  = "unseen"
	EphineaServerName = "ephinea"
)

type EphineaAccountMode int

const (
	Normal EphineaAccountMode = iota
	Hardcore
	Sandbox
)
