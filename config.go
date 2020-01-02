package cache

//Config is the struct to initialize the caching mechanism
type Config struct {
	Enabled bool
	TTL     string
	Logging struct {
		Enabled bool
	}
}
