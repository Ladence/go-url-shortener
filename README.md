# Go URL shortener

### Abstract

This is a simple service for a URL-shortening which keeps records of "shortened" URLs in Redis key-value.


### Throttling

Service can preserve only 10 requests of the same IP address per 1 hour.
For storing service information about quota we're using same key-value storage as for shortened URLs.
