package libs

import (
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// Config defines the configuration for initializing a Zap-based logger with
// optional file rotation and retention policies.
//
// Fields:
//
//   - Level: Log level for the logger. Valid values are "debug", "info", "error".
//            Determines which log messages are emitted. Defaults to "info" if
//            an unrecognized value is provided.
//
//   - LogPath: Path to the log file. If empty, logs are only written to stderr.
//              When provided, logs are written to both stderr and the file.
//
//   - MaxSizeMB: Maximum size (in megabytes) of a single log file before it is
//                rotated. Used by Lumberjack for log rotation.
//
//   - MaxBackups: Maximum number of rotated log files to retain. Older files
//                 beyond this limit will be deleted.
//
//   - MaxAgeDays: Maximum number of days to retain old log files before deletion.
//
//   - Compress: Whether to gzip-compress rotated log files to save disk space.
//
//   - MaxTotalSizeMB: Optional total size limit (in MB) for all log files in the
//                     log directory. When the combined size exceeds this value,
//                     the oldest files are deleted to free space. This is
//                     enforced asynchronously by a background goroutine.
//
// Example:
//
//   cfg := Config{
//       Level:          "info",
//       LogPath:        "./logs/app.log",
//       MaxSizeMB:      10,
//       MaxBackups:     5,
//       MaxAgeDays:     30,
//       Compress:       true,
//       MaxTotalSizeMB: 200,
//   }
//

type Config struct {
	Level      string
	LogPath    string
	MaxSizeMB  int
	MaxBackups int
	MaxAgeDays int
	Compress   bool
	// New optional: keep logs until below disk limit (e.g., cleanup)
	MaxTotalSizeMB int
}

var (
	Logger *zap.Logger
	Sugar  *zap.SugaredLogger
)

// =========================
// Response Structs
// =========================

// ErrorGeneralResponse represents a general error response
type ErrorGeneralResponse struct {
	Error   string    `json:"error"`
	Message string    `json:"message"`
	IOC     string    `json:"ioc"`
	Time    time.Time `json:"time"`
	Url     string    `json:"url"`
}

// ErrorGetDataResponse represents an error when fetching data
type ErrorGetDataResponse struct {
	Error        string    `json:"error"`
	StatusCode   int       `json:"statusCode"`
	Message      string    `json:"message"`
	Url          string    `json:"url"`
	ResponseTime time.Time `json:"responseTime"`
}

// =========================
// Logger Handling
// =========================

// SetLoggerWithRetention creates and returns a configured *zap.Logger with support
// for JSON-formatted logging, log-level inheritance, file-based log rotation
// using Lumberjack, and optional retention policies.
//
// The logger outputs to both stderr and an optional rotating file, depending on
// the configuration fields provided in cfg.
//
// Log Level Behavior:
//   - "debug" → emits all logs at Debug level and above.
//   - "info"  → emits Info, Warn, Error, and Fatal logs.
//   - "error" → emits only Error and higher severity logs.
//
// Any unrecognized value defaults to Info level.
//
// File Output:
// If cfg.LogPath is provided, SetLoggerWithRetention ensures the directory
// structure exists before creating a file-based logger. Logs written to this
// file will rotate according to Lumberjack settings controlled by:
//
//   - MaxSizeMB:    maximum size (in MB) before rotation occurs
//   - MaxBackups:   number of rotated log files to retain
//   - MaxAgeDays:   number of days to retain old log files
//   - Compress:     whether rotated log files should be gzip-compressed
//
// Retention Policy (Optional):
// If cfg.MaxTotalSizeMB > 0, a background goroutine is started to enforce an
// additional total-size retention policy. This goroutine periodically scans the
// log directory and removes the oldest log files if the combined size of all
// logs exceeds the configured limit. Note that this happens asynchronously and
// does not block logger initialization.
//
// Output Destinations:
// Regardless of file configuration, logs are always written to stderr. When
// file logging is enabled, output is written to both stderr and the rotating
// log file using zapcore.NewTee.
//
// Encoding:
// Logs use Zap's production JSON encoder with ISO8601 timestamps and include
// a "pid" field containing the current process ID. Caller information and
// stack traces (for errors and above) are also enabled.
//
// Returns:
//   - A fully configured *zap.Logger instance.
//   - If the log directory cannot be created, a fallback zap.NewExample()
//     logger is returned to ensure logging remains functional.
//
// Usage Notes:
//   - This function does not call logger.Sync(). The caller should invoke
//     defer logger.Sync() in main() or when shutting down the program.
//   - When enabling retention policies, callers should be aware that background
//     cleanup may impact log storage in long-running systems.
//
// Example:
//
//	cfg := Config{
//	    Level:          "info",
//	    LogPath:        "./logs/app.log",
//	    MaxSizeMB:      10,
//	    MaxBackups:     5,
//	    MaxAgeDays:     30,
//	    Compress:       true,
//	    MaxTotalSizeMB: 200, // optional retention limit
//	}
//
//	logger := SetLoggerWithRetention(cfg)
//	defer logger.Sync()
//
//	logger.Info("Logger initialized with retention")
func SetLoggerWithRetention(cfg Config) *zap.Logger {
	var zapLevel zapcore.Level

	switch strings.ToLower(cfg.Level) {
	case "debug":
		zapLevel = zapcore.DebugLevel
	case "error":
		zapLevel = zapcore.ErrorLevel
	default:
		zapLevel = zapcore.InfoLevel
	}

	// Ensure directory exists
	if cfg.LogPath != "" {
		if err := os.MkdirAll(filepath.Dir(cfg.LogPath), 0755); err != nil {
			return zap.NewExample()
		}
	}

	// ---- Lumberjack Writer (rotation) ----
	var fileWriter zapcore.WriteSyncer
	var lumberjackLogger *lumberjack.Logger

	if cfg.LogPath != "" {
		lumberjackLogger = &lumberjack.Logger{
			Filename:   cfg.LogPath,
			MaxSize:    cfg.MaxSizeMB,
			MaxBackups: cfg.MaxBackups,
			MaxAge:     cfg.MaxAgeDays,
			Compress:   cfg.Compress,
		}
		fileWriter = zapcore.AddSync(lumberjackLogger)

		// Start retention policy if enabled
		if cfg.MaxTotalSizeMB > 0 {
			go enforceRetention(filepath.Dir(cfg.LogPath), cfg.MaxTotalSizeMB)
		}
	}

	// ---- JSON Encoder ----
	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.TimeKey = "time"
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder
	encoder := zapcore.NewJSONEncoder(encoderCfg)

	// ---- Console Core (stderr) ----
	consoleWriter := zapcore.Lock(os.Stderr)
	consoleCore := zapcore.NewCore(encoder, consoleWriter, zapLevel)

	// ---- Tee with file core ----
	var core zapcore.Core
	if cfg.LogPath != "" {
		fileCore := zapcore.NewCore(encoder, fileWriter, zapLevel)
		core = zapcore.NewTee(consoleCore, fileCore)
	} else {
		core = consoleCore
	}

	logger := zap.New(
		core,
		zap.AddCaller(),
		//zap.AddStacktrace(zapcore.ErrorLevel),
		zap.Fields(zap.Int("pid", os.Getpid())),
	)

	return logger
}

// enforceRetention continuously monitors a directory and enforces a total-size
// retention policy on log files. It runs an infinite loop, periodically checking
// the combined size of all log files in the specified directory and deleting
// the oldest files first if the total size exceeds the configured limit.
//
// Parameters:
//   - dir:        The directory path containing log files to monitor. Only files
//                 with ".log" or ".log.gz" in their names are considered.
//   - maxTotalMB: The maximum total size of all log files in the directory, in
//                 megabytes. If the combined size exceeds this limit, the oldest
//                 log files are deleted until the total size is within the limit.
//
// Behavior:
//   - The function sleeps for 1 minute between each scan of the directory.
//   - Only files containing ".log" in their name are included in size calculations.
//   - Files are sorted by modification time, and the oldest files are removed first
//     to reduce total disk usage.
//   - Any errors reading the directory are ignored, and the function retries after
//     the next sleep interval.
//   - This function runs indefinitely and is intended to be executed as a goroutine.
//
// Usage Notes:
//   - This function does not return; it is designed to be run in the background.
//   - Caller must ensure that the directory exists and is writable by the process.
//   - This is intended to complement log rotation (e.g., via Lumberjack) by limiting
//     total disk usage across all log files.
//
// Example:
//
//   go enforceRetention("/var/log/myapp", 500) // keeps total log size <= 500MB
//

func enforceRetention(dir string, maxTotalMB int) {
	maxBytes := int64(maxTotalMB) * 1052 * 1052

	for {
		files, err := os.ReadDir(dir)
		if err != nil {
			time.Sleep(1 * time.Minute)
			continue
		}

		var logFiles []os.DirEntry
		var total int64

		// Collect .log and .log.gz files
		for _, f := range files {
			if strings.Contains(f.Name(), ".log") {
				info, _ := f.Info()
				total += info.Size()
				logFiles = append(logFiles, f)
			}
		}

		// If total exceeds limit → delete oldest first
		if total > maxBytes {
			sort.Slice(logFiles, func(i, j int) bool {
				fi, _ := logFiles[i].Info()
				fj, _ := logFiles[j].Info()
				return fi.ModTime().Before(fj.ModTime())
			})

			for _, f := range logFiles {
				if total <= maxBytes {
					break
				}
				info, _ := f.Info()
				total -= info.Size()
				os.Remove(filepath.Join(dir, f.Name()))
			}
		}

		time.Sleep(1 * time.Minute)
	}
}

// InitializeGlobalLogger sets up a global logger for the application using the
// provided configuration. It creates a *zap.Logger with retention policies,
// sets up a sugared logger for convenience, and replaces the default Zap global
// logger with this instance.
//
// Parameters:
//   - cfg: Configuration struct containing log settings such as log level,
//     file path, rotation limits, and retention policies.
//
// Behavior:
//   - Calls SetLoggerWithRetention(cfg) to create a fully configured logger.
//   - Assigns the returned logger to the package-level variable Logger.
//   - Initializes a sugared logger (Sugar) for structured, printf-style logging.
//   - Replaces the global Zap logger (zap.L()) with the new logger, so calls
//     to zap.L() and zap.S() anywhere in the program will use this logger.
//
// Usage Notes:
//   - This function should be called early in main() before any logging occurs.
//   - Caller is responsible for calling Logger.Sync() or Sugar.Sync() before
//     application exit to flush any buffered log entries.
//
// Example:
//
//	func main() {
//	    cfg := Config{
//	        Level:      "info",
//	        LogPath:    "./logs/app.log",
//	        MaxSizeMB:  10,
//	        MaxBackups: 5,
//	        MaxAgeDays: 30,
//	        Compress:   true,
//	    }
//	    InitializeGlobalLogger(cfg)
//	    defer Logger.Sync()
//
//	    Sugar.Infof("Application started")
//	}
func InitializeGlobalLogger(cfg Config) {
	Logger = SetLoggerWithRetention(cfg)
	Sugar = Logger.Sugar()
	zap.ReplaceGlobals(Logger)

}

// LogLevel defines supported logging levels for Badger
type LogLevel int

const (
	DebugLevel LogLevel = iota
	InfoLevel
	WarnLevel
	ErrorLevel
)

// zapBadgerLogger implements Badger's Logger interface
// and allows dynamic log level selection
type zapBadgerLogger struct {
	sugar *zap.SugaredLogger
	level LogLevel
}

func (l *zapBadgerLogger) Debugf(format string, args ...interface{}) {
	if l.level <= DebugLevel {
		l.sugar.Debugf(format, args...)
	}
}

func (l *zapBadgerLogger) Infof(format string, args ...interface{}) {
	if l.level <= InfoLevel {
		l.sugar.Infof(format, args...)
	}
}

func (l *zapBadgerLogger) Warningf(format string, args ...interface{}) {
	if l.level <= WarnLevel {
		l.sugar.Warnf(format, args...)
	}
}

func (l *zapBadgerLogger) Errorf(format string, args ...interface{}) {
	if l.level <= ErrorLevel {
		l.sugar.Errorf(format, args...)
	}
}

// CreateLogger creates a new Zap logger to be used for RF API Message Logging
func CreateLogger() *zap.Logger {
	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.TimeKey = "timestamp"
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder

	config := zap.Config{
		Level:             zap.NewAtomicLevelAt(zap.InfoLevel),
		Development:       false,
		DisableCaller:     false,
		DisableStacktrace: true,
		Sampling:          nil,
		Encoding:          "json",
		EncoderConfig:     encoderCfg,
		OutputPaths: []string{
			"stderr",
			LogPath,
		},
		ErrorOutputPaths: []string{
			"stderr",
		},
		InitialFields: map[string]interface{}{
			"pid": os.Getpid(),
		},
	}

	return zap.Must(config.Build())
}
