package blog

import "github.com/joho/godotenv"

func init() {
	godotenv.Load("../.env.test")
}
