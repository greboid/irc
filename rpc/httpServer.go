package rpc

import (
	"context"
	"errors"
	"fmt"
	grpcauth "github.com/grpc-ecosystem/go-grpc-middleware/auth"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"
)

type httpServer struct {
	WebPort int
	plugins []Plugin
	pathMap map[string]*descriptor
}

type descriptor struct {
	prefix string
	token string
	stream *HTTPPlugin_GetRequestServer
	receive chan *HttpResponse
}

func NewHttpServer(port int, plugin []Plugin) *httpServer {
	return &httpServer{
		WebPort: port,
		plugins: plugin,
		pathMap: make(map[string]*descriptor),
	}

}

func (h *httpServer) Start() {
	go func() {
		mux := http.NewServeMux()
		mux.HandleFunc("/", h.handleRequest)
		server := http.Server{
			Addr:    fmt.Sprintf(":%d", h.WebPort),
			Handler: mux,
		}
		go func() {
			log.Print(server.ListenAndServe())
		}()
		stop := make(chan os.Signal, 1)
		signal.Notify(stop, os.Interrupt, os.Kill)
		log.Printf("Waiting for stop")
		<-stop
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := server.Shutdown(ctx); err != nil {
			log.Fatalf("Unable to shutdown: %s", err.Error())
		}
	}()
}

func (h *httpServer) authPlugin(ctx context.Context) (context.Context, error) {
	token, err := grpcauth.AuthFromMD(ctx, "bearer")
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "invalid auth token: %s", err.Error())
	}
	if !h.checkPlugin(token) {
		return nil, status.Errorf(codes.Unauthenticated, "access denied")
	}
	return ctx, nil
}

func (h *httpServer) checkPlugin(token string) bool {
	for _, plugin := range h.plugins {
		if plugin.Token == token {
			return true
		}
	}
	return false
}

func (h *httpServer) RegisterRoute(ctx context.Context, route *Route) (*Error, error) {
	token, _ := grpcauth.AuthFromMD(ctx, "bearer")
	for k := range h.pathMap {
		if route.Prefix == h.pathMap[k].prefix {
			return &Error{
				Message: "Prefix already registered",
			}, errors.New("prefix already registered")
		}
	}
	h.pathMap[token] = &descriptor{
		token:  token,
		prefix: route.Prefix,
		receive: make(chan *HttpResponse, 1),
	}
	log.Printf("%s registered %s", token, route.Prefix)
	return &Error{}, nil
}

func (h *httpServer) handleRequest(writer http.ResponseWriter, request *http.Request) {
	log.Printf("Request: %s", request.URL)
	for k := range h.pathMap {
		if strings.HasPrefix(request.URL.Path, fmt.Sprintf("/%s", h.pathMap[k].prefix)) {
			log.Printf("Request matvhes prefix")
			log.Printf("Desriptor: %+v", h.pathMap[k])
			stream := *h.pathMap[k].stream
			if stream != nil {
				body, err := ioutil.ReadAll(request.Body)
				if err != nil {
					writer.WriteHeader(http.StatusInternalServerError)
				}
				log.Printf("Sending request to plugin %s", h.pathMap[k].token)
				err = stream.Send(&HttpRequest{
					Header: ConvertToRPCHeaders(request.Header),
					Body: body,
				})
				if err != nil {
					log.Printf("Unable to send to plugin")
					return
				}
				select {
					case response := <-h.pathMap[k].receive:
						for index := range response.Header {
							writer.Header().Add(response.Header[index].Key, response.Header[index].Value)
						}
						writer.WriteHeader(int(response.Status))
						_, _ = writer.Write(response.Body)
					case <-time.After(1 * time.Second):
				}
			}
		}
	}
}

func (h *httpServer) GetRequest(stream HTTPPlugin_GetRequestServer) error {
	token, _ := grpcauth.AuthFromMD(stream.Context(), "bearer")
	descriptor, ok := h.pathMap[token]
	if !ok {
		return errors.New("plugin not registered")
	}
	descriptor.stream = &stream
	for {
		in, err := stream.Recv()
		if err == io.EOF {
			return nil
		}
		if err != nil {
			return err
		}
		descriptor.receive <- in
	}
}

func ConvertFromRPCHeaders(headers []*HttpHeader) http.Header {
	httpHeaders := http.Header{}
	for index := range headers {
		httpHeaders.Add(headers[index].Key, headers[index].Value)
	}
	return httpHeaders
}

func ConvertToRPCHeaders(headers http.Header) []*HttpHeader {
	rpcHeaders := make([]*HttpHeader, 0)
	for key, value := range headers {
		for index := range value {
			rpcHeaders = append(rpcHeaders, &HttpHeader{
				Key:                  key,
				Value:                value[index],
			})
		}

	}
	return rpcHeaders
}