package scraper

import "testing"

func TestBaseScraper_CheckDomain(t *testing.T) {
	tests := []struct {
		name     string
		host     string
		inputUrl string
		want     bool
	}{
		{
			name:     "test sanook.com",
			host:     "sanook.com",
			inputUrl: "https://sanook.com/news/1234",
			want:     true,
		},
		{
			name:     "test abc.com",
			host:     "sanook.com",
			inputUrl: "https://abc.com/news/1234",
			want:     false,
		},
		{
			name:     "test www.sanook.com",
			host:     "sanook.com",
			inputUrl: "https://www.sanook.com/news/1234",
			want:     true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := BaseScraper{
				Host: tt.host,
			}
			got, err := s.CheckDomain(tt.inputUrl)
			if err != nil {
				t.Errorf("CheckDomain() got error = %v", err)
			}

			if got != tt.want {
				t.Errorf("CheckDomain() got = %v, want %v", got, tt.want)
			}
		})
	}
}
