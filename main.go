package main;

import (
    "flag"
    "log"
    "github.com/patmooney/captaingo/matcher"
);

var captain matcher.Matcher;

func main () {
    var filenamePtr *string = flag.String( "source", "", "filename from which to source data" );
    flag.Parse();

    if *filenamePtr == "" {
        log.Fatal( "--source is a required option" );
    }

    captain = matcher.NewMatcher( *filenamePtr );
}
