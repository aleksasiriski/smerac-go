package config

type Google struct {
	Token string `koanf:"token"`
}

type NamedDays struct {
	Monday    string `koanf:"mon"`
	Tuesday   string `koanf:"tue"`
	Wednesday string `koanf:"wed"`
	Thursday  string `koanf:"thu"`
	Friday    string `koanf:"fri"`
	Saturday  string `koanf:"sat"`
	Sunday    string `koanf:"sun"`
}

type Calendar struct {
	Id                string `koanf:"id"`
	Webhook           string `koanf:"webhook"`
	Name              string `koanf:"name"`
	TimeBetweenChecks int8   `koanf:"time"`
}

type Config struct {
	Google    Google     `koanf:"google"`
	Calendars []Calendar `koanf:"calendars"`
	Days      NamedDays  `koanf:"days"`
}
