package router

import (
	"net/http"
	"regexp"
)

var (
	curlUserAgentRegexp    *regexp.Regexp
	wgetUserAgentRegexp    *regexp.Regexp
	httpieUserAgentRegexp  *regexp.Regexp
	mozillaUserAgentRegexp *regexp.Regexp
	safariUserAgentRegexp  *regexp.Regexp
	chromeUserAgentRegexp  *regexp.Regexp
	operaUserAgentRegexp   *regexp.Regexp
)

func init() {
	// To verify:
	// ```
	//   % curl -sv http://checkip.amazonaws.com -o/dev/null/ 2>&1 | grep -i 'user-agent'
	// ```
	curlUserAgentRegexp = regexp.MustCompile(`^curl/[\d\.]+$`)
	httpieUserAgentRegexp = regexp.MustCompile(`^HTTPie/[\d\.]+$`)
	// To verify:
	// ```
	//   % wget -d http://checkip.amazonaws.com -O/dev/null 2>&1 | grep -i '^user-agent'
	// ```
	wgetUserAgentRegexp = regexp.MustCompile(`^Wget/[\d+\.]+$`)

	// See: https://developer.mozilla.org/en-US/docs/Web/HTTP/Browser_detection_using_the_user_agent
	mozillaUserAgentRegexp = regexp.MustCompile(`.*(Firefox|Seamonkey)/.*`)
	safariUserAgentRegexp = regexp.MustCompile(`.*Safari/.*`)
	chromeUserAgentRegexp = regexp.MustCompile(`.*Chrome?(ium)?/.*`)
	operaUserAgentRegexp = regexp.MustCompile(`.*(OPR|Opera)/.*`)
}

func userAgentIsCurlLike(req *http.Request) bool {
	var userAgent = extractUserAgent(req)
	return len(userAgent) == 0 || curlUserAgentRegexp.Match(userAgent) || wgetUserAgentRegexp.Match(userAgent) || httpieUserAgentRegexp.Match(userAgent)
}

func userAgentIsBrowser(req *http.Request) bool {
	var userAgent = extractUserAgent(req)
	return mozillaUserAgentRegexp.Match(userAgent) ||
		safariUserAgentRegexp.Match(userAgent) ||
		chromeUserAgentRegexp.Match(userAgent) ||
		operaUserAgentRegexp.Match(userAgent)
}

func extractUserAgent(req *http.Request) []byte {
	var userAgent = req.Header.Get("User-Agent")
	return []byte(userAgent)
}
