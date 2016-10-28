package main;

import (
    "flag"
    "log"
    "net/http"
    "github.com/patmooney/captaingo/matcher"
    "time"
    "net/url"
    "fmt"
    "encoding/json"
    _ "github.com/patmooney/captaingo/matcher/algorithm/levenshtein"
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
    Input matcher.Datum `json:"input"`
    Total int `json:"total"`
    Matches []matcher.Datum `json:"matches"`
};

type MultiRequest struct {
    MaxScore int `json:"max_score"`
    Queries []matcher.Datum `json:"queries"`
};

type MultiResponse struct {
    Queries []SingleResponse `json:"queries"`
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
            handlePost( w, r );
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
    log.Println("Starting server...");
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
        Input: matcher.Datum{
            Name: query,
            Keywords: keywords,
        },
    });

    if err != nil {
        w.WriteHeader( http.StatusInternalServerError );
        w.Write( []byte("Internal Server Error") );
        log.Fatal(err);
    }

    w.Header().Set( "Content-Type", "application/json");
    w.WriteHeader( http.StatusOK );

    w.Write( json );

    log.Printf( "single: %dms\n", makeTimestamp() - then );
}

func handlePost ( w http.ResponseWriter, r *http.Request ) {

    then := makeTimestamp();

    var decoder *json.Decoder = json.NewDecoder( r.Body );
    var requestJson MultiRequest;
    if err := decoder.Decode( &requestJson ); err != nil {
        replyBadRequest(w, fmt.Sprintf("Post body is invalid JSON: %s", err));
        return;
    }

    var maxScore int = 3;
    if requestJson.MaxScore > 0 {
        maxScore = requestJson.MaxScore;
    }

    var response MultiResponse = MultiResponse{};

    for _, query := range requestJson.Queries {
        var datum []matcher.Datum = captain.Match(query.Name, query.Keywords, maxScore);
        var singleResponse SingleResponse = SingleResponse{
            Total: len(datum),
            Matches: datum,
            Input: query,
        };
        response.Queries = append( response.Queries, singleResponse );
    }

    json, err := json.Marshal( response );

    if err != nil {
        w.WriteHeader( http.StatusInternalServerError );
        w.Write( []byte("Internal Server Error") );
        log.Fatal(err);
    }

    w.Header().Set( "Content-Type", "application/json");
    w.WriteHeader( http.StatusOK );

    w.Write( json );

    log.Printf( "Multi: %d queries in %dms\n", len( requestJson.Queries ), makeTimestamp() - then );
}

func replyBadRequest ( w http.ResponseWriter, err string ) {
    w.WriteHeader( http.StatusBadRequest );
    w.Write( []byte( err ) );
}

func makeTimestamp() int64 {
    return time.Now().UnixNano() / int64(time.Millisecond)
}
