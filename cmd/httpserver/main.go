package main

import (
	"crypto/sha256"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/rpstvs/httpfromtcp/internal/headers"
	"github.com/rpstvs/httpfromtcp/internal/request"
	"github.com/rpstvs/httpfromtcp/internal/response"
	"github.com/rpstvs/httpfromtcp/internal/server"
)

const port = 42069

func main() {
	server, err := server.Serve(port, handler)
	if err != nil {
		log.Fatalf("Error starting server %v", err)
	}

	defer server.Close()

	log.Println("Server started on port", port)

	sigChan := make(chan os.Signal, 1)

	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	<-sigChan
	log.Println("Server gracefully closed")
}

func handler(w *response.Writer, req *request.Request) {
	if strings.HasPrefix(req.RequestLine.RequestTarget, "/httpbin") {
		proxyHandler(w, req)
		return
	}
	if req.RequestLine.RequestTarget == "/yourproblem" {
		handler200(w, req)
		return
	}
	if req.RequestLine.RequestTarget == "/myproblem" {
		handler500(w, req)
		return
	}
	if req.RequestLine.RequestTarget == "/video" {
		handlervideo(w, req)
		return
	}
	handler200(w, req)
	return
}

func handler400(w *response.Writer, _ *request.Request) {
	w.WriteStatusLine(response.StatusBadRequest)
	body := []byte(`<html>
<head>
<title>400 Bad Request</title>
</head>
<body>
<h1>Bad Request</h1>
<p>Your request honestly kinda sucked.</p>
</body>
</html>
`)

	h := response.GetDefaultHeaders(len(body))
	h.Override("Content-Type", "text/html")
	w.WriteHeaders(h)
	w.WriteBody(body)
	return
}

func handler500(w *response.Writer, _ *request.Request) {
	w.WriteStatusLine(response.StatusInternalServerError)
	body := []byte(`<html>
<head>
<title>500 Internal Server Error</title>
</head>
<body>
<h1>Internal Server Error</h1>
<p>Okay, you know what? This one is on me.</p>
</body>
</html>
`)

	h := response.GetDefaultHeaders(len(body))
	h.Override("Content-Type", "text/html")
	w.WriteHeaders(h)
	w.WriteBody(body)
	return
}

func handler200(w *response.Writer, _ *request.Request) {
	w.WriteStatusLine(response.StatusOK)
	body := []byte(`<html>
<head>
<title>200 OK</title>
</head>
<body>
<h1>Success!</h1>
<p>Your request was an absolute banger.</p>
</body>
</html>
`)

	h := response.GetDefaultHeaders(len(body))
	h.Override("Content-Type", "text/html")
	w.WriteHeaders(h)
	w.WriteBody(body)
	return
}

func proxyHandler(w *response.Writer, r *request.Request) {
	target := strings.TrimPrefix(r.RequestLine.RequestTarget, "/httpbin/")
	url := "https://httpbin.org/" + target
	fmt.Println("Proxying to", url)
	resp, err := http.Get(url)

	if err != nil {
		handler500(w, r)
		return
	}

	defer resp.Body.Close()
	// Writing Headers
	w.WriteStatusLine(response.StatusOK)
	h := response.GetDefaultHeaders(0)
	h.Override("Transfer-Encoding", "chunked")
	h.Override("Trailer", "X-Content-SHA256, X-Content-Length")
	h.Remove("Content-Length")
	w.WriteHeaders(h)

	const maxChunksize = 1024
	buffer := make([]byte, maxChunksize)
	fullBody := make([]byte, 0)
	for {
		n, err := resp.Body.Read(buffer)

		fmt.Println("Read", n, "bytes")

		if n > 0 {
			_, err := w.WriteChunkedBody(buffer[:n])

			if err != nil {
				fmt.Println("Error writing chunked body", err)
				break
			}
			fullBody = append(fullBody, buffer[:n]...)
		}
		if err == io.EOF {
			break
		}

		if err != nil {
			fmt.Println("entrei aqui porque", err)
			fmt.Println("Error reading response body", err)
			break

		}
	}

	_, err = w.WriteChunkedBodyDone()

	if err != nil {
		fmt.Println("Error writing chunked body done", err)
	}

	trailers := headers.NewHeaders()

	sha256 := fmt.Sprintf("%x", sha256.Sum256(fullBody))

	trailers.Override("X-Content-SHA256", sha256)
	trailers.Override("X-Content-Length", fmt.Sprintf("%d", len(fullBody)))

	err = w.WriteTrailers(trailers)

	if err != nil {
		fmt.Println("couldnt write trailers", err)
	}

	fmt.Println("Wrote trailers")

}

func handlervideo(w *response.Writer, _ *request.Request) {
	w.WriteStatusLine(response.StatusOK)
	vidPath := "./assets/vim.mp4"
	videoByte, err := os.ReadFile(vidPath)
	if err != nil {
		handler500(w, nil)
	}
	h := response.GetDefaultHeaders(len(videoByte))
	h.Override("Content-Type", "video/mp4")
	w.WriteHeaders(h)
	w.WriteBody(videoByte)
	return
}
