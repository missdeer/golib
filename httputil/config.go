package httputil

var (
	insecureSkipVerify bool
	userAgent        = "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/74.0.3729.108 Safari/537.36"
)

func SetUserAgent(ua string) {
	userAgent = ua
}

func SetInsecureSkipVerify(b bool) {
	insecureSkipVerify = b
}
