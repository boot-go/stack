/*
 * Copyright (c) 2021-2023 boot-go
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all
 * copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 *
 */

package chi

import (
	"context"
	"net/http"

	"github.com/boot-go/boot"
	"github.com/go-chi/chi/v5"
)

type Server interface {
	// Middleware on Router
	Use(middlewares ...func(http.Handler) http.Handler)
	// Middleware for Endpointt Handler
	With(middlewares ...func(http.Handler) http.Handler) chi.Router
	// Generic routing
	Handle(pattern string, handler http.Handler)
	HandleFunc(pattern string, handlerFunc http.HandlerFunc)
	// HTTP-method routing`
	Connect(pattern string, handlerFunc http.HandlerFunc)
	Delete(pattern string, handlerFunc http.HandlerFunc)
	Get(pattern string, handlerFunc http.HandlerFunc)
	Head(pattern string, handlerFunc http.HandlerFunc)
	Options(pattern string, handlerFunc http.HandlerFunc)
	Patch(pattern string, handlerFunc http.HandlerFunc)
	Post(pattern string, handlerFunc http.HandlerFunc)
	Put(pattern string, handlerFunc http.HandlerFunc)
	Trace(pattern string, handlerFunc http.HandlerFunc)
	// Route
	Route(pattern string, fn func(r chi.Router)) chi.Router
	// Group
	Group(fn func(r chi.Router)) chi.Router
	// Mount
	Mount(pattern string, handlerFunc http.Handler)
	// Method
	Method(method, pattern string, handler http.Handler)
	MethodFunc(method, pattern string, handlerFunc http.HandlerFunc)
	// NotFound
	NotFound(handlerFunc http.HandlerFunc)
	// NotAllowed
	MethodNotAllowed(handlerFunc http.HandlerFunc)
	// Router Traversal
	Routes() []chi.Route
	Middlewares() chi.Middlewares
	Match(ctx *chi.Context, method, path string) bool
	// Server control
	Shutdown()
}

func (s *server) Match(ctx *chi.Context, method, path string) bool {
	boot.Logger.Debug.Printf("match for method %s on path %s", method, path)
	return s.router.Match(ctx, method, path)
}

func (s *server) Routes() []chi.Route {
	boot.Logger.Debug.Printf("returning all routes")
	return s.router.Routes()
}

func (s *server) Middlewares() chi.Middlewares {
	boot.Logger.Debug.Printf("returning all middlewares")
	return s.router.Middlewares()
}

func (s *server) With(middlewares ...func(http.Handler) http.Handler) chi.Router {
	for _, middleware := range middlewares {
		boot.Logger.Debug.Printf("attaching middleware %s\n", boot.QualifiedName(middleware))
	}
	return s.router.With(middlewares...)
}

func (s *server) Group(fn func(r chi.Router)) chi.Router {
	boot.Logger.Debug.Printf("group - new inline router along current routing path with new middlerware")
	return s.router.Group(fn)
}

func (s *server) Mount(pattern string, handler http.Handler) {
	boot.Logger.Debug.Printf("mount handler %s at %s", boot.QualifiedName(handler), pattern)
	s.router.Mount(pattern, handler)
}

func (s *server) Method(method, pattern string, handler http.Handler) {
	boot.Logger.Debug.Printf("method %s handler %s at %s", method, boot.QualifiedName(handler), pattern)
	s.router.Method(method, pattern, handler)
}

func (s *server) MethodFunc(method, pattern string, handlerFunc http.HandlerFunc) {
	boot.Logger.Debug.Printf("method %s handlerFunc %s at %s", method, boot.QualifiedName(handlerFunc), pattern)
	s.router.MethodFunc(method, pattern, handlerFunc)
}

func (s *server) MethodNotAllowed(handlerFunc http.HandlerFunc) {
	boot.Logger.Debug.Printf("method not allowed handlerFunc %s", boot.QualifiedName(handlerFunc))
	s.router.MethodNotAllowed(handlerFunc)
}

