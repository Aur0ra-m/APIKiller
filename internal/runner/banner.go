package runner

import "fmt"

var banner = fmt.Sprintf(`
 █████╗ ██████╗ ██╗██╗  ██╗██╗██╗     ██╗     ███████╗██████╗
██╔══██╗██╔══██╗██║██║ ██╔╝██║██║     ██║     ██╔════╝██╔══██╗
███████║██████╔╝██║█████╔╝ ██║██║     ██║     █████╗  ██████╔╝
██╔══██║██╔═══╝ ██║██╔═██╗ ██║██║     ██║     ██╔══╝  ██╔══██╗
██║  ██║██║     ██║██║  ██╗██║███████╗███████╗███████╗██║  ██║
╚═╝  ╚═╝╚═╝     ╚═╝╚═╝  ╚═╝╚═╝╚══════╝╚══════╝╚══════╝╚═╝  ╚═╝
Version: %s`+"\n", VERSION)

//
// showBanner
//  @Description: show banner on the terminal
//
func showBanner() {
	fmt.Println(banner)
}
