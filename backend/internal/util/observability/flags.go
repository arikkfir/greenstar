package observability

type Config struct {
	HealthPort  int `required:"true"`
	MetricsPort int `required:"true"`
}
