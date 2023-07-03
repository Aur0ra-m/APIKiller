package authorize

type authGroup struct {
	Domain        []string       `mapstructure:"domain"`
	ReplaceGroups []replaceGroup `mapstructure:"replaceGroup"`
}

type replaceGroup struct {
	Position int    `mapstructure:"position"`
	Key      string `mapstructure:"key"`
	Value    string `mapstructure:"value"`
}

// position codes
const (
	Replace_Position_Code_Header = 0
	Replace_Position_Code_Query  = 1
	Replace_Position_Code_Body   = 2
)
