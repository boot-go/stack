/*
 * Copyright (c) 2021 boot-go
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

/*
 * Copyright (c) 2021 boot-go
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

package httpcmp

import (
	"context"
	"github.com/boot-go/boot"
	"net/http"
)

type Server interface {
	// Middleware
	Use(middlewares ...func(http.Handler) http.Handler)
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
	// Server control
	Shutdown()
}

func (s *server) Use(middlewares ...func(http.Handler) http.Handler) {
	for _, middleware := range middlewares {
		boot.Logger.Debug.Printf("Attaching middleware " + boot.QualifiedName(middleware))
	}
	s.router.Use(middlewares...)
}
func (s *server) Handle(pattern string, handler http.Handler) {
	boot.Logger.Debug.Printf("Attaching handler %s at %s", boot.QualifiedName(handler), pattern)
	s.router.Handle(pattern, handler)
}

func (s *server) HandleFunc(pattern string, handlerFunc http.HandlerFunc) {
	boot.Logger.Debug.Printf("Attaching handler function %s at %s", boot.QualifiedName(handlerFunc), pattern)
	s.router.HandleFunc(pattern, handlerFunc)
}

// HTTP-method routing along `pattern`
func (s *server) Connect(pattern string, handlerFunc http.HandlerFunc) {
	boot.Logger.Debug.Printf("Attaching <Connect> handler %s at %s", boot.QualifiedName(handlerFunc), pattern)
	s.router.Connect(pattern, handlerFunc)
}
func (s *server) Delete(pattern string, handlerFunc http.HandlerFunc) {
	boot.Logger.Debug.Printf("Attaching <Delete> handler %s at %s", boot.QualifiedName(handlerFunc), pattern)
	s.router.Delete(pattern, handlerFunc)
}
func (s *server) Get(pattern string, handlerFunc http.HandlerFunc) {
	boot.Logger.Debug.Printf("Attaching <Get> handler %s at %s", boot.QualifiedName(handlerFunc), pattern)
	s.router.Get(pattern, handlerFunc)
}
func (s *server) Head(pattern string, handlerFunc http.HandlerFunc) {
	boot.Logger.Debug.Printf("Attaching <Head> handler %s at %s", boot.QualifiedName(handlerFunc), pattern)
	s.router.Head(pattern, handlerFunc)
}
func (s *server) Options(pattern string, handlerFunc http.HandlerFunc) {
	boot.Logger.Debug.Printf("Attaching <Options> handler %s at %s", boot.QualifiedName(handlerFunc), pattern)
	s.router.Options(pattern, handlerFunc)
}
func (s *server) Patch(pattern string, handlerFunc http.HandlerFunc) {
	boot.Logger.Debug.Printf("Attaching <Patch> handler %s at %s", boot.QualifiedName(handlerFunc), pattern)
	s.router.Patch(pattern, handlerFunc)
}
func (s *server) Post(pattern string, handlerFunc http.HandlerFunc) {
	boot.Logger.Debug.Printf("Attaching <Post> handler %s at %s", boot.QualifiedName(handlerFunc), pattern)
	s.router.Post(pattern, handlerFunc)
}
func (s *server) Put(pattern string, handlerFunc http.HandlerFunc) {
	boot.Logger.Debug.Printf("Attaching <Put> handler %s at %s", boot.QualifiedName(handlerFunc), pattern)
	s.router.Put(pattern, handlerFunc)
}
func (s *server) Trace(pattern string, handlerFunc http.HandlerFunc) {
	boot.Logger.Debug.Printf("Attaching <Trace> handler %s at %s", boot.QualifiedName(handlerFunc), pattern)
	s.router.Trace(pattern, handlerFunc)
}

// Shutdown gracefully shuts down the server
func (s *server) Shutdown() {
	go func() {
		s.state = ServerShuttingDown
		ctx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer cancel()
		s.shutdown <- s.httpServer.Shutdown(ctx)
	}()
}
