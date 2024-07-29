package middleware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"go-rest-api/internal/config"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"io"
	"net/http"
	"time"
)

type ResponseWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (rw *ResponseWriter) Write(b []byte) (int, error) {
	rw.body.Write(b)
	return rw.ResponseWriter.Write(b)
}

func (rw *ResponseWriter) WriteString(s string) (int, error) {
	rw.body.WriteString(s)
	return rw.ResponseWriter.WriteString(s)
}

func LoggerHandler() gin.HandlerFunc {
	logger := config.Logger

	return func(c *gin.Context) {
		start := time.Now()

		// Read the requests body
		bodyBytes, err := io.ReadAll(c.Request.Body)
		if err != nil {
			logger.Error("Failed to read requests body", zap.Error(err))
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Internal Server Error"})
			return
		}

		// Reset the requests body so it can be read again
		c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

		var bodyMap map[string]interface{}
		err = json.Unmarshal(bodyBytes, &bodyMap)
		if err != nil {
			fmt.Println("Error unmarshalling JSON:", err)
			return
		}

		// List of keys to replace with empty string
		keysToReplace := []string{"password", "token"}

		// Replacing the values with an empty string if the keys exist
		for _, key := range keysToReplace {
			if _, exists := bodyMap[key]; exists {
				bodyMap[key] = ""
			}
		}

		modifiedBodyBytes, err := json.Marshal(bodyMap)
		if err != nil {
			fmt.Println("Error marshalling JSON:", err)
			return
		}

		// Wrap the response writer to capture the response body
		rw := &ResponseWriter{ResponseWriter: c.Writer, body: &bytes.Buffer{}}
		c.Writer = rw

		// Process requests
		c.Next()

		// Getting requests and response details
		req := c.Request
		res := c.Writer

		id := req.Header.Get("X-Request-ID")
		if id == "" {
			id = res.Header().Get("X-Request-ID")
		}

		fields := []zapcore.Field{
			zap.String("remote_ip", c.ClientIP()),
			zap.String("latency", time.Since(start).String()),
			zap.String("host", req.Host),
			zap.String("requests", fmt.Sprintf("%s %s", req.Method, req.RequestURI)),
			zap.Int("status", res.Status()),
			zap.Int("size", res.Size()),
			zap.String("user_agent", req.UserAgent()),
			zap.String("request_body", string(modifiedBodyBytes)),
			zap.String("response_body", rw.body.String()),
			zap.String("request_id", id),
			//zap.Errors("errors", c.Errors.Errors()),
		}

		n := res.Status()
		switch {
		case n >= http.StatusInternalServerError:
			logger.Error("Server error", fields...)
		case n >= http.StatusBadRequest:
			logger.Warn("Client error", fields...)
		case n >= http.StatusMultipleChoices:
			logger.Info("Redirection", fields...)
		default:
			logger.Info("Success", fields...)
		}
	}
}
