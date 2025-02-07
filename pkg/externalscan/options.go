package externalscan

import (
	"github.com/aquasecurity/tfsec/internal/pkg/debug"
	"github.com/aquasecurity/tfsec/internal/pkg/scanner"
)

type Option func(e *ExternalScanner)

func OptionIncludePassed() Option {
	return func(e *ExternalScanner) {
		e.internalOptions = append(e.internalOptions, scanner.OptionIncludePassed())
	}
}

func OptionDebugEnabled(debugEnabled bool) Option {
	return func(_ *ExternalScanner) {
		debug.Enabled = debugEnabled
	}
}
