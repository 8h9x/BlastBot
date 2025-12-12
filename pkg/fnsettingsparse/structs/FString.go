package structs

type FString struct {
	Length int32
	Bytes  []byte
}

func (f *FString) String() string {
	return ""
}
