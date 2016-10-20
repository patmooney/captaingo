package matcher

import (
    "../matcher/"
    "testing"
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
