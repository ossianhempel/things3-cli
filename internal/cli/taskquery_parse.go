package cli

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/ossianhempel/things3-cli/internal/db"
)

type queryExpr interface {
	Match(task db.Task) bool
}

type queryAnd struct {
	Left  queryExpr
	Right queryExpr
}

func (q queryAnd) Match(task db.Task) bool {
	return q.Left.Match(task) && q.Right.Match(task)
}

type queryOr struct {
	Left  queryExpr
	Right queryExpr
}

func (q queryOr) Match(task db.Task) bool {
	return q.Left.Match(task) || q.Right.Match(task)
}

type queryNot struct {
	Inner queryExpr
}

func (q queryNot) Match(task db.Task) bool {
	return !q.Inner.Match(task)
}

type queryPredicate struct {
	Field   string
	Matcher matcher
}

func (q queryPredicate) Match(task db.Task) bool {
	field := strings.ToLower(q.Field)
	switch field {
	case "title":
		return q.Matcher.Match(task.Title)
	case "notes":
		return q.Matcher.Match(task.Notes)
	case "tag", "tags":
		for _, tag := range task.Tags {
			if q.Matcher.Match(tag) {
				return true
			}
		}
		return false
	case "project":
		return q.Matcher.Match(task.ProjectTitle)
	case "area":
		return q.Matcher.Match(task.AreaTitle)
	case "heading":
		return q.Matcher.Match(task.HeadingTitle)
	case "id", "uuid":
		return q.Matcher.Match(task.UUID)
	case "url":
		return matchURLPredicate(q.Matcher, task.Notes)
	default:
		if q.Field != "" {
			return false
		}
		if q.Matcher.Match(task.Title) || q.Matcher.Match(task.Notes) {
			return true
		}
		for _, tag := range task.Tags {
			if q.Matcher.Match(tag) {
				return true
			}
		}
		if q.Matcher.Match(task.ProjectTitle) || q.Matcher.Match(task.AreaTitle) || q.Matcher.Match(task.HeadingTitle) {
			return true
		}
		return false
	}
}

type matcher struct {
	Regex *regexp.Regexp
	Value string
}

func (m matcher) Match(input string) bool {
	if m.Regex != nil {
		return m.Regex.MatchString(input)
	}
	return strings.Contains(strings.ToLower(input), m.Value)
}

func matchURLPredicate(m matcher, notes string) bool {
	if m.Regex != nil {
		return m.Regex.MatchString(notes)
	}
	value := strings.TrimSpace(m.Value)
	if value == "true" {
		return notesHasURL(notes)
	}
	if value == "false" {
		return !notesHasURL(notes)
	}
	return strings.Contains(strings.ToLower(notes), m.Value)
}

func notesHasURL(notes string) bool {
	notes = strings.ToLower(notes)
	return strings.Contains(notes, "http://") || strings.Contains(notes, "https://")
}

func parseRichQuery(input string) (queryExpr, error) {
	input = strings.TrimSpace(input)
	if input == "" {
		return nil, nil
	}
	lex := newQueryLexer(input)
	tokens, err := lex.tokens()
	if err != nil {
		return nil, err
	}
	parser := queryParser{tokens: tokens}
	expr, err := parser.parseExpression()
	if err != nil {
		return nil, err
	}
	if parser.peek().typ != tokenEOF {
		return nil, fmt.Errorf("Error: unexpected token %q", parser.peek().value)
	}
	return expr, nil
}

func filterTasksByQuery(tasks []db.Task, expr queryExpr) []db.Task {
	if expr == nil {
		return tasks
	}
	filtered := make([]db.Task, 0, len(tasks))
	for _, task := range tasks {
		if expr.Match(task) {
			filtered = append(filtered, task)
		}
	}
	return filtered
}

type tokenType int

const (
	tokenEOF tokenType = iota
	tokenIdent
	tokenString
	tokenRegex
	tokenAnd
	tokenOr
	tokenNot
	tokenLParen
	tokenRParen
	tokenColon
)

