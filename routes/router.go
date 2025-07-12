package routes

import (
	"context"
	"net/http"
	"regexp"
	"strings"
)

type Route struct {
	Method      string
	Pattern     string
	Handler     http.HandlerFunc
	Regex       *regexp.Regexp
	ParamNames  []string
	Middlewares []func(http.Handler) http.Handler
}

type Router struct {
	routes      []Route
	middlewares []func(http.Handler) http.Handler
}

type contextKey string

const ParamsKey contextKey = "params"

func NewRouter() *Router {
	return &Router{}
}

// Use adds middleware to all routes
func (r *Router) Use(middleware func(http.Handler) http.Handler) {
	r.middlewares = append(r.middlewares, middleware)
}

func (r *Router) AddRoute(method, pattern string, handler http.HandlerFunc, middlewares ...func(http.Handler) http.Handler) {
	// Extract parameter names and convert pattern to regex
	paramNames := []string{}
	regexPattern := "^" + pattern + "$"

	// Find all {param} patterns
	paramRegex := regexp.MustCompile(`\{([^}]+)\}`)
	matches := paramRegex.FindAllStringSubmatch(pattern, -1)

	for _, match := range matches {
		paramNames = append(paramNames, match[1])
		// Replace {param} with capture group
		regexPattern = strings.Replace(regexPattern, match[0], "([^/]+)", 1)
	}

	regex := regexp.MustCompile(regexPattern)

	route := Route{
		Method:      method,
		Pattern:     pattern,
		Handler:     handler,
		Regex:       regex,
		ParamNames:  paramNames,
		Middlewares: middlewares,
	}
	r.routes = append(r.routes, route)
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	for _, route := range r.routes {
		if route.Method == req.Method && route.Regex.MatchString(req.URL.Path) {
			// Extract parameters
			matches := route.Regex.FindStringSubmatch(req.URL.Path)
			params := make(map[string]string)

			for i, name := range route.ParamNames {
				if i+1 < len(matches) {
					params[name] = matches[i+1]
				}
			}

			// Add params to request context
			ctx := context.WithValue(req.Context(), ParamsKey, params)
			req = req.WithContext(ctx)

			// Build handler chain with middlewares
			handler := http.Handler(route.Handler)

			// Apply route-specific middlewares
			for i := len(route.Middlewares) - 1; i >= 0; i-- {
				handler = route.Middlewares[i](handler)
			}

			// Apply global middlewares
			for i := len(r.middlewares) - 1; i >= 0; i-- {
				handler = r.middlewares[i](handler)
			}

			handler.ServeHTTP(w, req)
			return
		}
	}
	http.NotFound(w, req)
}

// Helper methods for HTTP verbs
func (r *Router) GET(pattern string, handler http.HandlerFunc, middlewares ...func(http.Handler) http.Handler) {
	r.AddRoute("GET", pattern, handler, middlewares...)
}

func (r *Router) POST(pattern string, handler http.HandlerFunc, middlewares ...func(http.Handler) http.Handler) {
	r.AddRoute("POST", pattern, handler, middlewares...)
}

func (r *Router) PUT(pattern string, handler http.HandlerFunc, middlewares ...func(http.Handler) http.Handler) {
	r.AddRoute("PUT", pattern, handler, middlewares...)
}

func (r *Router) DELETE(pattern string, handler http.HandlerFunc, middlewares ...func(http.Handler) http.Handler) {
	r.AddRoute("DELETE", pattern, handler, middlewares...)
}

func (r *Router) PATCH(pattern string, handler http.HandlerFunc, middlewares ...func(http.Handler) http.Handler) {
	r.AddRoute("PATCH", pattern, handler, middlewares...)
}

// Group allows grouping routes with common prefix and middlewares
func (r *Router) Group(prefix string, middlewares ...func(http.Handler) http.Handler) *RouteGroup {
	return &RouteGroup{
		router:      r,
		prefix:      prefix,
		middlewares: middlewares,
	}
}

type RouteGroup struct {
	router      *Router
	prefix      string
	middlewares []func(http.Handler) http.Handler
}

func (rg *RouteGroup) GET(pattern string, handler http.HandlerFunc, middlewares ...func(http.Handler) http.Handler) {
	allMiddlewares := append(rg.middlewares, middlewares...)
	rg.router.AddRoute("GET", rg.prefix+pattern, handler, allMiddlewares...)
}

func (rg *RouteGroup) POST(pattern string, handler http.HandlerFunc, middlewares ...func(http.Handler) http.Handler) {
	allMiddlewares := append(rg.middlewares, middlewares...)
	rg.router.AddRoute("POST", rg.prefix+pattern, handler, allMiddlewares...)
}

func (rg *RouteGroup) PUT(pattern string, handler http.HandlerFunc, middlewares ...func(http.Handler) http.Handler) {
	allMiddlewares := append(rg.middlewares, middlewares...)
	rg.router.AddRoute("PUT", rg.prefix+pattern, handler, allMiddlewares...)
}

func (rg *RouteGroup) DELETE(pattern string, handler http.HandlerFunc, middlewares ...func(http.Handler) http.Handler) {
	allMiddlewares := append(rg.middlewares, middlewares...)
	rg.router.AddRoute("DELETE", rg.prefix+pattern, handler, allMiddlewares...)
}

// Helper function to get parameters from request context
func GetParams(r *http.Request) map[string]string {
	if params, ok := r.Context().Value(ParamsKey).(map[string]string); ok {
		return params
	}
	return make(map[string]string)
}

// Helper function to get a specific parameter
func GetParam(r *http.Request, name string) string {
	params := GetParams(r)
	return params[name]
}

// Static file serving
func (r *Router) Static(prefix, dir string) {
	fileServer := http.FileServer(http.Dir(dir))
	r.GET(prefix+"/*", func(w http.ResponseWriter, req *http.Request) {
		req.URL.Path = strings.TrimPrefix(req.URL.Path, prefix)
		fileServer.ServeHTTP(w, req)
	})
}
