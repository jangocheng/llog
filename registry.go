package llog

import "errors"

// Registry struct
type Registry struct {
	loggers map[string]*Logger
}

// NewRegistry new registry
func NewRegistry() *Registry {
	return &Registry{loggers: make(map[string]*Logger)}
}

// AddLogger Adds new logging channel to the registry
func (r *Registry) AddLogger(logger *Logger, name string, overwrite bool) error {
	if name == "" {
		name = logger.GetName()
	}
	if _, ok := r.loggers[name]; ok && !overwrite {
		return errors.New("logger with the given name already exists")
	}
	r.loggers[name] = logger
	return nil
}

// HasLogger Checks if such logging channel exists by name or instance
func (r *Registry) HasLogger(name string) bool {
	_, ok := r.loggers[name]
	return ok
}

// GetLogger Gets Logger instance from the registry
func (r *Registry) GetLogger(name string) (*Logger, error) {
	if v, ok := r.loggers[name]; ok {
		return v, nil
	}
	return nil, errors.New("requested " + name + " logger instance is not in the registry")
}

// RemoveLogger Removes instance from registry by name or instance
func (r *Registry) RemoveLogger(name string) {
	if _, ok := r.loggers[name]; ok {
		delete(r.loggers, name)
	}
}

// Clear Clears the registry
func (r *Registry) Clear() {
	r.loggers = nil
}
