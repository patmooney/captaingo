package main;

import (
    "flag"
    "log"
    "net/http"
    "github.com/patmooney/captaingo/matcher"
    "time"
    "fmt"
    "strings"
);

/*
    Captain Go

    For creating a very simple keyword search for mapping input strings against a
    predefined set. Uses Levenshtein distance (https://en.wikipedia.org/wiki/Levenshtein_distance).

    ./captaingo --source=my_data.json --port=8080

    Usage ( Single ):

        GET /match?q=London&keywords=England&keywords=United+Kingdom
        Content-Type: application/x-www-form-urlencoded

    Out:

        {
            "total": 2,
            "matches": [
                {
                    "name": "City of London",
                    "normalised_name": "city of london",
                    "id": 12345,
                    "keywords": "City of London, Greater London, England, United Kingdom",
                    "score": 0
                },
                ...
            ]
        }

    Usage ( Multi ):

        POST /match
        Content-Type: application/json

        {
            "max_score": 1,
            "queries": [
                {
                    "q": "London",
                    "keywords": [ "Greater London", "England" ],
                    "id": 554433
                },
                ...
            ]
        }

    Out:

        {
            "queries": [
                {
                    "input": {
                        "q": "London",
                        "keywords": [ "Greater London", "England" ],
                        "id": 554433
                    },
                    "total": 2,
                    "matches": [
                        {
                            "name": "City of London",
                            "id": 12345,
                            "keywords": "City of London, Greater London, England, United Kingdom",
                            "score": 0
                        },
                        ...
                    ]
                }
            ]
        }

    Input source should be in format

    [
        {
            "name": "City of London",
            "id": 12345,
            "keywords": ["City of London", "Greater London", "England", "United Kingdom"]
        },
        ...
    ]
*/

var captain matcher.Matcher;

func main () {
    var filenamePtr *string = flag.String( "source", "", "filename from which to source data" );
    flag.Parse();

    if *filenamePtr == "" {
        log.Fatal( "--source is a required option" );
    }

    captain = matcher.NewMatcher( *filenamePtr );

    http.HandleFunc( "/bar", func(w http.ResponseWriter, r *http.Request) {
        keyword := strings.ToLower( r.FormValue("q") );
        var datum []matcher.Datum = captain.Match(keyword, []string{});
        w.Write( []byte(fmt.Sprintf( "%s - %s", datum[0].Name, datum[0].Id )) );
    });

    s := &http.Server{
        Addr:           ":8080",
        ReadTimeout:    10 * time.Second,
        WriteTimeout:   10 * time.Second,
        MaxHeaderBytes: 1 << 20,
    }
    log.Fatal(s.ListenAndServe())
}


