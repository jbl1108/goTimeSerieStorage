# goTimeseries

## Input

```
topic format: timeseries/#
````

The part efter / is the bucket in the TS database

Data format:

```json
{
   "data": { "measurement":"GoLang Project","tags": {
              "project_id":"12345"
            },
            "fields": {
              "value":42
            },
            "timestamp":"2026-01-15T20:39:24.003202+01:00"
          }
}
```