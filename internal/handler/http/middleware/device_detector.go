package middleware

import (
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/mssola/user_agent"
)

type DeviceInfo struct {
	Platform       string
	Model          string
	Browser        string
	BrowserVersion string
	IsMobile       bool
	Source         string
	RawUA          string
	RawCH          string
}

func DeviceDetector() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		h := ctx.Request.Header
		ch := map[string]string{
			"Sec-CH-UA":                 h.Get("Sec-CH-UA"),
			"Sec-CH-UA-Mobile":          h.Get("Sec-CH-UA-Mobile"),
			"Sec-CH-UA-Platform":        h.Get("Sec-CH-UA-Platform"),
			"Sec-CH-UA-Model":           h.Get("Sec-CH-UA-Model"),
			"Sec-CH-UA-Browser":         h.Get("Sec-CH-UA-Browser"),
			"Sec-CH-UA-Browser-Version": h.Get("Sec-CH-UA-Browser-Version"),
		}

		if ch["Sec-CH-UA"] != "" || ch["Sec-CH-UA-Mobile"] != "" {
			info := &DeviceInfo{
				Platform:       strings.Trim(ch["Sec-CH-UA-Platform"], `"`),
				Model:          strings.Trim(ch["Sec-CH-UA-Model"], `"`),
				Browser:        strings.Trim(ch["Sec-CH-UA-Browser"], `"`),
				BrowserVersion: strings.Trim(ch["Sec-CH-UA-Browser-Version"], `"`),
				IsMobile:       ch["Sec-CH-UA-Mobile"] == "true",
				Source:         "client-hint",
				RawUA:          "",
				RawCH: strings.Join([]string{
					ch["Sec-CH-UA"],
					ch["Sec-CH-UA-Mobile"],
					ch["Sec-CH-UA-Platform"],
					ch["Sec-CH-UA-Model"],
					ch["Sec-CH-UA-Browser"],
					ch["Sec-CH-UA-Browser-Version"],
				}, "; "),
			}
			ctx.Set("deviceInfo", &info)
			ctx.Next()
		}

		uaStr := h.Get("User-Agent")
		ua := user_agent.New(uaStr)
		name, version := ua.Browser()
		platform := ""
		switch {
		case ua.Platform() != "":
			platform = ua.Platform()
		case ua.OS() != "":
			platform = ua.OS()
		}
		info := DeviceInfo{
			Platform:       platform,
			Model:          "",
			Browser:        name,
			BrowserVersion: version,
			IsMobile:       ua.Mobile(),
			Source:         "user-agent",
			RawUA:          uaStr,
			RawCH:          "",
		}
		ctx.Set("deviceInfo", &info)
		ctx.Next()
	}
}

// helper func
func GetDeviceInfo(c *gin.Context) *DeviceInfo {
	if v, ok := c.Get("deviceInfo"); ok {
		if di, ok := v.(*DeviceInfo); ok {
			return di
		}
	}
	return nil
}
