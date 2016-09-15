package logging

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/robjporter/GO-Color"
	"github.com/zenazn/goji/web"
)

// Logger is a middleware handler that logs the request as it goes in and the response as it goes out.
type Logger struct {
	// Logger inherits from log.Logger used to log messages with the Logger middleware
	*log.Logger
}

// Negroni compatible interface
func ServeHTTP(w http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	start := time.Now()
	ww := &StatusTrackingResponseWriter{w, http.StatusOK}
	next(ww, r)

	var remoteAddr string
	fwd := r.Header.Get("X-Forwarded-For")
	if fwd == "" {
		remoteAddr = r.RemoteAddr
	} else {
		remoteAddr = fwd + ":" + r.Header.Get("X-Forwarded-Port")
	}

	id := r.Header.Get("X-Request-Id")

	//log.Printf("%s | %s | %s | %s | %s | %s | %s", start.Format("02/01/2006 - 15:04:05.999999999"), codeColor(ww.Status), constantLength(time.Since(start).String(), 14), remoteAddr, methodColor(r.Method), r.Method, r.RequestURI)
	fmt.Printf("%s | %s | %s | %s | %s | %s | %s | %s\n", constantLengthRight(start.Format("02/01/2006 - 15:04:05.999999999"), 32), id, codeColor(ww.Status), constantLength(time.Since(start).String(), 14), constantLengthRight(remoteAddr, 24), methodColor(r.Method), r.Method, r.RequestURI)
}

func LoggingMiddleWare(c *web.C, h http.Handler) http.Handler {
	handler := func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		ww := &StatusTrackingResponseWriter{w, http.StatusOK}
		h.ServeHTTP(ww, r)

		var remoteAddr string
		fwd := r.Header.Get("X-Forwarded-For")
		if fwd == "" {
			remoteAddr = r.RemoteAddr
		} else {
			remoteAddr = fwd + ":" + r.Header.Get("X-Forwarded-Port")
		}

		log.Printf("%s | %s | %s | %s | %s | %s | %s", start.Format("02/01/2006 - 15:04:05.999999999"), codeColor(ww.Status), constantLength(time.Since(start).String(), 14), remoteAddr, methodColor(r.Method), r.Method, r.RequestURI)
	}
	return http.HandlerFunc(handler)
}

func codeColor(code2 int) string {
	code := strconv.Itoa(code2)
	if code == "404" {
		return color.Red(" " + code + " ")
	}
	if code == "301" {
		return color.Yellow(" " + code + " ")
	}
	if code == "304" {
		return color.Magenta(" " + code + " ")
	}
	if code == "200" {
		return color.Green(" " + code + " ")
	}
	return " " + code + " "
}

func methodColor(method string) string {
	if method == "GET" {
		return color.BlueBg("  ")
	}
	if method == "POST" {
		return color.GreenBg("  ")
	}
	if method == "DELETE" {
		return color.RedBg("  ")
	}
	if method == "HEAD" {
		return color.YellowBg("  ")
	}
	return "  "
}

func constantLengthRight(value string, length int) string {
	if len(value) < length {
		remains := length - len(value)
		space := ""
		for i := 0; i < remains; i++ {
			space = space + " "
		}
		toReturn := value + space
		return toReturn
	} else {
		return value[:10]
	}
}

func constantLength(value string, length int) string {
	if len(value) < length {
		remains := length - len(value)
		space := ""
		for i := 0; i < remains; i++ {
			space = space + " "
		}
		toReturn := space + value
		return toReturn
	} else {
		return value[:10]
	}
}

type StatusTrackingResponseWriter struct {
	http.ResponseWriter
	// http status code written
	Status int
}

func (w *StatusTrackingResponseWriter) WriteHeader(status int) {
	w.Status = status
	w.ResponseWriter.WriteHeader(status)
}

//[GIN] 2016/01/31 - 16:02:24 | 404 |        9.42µs | 127.0.0.1 |   GET     /
