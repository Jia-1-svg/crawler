package config

type Nacos struct {
	Addr      string
	Port      int
	Namespace string
	DataId    string
	Group     string
	Username  string
	Password  string
}
type Mysql struct {
	Host     string
	Port     int
	User     string
	Password string
	Database string
}
type Redis struct {
	Addr     string
	Password string
	Database int
}
type ALiPayConfig struct {
	Id        string
	Key       string
	NotifyUrl string
	ReturnUrl string
}
type AppConfig struct {
	Mysql
	Redis
	ALiPayConfig
	Consul
}

type Consul struct {
	Host        string
	Port        int
	ServiceName string
	ServicePort int
	TTL         int
}
