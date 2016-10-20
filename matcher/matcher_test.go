package matcher_test;

import (
    "github.com/patmooney/captaingo/matcher"
    "testing"
    "fmt"
);

func TestMatch(t *testing.T) {
    var captain matcher.Matcher = matcher.Matcher{};
    captain.SetSource( []byte("[{\"name\":\"London\",\"id\":\"aabbcc\",\"keywords\":[\"England\"]}]") );
    var matches []matcher.Datum = captain.Match("London", []string{});

    if len(matches) <= 0 {
        t.FailNow();
    }

    if matches[0].Name != "London" {
        t.Error( "Expected London to be the first match" );
    }

}

func ExampleMatch() {
    var captain matcher.Matcher = matcher.Matcher{};
    captain.SetSource( []byte("[{\"name\":\"London\",\"id\":\"aabbcc\",\"keywords\":[\"England\"]}]") );
    var matches []matcher.Datum = captain.Match("London", []string{});

    fmt.Println( matches[0].Name );
    // Output: London
}