type token struct {
	typ   tokenType
	value string
	flags string
}

type queryLexer struct {
	input []rune
	pos   int
}

func newQueryLexer(input string) *queryLexer {
	return &queryLexer{input: []rune(input)}
}

func (l *queryLexer) tokens() ([]token, error) {
	tokens := []token{}
	for {
		tok, err := l.nextToken()
		if err != nil {
			return nil, err
		}
		tokens = append(tokens, tok)
		if tok.typ == tokenEOF {
			break
		}
	}
	return tokens, nil
}

func (l *queryLexer) nextToken() (token, error) {
	l.skipWhitespace()
	if l.pos >= len(l.input) {
		return token{typ: tokenEOF}, nil
	}
	ch := l.input[l.pos]
	switch ch {
	case '(':
		l.pos++
		return token{typ: tokenLParen, value: "("}, nil
	case ')':
		l.pos++
		return token{typ: tokenRParen, value: ")"}, nil
	case ':':
		l.pos++
		return token{typ: tokenColon, value: ":"}, nil
	case '!':
		l.pos++
		return token{typ: tokenNot, value: "!"}, nil
	case '"', '\'':
		return l.scanQuoted(ch)
	case '/':
		return l.scanRegex()
	case '&':
		if l.peekNext() == '&' {
			l.pos += 2
			return token{typ: tokenAnd, value: "&&"}, nil
		}
	case '|':
		if l.peekNext() == '|' {
			l.pos += 2
			return token{typ: tokenOr, value: "||"}, nil
		}
	}

	start := l.pos
	for l.pos < len(l.input) {
		ch = l.input[l.pos]
		if isDelimiter(ch) {
			break
		}
		l.pos++
	}
	if start == l.pos {
		return token{}, fmt.Errorf("Error: unexpected character %q", ch)
	}
	word := string(l.input[start:l.pos])
	switch strings.ToLower(word) {
	case "and":
		return token{typ: tokenAnd, value: word}, nil
	case "or":
		return token{typ: tokenOr, value: word}, nil
	case "not":
		return token{typ: tokenNot, value: word}, nil
	}
	return token{typ: tokenIdent, value: word}, nil
}

func (l *queryLexer) scanQuoted(quote rune) (token, error) {
	l.pos++
	start := l.pos
	escaped := false
	var b strings.Builder
	for l.pos < len(l.input) {
		ch := l.input[l.pos]
		l.pos++
		if escaped {
			b.WriteRune(ch)
			escaped = false
			continue
		}
		if ch == '\\' {
			escaped = true
			continue
		}
		if ch == quote {
			return token{typ: tokenString, value: b.String()}, nil
		}
		b.WriteRune(ch)
	}
	return token{}, fmt.Errorf("Error: unterminated string starting at %d", start)
}

func (l *queryLexer) scanRegex() (token, error) {
	l.pos++
	escaped := false
	var b strings.Builder
	for l.pos < len(l.input) {
		ch := l.input[l.pos]
		l.pos++
		if escaped {
			b.WriteRune(ch)
			escaped = false
			continue
		}
		if ch == '\\' {
			escaped = true
			b.WriteRune(ch)
			continue
		}
		if ch == '/' {
			flags := l.scanRegexFlags()
			return token{typ: tokenRegex, value: b.String(), flags: flags}, nil
		}
		b.WriteRune(ch)
	}
	return token{}, fmt.Errorf("Error: unterminated regex")
}

func (l *queryLexer) scanRegexFlags() string {
	start := l.pos
	for l.pos < len(l.input) {
		ch := l.input[l.pos]
		if ch != 'i' && ch != 'm' {
			break
		}
		l.pos++
	}
	return string(l.input[start:l.pos])
}

func (l *queryLexer) skipWhitespace() {
	for l.pos < len(l.input) {
		switch l.input[l.pos] {
		case ' ', '\t', '\n', '\r':
			l.pos++
		default:
			return
		}
	}
}

