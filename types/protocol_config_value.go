package types

// ProtocolConfigValue is an interface that defines a protocol config value type.
type ProtocolConfigValue interface {
	isProtocolConfigValue()
}

// ProtocolConfigValueU32 defines a protocol config value of type U32.
type ProtocolConfigValueU32 struct {
	U32 string `json:"u32"`
}

// ProtocolConfigValueU64 defines a protocol config value of type U64.
type ProtocolConfigValueU64 struct {
	U64 string `json:"u64"`
}

// ProtocolConfigValueF64 defines a protocol config value of type F64.
type ProtocolConfigValueF64 struct {
	F64 string `json:"f64"`
}

// isProtocolConfigValue implements the ProtocolConfigValue interface for ProtocolConfigValueU32.
func (ProtocolConfigValueU32) isProtocolConfigValue() {}

// isProtocolConfigValue implements the ProtocolConfigValue interface for ProtocolConfigValueU64.
func (ProtocolConfigValueU64) isProtocolConfigValue() {}

// isProtocolConfigValue implements the ProtocolConfigValue interface for ProtocolConfigValueF64.
func (ProtocolConfigValueF64) isProtocolConfigValue() {}
