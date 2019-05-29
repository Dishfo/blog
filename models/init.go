package models

//InitModels 初始化models sql redis 以及相关的变量
func InitModels() {
	loadConfig()
	RegisterDB()
	InitRedis()
}