func (s *server) NotFound(handlerFunc http.HandlerFunc) {
	boot.Logger.Debug.Printf("not found handlerFunc %s", boot.QualifiedName(handlerFunc))
	s.router.NotFound(handlerFunc)
}

func (s *server) Route(pattern string, fn func(r chi.Router)) chi.Router {
	boot.Logger.Debug.Printf("attaching route %s at %s", boot.QualifiedName(fn), pattern)
	return s.router.Route(pattern, fn)
}

func (s *server) Use(middlewares ...func(http.Handler) http.Handler) {
	for _, middleware := range middlewares {
		boot.Logger.Debug.Printf("attaching middleware %s\n", boot.QualifiedName(middleware))
	}
	s.router.Use(middlewares...)
}

func (s *server) Handle(pattern string, handler http.Handler) {
	boot.Logger.Debug.Printf("attaching handler %s at %s", boot.QualifiedName(handler), pattern)
	s.router.Handle(pattern, handler)
}

func (s *server) HandleFunc(pattern string, handlerFunc http.HandlerFunc) {
	boot.Logger.Debug.Printf("attaching handler function %s at %s", boot.QualifiedName(handlerFunc), pattern)
	s.router.HandleFunc(pattern, handlerFunc)
}

// HTTP-method routing along `pattern`
func (s *server) Connect(pattern string, handlerFunc http.HandlerFunc) {
	boot.Logger.Debug.Printf("attaching <Connect> handler %s at %s", boot.QualifiedName(handlerFunc), pattern)
	s.router.Connect(pattern, handlerFunc)
}

func (s *server) Delete(pattern string, handlerFunc http.HandlerFunc) {
	boot.Logger.Debug.Printf("attaching <Delete> handler %s at %s", boot.QualifiedName(handlerFunc), pattern)
	s.router.Delete(pattern, handlerFunc)
}

func (s *server) Get(pattern string, handlerFunc http.HandlerFunc) {
	boot.Logger.Debug.Printf("attaching <Get> handler %s at %s", boot.QualifiedName(handlerFunc), pattern)
	s.router.Get(pattern, handlerFunc)
}

func (s *server) Head(pattern string, handlerFunc http.HandlerFunc) {
	boot.Logger.Debug.Printf("attaching <Head> handler %s at %s", boot.QualifiedName(handlerFunc), pattern)
	s.router.Head(pattern, handlerFunc)
}

func (s *server) Options(pattern string, handlerFunc http.HandlerFunc) {
	boot.Logger.Debug.Printf("attaching <Options> handler %s at %s", boot.QualifiedName(handlerFunc), pattern)
	s.router.Options(pattern, handlerFunc)
}

func (s *server) Patch(pattern string, handlerFunc http.HandlerFunc) {
	boot.Logger.Debug.Printf("attaching <Patch> handler %s at %s", boot.QualifiedName(handlerFunc), pattern)
	s.router.Patch(pattern, handlerFunc)
}

func (s *server) Post(pattern string, handlerFunc http.HandlerFunc) {
	boot.Logger.Debug.Printf("attaching <Post> handler %s at %s", boot.QualifiedName(handlerFunc), pattern)
	s.router.Post(pattern, handlerFunc)
}

func (s *server) Put(pattern string, handlerFunc http.HandlerFunc) {
	boot.Logger.Debug.Printf("attaching <Put> handler %s at %s", boot.QualifiedName(handlerFunc), pattern)
	s.router.Put(pattern, handlerFunc)
}

func (s *server) Trace(pattern string, handlerFunc http.HandlerFunc) {
	boot.Logger.Debug.Printf("attaching <Trace> handler %s at %s", boot.QualifiedName(handlerFunc), pattern)
	s.router.Trace(pattern, handlerFunc)
}

// Shutdown gracefully shuts down the net
func (s *server) Shutdown() {
	go func() {
		s.state = ServerShuttingDown
		ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer cancel()
		s.shutdown <- s.httpServer.Shutdown(ctx)
	}()
}
