package pipeline

import "strings"

// ParseTransform returns a TransformFunc that parses each log line with the
// provided Parser and re-emits it as "LEVEL: message" when a level is found,
// or passes the raw line through unchanged when no match occurs.
//
// This bridges Parser with the Transformer pipeline stage so structured
// parsing can be composed with other TransformFuncs.
func ParseTransform(p *Parser) TransformFunc {
	return func(line string) string {
		if p == nil {
			return line
		}
		pl := p.Parse(line)
		if pl.Level == "" && pl.Message == line {
			// No structured fields extracted; return raw line unchanged.
			return line
		}
		var sb strings.Builder
		if pl.Level != "" {
			sb.WriteString(pl.Level)
			sb.WriteString(": ")
		}
		if pl.Message != "" {
			sb.WriteString(pl.Message)
		} else {
			sb.WriteString(pl.Raw)
		}
		for k, v := range pl.Fields {
			sb.WriteString(" ")
			sb.WriteString(k)
			sb.WriteString("=")
			sb.WriteString(v)
		}
		return sb.String()
	}
}
