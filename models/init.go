package models

func InitModels() {
	loadConfig()
	RegisterDB()
	InitRedis()
}
