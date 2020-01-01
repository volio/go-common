package config

func LocalClient() LocalConfigClient {
	return new(localConfigClient)
}
