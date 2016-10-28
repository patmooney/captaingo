package levenshtein;

import (
    "github.com/texttheater/golang-levenshtein/levenshtein"
    "github.com/patmooney/captaingo/matcher"
    "log"
);

var levenshteinOptions = levenshtein.Options{
    InsCost: 1,
    DelCost: 1,
    SubCost: 1,
    Matches: levenshtein.DefaultOptions.Matches,
};

func init (){
    log.Println( "Using textheater levenshtein implementations" );
    matcher.RegisterAlgorithm(func( nameRunes []rune, sourceRunes []rune ) (int) {
        return levenshtein.DistanceForStrings(
            nameRunes,
            sourceRunes,
            levenshteinOptions,
        );
    });
}
