package libconfig

// ----------------------------------------------------------------
// LIBCONFIG data structure
// ----------------------------------------------------------------

type LIBCONFIG struct {
	Settings []*SettingT `@@*`
}

type SettingT struct {
	SetingName   string         `@Name ("="|":")`
	SettingValue *SettingValueT `@@`
}

type SettingValueT struct {
	Primitive *PrimitiveT `( @@ (";"?","?)`
	Group     *GroupT     ` | @@ (","?)`
	Array     *ArrayT     ` | @@ (","?)`
	List      *ListT      ` | @@ (","?))`
}

type PrimitiveT struct {
	Value string `@Value (","?)`
}

type ArrayT struct {
	Primitives []*PrimitiveT `"[" @@* "]"`
}

type GroupT struct {
	Settings []*SettingT `"{" @@* "}"`
}

type ListT struct {
	List []*SettingValueT `"(" @@* ")"`
}
