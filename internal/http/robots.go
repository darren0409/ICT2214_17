package http

import (
	"fmt"
	"strings"

	"github.com/valyala/fasthttp"

	"github.com/yunginnanet/HellPot/internal/config"
)

var logMessageChannel = make(chan map[string]interface{})
var null = map[string]interface{}{}

func robotsTXT(ctx *fasthttp.RequestCtx) {
	slog := log.With().
		Str("USERAGENT", string(ctx.UserAgent())).
		Str("REMOTE_ADDR", getRealRemote(ctx)).
		Interface("URL", string(ctx.RequestURI())).Logger()
	paths := &strings.Builder{}
	paths.WriteString("User-agent: *\r\n")
	for _, p := range config.Paths {
		paths.WriteString("Disallow: ")
		paths.WriteString(p)
		paths.WriteString("\r\n")
	}
	paths.WriteString("\r\n")

	slog.Debug().
		Strs("PATHS", config.Paths).
		Msg("SERVE_ROBOTS")

	if _, err := fmt.Fprintf(ctx, paths.String()); err != nil {
		slog.Error().Err(err).Msg("SERVE_ROBOTS_ERROR")
	}

	// Construct log message with relevant information
	logMessage := map[string]interface{}{
		"Event":      "New request",
		"User Agent": string(ctx.UserAgent()),
		"IP address": getRealRemote(ctx),
		"URL":        string(ctx.RequestURI()),
	}
	logMessageChannel <- null
	logMessageChannel <- logMessage
}