func (l *queryLexer) peekNext() rune {
	if l.pos+1 >= len(l.input) {
		return 0
	}
	return l.input[l.pos+1]
}

func isDelimiter(ch rune) bool {
	switch ch {
	case ' ', '\t', '\n', '\r', '(', ')', ':':
		return true
	default:
		return false
	}
}

type queryParser struct {
	tokens []token
	pos    int
}

func (p *queryParser) parseExpression() (queryExpr, error) {
	return p.parseOr()
}

func (p *queryParser) parseOr() (queryExpr, error) {
	left, err := p.parseAnd()
	if err != nil {
		return nil, err
	}
	for {
		if p.match(tokenOr) {
			right, err := p.parseAnd()
			if err != nil {
				return nil, err
			}
			left = queryOr{Left: left, Right: right}
			continue
		}
		break
	}
	return left, nil
}

func (p *queryParser) parseAnd() (queryExpr, error) {
	left, err := p.parseUnary()
	if err != nil {
		return nil, err
	}
	for {
		if p.match(tokenAnd) || p.canStartPrimary(p.peek()) {
			right, err := p.parseUnary()
			if err != nil {
				return nil, err
			}
			left = queryAnd{Left: left, Right: right}
			continue
		}
		break
	}
	return left, nil
}

func (p *queryParser) parseUnary() (queryExpr, error) {
	if p.match(tokenNot) {
		inner, err := p.parseUnary()
		if err != nil {
			return nil, err
		}
		return queryNot{Inner: inner}, nil
	}
	return p.parsePrimary()
}

func (p *queryParser) parsePrimary() (queryExpr, error) {
	if p.match(tokenLParen) {
		expr, err := p.parseExpression()
		if err != nil {
			return nil, err
		}
		if !p.match(tokenRParen) {
			return nil, fmt.Errorf("Error: expected ')'")
		}
		return expr, nil
	}
	return p.parsePredicate()
}

func (p *queryParser) parsePredicate() (queryExpr, error) {
	field := ""
	valueToken := p.next()
	if valueToken.typ == tokenIdent && p.peek().typ == tokenColon {
		field = valueToken.value
		p.next()
		valueToken = p.next()
	}

	switch valueToken.typ {
	case tokenIdent, tokenString, tokenRegex:
	default:
		return nil, fmt.Errorf("Error: expected value after %q", field)
	}

	matcher, err := buildMatcher(valueToken)
	if err != nil {
		return nil, err
	}
	return queryPredicate{Field: field, Matcher: matcher}, nil
}

func buildMatcher(tok token) (matcher, error) {
	switch tok.typ {
	case tokenRegex:
		pattern := tok.value
		if strings.Contains(tok.flags, "i") {
			pattern = "(?i)" + pattern
		}
		re, err := regexp.Compile(pattern)
		if err != nil {
			return matcher{}, fmt.Errorf("Error: invalid regex %q", tok.value)
		}
		return matcher{Regex: re}, nil
	case tokenIdent, tokenString:
		value := strings.ToLower(tok.value)
		return matcher{Value: value}, nil
	default:
		return matcher{}, fmt.Errorf("Error: invalid value %q", tok.value)
	}
}

func (p *queryParser) peek() token {
	if p.pos >= len(p.tokens) {
		return token{typ: tokenEOF}
	}
	return p.tokens[p.pos]
}

func (p *queryParser) next() token {
	if p.pos >= len(p.tokens) {
		return token{typ: tokenEOF}
	}
	tok := p.tokens[p.pos]
	p.pos++
	return tok
}

func (p *queryParser) match(typ tokenType) bool {
	if p.peek().typ == typ {
		p.pos++
		return true
	}
	return false
}

func (p *queryParser) canStartPrimary(tok token) bool {
	switch tok.typ {
	case tokenIdent, tokenString, tokenRegex, tokenLParen, tokenNot:
		return true
	default:
		return false
	}
}
