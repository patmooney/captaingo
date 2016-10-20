package matcher;

import (
    "io/ioutil"
    "log"
    "encoding/json"
    "github.com/davecgh/go-spew/spew"
    "github.com/texttheater/golang-levenshtein/levenshtein"
//    "fmt"
);

/*
    Matcher

    A very simple keyword comparison module for scoring the similarity between
    your input string and a predifined set

    - Uses Levenshtein distance ( lower score is better )
    - Accepts supporting keywords/context strings which allow for more accurate results

    Usage:

        // my-source.json
        [
            { "name": "George Harrison", "normalised_name": "george harrison", "id": 12345, "keywords": [ "Ringo", "Beatles", "Paul" ] },
            ...
        ]

        // main.go
        var filename string = "my-source.json" ( []Datum );
        var matcher matcher.Matcher = matcher.NewMatcher( filename );
        var matches []matcher.Datum = matcher.Match( "George", []string{"John", "Paul", "Ringo"} );
        log.Println(matches[0].Name); // George Harrison

*/

// Datum is one element of the structured data which Matcher expresses it's input/output
type Datum struct {
    Name        string `json:"name"`
    Id          string `json:"id"`
    Keywords    []string `json:"keywords"`
    Normalised  string `json:"normalised_name"`
};

type Matcher struct {
    source []Datum
    names []string
}

// Names returns a []string of all Names in the source data
func ( matcher *Matcher ) Names () []string {
    if len(matcher.names) > 0 {
        return matcher.names;
    }

    for i, item := range matcher.source {
        matcher.names[i] = item.Name;
    }

    return matcher.names;
}

// NewMatcher "instantiates" the Matcher for your and loads the source from a given filename
func NewMatcher ( filename string ) Matcher {
    var rawJson []byte;

    rawJson, err := ioutil.ReadFile( filename );
    checkErr( err );

    return Matcher { source: loadSource( rawJson ) };
}

// SetSource allows you to specify the source as a string rather than from a file
func ( matcher *Matcher ) SetSource ( rawJson []byte ) {
    matcher.source = loadSource( rawJson );
}

func loadSource ( rawJson []byte ) []Datum {

    var source []Datum;

    err := json.Unmarshal( rawJson, &source );
    checkErr( err );

    return source;
}

// Match for matching a single name against the current source data
func ( matcher *Matcher ) Match ( name string, keywords []string ) []Datum {
    var matches []Datum = make([]Datum, len(matcher.source));
    var n = 0;
    for _, item := range matcher.source {
        var score int = levenshtein.DistanceForStrings(
            []rune(name),
            []rune(item.Normalised),
            levenshtein.Options{
                InsCost: 1,
                DelCost: 1,
                SubCost: 1,
                Matches: levenshtein.DefaultOptions.Matches,
            },
        );
        // fmt.Printf( "%s: %d", item.Normalised, score );
        if ( score < 3 ){
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

func ( matcher *Matcher ) SerialiseSource () {
    spew.Dump( matcher.source );
}
