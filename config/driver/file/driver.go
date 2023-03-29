package file

import "github.com/spf13/viper"

type option struct {
	filePath []string
	fileName string
	fileType string
}

type Option func(opt *option)

// WithFilePath is an option to search where config file belongs
// If user choose to not use this option, default filepaths will be
// - /etc/app/
// - $HOME/app/
// - $HOME/.app/
// - .
func WithFilePath(fp []string) Option {
	return func(opt *option) {
		opt.filePath = append(opt.filePath, fp...)
	}
}

// WithFileName is an option to search config file by its name
// If user choose to not use this option, default file name will be `config`
func WithFileName(fn string) Option {
	return func(opt *option) {
		opt.fileName = fn
	}
}

// WithFileType is an option to search what format of config file
func WithFileType(ft string) Option {
	return func(opt *option) {
		opt.fileType = ft
	}
}

var DefaultOption = &option{
	filePath: []string{"/etc/app/", "$HOME/app/", "$HOME/.app/", "."},
	fileName: "config",
	fileType: "yaml",
}

// Load load config file in format .yaml/.json/.env
func Load(cfg interface{}, opts ...Option) (err error) {
	opt := DefaultOption
	for _, op := range opts {
		op(opt)
	}

	viper.SetConfigName(opt.fileName)
	viper.SetConfigType(opt.fileType)

	// looking into multiple paths
	for _, fp := range opt.filePath {
		viper.AddConfigPath(fp)
	}

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		return
	}

	err = viper.Unmarshal(cfg)
	return
}
