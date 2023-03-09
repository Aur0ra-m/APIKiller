package runner

import "APIKiller/v1/pkg/types"

// ParseOptions parses the command line flags provided by user
func ParseOptions(options *types.Options) {
	// Show the user banner
	showBanner()
}
