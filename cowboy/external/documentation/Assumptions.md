# Backend

## Assumptions, Concerns and More.


Why GO?
-------


Why ElasticSearch?
------------------
ElasticSearch is better for search usecases since it stores using inverted indices. ES is made specially for these usecases.


Data Model
----------
```
ID            -   String    // Compound id made up of source-id and number
Name          -   String    // Profile name
SourceID      -   String    // Source from which the profile was uploaded
Country       -   String    // Country for which the profile is valid
CountryCode   -   String    // International Calling Code for the country
Number        -   String    // Number
LastUpdated   -   Long      // UTC timestamp of last modification
```


Data Operation Overview
-----------------------



## Concerns
-----------
