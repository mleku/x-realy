package normalize

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type urlTest struct {
	url, expected string
}

var urlTests = []urlTest{
	{"", ""},
	{"wss://x.com/y", "wss://x.com/y"},
	{"wss://x.com/y/", "wss://x.com/y"},
	{"http://x.com/y", "ws://x.com/y"},
	{Url("http://x.com/y"), "ws://x.com/y"},
	{Url("wss://x.com"), "wss://x.com"},
	{Url("wss://x.com/"), "wss://x.com"},
	{Url(Url(Url("wss://x.com/"))), "wss://x.com"},
	{"wss://x.com", "wss://x.com"},
	{"wss://x.com/", "wss://x.com"},
	{"x.com////", "wss://x.com"},
	{"x.com/?x=23", "wss://x.com?x=23"},
	{"localhost:4036", "ws://localhost:4036"},
	{"localhost:4036/relay", "ws://localhost:4036/relay"},
	{"localhostmagnanimus.com", "wss://localhostmagnanimus.com"},
	{Url("localhost:4036/relay"), "ws://localhost:4036/relay"},
}

func TestUrl(t *testing.T) {
	for _, test := range urlTests {
		output := Url(test.url)
		assert.Equal(t, test.expected, output)
	}
}
