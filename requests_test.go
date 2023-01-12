package requests

import (
	"testing"
	"time"
)

func TestRequest(t *testing.T) {
	// request
	url := "http://x.x.x.x"
	r := NewRequest(nil)
	r.Get(url, nil)
	r.Get(url, map[string]string{"a": "1", "b": "2"})
	r.Post(url, map[string]interface{}{"a": "1", "b": "2"})
	r.PostForm(url, map[string]string{"a": "1", "b": "2"})
	r.Put(url, map[string]interface{}{"a": "1", "b": "2"})
	r.Delete(url)
	// session
	loginUrl := "http://x.x.x.x/login"
	s := NewSession(nil)
	s.Post(loginUrl, map[string]interface{}{"user": "username", "password": "password"})
	s.Get(url, nil)
	s.Get(url, map[string]string{"a": "1", "b": "2"})
	s.Post(url, map[string]interface{}{"a": "1", "b": "2"})
	s.PostForm(url, map[string]string{"a": "1", "b": "2"})
	s.Put(url, map[string]interface{}{"a": "1", "b": "2"})
	s.Delete(url)
}

func TestDownload(t *testing.T) {
	r := NewRequest(nil)
	r.SetTimeout(time.Minute)
	r.DownloadWithRateLimit("/tmp/2", "http://127.0.0.1:60080/xxx", 1024*1024)
}
