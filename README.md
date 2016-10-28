## Captain Go - Simple Fuzzy Text Searching API

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
