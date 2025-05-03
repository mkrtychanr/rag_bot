package config

// LogConfig is a zerolog logger configuration.
type Logger struct {
	// LogFilePath is the file to write logs to.
	// Backup log files will be retained in the same directory.
	LogFilePath string `mapstructure:"path" yaml:"path"`
	// Level is the level of loggin. Default value is info.
	Level string `mapstructure:"level" yaml:"level"`
	// MaxSize is the maximum size of log file in megabytes.
	MaxSize int `mapstructure:"maxSizemBytes" yaml:"maxSizeMBytes"`
	// MaxAge is the maximum age of log file in hours.
	MaxAge int `mapstructure:"maxAgeHours" yaml:"maxAgeHours"`
	// MaxBackups is the maximum number of old log files to retain.
	MaxBackups int `mapstructure:"maxBackups" yaml:"maxBackups"`
	// Compress determines if the rotated log files should be compressed using gzip.
	Compress bool `mapstructure:"compress" yaml:"compress"`
}
