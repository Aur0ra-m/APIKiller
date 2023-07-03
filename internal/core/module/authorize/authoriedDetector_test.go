package authorize

import (
	"fmt"
	"github.com/spf13/viper"
	"testing"
)

func Test(t *testing.T) {
	tests := []struct {
		name string
	}{
		// TODO: Add test cases.
		{},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			viper.SetConfigFile("D:\\Projects\\GO\\APIKiller\\config\\config.dev.yaml")
			viper.ReadInConfig()

			var authGroups = []authGroup{}
			viper.UnmarshalKey("app.module.authorizedDetector.authGroup", &authGroups)

			fmt.Println(authGroups)
		})
	}
}
