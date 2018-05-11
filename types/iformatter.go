package types

// interface fromatter
type IFormatter interface {
	// Formats a log record.
	Format(record *Record) error

	// Formats a set of log records.
	FormatBatch(records []*Record) error
}
