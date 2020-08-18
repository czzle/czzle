package transport

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/czzle/czzle"
	"github.com/czzle/czzle/internal/app/endpoints"
	"github.com/czzle/czzle/pkg/multierr"
	"github.com/czzle/czzle/pkg/uuid"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	httptp "github.com/go-kit/kit/transport/http"
)

var (
	ErrEncodeTypeAssertion = errors.New("type assertion while encoding transport")
)

func NewHTTPServer(eps endpoints.Set, opts ...httptp.ServerOption) http.Handler {
	opts = append(opts, httptp.ServerErrorEncoder(encodeHTTPError))
	r := chi.NewRouter()
	r.Use(
		middleware.Timeout(time.Second*10),
		middleware.AllowContentType("application/json"),
		cors.Handler(cors.Options{
			AllowedOrigins:   []string{"*"},
			AllowedMethods:   []string{"POST"},
			AllowedHeaders:   []string{"Accept", "Content-Type", "CZZRL-CLIENT-ID"},
			ExposedHeaders:   []string{""},
			AllowCredentials: true,
			MaxAge:           300, // Maximum value not ignored by any of major browsers
		}),
	)
	r.Method("POST", "/v1/begin",
		httptp.NewServer(
			eps.BeginEndpoint,
			decodeHTTPBeginReq,
			encodeHTTPBeginRes,
			opts...,
		),
	)
	r.Method("POST", "/v1/solve",
		httptp.NewServer(
			eps.SolveEndpoint,
			decodeHTTPSolveReq,
			encodeHTTPSolveRes,
			opts...,
		),
	)
	r.Method("POST", "/v1/validate",
		httptp.NewServer(
			eps.ValidateEndpoint,
			decodeHTTPValidateReq,
			encodeHTTPValidateRes,
			opts...,
		),
	)
	return r
}

func encodeHTTPError(ctx context.Context, err error, w http.ResponseWriter) {
	multierr.ToHTTP(w, err)
}

func decodeHTTPBeginReq(ctx context.Context, r *http.Request) (interface{}, error) {
	ip := getIPAdress(r)
	return &czzle.BeginReq{
		Client: &czzle.ClientInfo{
			ID:        getClientID(r),
			IP:        ip,
			UserAgent: getUserAgent(r),
			Time:      time.Now().UnixNano() / int64(time.Millisecond),
		},
	}, nil
}

func encodeHTTPBeginRes(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	res, ok := response.(*czzle.BeginRes)
	if !ok {
		return ErrEncodeTypeAssertion
	}
	data, err := json.Marshal(res)
	if err != nil {
		return err
	}
	c := &http.Cookie{
		Name:     "czzle-device-id",
		Value:    res.GetPuzzle().GetClient().GetID().String(),
		HttpOnly: true,
	}
	http.SetCookie(w, c)
	w.Header().Add("content-type", "application/json")
	w.Write(data)
	return nil
}

func decodeHTTPSolveReq(ctx context.Context, r *http.Request) (interface{}, error) {
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}
	var req czzle.SolveReq
	err = json.Unmarshal(data, &req)
	if err != nil {
		return nil, err
	}
	return &req, nil
}

func encodeHTTPSolveRes(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	res, ok := response.(*czzle.SolveRes)
	if !ok {
		return ErrEncodeTypeAssertion
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Add("content-type", "application/json")
	return json.NewEncoder(w).Encode(res)
}

func decodeHTTPValidateReq(ctx context.Context, r *http.Request) (interface{}, error) {
	return &czzle.ValidateReq{}, nil
}

func encodeHTTPValidateRes(ctx context.Context, w http.ResponseWriter, response interface{}) error {
	res, ok := response.(*czzle.ValidateRes)
	if !ok {
		return ErrEncodeTypeAssertion
	}
	w.WriteHeader(http.StatusOK)
	w.Header().Add("content-type", "application/json")
	return json.NewEncoder(w).Encode(res)
}

type ipRange struct {
	start net.IP
	end   net.IP
}

// inRange - check to see if a given ip address is within a range given
func inRange(r ipRange, ipAddress net.IP) bool {
	// strcmp type byte comparison
	if bytes.Compare(ipAddress, r.start) >= 0 && bytes.Compare(ipAddress, r.end) < 0 {
		return true
	}
	return false
}

var privateRanges = []ipRange{
	ipRange{
		start: net.ParseIP("10.0.0.0"),
		end:   net.ParseIP("10.255.255.255"),
	},
	ipRange{
		start: net.ParseIP("100.64.0.0"),
		end:   net.ParseIP("100.127.255.255"),
	},
	ipRange{
		start: net.ParseIP("172.16.0.0"),
		end:   net.ParseIP("172.31.255.255"),
	},
	ipRange{
		start: net.ParseIP("192.0.0.0"),
		end:   net.ParseIP("192.0.0.255"),
	},
	ipRange{
		start: net.ParseIP("192.168.0.0"),
		end:   net.ParseIP("192.168.255.255"),
	},
	ipRange{
		start: net.ParseIP("198.18.0.0"),
		end:   net.ParseIP("198.19.255.255"),
	},
}

// isPrivateSubnet - check to see if this ip is in a private subnet
func isPrivateSubnet(ipAddress net.IP) bool {
	// my use case is only concerned with ipv4 atm
	if ipCheck := ipAddress.To4(); ipCheck != nil {
		// iterate over all our ranges
		for _, r := range privateRanges {
			// check if this ip is in a private range
			if inRange(r, ipAddress) {
				return true
			}
		}
	}
	return false
}

func getIPAdress(r *http.Request) string {
	for _, h := range []string{"x-original-forwarded-for", "x-forwarded-for", "x-real-ip"} {
		addresses := strings.Split(r.Header.Get(h), ",")
		// march from right to left until we get a public address
		// that will be the address right before our proxy.
		for i := len(addresses) - 1; i >= 0; i-- {
			ip := strings.TrimSpace(addresses[i])
			// header can contain spaces too, strip those out.
			realIP := net.ParseIP(ip)
			if !realIP.IsGlobalUnicast() || isPrivateSubnet(realIP) {
				// bad address, go to next
				continue
			}
			return ip
		}
	}
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return ""
	}
	return ip
}

func getUserAgent(r *http.Request) string {
	return r.Header.Get("user-agent")
}

func getClientID(r *http.Request) uuid.UUID {
	c, err := r.Cookie("czzle-device-id")
	if err != nil {
		return uuid.Null()
	}
	return uuid.FromString(c.String())
}
