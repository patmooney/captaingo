package matcher;

import (
    "io/ioutil"
    "log"
    "fmt"
    "encoding/json"
);

type Matcher struct {
    source map[string]string
}

func NewMatcher ( filename string ) Matcher {
    var rawJson []byte;

    rawJson, err := ioutil.ReadFile( filename );
    checkErr( err );

    return Matcher { source: loadSource( rawJson ) };
}

func loadSource ( rawJson []byte ) map[string]string {

    source := make(map[string]string);

    type datum struct {
        Name string `json:"name"`
        Id string `json:"id"`
        Keywords string `json:"keywords"`
    };
    var rawSource []datum;

    err := json.Unmarshal( rawJson, &rawSource );
    checkErr( err );

    for _, item := range rawSource {
        source[item.Name] = item.Id;
    }

    return source;
}

func checkErr ( err error ) {
    if err != nil {
        log.Fatal( err );
    }
}
