package config

func GetCommandererUrl() string {
	return readEnvWithFallback("COMMANDERER_URL", "https://commanderer.hoenle.xyz")
}
