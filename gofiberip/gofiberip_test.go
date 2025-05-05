package gofiberip

import (
	"io"
	"net"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func TestIPRestrictAdvanced(t *testing.T) {
	allowedIPs := []string{"192.168.1.100", "10.0.0.0/24", "2001:db8::/32"}
	middleware := New(Config{AllowedIPs: allowedIPs})

	tests := []struct {
		name          string
		ip            string
		xForwardedFor string
		expectedCode  int
	}{
		{
			name:         "IP exacta permitida",
			ip:           "192.168.1.100",
			expectedCode: fiber.StatusOK,
		},
		{
			name:         "IP en rango CIDR permitido (IPv4)",
			ip:           "10.0.0.50",
			expectedCode: fiber.StatusOK,
		},
		{
			name:         "IP en rango CIDR permitido (IPv6)",
			ip:           "2001:db8::1",
			expectedCode: fiber.StatusOK,
		},
		{
			name:         "IP no permitida",
			ip:           "192.168.1.101",
			expectedCode: fiber.StatusForbidden,
		},
		{
			name:          "IP en X-Forwarded-For permitida",
			ip:            "1.2.3.4",
			xForwardedFor: "192.168.1.100, 1.2.3.4",
			expectedCode:  fiber.StatusOK,
		},
		{
			name:         "IP malformada",
			ip:           "not.an.ip.address",
			expectedCode: fiber.StatusForbidden,
		},
		{
			name:         "IP vacía",
			ip:           "",
			expectedCode: fiber.StatusForbidden,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := fiber.New()
			app.Use(middleware)
			app.Get("/", func(c *fiber.Ctx) error {
				return c.SendString("OK")
			})

			req := httptest.NewRequest("GET", "/", nil)
			req.RemoteAddr = tt.ip + ":1234" 

			if tt.xForwardedFor != "" {
				req.Header.Set("X-Forwarded-For", tt.xForwardedFor)
			}

			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedCode, resp.StatusCode)

			if resp.StatusCode == fiber.StatusOK {
				body, err := io.ReadAll(resp.Body)
				assert.NoError(t, err)
				assert.Equal(t, "OK", body)
			}
		})
	}
}

func TestIPRestrictAdvancedCIDR(t *testing.T) {
	// Specific test for CIDR verification
	middleware := New(Config{AllowedIPs: []string{"192.168.1.0/24"}})

	app := fiber.New()
	app.Use(middleware)
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("OK")
	})

	// Generate all the IPs in the range 192.168.1.0/24 and verify
	ip, ipnet, err := net.ParseCIDR("192.168.1.0/24")
	assert.NoError(t, err)

	for ip := ip.Mask(ipnet.Mask); ipnet.Contains(ip); inc(ip) {
		req := httptest.NewRequest("GET", "/", nil)
		req.RemoteAddr = ip.String() + ":1234"

		resp, err := app.Test(req)
		assert.NoError(t, err)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	}

	// Verify if an IP out of range is rejected
	req := httptest.NewRequest("GET", "/", nil)
	req.RemoteAddr = "192.168.2.1:1234"
	resp, err := app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, fiber.StatusForbidden, resp.StatusCode)
}

// increments an IP address
func inc(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}

func TestIPRestrictAdvancedEdgeCases(t *testing.T) {
	tests := []struct {
		name         string
		allowedIPs   []string
		ip           string
		expectedCode int
	}{
		{
			name:         "Lista vacía de IPs permitidas",
			allowedIPs:   []string{},
			ip:           "192.168.1.1",
			expectedCode: fiber.StatusForbidden,
		},
		{
			name:         "CIDR inválido en configuración",
			allowedIPs:   []string{"192.168.1.100", "invalid.cidr/24"},
			ip:           "192.168.1.100",
			expectedCode: fiber.StatusOK,
		},
		{
			name:         "IP con puerto",
			allowedIPs:   []string{"192.168.1.100"},
			ip:           "192.168.1.100:12345",
			expectedCode: fiber.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			middleware := New(Config{AllowedIPs: tt.allowedIPs})

			app := fiber.New()
			app.Use(middleware)
			app.Get("/", func(c *fiber.Ctx) error {
				return c.SendString("OK")
			})

			req := httptest.NewRequest("GET", "/", nil)
			req.RemoteAddr = tt.ip

			resp, err := app.Test(req)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedCode, resp.StatusCode)
		})
	}
}
