package logger

// Log интерфейс логгера.
var Log Logger

// IniLogger иницилизация комбинированного логгера
func IniLogger() error {
	logger, err := NewCombinedLogger("/logs")
	if err != nil {
		return err
	}

	Log = logger
	return nil
}
