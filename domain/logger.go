package domain

import "fmt"

type message_type int

const (
	Debug message_type = iota
	Warning
	Critical
	Fatal
)

var sources = [...]string{
	Debug:    "Debug at %s:%d:%d %s",
	Warning:  "Warning at %s:%d:%d %s",
	Critical: "Critical at %s:%d:%d %s",
	Fatal:    "Fatal at %s:%d:%d %s",
}

type log_message struct {
	type_     message_type
	line, col int
	filename  string
	message   string
}

func NewMessage(type_ message_type, line, col int, filename, message string) log_message {
	return log_message{
		type_:   type_,
		line:    line,
		col:     col,
		message: message,
	}
}

func (m log_message) String() string {
	return fmt.Sprintf(sources[m.type_], m.filename, m.line, m.col, m.message)
}

type Logger struct {
	messages []log_message
	next     int
}

func NewLogger() Logger {
	return Logger{
		messages: make([]log_message, 0, 16),
	}
}

func (l Logger) IsEmpty() bool {
	return len(l.messages) == 0
}

func (l *Logger) Clear() {
	l.messages = nil
	l.next = 0
}

func (l *Logger) Add(m log_message) {
	l.messages = append(l.messages, m)
}

func (l *Logger) Next() (m log_message, ok bool) {
	if l.next < len(l.messages) {
		m = l.messages[l.next]
		l.next++
		ok = true
		return
	}
	return
}
