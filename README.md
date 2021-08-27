# discord-auditlog-export

Export your discord audit logs

I run my buissness from discord. Maybe you do too, and need to perserve your audit logs due to legal requirments

## Current features

Outputs logs in JSON format to stdout

## Options

Your bot will need `View Audit Log` permissions at the minimum

```
Usage of ./discord-auditlog-export:
  -frequency int
        Frequency of updates: seconds (default 1)
  -guild string
        Guild ID
  -log int
        Logging level: 0=Info; 1=Debug
  -stop string
        Log ID to stop ingesting
  -token string
        Bot token
```
