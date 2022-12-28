package ir

type TypeFlags uint32

const (
	FlagPrivate TypeFlags = 1 << iota
	FlagProtected
	FlagPublic
	FlagIntGtZero
	FlagStringNonEmpty
)

func (flags TypeFlags) IsPrivate() bool   { return flags&FlagPrivate != 0 }
func (flags TypeFlags) IsProtected() bool { return flags&FlagProtected != 0 }
func (flags TypeFlags) IsPublic() bool    { return flags&FlagPublic != 0 }

func (flags TypeFlags) IsGtZeroInt() bool      { return flags&FlagIntGtZero != 0 }
func (flags TypeFlags) IsNonEmptyString() bool { return flags&FlagStringNonEmpty != 0 }
