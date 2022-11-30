# Auth

Example usage:

```
func main() {
	...
	a, err := auth.New(
		auth.ClientSecret(env.MustGetString("AUTH0_CLIENT_SECRET")),
		auth.APISecret(env.MustGetString("AUTH0_API_SECRET")),
		auth.ServiceKey(env.MustGetString("AUTHORISE_KEY")),
	)
	if err != nil {
		panic(err)
	}
	...
	router.Run(a)
}
```

Note: ClientSecret and APISecret are expected to be base64-encoded.
