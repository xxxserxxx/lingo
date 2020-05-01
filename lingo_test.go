package lingo

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestFallback(t *testing.T) {
	l, err := New("en_US", "translations", nil)
	if err != nil {
		t.Error(err)
	}
	t1 := l.TranslationsForLocale("de_DE")
	actual := t1.Value("only.english")
	expected := "This is only in the English file"
	if actual != expected {
		t.Errorf("Expected %v, got %v", expected, actual)
		t.Fail()
	}
}

func TestLingo(t *testing.T) {
	tests := []struct {
		key      string
		expected string
	}{
		{"main.subtitle", "Knives that put cut in cutlery."},
		{"home.title", "Welcome to CutleryPlus!"},
		{"menu.products.self", "Products"},
		{"menu.non.existant", "menu.non.existant"},
	}

	l, err := New("de_DE", "translations", nil)
	if err != nil {
		t.Error(err)
	}
	t1 := l.TranslationsForLocale("en_US")

	for _, tc := range tests {
		actual := t1.Value(tc.key)
		if actual != tc.expected {
			t.Errorf("Expected %v, got %v", tc.expected, actual)
			t.Fail()
		}
	}

	actual := t1.Value("error.404", "idnex.html")
	expected := "Page idnex.html not found!"
	if actual != expected {
		t.Errorf("Expected %v, got %v", expected, actual)
		t.Fail()
	}
}

func TestLingoHttp(t *testing.T) {
	l, err := New("en_US", "translations", nil)
	if err != nil {
		t.Error(err)
	}
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		expected := r.Header.Get("Expected-Results")
		t1 := l.TranslationsForRequest(r)
		actual := t1.Value("error.500")
		if actual != expected {
			t.Errorf("Expected %v, got %v", expected, actual)
			t.Fail()
		}
	}))
	defer srv.Close()
	url, _ := url.Parse(srv.URL)

	req1 := &http.Request{
		Method: "GET",
		Header: map[string][]string{
			"Accept-Language":  {"sr, en-gb;q=0.8, en;q=0.7"},
			"Expected-Results": {"Greska sa nase strane, pokusajte ponovo."},
		},
		URL: url,
	}
	req2 := &http.Request{
		Method: "GET",
		Header: map[string][]string{
			"Accept-Language":  {"en-US, en-gb;q=0.8, en;q=0.7"},
			"Expected-Results": {"Something is wrong on our side, please try again."},
		},
		URL: url,
	}
	req3 := &http.Request{
		Method: "GET",
		Header: map[string][]string{
			"Accept-Language":  {"de-at, en-gb;q=0.8, en;q=0.7"},
			"Expected-Results": {"Stimmt etwas nicht auf unserer Seite ist, versuchen Sie es erneut."},
		},
		URL: url,
	}

	http.DefaultClient.Do(req1)
	http.DefaultClient.Do(req2)
	http.DefaultClient.Do(req3)
}
