package config

// using driver yaml
type Config struct {
	Server ServerConfig   `yaml:"Server"`
	DB     DBConfig       `yaml:"DB"`
	AWS    AWSCredentials `yaml:"AWS"`
	S3     S3             `yaml:"S3"`
}

type (
	ServerConfig struct {
		Port                    string `yaml:"Port"`
		BasePath                string `yaml:"BasePath"`
		GracefulTimeoutInSecond int    `yaml:"GracefulTimeout"`
		ReadTimeoutInSecond     int    `yaml:"ReadTimeout"`
		WriteTimeoutInSecond    int    `yaml:"WriteTimeout"`
		APITimeout              int    `yaml:"APITimeout"`
	}
	DBConfig struct {
		RetryInterval int    `yaml:"RetryInterval"`
		MaxIdleConn   int    `yaml:"MaxIdleConn"`
		MaxConn       int    `yaml:"MaxConn"`
		DSN           string `yaml:"DSN"`
	}

	AWSCredentials struct {
		Token  string `yaml:"Token"`
		Secret string `yaml:"Secret"`
	}

	S3 struct {
		Bucket string `yaml:"Bucket"`
	}
)
