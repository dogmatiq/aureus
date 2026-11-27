package markdownloader_test

import (
	"strings"
	"testing"

	. "github.com/dogmatiq/aureus/internal/loader/markdownloader"
)

func TestLoadContent(t *testing.T) {
	t.Run("when a fixture declaration is invalid", func(t *testing.T) {
		cases := []struct {
			Name string
			Info string
			Want string
		}{
			{
				Name: "also an input",
				Info: "sql au:fixture au:input au:group=example",
				Want: "only one of",
			},
			{
				Name: "also an output",
				Info: "sql au:fixture au:output au:group=example",
				Want: "only one of",
			},
			{
				Name: "fixture flag has a value",
				Info: "sql au:fixture=yes au:group=example",
				Want: `"au:fixture" attribute must not have a value`,
			},
		}

		for _, tc := range cases {
			t.Run(tc.Name, func(t *testing.T) {
				_, _, err := LoadContent(nil, tc.Info, "fixture")
				if err == nil || !strings.Contains(err.Error(), tc.Want) {
					t.Fatalf("expected an error containing %q, got %v", tc.Want, err)
				}
			})
		}
	})
}
