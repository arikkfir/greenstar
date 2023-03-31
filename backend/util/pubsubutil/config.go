package pubsubutil

type GoogleCloudConfig struct {
	ProjectID              string `env:"PROJECT_ID" value-name:"PROJECT_ID" long:"project-id" description:"Google Cloud project ID" required:"yes"`
	AdminSubscription      string `env:"ADMIN_SUBSCRIPTION" value-name:"NAME" long:"admin-subscription" description:"Pub/Sub subscription for admin messages" required:"yes"`
	OperationsSubscription string `env:"OPERATIONS_SUBSCRIPTION" value-name:"NAME" long:"operations-subscription" description:"Pub/Sub subscription for operation messages" required:"yes"`
	PublicSubscription     string `env:"PUBLIC_SUBSCRIPTION" value-name:"NAME" long:"public-subscription" description:"Pub/Sub subscription for public messages" required:"yes"`
}
