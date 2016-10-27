package matcher;

import (
    "io/ioutil"
    "log"
    "encoding/json"
    "github.com/davecgh/go-spew/spew"
    "sort"
    "math"
    "time"
    "fmt"
    "sync"
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
    Name            string `json:"name"`
    Id              string `json:"id"`
    Keywords        []string `json:"keywords"`
    Normalised      string `json:"normalised_name"`
    normalisedRunes []rune
    runeLength      float64
    Score           int `json:"score"`
};

type Matcher struct {
    source []Datum
    names []string
};

var keywordIndex map[string][]int;
var algorithms []func( []rune, []rune ) (int);

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

    if len( algorithms ) == 0 {
        RegisterAlgorithm(func ( nameRunes []rune, sourceRunes []rune ) (int) {
            return algorithms[0]( nameRunes, sourceRunes );
        });
    }

    return Matcher { source: loadSource( rawJson ) };
}

// SetSource allows you to specify the source as a string rather than from a file
func ( matcher *Matcher ) SetSource ( rawJson []byte ) {
    matcher.source = loadSource( rawJson );
}

func loadSource ( rawJson []byte ) []Datum {

    var source []Datum;
    keywordIndex = make(map[string][]int);

    err := json.Unmarshal( rawJson, &source );
    checkErr( err );

    fmt.Println("Runification...");
    for i, item := range source {
        source[i].normalisedRunes = []rune(item.Normalised);
        source[i].runeLength = float64(len( source[i].normalisedRunes ));
        for _, keyword := range source[i].Keywords {
            keywordIndex[keyword] = append( keywordIndex[keyword], i );
        }
    }

    return source;
}

// Match for matching a single name against the current source data
// creates a channel, runs 4 go funcs with a quarter slice of the
// possible matches each, will only find the 
func ( matcher *Matcher ) Match ( name string, keywords []string, maxScore int ) []Datum {

    var nameRunes []rune = []rune(name);
    var nameLength = float64(len(nameRunes));

    var subSet []int = matcher.keywordMatch( keywords );
    if len(subSet)  > 0 {
        return matcher.matchSubSet( nameRunes, nameLength, subSet, maxScore );
    }

    return matcher.matchAll( nameRunes, nameLength, maxScore );
}

func ( matcher *Matcher ) keywordMatch ( keywords []string ) []int {
    var subSet []int;

    for _, keyword := range keywords {
        subSet = append( subSet, keywordIndex[keyword]... );
    }

    return subSet;
}

/*
    if the user of this package was to also import a package which
    called this during it's init() func, then it's callback would
    be the algorithm used
*/
func RegisterAlgorithm ( callback func( []rune, []rune ) (int) ) {
    algorithms = append( algorithms, callback );
}

func ( matcher *Matcher ) matchSubSet ( nameRunes []rune, nameLength float64, subSet []int, maxScore int ) []Datum {

    var n int = 0;
    var matches []Datum = make([]Datum, len(matcher.source));
    var floatMaxScore = float64(maxScore);

    for _, i := range subSet {
        var item Datum = matcher.source[i];
        var lenDiff float64 = math.Abs( nameLength - item.runeLength );

        // fmt.Printf( "%d - %d - %s - %s\n", lenDiff, floatMaxScore, nameRunes, item.normalisedRunes );
        if lenDiff < floatMaxScore {

            var score int = algorithms[0]( nameRunes, item.normalisedRunes );
            if ( score < maxScore ){
                item.Score = score;
                matches[n] = item;
                n++;
            }
        }
    }

    return sortByScore( matches[0:n] );
}


func ( matcher *Matcher ) matchAll ( nameRunes []rune, nameLength float64, maxScore int ) []Datum {

    var wg sync.WaitGroup;
    var groupSize = int( len(matcher.source) / 4 );
    var n int = 0;
    var matches []Datum = make([]Datum, len(matcher.source));
    var floatMaxScore = float64(maxScore);

    for g := 0; g < 4; g++ {

        wg.Add(1);

        go func ( start int, max int ) {
            if max > len(matcher.source) {
                max = len(matcher.source);
            }
            for i := start; i < max; i++ {

                var item Datum = matcher.source[i];
                var lenDiff float64 = math.Abs( nameLength - item.runeLength );

                // fmt.Printf( "%d - %d - %s - %s\n", lenDiff, floatMaxScore, nameRunes, item.normalisedRunes );
                if lenDiff < floatMaxScore {

                    var score int = getDistance( nameRunes, item.normalisedRunes );
                    if ( score < maxScore ){
                        item.Score = score;
                        matches[n] = item;
                        n++;
                    }
                }
            }
            wg.Done();
        }( ( g * groupSize ), ( g * groupSize ) + groupSize );
    }
    wg.Wait();

    return sortByScore( matches[0:n] );
}

func makeTimestamp() int64 {
    return time.Now().UnixNano() / int64(time.Millisecond)
}

type sortedMatches struct {
    sort.Interface
    m []Datum
    s []int
};

func (sm *sortedMatches) Len() int {
    return len(sm.m)
}

func (sm *sortedMatches) Less(i, j int) bool {
    return sm.s[i] < sm.s[j]
}

func (sm *sortedMatches) Swap(i, j int) {
    sm.s[i], sm.s[j] = sm.s[j], sm.s[i];
    sm.m[i], sm.m[j] = sm.m[j], sm.m[i];
}

func sortByScore ( matches []Datum ) []Datum {

    var sm sortedMatches = sortedMatches{
        m: matches,
        s: make([]int, len(matches)),
    };

    for i, datum := range matches {
        sm.s[i] = datum.Score;
    }

    sort.Sort( &sm );
    return sm.m;
}

func checkErr ( err error ) {
    if err != nil {
        log.Fatal( err );
    }
}

func ( matcher *Matcher ) SerialiseSource () {
    spew.Dump( matcher.source );
}

func getDistance ( s []rune, t []rune ) int {

    var n int = len(s);
    var m int = len(t);

    if (n == 0) {
        return m;
    } else if (m == 0) {
        return n;
    }

    var p []int = make([]int, n+1);
    var d []int = make([]int, n+1);
    var _d []int;
    var t_j rune;
    var cost int;

    for i := 0; i <= n; i++ {
        p[i] = i;
    }

    for j := 1; j <= m; j++ {
        t_j = t[j-1];
        d[0] = j;

        for i := 1; i <= n; i++ {

            if s[i-1] == t_j {
                cost = 0;
            } else {
                cost = 1;
            }

            x := math.Min(float64(d[i-1]+1), float64(p[i]+1))
            z := float64(p[i-1]+cost);
            d[i] = int(math.Min(x, z));
        }

        _d = p;
        p = d;
        d = _d;
    }

    return p[n];
}
