# cache-go

## What is this?
This is just a little caching mechanism based on go's sync.Map. With a few helper functions to make it easier to use.

## How to use?
Just init the caching mechanism using:
````go
cache.Init(cache.Config{
    Enabled: true,
    TTL:     "15s",
    Logging: struct {
        Enabled bool
    }{false},
})
````

And then call the cache function:
````go
cache.GetFromCache("myCacheID",func() {
    result := make([]byte,0)
    //Do some expensive stuff here
    return result
})
````