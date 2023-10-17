package router

import (
	"encoding/json"
	"io"
	"log"
	"path/filepath"
	"testing"
)

type userAgent struct {
	UserAgent string  `json:"ua"`
	Percent   float32 `json:"pct"`
}

func loadUserAgents(filename string) []userAgent {
	contents, err := io.ReadAll(filename)
	if err != nil {
		log.Fatal(err)
	}
	var dst []userAgent
	if err = json.Unmarshal(contents, &dst); err != nil {
		log.Fatal(err)
	}
	return dst
}

func TestDetectDesktopBrowsers(t *testing.T) {
	uas := loadUserAgents(filepath.Join("testdata", "desktop_user_agents.json"))
	for _, ua := range uas {
		result := isBrowser([]byte(ua.UserAgent))
		if result == false {
			t.Errorf(`Failure to detect this desktop user agent (with %0.2f%% usage) as a browser-type user agent: %q`, ua.Percent, ua.UserAgent)
		}
	}
}

func TestDetectMobileBrowsers(t *testing.T) {
	uas := loadUserAgents(filepath.Join("testdata", "mobile_user_agents.json"))
	for _, ua := range uas {
		result := isBrowser([]byte(ua.UserAgent))
		if result == false {
			t.Errorf(`Failure to detect this mobile user agent (with %0.2f%% usage) as a browser-type user agent: %q`, ua.Percent, ua.UserAgent)
		}
	}
}

func TestDetectCurl(t *testing.T) {
	curlUserAgent := "curl/7.88.1"
	result := isCurlLike([]byte(curlUserAgent))
	if result == false {
		t.Fatalf("Fail to detect curl user agent")
	}
}
