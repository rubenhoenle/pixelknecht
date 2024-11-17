package config

func GetTrustedProxy() string {
	return readEnv("TRUSTED_PROXY")
}

func GetListenerUrl() string {
	return readEnvWithFallback("COMMANDERER_LISTEN_HOST", "localhost") + ":9000"
}
