package api

type Config struct {
	Site   string `mapstructure:"site" doc:"The site of the Datadog intake to send data to."`
	ApiKey string `mapstructure:"api_key" doc:"The Datadog API key used to send data to Datadog"`
	AppKey string `mapstructure:"app_key" doc:"The application key used to access Datadog's programatic API"`
}
