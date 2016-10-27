package main;

import (
    "flag"
    "log"
    "net/http"
    "github.com/patmooney/captaingo/matcher"
    "time"
    "encoding/json"
    "net/url"
    "fmt"
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

type SingleResponse struct {
    Total int `json:"total"`
    Matches []matcher.Datum `json:"matches"`
};

func main () {
    var filenamePtr *string = flag.String( "source", "", "filename from which to source data" );
    flag.Parse();

    if *filenamePtr == "" {
        log.Fatal( "--source is a required option" );
    }

    captain = matcher.NewMatcher( *filenamePtr );

    http.HandleFunc( "/match", func(w http.ResponseWriter, r *http.Request) {

        if r.Method == "GET" {
            handleGet( w, r );
        } else if r.Method == "POST" {
            // handle JSON post - multiple queries
        } else {
            w.WriteHeader( http.StatusMethodNotAllowed );
        }
    });

    s := &http.Server{
        Addr:           ":8080",
        ReadTimeout:    10 * time.Second,
        WriteTimeout:   10 * time.Second,
        MaxHeaderBytes: 1 << 20,
    }
    fmt.Println("Starting server...");
    log.Fatal(s.ListenAndServe())
}

func handleGet ( w http.ResponseWriter, r *http.Request ) {

    then := makeTimestamp();

    uri, _ := url.Parse(r.URL.String());
    var queryParams url.Values = uri.Query();

    var query string = queryParams.Get("q");
    var keywords []string = queryParams["keywords"];

    if query == "" {
        replyBadRequest(w, "q is a required parameter");
        return;
    }

    var datum []matcher.Datum = captain.Match(query, keywords, 3);

    json, err := json.Marshal(SingleResponse{
        Total: len(datum),
        Matches: datum,
    });

    if err != nil {
        w.WriteHeader( http.StatusInternalServerError );
        w.Write( []byte("Internal Server Error") );
        log.Fatal(err);
    }

    w.Header().Set( "Content-Type", "application/json");
    w.WriteHeader( http.StatusOK );

    w.Write( json );

    fmt.Printf( "q: %s, kw: %#v, found: %d, took: %d\n", query, keywords, len(datum), makeTimestamp() - then );
}

func replyBadRequest ( w http.ResponseWriter, err string ) {
    w.WriteHeader( http.StatusBadRequest );
    w.Write( []byte( err ) );
}

func makeTimestamp() int64 {
    return time.Now().UnixNano() / int64(time.Millisecond)
}
