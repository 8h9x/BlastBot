package database

type Flag uint32

func (f *Flag) Has(flag Flag) bool { return *f&flag != 0 }
func (f *Flag) Set(flag Flag)      { *f |= flag }
func (f *Flag) Toggle(flag Flag)   { *f ^= flag }
func (f *Flag) Clear(flag Flag)    { *f &= ^flag }

const (
	FlagUser Flag = 0

	FlagDeveloper Flag = 1 << iota
	FlagBetaTester
	FlagVIP
	FlagPremiumUser
)
