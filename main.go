package main

import (
	"enzovu/bootstrap"
	"enzovu/routes"
	"fmt"
	"net/http"
)

func main() {
	elephant := `
	_..--""-.                  .-""--.._
 .-'         \ __...----...__ /         '-.
.'      .:::...,'              ',...:::.      '.
(     .''''''::;                  ;::''''''.     )
\             '-)              (-'             /
\             /                \             /
 \          .'.-.            .-.'.          /
  \         | \0|            |0/ |         /
   |         \  |   .-==-.   |  /         |
   \         '/';          ;'\'         /
	'.._      (_ |  .-==-.  | _)      _..'
		'""'-./ /'        '\ \.-'"'"
			 / /';   .==.   ;'\ \
		.---/ /   \  .==.  /   \ \---.
		|   | |   / .''''. \   | |   |
		|   | |   \ \    / /   | |   |
		|   \ \   /  '""'  \   / /   |
		\    \ \_/          \_/ /    /
		 \    \  -._      _. -  /    /
		  \    \    '""""'    /    /
		   \    \     _    /    /
			\    \   /----\   /    
`
	fmt.Println(elephant)
	bootstrap.InitializeApp()  // Initialize the app
	routes.RegisterWebRoutes() // Register routes
	http.ListenAndServe(":8000", nil)
}
