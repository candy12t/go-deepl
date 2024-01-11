package deepl

import "testing"

func Test_NewClient(t *testing.T) {
	tests := map[string]struct {
		authKey     string
		wantBaseURL string
	}{
		"free acount": {
			authKey:     "free-auth-key:fx",
			wantBaseURL: "https://api-free.deepl.com/v2",
		},
		"pro account": {
			authKey:     "pro-auth-key",
			wantBaseURL: "https://api.deepl.com/v2",
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			c := NewClient(tt.authKey)
			if got := c.BaseURL.String(); got != tt.wantBaseURL {
				t.Errorf("NewClient(%q).BaseURL is %v, want %v", tt.authKey, got, tt.wantBaseURL)
			}
		})
	}
}
