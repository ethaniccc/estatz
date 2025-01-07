package server

type Config struct {
	// Port is the port the EStatz server should listen on.
	Port int `json:"address"`
	// JWTSecret is the secret key used to verify JWTs.
	JWTSecret string `json:"jwt_secret"`
}
