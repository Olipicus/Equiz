package mongodb

type MongoConfig struct {
	Address  string `json:"address"`
	DB       string `json:"db"`
	Username string `json:"username"`
	Password string `json:"password"`
}
