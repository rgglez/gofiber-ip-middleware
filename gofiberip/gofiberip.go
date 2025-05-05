/*
Copyright 2025 Rodolfo González González

   Licensed under the Apache License, Version 2.0 (the "License");
   you may not use this file except in compliance with the License.
   You may obtain a copy of the License at

       http://www.apache.org/licenses/LICENSE-2.0

   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
*/

package gofiberip

import (
	"net"
	"strings"

	fiber "github.com/gofiber/fiber/v2"
)

//-----------------------------------------------------------------------------

type Config struct {
	// Next defines a function to skip this middleware when returned true.
	//
	// Optional. Default: nil
	Next func(c *fiber.Ctx) bool

	// AllowedIPs define the allowed IPs/CIDRs
	//
	// Required
	AllowedIPs []string
}

//-----------------------------------------------------------------------------

// Set the default configuration.
var ConfigDefault = Config{
	Next:       nil,
	AllowedIPs: []string{"localhost"},
}

//-----------------------------------------------------------------------------

func New(config ...Config) fiber.Handler {
	cfg := ConfigDefault

	if len(config) > 0 {
		cfg = config[0]
	}

	// Preparse las IPs/CIDRs permitidas
	var ipnets []*net.IPNet
	for _, ip := range cfg.AllowedIPs {
		if strings.Contains(ip, "/") {
			_, ipnet, err := net.ParseCIDR(ip)
			if err == nil {
				ipnets = append(ipnets, ipnet)
			}
		}
	}

	return func(c *fiber.Ctx) error {
		// Should we pass?
		if cfg.Next != nil && cfg.Next(c) {
			return c.Next()
		}

		allPossibleIPs := append([]string{c.IP()}, c.IPs()...)

		for _, ip := range allPossibleIPs {
			cleanIP := strings.Split(strings.TrimSpace(ip), ":")[0]
			parsedIP := net.ParseIP(cleanIP)
			if parsedIP == nil {
				continue
			}

			// Verify exact IPs
			for _, allowedIP := range cfg.AllowedIPs {
				if !strings.Contains(allowedIP, "/") && cleanIP == allowedIP {
					return c.Next()
				}
			}

			// Verify CIDR ranges
			for _, ipnet := range ipnets {
				if ipnet.Contains(parsedIP) {
					return c.Next()
				}
			}
		}

		return c.Status(fiber.StatusForbidden).SendString("Forbidden access")
	}
}
