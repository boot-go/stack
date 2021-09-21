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
	"github.com/go-chi/chi"
	"net/http"
	"net/http/httptest"
	"strconv"
	"time"
)

type lifeState uint8

// server provides the default implementation using the httpcmp.server. Other components
// can register context paths to process certain requests.
type server struct {
	Eventbus   boot.EventBus `boot:"wire"`
	Runtime    boot.Runtime  `boot:"wire"`
	Port       int           `boot:"config,key:${HTTP_SERVER_PORT},default:8080"`
	router     chi.Router
	httpServer *http.Server
	testServer *httptest.Server
	// lifecycle
	shutdown chan error
	state    lifeState
}

const (
	shutdownTimeout = 5 * time.Second
)

// lifecycle states
const (
	ServerLive lifeState = iota
	ServerReady
	ServerShuttingDown
	ServerShuttedDown
)

func init() {
	boot.Register(func() boot.Component {
		return &server{}
	})
}

func (s *server) Init() {
	s.router = chi.NewRouter()
	if s.Runtime.HasFlag(boot.StandardFlag) {
		s.initHttpServer()
	} else if s.Runtime.HasFlag(boot.UnitTestFlag) {
		s.initTestServer()
	}
}

func (s *server) initHttpServer() {
	s.httpServer = &http.Server{
		Addr:    ":" + strconv.Itoa(s.Port),
		Handler: s.router,
	}
	s.Use(func(handler http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			boot.Logger.Debug.Printf("called request uri: %s", r.RequestURI)
			handler.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	})
	s.state = ServerReady
}

func (s *server) initTestServer() {
	s.testServer = httptest.NewUnstartedServer(s.router)
	s.state = ServerReady
}

func (s *server) Start() {
	if s.Runtime.HasFlag(boot.StandardFlag) {
		s.startHttpServer()
	} else if s.Runtime.HasFlag(boot.UnitTestFlag) {
		s.startTestServer()
	}
	// wait for shutdown signal then return
	<-s.shutdown
}

func (s *server) startHttpServer() {
	s.router.HandleFunc("/", logRequestHandler)
	s.Eventbus.Publish(InitializedEvent{})
	boot.Logger.Info.Printf("http server listening on %s", s.httpServer.Addr)
	s.shutdown = make(chan error)
	s.state = ServerLive
	s.httpServer.RegisterOnShutdown(func() {
		s.Eventbus.Publish(ShutDownInitiatedEvent{})
	})
	if err := s.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		boot.Logger.Error.Printf("http server closed unexpectedly: %s", err.Error())
	}
	s.httpServer.Close()
}

func (s *server) startTestServer() {
	s.Eventbus.Publish(InitializedEvent{})
	s.state = ServerLive
	s.testServer.Start()
}

func (s *server) Stop() {
	if s.Runtime.HasFlag(boot.StandardFlag) {
		s.stopHttpServer()
	} else if s.Runtime.HasFlag(boot.UnitTestFlag) {
		s.stopTestServer()
	}
}

func (s *server) stopHttpServer() {
	s.state = ServerShuttingDown
	if s.httpServer != nil {
		boot.Logger.Info.Printf("shutting down server")
		s.Eventbus.Publish(ShutDownInitiatedEvent{})
		s.httpServer.Shutdown(context.Background())
		s.Eventbus.Publish(ShutDownCompletedEvent{})
		s.httpServer.Close()
	} else {
		boot.Logger.Warn.Printf("server is not in shutdown mode! shutdown first before stopping it...")
	}
	s.state = ServerShuttedDown
}
func (s *server) stopTestServer() {
	s.state = ServerShuttingDown
	s.testServer.Close()
	s.state = ServerShuttedDown
}
