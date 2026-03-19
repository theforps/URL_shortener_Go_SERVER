package tests

import (
	"testing"
	"time"

	"url_shortener/internal/data/repository"

	_ "github.com/lib/pq"
)

// add new redirect to db
func TestAdd(t *testing.T) {
	db, err := waitForDB()
	if err != nil {
		t.Fatalf("db not ready: %v", err)
	}
	defer db.Close()
	var repository repository.StorageRepository = repository.NewStorageRepository(db)

	var testcases = []struct {
		name string
		code string
		url  string
		days int
		want string
	}{
		{"positive-1", "r43eormervu", "https://youtube.com", 2, "https://youtube.com"},
		{"positive-2", "rkr65evurlr", "https://vk.com", 1, "https://vk.com"},
		{"positive-3", "rkr65eV23lr", "https://mail.google.com", 6, "https://mail.google.com"},
	}

	for _, tt := range testcases {
		t.Run(tt.name, func(t *testing.T) {
			date := time.Now().AddDate(0, 0, tt.days).UTC()
			dateFormat := date.Format("2006-01-02 15:04:05")

			err = repository.AddCode(tt.code, tt.url, dateFormat)
			if err != nil {
				t.Fatalf("couldn't add code: %v", err)
			}

			url, err := repository.GetBaseUrl(tt.code)
			if err != nil {
				t.Fatalf("couldn't get url: %v", err)
			}

			if url != tt.want {
				t.Fatalf("url mismatch: expected %s, got %s", tt.want, url)
			}
		})
	}
}
