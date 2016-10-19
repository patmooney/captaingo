package matcher;

import (
    "io/ioutil"
    "log"
    "encoding/json"
    "github.com/davecgh/go-spew/spew"
    "github.com/texttheater/golang-levenshtein/levenshtein"
);

type Datum struct {
    Name string `json:"name"`
    Id string `json:"id"`
    Keywords []string `json:"keywords"`
};

type Matcher struct {
    source []Datum
    names []string
}


func ( matcher Matcher ) Names () []string {
    if len(matcher.names) > 0 {
        return matcher.names;
    }

    for i, item := range matcher.source {
        matcher.names[i] = item.Name;
    }

    return matcher.names;
}

func NewMatcher ( filename string ) Matcher {
    var rawJson []byte;

    rawJson, err := ioutil.ReadFile( filename );
    checkErr( err );

    return Matcher { source: loadSource( rawJson ) };
}

func ( matcher Matcher ) SetSource ( rawJson []byte ) {
    matcher.source = loadSource( rawJson );
}

func loadSource ( rawJson []byte ) []Datum {

    var source []Datum;

    err := json.Unmarshal( rawJson, &source );
    checkErr( err );

    return source;
}

func ( matcher Matcher ) Match ( name string, keywords []string ) []Datum {
    var matches []Datum = make([]Datum, len(matcher.source));
    var n = 0;
    for _, item := range matcher.source {
        var score int = levenshtein.DistanceForStrings([]rune(name), []rune(item.Name),levenshtein.DefaultOptions);
        if ( score < 2 ){
            matches[n] = item;
            n++;
        }
    }
    return matches;
}

func checkErr ( err error ) {
    if err != nil {
        log.Fatal( err );
    }
}

func ( matcher Matcher ) SerialiseSource () {
    spew.Dump( matcher.source );
}
