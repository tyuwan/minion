package minion

import (
	"net/http"
	"path"
)

// Router TODO
type Router struct {
	Handlers     []HandlerFunc
	absolutePath string
	engine       *Engine
}

// Use Adds middlewares to the group
func (c *Router) Use(middlewares ...HandlerFunc) {
	c.Handlers = append(c.Handlers, middlewares...)
}

// Group Creates a new router group. You should add all the routes that have common middlwares or the same path prefix.
// For example, all the routes that use a common middlware for authorization could be grouped.
func (c *Router) Group(relativePath string, fn func(*Router), handlers ...HandlerFunc) *Router {
	router := &Router{
		Handlers:     c.combineHandlers(handlers),
		absolutePath: c.calculateAbsolutePath(relativePath),
		engine:       c.engine,
	}
	fn(router)
	return router
}

// Handle registers a new request handle and middlewares with the given path and method.
// The last handler should be the real handler, the other ones should be middlewares that can and should be shared among different routes.
//
// For GET, POST, PUT, PATCH and DELETE requests the respective shortcut
// functions can be used.
//
// This function is intended for bulk loading and to allow the usage of less
// frequently used, non-standardized or custom methods (e.g. for internal
// communication with a proxy).
func (c *Router) Handle(httpMethod, relativePath string, handlers []HandlerFunc) {
	absolutePath := c.calculateAbsolutePath(relativePath)
	handlers = c.combineHandlers(handlers)
	c.engine.router.Get(absolutePath, func(w http.ResponseWriter, req *http.Request) {
		ctx := c.engine.createContext(w, req, handlers)
		ctx.Next()
		ctx.Writer.WriteHeaderNow()
		c.engine.reuseContext(ctx)
	})
}

// Post is a shortcut for router.Handle("POST", path, handle)
func (c *Router) Post(relativePath string, handlers ...HandlerFunc) {
	absolutePath := c.calculateAbsolutePath(relativePath)
	handlers = c.combineHandlers(handlers)
	c.engine.router.Post(absolutePath, func(w http.ResponseWriter, req *http.Request) {
		ctx := c.engine.createContext(w, req, handlers)
		ctx.Next()
		ctx.Writer.WriteHeaderNow()
		c.engine.reuseContext(ctx)
	})
}

// Get is a shortcut for router.Handle("GET", path, handle)
func (c *Router) Get(relativePath string, handlers ...HandlerFunc) {
	absolutePath := c.calculateAbsolutePath(relativePath)
	handlers = c.combineHandlers(handlers)
	c.engine.router.Get(absolutePath, func(w http.ResponseWriter, req *http.Request) {
		ctx := c.engine.createContext(w, req, handlers)
		ctx.Next()
		ctx.Writer.WriteHeaderNow()
		c.engine.reuseContext(ctx)
	})
}

// Delete is a shortcut for router.Handle("DELETE", path, handle)
func (c *Router) Delete(relativePath string, handlers ...HandlerFunc) {
	absolutePath := c.calculateAbsolutePath(relativePath)
	handlers = c.combineHandlers(handlers)
	c.engine.router.Delete(absolutePath, func(w http.ResponseWriter, req *http.Request) {
		ctx := c.engine.createContext(w, req, handlers)
		ctx.Next()
		ctx.Writer.WriteHeaderNow()
		c.engine.reuseContext(ctx)
	})
}

// Patch is a shortcut for router.Handle("PATCH", path, handle)
func (c *Router) Patch(relativePath string, handlers ...HandlerFunc) {
	absolutePath := c.calculateAbsolutePath(relativePath)
	handlers = c.combineHandlers(handlers)
	c.engine.router.Patch(absolutePath, func(w http.ResponseWriter, req *http.Request) {
		ctx := c.engine.createContext(w, req, handlers)
		ctx.Next()
		ctx.Writer.WriteHeaderNow()
		c.engine.reuseContext(ctx)
	})
}

// Put is a shortcut for router.Handle("PUT", path, handle)
func (c *Router) Put(relativePath string, handlers ...HandlerFunc) {
	absolutePath := c.calculateAbsolutePath(relativePath)
	handlers = c.combineHandlers(handlers)
	c.engine.router.Put(absolutePath, func(w http.ResponseWriter, req *http.Request) {
		ctx := c.engine.createContext(w, req, handlers)
		ctx.Next()
		ctx.Writer.WriteHeaderNow()
		c.engine.reuseContext(ctx)
	})
}

// Options is a shortcut for router.Handle("OPTIONS", path, handle)
func (c *Router) Options(relativePath string, handlers ...HandlerFunc) {
	absolutePath := c.calculateAbsolutePath(relativePath)
	handlers = c.combineHandlers(handlers)
	c.engine.router.Options(absolutePath, func(w http.ResponseWriter, req *http.Request) {
		ctx := c.engine.createContext(w, req, handlers)
		ctx.Next()
		ctx.Writer.WriteHeaderNow()
		c.engine.reuseContext(ctx)
	})
}

// Head is a shortcut for router.Handle("HEAD", path, handle)
func (c *Router) Head(relativePath string, handlers ...HandlerFunc) {
	absolutePath := c.calculateAbsolutePath(relativePath)
	handlers = c.combineHandlers(handlers)
	c.engine.router.Head(absolutePath, func(w http.ResponseWriter, req *http.Request) {
		ctx := c.engine.createContext(w, req, handlers)
		ctx.Next()
		ctx.Writer.WriteHeaderNow()
		c.engine.reuseContext(ctx)
	})
}

// Static serves files from the given file system root.
// use : router.Static("/static", "/var/www")
func (c *Router) Static(path, dir string) {
	c.engine.router.FileServer(path, http.Dir(dir))
}

func (c *Router) combineHandlers(handlers []HandlerFunc) []HandlerFunc {
	finalSize := len(c.Handlers) + len(handlers)
	mergedHandlers := make([]HandlerFunc, 0, finalSize)
	mergedHandlers = append(mergedHandlers, c.Handlers...)
	return append(mergedHandlers, handlers...)
}

func (c *Router) calculateAbsolutePath(relativePath string) string {
	if len(relativePath) == 0 {
		return c.absolutePath
	}
	absolutePath := path.Join(c.absolutePath, relativePath)
	appendSlash := lastChar(relativePath) == '/' && lastChar(absolutePath) != '/'
	if appendSlash {
		return absolutePath + "/"
	}
	return absolutePath
}