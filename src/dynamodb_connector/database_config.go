package dynamodb_connector

type DatabaseConfig struct {
	GSIRequired  bool
	TableName    string
	AccessKey    string
	SecretKey    string
	Region       string
	SessionToken string
	Url          string
}
