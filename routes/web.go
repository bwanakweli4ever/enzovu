package routes

import (
	controllers "enzovu/app/Http/Controllers"
	"net/http"
)

func RegisterWebRoutes() {
	http.HandleFunc("/", controllers.Home)
	http.HandleFunc("/about", controllers.About)
	http.HandleFunc("/test-model", controllers.TestModel)
	http.HandleFunc("/user", (&controllers.UserController{}).RenderView)
	
}
