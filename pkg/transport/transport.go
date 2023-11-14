package transport

import (
	"time"

	"golang.org/x/net/http/httpproxy"
)

// CustomTransport Allows to create a global transport configuration
// re-usable across providers
type CustomTransport struct {

	// Allows to specify a proxy via the standard ENV variables
	// HTTPProxy represents the value of the HTTP_PROXY or
	// http_proxy environment variable. It will be used as the proxy
	// URL for HTTP requests unless overridden by NoProxy.
	// HTTPProxy string

	// HTTPSProxy represents the HTTPS_PROXY or https_proxy
	// environment variable. It will be used as the proxy URL for
	// HTTPS requests unless overridden by NoProxy.
	// HTTPSProxy string

	// NoProxy represents the NO_PROXY or no_proxy environment
	// variable. It specifies a string that contains comma-separated values
	// specifying hosts that should be excluded from proxying. Each value is
	// represented by an IP address prefix (1.2.3.4), an IP address prefix in
	// CIDR notation (1.2.3.4/8), a domain name, or a special DNS label (*).
	// An IP address prefix and domain name can also include a literal port
	// number (1.2.3.4:80).
	// A domain name matches that name and all subdomains. A domain name with
	// a leading "." matches subdomains only. For example "foo.com" matches
	// "foo.com" and "bar.foo.com"; ".y.com" matches "x.y.com" but not "y.com".
	// A single asterisk (*) indicates that no proxying should be done.
	// A best effort is made to parse the string and errors are
	// ignored.
	// NoProxy string
	Proxy httpproxy.Config

	// This setting represents the maximum amount of time to keep an idle network connection alive between HTTP requests.
	//
	// Set to 0 for no limit.
	//
	// See https://golang.org/pkg/net/http/#Transport.IdleConnTimeout
	IdleConnTimeout time.Duration

	// This setting represents the maximum number of idle (keep-alive) connections
	// across all hosts. One use case for increasing this value is when you are seeing many connections in a short period from the same clients
	//
	// 0 means no limit.
	//
	// See https://golang.org/pkg/net/http/#Transport.MaxIdleConns
	MaxIdleConns int

	// This setting represents the maximum number of idle (keep-alive) connections
	// to keep per-host.
	// One use case for increasing this value is when you are seeing many connections
	// in a short period from the same clients
	//
	// Default is two idle connections per host.
	//
	// Set to 0 to use DefaultMaxIdleConnsPerHost (2).
	//
	// See https://golang.org/pkg/net/http/#Transport.MaxIdleConnsPerHost
	MaxIdleConnsPerHost int

	// This setting represents the maximum amount of time to wait for a client to
	// read the response header.
	// If the client isn’t able to read the response’s header within this duration,
	// the request fails with a timeout error.
	// Be careful setting this value when using long-running Lambda functions,
	// as the operation does not return any response headers until the Lambda
	// function has finished or timed out.
	// However, you can still use this option with the ** InvokeAsync** API operation.
	//
	// Default is no timeout; wait forever.
	//
	// See https://golang.org/pkg/net/http/#Transport.ResponseHeaderTimeout
	ResponseHeaderTimeout time.Duration

	// This setting represents the maximum amount of time waiting for a
	// TLS handshake to be completed.
	//
	// Default is 10 seconds.
	//
	// Zero means no timeout.
	//
	// See https://golang.org/pkg/net/http/#Transport.TLSHandshakeTimeout
	TLSHandshakeTimeout time.Duration
}

func CustomTransportFromConfig() *CustomTransport {
	return nil
}
