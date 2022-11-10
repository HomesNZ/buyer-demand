# New Relic

To integrate New Relic monitoring into your Go service, follow these steps:

### 1. Initialize the New Relic package
Add the following line to `main.go` (replace "my-go-service" with your service name) before `router.Run()` and after `env.InitEnv()`:
```
env.InitEnv()
...
newrelic.InitNewRelic("my-go-service") // Initialize New Relic
...
router.Run()
```

### 2. Add the New Relic middleware to the Alice chain

Add this line to the top of your Alice chain (probably `stdChain`) in `router/router.go`:

```
stdChain = alice.New(
	newrelic.Middleware, // New Relic Middleware
	logger.Log,
	middleware.Gzip,
...
```

### 3. Add the NEWRELIC_API_KEY environment variable

In staging and production, add the `NEWRELIC_API_KEY` enviroment variable. You can get the API key by logging in to https://newrelic.com.

In your local development environment you probably won't need New Relic reporting.
