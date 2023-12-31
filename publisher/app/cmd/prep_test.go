package cmd

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestPrep_MakeShow(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/lastmd/12", r.URL.Path)
		w.Write([]byte("- blah1\n- blah2"))
	}))
	defer ts.Close()

	p := Prep{
		Client:  http.Client{Timeout: 100 * time.Millisecond},
		NewsHrs: 12,
		NewsAPI: ts.URL,
		Dest:    "/tmp",
		now:     func() time.Time { return time.Date(2020, 2, 3, 20, 18, 53, 0, time.Local) },
	}

	err := p.MakeShow(123)
	require.NoError(t, err)
	defer os.Remove("/tmp/podcast-123.md")

	b, err := os.ReadFile("/tmp/podcast-123.md")
	require.NoError(t, err)
	exp := `+++
title = "Радио-Т 123"
date = 2020-02-03T20:18:53
categories = ["podcast"]
image = "https://radio-t.com/images/radio-t/rt123.jpg"
filename = "rt_podcast123"
+++

![](https://radio-t.com/images/radio-t/rt123.jpg)

- blah1
- blah2

[аудио](https://cdn.radio-t.com/rt_podcast123.mp3) • [лог чата](https://chat.radio-t.com/logs/radio-t-123.html)
<audio src="https://cdn.radio-t.com/rt_podcast123.mp3" preload="none"></audio>
`
	assert.Equal(t, exp, string(b))
}

func TestPrep_MakePrep(t *testing.T) {
	p := Prep{
		Client:  http.Client{Timeout: 100 * time.Millisecond},
		NewsHrs: 12,
		Dest:    "/tmp",
		now:     func() time.Time { return time.Date(2020, 2, 3, 20, 18, 53, 0, time.Local) },
	}

	err := p.MakePrep(123)
	require.NoError(t, err)
	defer os.Remove("/tmp/prep-123.md")

	b, err := os.ReadFile("/tmp/prep-123.md")
	require.NoError(t, err)
	exp := `+++
title = "Темы для 123"
date = 2020-02-03T20:18:53
categories = ["prep"]
+++
`
	assert.Equal(t, exp, string(b))
}
