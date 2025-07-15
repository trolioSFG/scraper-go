package main

import (
	"reflect"
	"testing"
)

func TestGetURLS(t *testing.T) {
	tests := []struct {
		name string
		input string
		base string
		expectedURLS []string
	}{
		{
			name: "basic",
			input: `Here is <a href="http://boot.dev">link</a>`,
			base: "http://boot.dev",
			expectedURLS: []string{"http://boot.dev"},
		},
		{
			name: "multilink basic",
			input: `Here is <a href="/book">link</a> and <a href="/about/me">another</a>`,
			base: "http://boot.dev",
			expectedURLS: []string{"http://boot.dev/book", "http://boot.dev/about/me"},
		},
		{
			name: "multilink relative absolute",
			input: `Here is <a href="/book">link</a> and <a href="/about/me">another</a>
			 <p>The last one: <a href="https://apple.com">Apple</a></p>`,
			base: "http://boot.dev",
			expectedURLS: []string{"http://boot.dev/book",
				"http://boot.dev/about/me",
				"https://apple.com",
			},
		},
	}

	for i, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			links, err := getURLSFromHTML(tc.input, tc.base)
			if err != nil {
				t.Errorf("Test %v - '%s' FAILED: %v", i, tc.name, err)
				return
			}
			if !reflect.DeepEqual(tc.expectedURLS, links) {
				t.Errorf("Test %v - %s FAIL:\nExpected: %v\nActual: %v",
					i, tc.name, tc.expectedURLS, links)
			}
		})
	}

}

