package tests

import (
	"testing"
	"time"

	"url_shortener/internal/data/repository"
)

// delete old urls
func TestDeleteOld(t *testing.T) {
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
		want bool
	}{
		{"positive", "reormervu", "https://youtube.com", -7, false},
		{"negative", "rkrevurlr", "https://vk.com", 1, true},
	}

	for _, tt := range testcases {
		t.Run(tt.name, func(t *testing.T) {
			date := time.Now().AddDate(0, 0, tt.days).UTC()
			dateFormat := date.Format("2006-01-02 15:04:05")

			err = repository.AddCode(tt.code, tt.url, dateFormat)
			if err != nil {
				t.Fatalf("couldn't add code: %v", err)
			}

			isExists, err := repository.IsExists(tt.code)
			if err != nil {
				t.Fatalf("couldn't check code %s: %v", tt.code, err)
			} else if !isExists {
				t.Fatalf("code doesn't exist: expected %v, got %v", true, false)
			}

			err = repository.ClearOld()
			if err != nil {
				t.Fatalf("couldn't delete old url: %v", err)
			}

			isExists, err = repository.IsExists(tt.code)
			if err != nil {
				t.Fatalf("couldn't check code %s: %v", tt.code, err)
			} else if isExists != tt.want {
				t.Fatalf("code exists: expected %v, got %v", tt.want, isExists)
			}
		})
	}
}
