package config

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetSecurityTokenByOIDCV5(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" || r.URL.Path != "/v5/agencies/assume-with-oidc" {
			t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
		var body map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatal(err)
		}
		want := map[string]interface{}{
			"provider_urn":        "iam::domain:oidcProvider:ci",
			"agency_urn":          "iam::domain:agency:deployer",
			"agency_session_name": "deployer",
			"id_token":            "jwt",
			"duration_seconds":    float64(900),
		}
		for k, v := range want {
			if body[k] != v {
				t.Errorf("%s = %v, want %v", k, body[k], v)
			}
		}
		w.Write([]byte(`{"credentials":{"access_key_id":"ak","secret_access_key":"sk","security_token":"st"}}`))
	}))
	defer server.Close()

	c := &Config{
		Region:                "cn-north-4",
		Endpoints:             map[string]string{"sts": server.URL + "/"},
		AssumeRoleAgency:      "deployer",
		AssumeRoleDomainID:    "domain",
		AssumeRoleDuration:    900,
		AssumeRoleProviderURN: "iam::domain:oidcProvider:ci",
		AssumeRoleIdToken:     "jwt",
	}
	if err := getSecurityTokenByOIDCV5(c); err != nil {
		t.Fatal(err)
	}
	if c.AccessKey != "ak" || c.SecretKey != "sk" || c.SecurityToken != "st" {
		t.Errorf("unexpected credentials: %s %s %s", c.AccessKey, c.SecretKey, c.SecurityToken)
	}
}
