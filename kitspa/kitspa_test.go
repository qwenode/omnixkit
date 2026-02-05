package kitspa

import (
	"embed"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

//go:embed testdata/dist/* testdata/dist/assets/*
var testFS embed.FS

func TestMount_CacheControlDefaults(t *testing.T) {
	gin.SetMode(gin.TestMode)

	r := gin.New()
	err := Mount(r, Config{
		FS:      testFS,
		DistDir: "testdata/dist",
	})
	if err != nil {
		t.Fatalf("Mount() error = %v", err)
	}

	t.Run("assets cache", func(t *testing.T) {
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/assets/app.js", nil)
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("status = %d, want %d", w.Code, http.StatusOK)
		}
		if got, want := w.Header().Get("Cache-Control"), "public, max-age=31536000, immutable"; got != want {
			t.Fatalf("Cache-Control = %q, want %q", got, want)
		}
	})

	t.Run("index cache", func(t *testing.T) {
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("status = %d, want %d", w.Code, http.StatusOK)
		}
		if got, want := w.Header().Get("Cache-Control"), "no-cache"; got != want {
			t.Fatalf("Cache-Control = %q, want %q", got, want)
		}
	})

	t.Run("spa fallback cache", func(t *testing.T) {
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/about", nil)
		r.ServeHTTP(w, req)

		if w.Code != http.StatusOK {
			t.Fatalf("status = %d, want %d", w.Code, http.StatusOK)
		}
		if got, want := w.Header().Get("Cache-Control"), "no-cache"; got != want {
			t.Fatalf("Cache-Control = %q, want %q", got, want)
		}
	})

	t.Run("blocked prefix", func(t *testing.T) {
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/.env", nil)
		r.ServeHTTP(w, req)

		if w.Code != http.StatusNotFound {
			t.Fatalf("status = %d, want %d", w.Code, http.StatusNotFound)
		}
	})

	t.Run("non-GET fallback", func(t *testing.T) {
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/anything", nil)
		r.ServeHTTP(w, req)

		if w.Code != http.StatusNotFound {
			t.Fatalf("status = %d, want %d", w.Code, http.StatusNotFound)
		}
	})
}

func TestMount_CacheControlCustom(t *testing.T) {
	gin.SetMode(gin.TestMode)

	r := gin.New()
	err := Mount(r, Config{
		FS:                 testFS,
		DistDir:            "testdata/dist",
		AssetsCacheControl: "public, max-age=3600",
		IndexCacheControl:  "no-store",
	})
	if err != nil {
		t.Fatalf("Mount() error = %v", err)
	}

	t.Run("assets cache", func(t *testing.T) {
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/assets/app.js", nil)
		r.ServeHTTP(w, req)

		if got, want := w.Header().Get("Cache-Control"), "public, max-age=3600"; got != want {
			t.Fatalf("Cache-Control = %q, want %q", got, want)
		}
	})

	t.Run("index cache", func(t *testing.T) {
		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		r.ServeHTTP(w, req)

		if got, want := w.Header().Get("Cache-Control"), "no-store"; got != want {
			t.Fatalf("Cache-Control = %q, want %q", got, want)
		}
	})
}
