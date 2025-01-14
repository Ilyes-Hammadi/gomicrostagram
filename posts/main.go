package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/Ilyes-Hammadi/gomicrostagram/posts/services"
	"github.com/improbable-eng/grpc-web/go/grpcweb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
)

var serverAddress = os.Getenv("POSTS_SERVER_ADDRESS")
var serverPort = os.Getenv("POSTS_SERVER_PORT")

var host = fmt.Sprintf("%s:%s", serverAddress, serverPort)

func startServer() {
	grpcServer := grpc.NewServer()
	services.RegisterServices(grpcServer)

	grpclog.SetLogger(log.New(os.Stdout, "Posts service: ", log.LstdFlags))
	grpclog.Printf("Starting server. http port: %s", host)

	wrappedServer := grpcweb.WrapServer(grpcServer)
	handler := func(resp http.ResponseWriter, req *http.Request) {
		wrappedServer.ServeHTTP(resp, req)
	}

	httpServer := http.Server{
		Addr:    host,
		Handler: http.HandlerFunc(handler),
	}

	httpErr := httpServer.ListenAndServe()
	if httpErr != nil {
		log.Fatalf("Failed to listen: %v", httpErr)
	}
}

func main() {
	startServer()
}
