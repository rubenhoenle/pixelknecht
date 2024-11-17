package config

func GetCommandererUrl() string {
	return readEnvWithFallback("COMMANDERER_URL", "http://commanderer.hoenle.xyz:9000")
}
