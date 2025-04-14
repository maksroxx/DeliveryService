package handlers

import (
	"net/http"
	"strings"

	"github.com/maksroxx/DeliveryService/gateway/internal/utils"
	"github.com/sirupsen/logrus"
)

type RouteConfig struct {
	Prefix      string
	TargetURL   string
	PathRewrite string
}

type Router struct {
	routes []RouteConfig
	logger *logrus.Logger
}

func NewRouter(routes []RouteConfig, logger *logrus.Logger) *Router {
	return &Router{
		routes: routes,
		logger: logger,
	}
}

func (rt *Router) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	for _, route := range rt.routes {
		if strings.HasPrefix(r.URL.Path, route.Prefix) {
			newPath := strings.Replace(r.URL.Path, route.Prefix, route.PathRewrite, 1)
			targetURL := strings.TrimSuffix(route.TargetURL, "/") + newPath

			if r.URL.RawQuery != "" {
				targetURL += "?" + r.URL.RawQuery
			}

			rt.logger.Debugf("Routing %s -> %s", r.URL.Path, targetURL)

			if err := utils.ProxyRequest(w, r, targetURL); err != nil {
				rt.logger.Errorf("Proxy error: %v", err)
				http.Error(w, "Internal server error", http.StatusInternalServerError)
			}
			return
		}
	}

	http.NotFound(w, r)
}
