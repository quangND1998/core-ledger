package fingerprint

import (
	"encoding/json"
	"fmt"
	"github.com/zeebo/xxh3"
)

type AcceptHeader struct {
	Accept          string `json:"accept"`
	AcceptLanguage  string `json:"accept_language"`
	UserAgentHeader string `json:"user_agent"`
}

type IpAddress struct {
	Value string `json:"value"`
}

type Fingerprint struct {
	ID        string       `json:"id"`
	Headers   AcceptHeader `json:"headers"`
	UserAgent UserAgent    `json:"user_agent"`
	IpAddress IpAddress    `json:"ip_address"`
}

func (f *Fingerprint) BrowserID() (string, error) {
	ids, err := GenerateID(*f)
	if err != nil {
		return "", err
	}
	return ids.BrowserID, nil
}
func (f *Fingerprint) DeviceID() (string, error) {
	ids, err := GenerateID(*f)
	if err != nil {
		return "", err
	}
	return ids.DeviceID, nil
}

type IPInfo struct {
	Value string `json:"value"`
}

type UserAgent struct {
	Device  DeviceInfo  `json:"device"`
	OS      OSInfo      `json:"os"`
	Browser BrowserInfo `json:"browser"`
}

type DeviceInfo struct {
	Family  string   `json:"family"`
	Version []string `json:"version"`
}

type OSInfo struct {
	Family string `json:"family"`
	Major  string `json:"major"`
	Minor  string `json:"minor"`
}

type BrowserInfo struct {
	Family  string `json:"family"`
	Version string `json:"version"`
}

type FingerprintID struct {
	HeaderID  string
	DeviceID  string
	BrowserID string
}

// --- REFINED CODE ---

// generateHashedID is a helper function that serializes any given data structure to JSON,
// computes its 128-bit xxhash, and returns it as a hex string.
// This consolidates the repetitive marshal-hash-format logic.
func generateHashedID(data interface{}) (string, error) {
	bytes, err := json.Marshal(data)
	if err != nil {
		// Wrap the error with context for better debugging.
		return "", fmt.Errorf("failed to marshal data for hashing: %w", err)
	}
	hash := xxh3.Hash128(bytes)
	return fmt.Sprintf("%016x%016x", hash.Hi, hash.Lo), nil
}

// GenerateID creates a set of stable identifiers from a request fingerprint.
// It is now more concise and easier to follow the logic flow.
func GenerateID(fingerprint Fingerprint) (FingerprintID, error) {
	var err error

	// 1. Generate HeaderID
	headerID, err := generateHashedID(fingerprint.Headers)
	if err != nil {
		return FingerprintID{}, err
	}

	// 2. Generate DeviceID
	// Using an anonymous struct with json tags is more type-safe than a map[string]interface{}
	// and ensures the JSON output keys are consistent.
	deviceID, err := generateHashedID(struct {
		DeviceFamily  string   `json:"device_family"`
		DeviceVersion []string `json:"device_version"`
		OSFamily      string   `json:"os_family"`
		OSMajor       string   `json:"os_major"`
		OSMinor       string   `json:"os_minor"`
	}{
		DeviceFamily:  fingerprint.UserAgent.Device.Family,
		DeviceVersion: fingerprint.UserAgent.Device.Version,
		OSFamily:      fingerprint.UserAgent.OS.Family,
		OSMajor:       fingerprint.UserAgent.OS.Major,
		OSMinor:       fingerprint.UserAgent.OS.Minor,
	})
	if err != nil {
		return FingerprintID{}, err
	}

	// 3. Generate IP-based ID (used internally for the browser ID)
	ipID, err := generateHashedID(struct {
		IPAddress string `json:"ipAddress"`
	}{
		IPAddress: fingerprint.IpAddress.Value,
	})
	if err != nil {
		return FingerprintID{}, err
	}

	// 4. Generate the final BrowserID, which aggregates other IDs
	browserID, err := generateHashedID(struct {
		IPID           string `json:"ip_id"`
		HeaderID       string `json:"header_id"`
		DeviceID       string `json:"device_id"`
		BrowserFamily  string `json:"browser_family"`
		BrowserVersion string `json:"browser_version"`
	}{
		IPID:           ipID,
		HeaderID:       headerID,
		DeviceID:       deviceID,
		BrowserFamily:  fingerprint.UserAgent.Browser.Family,
		BrowserVersion: fingerprint.UserAgent.Browser.Version,
	})
	if err != nil {
		return FingerprintID{}, err
	}

	return FingerprintID{
		HeaderID:  headerID,
		DeviceID:  deviceID,
		BrowserID: browserID,
	}, nil
}

func hashToString(u xxh3.Uint128) string {
	return fmt.Sprintf("%016x%016x", u.Hi, u.Lo)
}
