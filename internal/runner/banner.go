package runner

import (
	"APIKiller/pkg/config"
	"fmt"
)

var banner = fmt.Sprintf(`
█████╗ ██████╗ ██╗██╗  ██╗██╗██╗     ██╗     ███████╗██████╗
██╔══██╗██╔══██╗██║██║ ██╔╝██║██║     ██║     ██╔════╝██╔══██╗
███████║██████╔╝██║█████╔╝ ██║██║     ██║     █████╗  ██████╔╝
██╔══██║██╔═══╝ ██║██╔═██╗ ██║██║     ██║     ██╔══╝  ██╔══██╗
██║  ██║██║     ██║██║  ██╗██║███████╗███████╗███████╗██║  ██║
╚═╝  ╚═╝╚═╝     ╚═╝╚═╝  ╚═╝╚═╝╚══════╝╚══════╝╚══════╝╚═╝  ╚═╝
Version: %s`+"\n",
	config.VERSION)

// showBanner is used to show banner to user
func showBanner() {
	fmt.Println(banner)
}
