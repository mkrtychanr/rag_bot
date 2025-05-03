package model

type RightsPolicy int64

var (
	ReadOnlyRightPolicy  = RightsPolicy(0)
	ReadWriteRightPolicy = RightsPolicy(1)
)
