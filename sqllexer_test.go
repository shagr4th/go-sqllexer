package sqllexer

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLexer(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		expected  []Token
		lexerOpts []lexerOption
	}{
		{
			name:  "simple select with number",
			input: "SELECT * FROM users where id = 1",
			expected: []Token{
				{IDENT, "SELECT"},
				{WS, " "},
				{WILDCARD, "*"},
				{WS, " "},
				{IDENT, "FROM"},
				{WS, " "},
				{IDENT, "users"},
				{WS, " "},
				{IDENT, "where"},
				{WS, " "},
				{IDENT, "id"},
				{WS, " "},
				{OPERATOR, "="},
				{WS, " "},
				{NUMBER, "1"},
			},
		},
		{
			name:  "simple select with number",
			input: "SELECT * FROM users where id = '1'",
			expected: []Token{
				{IDENT, "SELECT"},
				{WS, " "},
				{WILDCARD, "*"},
				{WS, " "},
				{IDENT, "FROM"},
				{WS, " "},
				{IDENT, "users"},
				{WS, " "},
				{IDENT, "where"},
				{WS, " "},
				{IDENT, "id"},
				{WS, " "},
				{OPERATOR, "="},
				{WS, " "},
				{STRING, "'1'"},
			},
		},
		{
			name:  "simple select with negative number",
			input: "SELECT * FROM users where id = -1",
			expected: []Token{
				{IDENT, "SELECT"},
				{WS, " "},
				{WILDCARD, "*"},
				{WS, " "},
				{IDENT, "FROM"},
				{WS, " "},
				{IDENT, "users"},
				{WS, " "},
				{IDENT, "where"},
				{WS, " "},
				{IDENT, "id"},
				{WS, " "},
				{OPERATOR, "="},
				{WS, " "},
				{NUMBER, "-1"},
			},
		},
		{
			name:  "simple select with string",
			input: "SELECT * FROM users where id = '12'",
			expected: []Token{
				{IDENT, "SELECT"},
				{WS, " "},
				{WILDCARD, "*"},
				{WS, " "},
				{IDENT, "FROM"},
				{WS, " "},
				{IDENT, "users"},
				{WS, " "},
				{IDENT, "where"},
				{WS, " "},
				{IDENT, "id"},
				{WS, " "},
				{OPERATOR, "="},
				{WS, " "},
				{STRING, "'12'"},
			},
		},
		{
			name:  "simple select with double quoted identifier",
			input: "SELECT * FROM \"users table\" where id = 1",
			expected: []Token{
				{IDENT, "SELECT"},
				{WS, " "},
				{WILDCARD, "*"},
				{WS, " "},
				{IDENT, "FROM"},
				{WS, " "},
				{QUOTED_IDENT, "\"users table\""},
				{WS, " "},
				{IDENT, "where"},
				{WS, " "},
				{IDENT, "id"},
				{WS, " "},
				{OPERATOR, "="},
				{WS, " "},
				{NUMBER, "1"},
			},
		},
		{
			name:  "simple select with single line comment",
			input: "SELECT * FROM users where id = 1 -- comment here",
			expected: []Token{
				{IDENT, "SELECT"},
				{WS, " "},
				{WILDCARD, "*"},
				{WS, " "},
				{IDENT, "FROM"},
				{WS, " "},
				{IDENT, "users"},
				{WS, " "},
				{IDENT, "where"},
				{WS, " "},
				{IDENT, "id"},
				{WS, " "},
				{OPERATOR, "="},
				{WS, " "},
				{NUMBER, "1"},
				{WS, " "},
				{COMMENT, "-- comment here"},
			},
		},
		{
			name:  "simple select with multi line comment",
			input: "SELECT * /* comment here */ FROM users where id = 1",
			expected: []Token{
				{IDENT, "SELECT"},
				{WS, " "},
				{WILDCARD, "*"},
				{WS, " "},
				{MULTILINE_COMMENT, "/* comment here */"},
				{WS, " "},
				{IDENT, "FROM"},
				{WS, " "},
				{IDENT, "users"},
				{WS, " "},
				{IDENT, "where"},
				{WS, " "},
				{IDENT, "id"},
				{WS, " "},
				{OPERATOR, "="},
				{WS, " "},
				{NUMBER, "1"},
			},
		},
		{
			name:  "simple malformed select",
			input: "SELECT * FROM users where id = 1 and name = 'j",
			expected: []Token{
				{IDENT, "SELECT"},
				{WS, " "},
				{WILDCARD, "*"},
				{WS, " "},
				{IDENT, "FROM"},
				{WS, " "},
				{IDENT, "users"},
				{WS, " "},
				{IDENT, "where"},
				{WS, " "},
				{IDENT, "id"},
				{WS, " "},
				{OPERATOR, "="},
				{WS, " "},
				{NUMBER, "1"},
				{WS, " "},
				{IDENT, "and"},
				{WS, " "},
				{IDENT, "name"},
				{WS, " "},
				{OPERATOR, "="},
				{WS, " "},
				{INCOMPLETE_STRING, "'j"},
			},
		},
		{
			name:  "truncated sql",
			input: "SELECT * FROM users where id = ",
			expected: []Token{
				{IDENT, "SELECT"},
				{WS, " "},
				{WILDCARD, "*"},
				{WS, " "},
				{IDENT, "FROM"},
				{WS, " "},
				{IDENT, "users"},
				{WS, " "},
				{IDENT, "where"},
				{WS, " "},
				{IDENT, "id"},
				{WS, " "},
				{OPERATOR, "="},
				{WS, " "},
			},
		},
		{
			name:  "simple select with array of literals",
			input: "SELECT * FROM users where id in (1, '2')",
			expected: []Token{
				{IDENT, "SELECT"},
				{WS, " "},
				{WILDCARD, "*"},
				{WS, " "},
				{IDENT, "FROM"},
				{WS, " "},
				{IDENT, "users"},
				{WS, " "},
				{IDENT, "where"},
				{WS, " "},
				{IDENT, "id"},
				{WS, " "},
				{IDENT, "in"},
				{WS, " "},
				{PUNCTUATION, "("},
				{NUMBER, "1"},
				{PUNCTUATION, ","},
				{WS, " "},
				{STRING, "'2'"},
				{PUNCTUATION, ")"},
			},
		},
		{
			name:  "dollar quoted function",
			input: "SELECT $func$INSERT INTO table VALUES ('a', 1, 2)$func$ FROM users",
			expected: []Token{
				{IDENT, "SELECT"},
				{WS, " "},
				{DOLLAR_QUOTED_FUNCTION, "$func$INSERT INTO table VALUES ('a', 1, 2)$func$"},
				{WS, " "},
				{IDENT, "FROM"},
				{WS, " "},
				{IDENT, "users"},
			},
		},
		{
			name:  "dollar quoted string",
			input: "SELECT * FROM users where id = $tag$test$tag$",
			expected: []Token{
				{IDENT, "SELECT"},
				{WS, " "},
				{WILDCARD, "*"},
				{WS, " "},
				{IDENT, "FROM"},
				{WS, " "},
				{IDENT, "users"},
				{WS, " "},
				{IDENT, "where"},
				{WS, " "},
				{IDENT, "id"},
				{WS, " "},
				{OPERATOR, "="},
				{WS, " "},
				{DOLLAR_QUOTED_STRING, "$tag$test$tag$"},
			},
		},
		{
			name:  "dollar quoted string",
			input: "SELECT * FROM users where id = $$test$$",
			expected: []Token{
				{IDENT, "SELECT"},
				{WS, " "},
				{WILDCARD, "*"},
				{WS, " "},
				{IDENT, "FROM"},
				{WS, " "},
				{IDENT, "users"},
				{WS, " "},
				{IDENT, "where"},
				{WS, " "},
				{IDENT, "id"},
				{WS, " "},
				{OPERATOR, "="},
				{WS, " "},
				{DOLLAR_QUOTED_STRING, "$$test$$"},
			},
		},
		{
			name:  "numbered parameter",
			input: "SELECT * FROM users where id = $1",
			expected: []Token{
				{IDENT, "SELECT"},
				{WS, " "},
				{WILDCARD, "*"},
				{WS, " "},
				{IDENT, "FROM"},
				{WS, " "},
				{IDENT, "users"},
				{WS, " "},
				{IDENT, "where"},
				{WS, " "},
				{IDENT, "id"},
				{WS, " "},
				{OPERATOR, "="},
				{WS, " "},
				{POSITIONAL_PARAMETER, "$1"},
			},
		},
		{
			name:  "identifier with underscore and period",
			input: "SELECT * FROM users where user_id = 2 and users.name = 'j'",
			expected: []Token{
				{IDENT, "SELECT"},
				{WS, " "},
				{WILDCARD, "*"},
				{WS, " "},
				{IDENT, "FROM"},
				{WS, " "},
				{IDENT, "users"},
				{WS, " "},
				{IDENT, "where"},
				{WS, " "},
				{IDENT, "user_id"},
				{WS, " "},
				{OPERATOR, "="},
				{WS, " "},
				{NUMBER, "2"},
				{WS, " "},
				{IDENT, "and"},
				{WS, " "},
				{IDENT, "users.name"},
				{WS, " "},
				{OPERATOR, "="},
				{WS, " "},
				{STRING, "'j'"},
			},
		},
		{
			name:  "select with hex and octal numbers",
			input: "SELECT * FROM users where id = 0x123 and id = 0X123 and id = 0123",
			expected: []Token{
				{IDENT, "SELECT"},
				{WS, " "},
				{WILDCARD, "*"},
				{WS, " "},
				{IDENT, "FROM"},
				{WS, " "},
				{IDENT, "users"},
				{WS, " "},
				{IDENT, "where"},
				{WS, " "},
				{IDENT, "id"},
				{WS, " "},
				{OPERATOR, "="},
				{WS, " "},
				{NUMBER, "0x123"},
				{WS, " "},
				{IDENT, "and"},
				{WS, " "},
				{IDENT, "id"},
				{WS, " "},
				{OPERATOR, "="},
				{WS, " "},
				{NUMBER, "0X123"},
				{WS, " "},
				{IDENT, "and"},
				{WS, " "},
				{IDENT, "id"},
				{WS, " "},
				{OPERATOR, "="},
				{WS, " "},
				{NUMBER, "0123"},
			},
		},
		{
			name:  "select with float numbers and scientific notation",
			input: "SELECT 1.2,1.2e3,1.2e-3,1.2E3,1.2E-3 FROM users",
			expected: []Token{
				{IDENT, "SELECT"},
				{WS, " "},
				{NUMBER, "1.2"},
				{PUNCTUATION, ","},
				{NUMBER, "1.2e3"},
				{PUNCTUATION, ","},
				{NUMBER, "1.2e-3"},
				{PUNCTUATION, ","},
				{NUMBER, "1.2E3"},
				{PUNCTUATION, ","},
				{NUMBER, "1.2E-3"},
				{WS, " "},
				{IDENT, "FROM"},
				{WS, " "},
				{IDENT, "users"},
			},
		},
		{
			name:  "select with double quoted identifier",
			input: `SELECT * FROM "users table"`,
			expected: []Token{
				{IDENT, "SELECT"},
				{WS, " "},
				{WILDCARD, "*"},
				{WS, " "},
				{IDENT, "FROM"},
				{WS, " "},
				{QUOTED_IDENT, `"users table"`},
			},
		},
		{
			name:  "select with double quoted identifier",
			input: `SELECT * FROM "public"."users table"`,
			expected: []Token{
				{IDENT, "SELECT"},
				{WS, " "},
				{WILDCARD, "*"},
				{WS, " "},
				{IDENT, "FROM"},
				{WS, " "},
				{QUOTED_IDENT, `"public"."users table"`},
			},
		},
		{
			name:  "select with escaped string",
			input: "SELECT * FROM users where id = 'j\\'s'",
			expected: []Token{
				{IDENT, "SELECT"},
				{WS, " "},
				{WILDCARD, "*"},
				{WS, " "},
				{IDENT, "FROM"},
				{WS, " "},
				{IDENT, "users"},
				{WS, " "},
				{IDENT, "where"},
				{WS, " "},
				{IDENT, "id"},
				{WS, " "},
				{OPERATOR, "="},
				{WS, " "},
				{STRING, "'j\\'s'"},
			},
		},
		{
			name:  "select with escaped string",
			input: "SELECT * FROM users where id =?",
			expected: []Token{
				{IDENT, "SELECT"},
				{WS, " "},
				{WILDCARD, "*"},
				{WS, " "},
				{IDENT, "FROM"},
				{WS, " "},
				{IDENT, "users"},
				{WS, " "},
				{IDENT, "where"},
				{WS, " "},
				{IDENT, "id"},
				{WS, " "},
				{OPERATOR, "="},
				{OPERATOR, "?"},
			},
		},
		{
			name:  "select with bind parameter",
			input: "SELECT * FROM users where id = :id and name = :1",
			expected: []Token{
				{IDENT, "SELECT"},
				{WS, " "},
				{WILDCARD, "*"},
				{WS, " "},
				{IDENT, "FROM"},
				{WS, " "},
				{IDENT, "users"},
				{WS, " "},
				{IDENT, "where"},
				{WS, " "},
				{IDENT, "id"},
				{WS, " "},
				{OPERATOR, "="},
				{WS, " "},
				{BIND_PARAMETER, ":id"},
				{WS, " "},
				{IDENT, "and"},
				{WS, " "},
				{IDENT, "name"},
				{WS, " "},
				{OPERATOR, "="},
				{WS, " "},
				{BIND_PARAMETER, ":1"},
			},
			lexerOpts: []lexerOption{WithDBMS(DBMSOracle)},
		},
		{
			name:  "select with bind parameter",
			input: "SELECT * FROM users where id = @id and name = @1",
			expected: []Token{
				{IDENT, "SELECT"},
				{WS, " "},
				{WILDCARD, "*"},
				{WS, " "},
				{IDENT, "FROM"},
				{WS, " "},
				{IDENT, "users"},
				{WS, " "},
				{IDENT, "where"},
				{WS, " "},
				{IDENT, "id"},
				{WS, " "},
				{OPERATOR, "="},
				{WS, " "},
				{BIND_PARAMETER, "@id"},
				{WS, " "},
				{IDENT, "and"},
				{WS, " "},
				{IDENT, "name"},
				{WS, " "},
				{OPERATOR, "="},
				{WS, " "},
				{BIND_PARAMETER, "@1"},
			},
		},
		{
			name:  "select with system variable",
			input: "SELECT @@VERSION AS SqlServerVersion",
			expected: []Token{
				{IDENT, "SELECT"},
				{WS, " "},
				{SYSTEM_VARIABLE, "@@VERSION"},
				{WS, " "},
				{IDENT, "AS"},
				{WS, " "},
				{IDENT, "SqlServerVersion"},
			},
		},
		{
			name:  "SQL Server quoted identifier",
			input: "SELECT [user] FROM [test].[table] WHERE [id] = 1",
			expected: []Token{
				{IDENT, "SELECT"},
				{WS, " "},
				{QUOTED_IDENT, "[user]"},
				{WS, " "},
				{IDENT, "FROM"},
				{WS, " "},
				{QUOTED_IDENT, "[test].[table]"},
				{WS, " "},
				{IDENT, "WHERE"},
				{WS, " "},
				{QUOTED_IDENT, "[id]"},
				{WS, " "},
				{OPERATOR, "="},
				{WS, " "},
				{NUMBER, "1"},
			},
			lexerOpts: []lexerOption{WithDBMS(DBMSSQLServer)},
		},
		{
			name:  "MySQL backtick quoted identifier",
			input: "SELECT `user` FROM `test`.`table` WHERE `id` = 1",
			expected: []Token{
				{IDENT, "SELECT"},
				{WS, " "},
				{QUOTED_IDENT, "`user`"},
				{WS, " "},
				{IDENT, "FROM"},
				{WS, " "},
				{QUOTED_IDENT, "`test`.`table`"},
				{WS, " "},
				{IDENT, "WHERE"},
				{WS, " "},
				{QUOTED_IDENT, "`id`"},
				{WS, " "},
				{OPERATOR, "="},
				{WS, " "},
				{NUMBER, "1"},
			},
			lexerOpts: []lexerOption{WithDBMS(DBMSMySQL)},
		},
		{
			name:  "Tokenize function",
			input: "SELECT count(*) FROM users",
			expected: []Token{
				{IDENT, "SELECT"},
				{WS, " "},
				{FUNCTION, "count"},
				{PUNCTUATION, "("},
				{WILDCARD, "*"},
				{PUNCTUATION, ")"},
				{WS, " "},
				{IDENT, "FROM"},
				{WS, " "},
				{IDENT, "users"},
			},
		},
		{
			name:  "Tokenize temp table",
			input: `SELECT * FROM #temp`,
			expected: []Token{
				{IDENT, "SELECT"},
				{WS, " "},
				{WILDCARD, "*"},
				{WS, " "},
				{IDENT, "FROM"},
				{WS, " "},
				{IDENT, "#temp"},
			},
			lexerOpts: []lexerOption{WithDBMS(DBMSSQLServer)},
		},
		{
			name:  "MySQL comment",
			input: `SELECT * FROM users # comment`,
			expected: []Token{
				{IDENT, "SELECT"},
				{WS, " "},
				{WILDCARD, "*"},
				{WS, " "},
				{IDENT, "FROM"},
				{WS, " "},
				{IDENT, "users"},
				{WS, " "},
				{COMMENT, "# comment"},
			},
			lexerOpts: []lexerOption{WithDBMS(DBMSMySQL)},
		},
		{
			name:  "drop table if exists",
			input: `DROP TABLE IF EXISTS users`,
			expected: []Token{
				{IDENT, "DROP"},
				{WS, " "},
				{IDENT, "TABLE"},
				{WS, " "},
				{IDENT, "IF"},
				{WS, " "},
				{IDENT, "EXISTS"},
				{WS, " "},
				{IDENT, "users"},
			},
		},
		{
			name:  "select only",
			input: "SELECT * FROM ONLY tab1 where id = 1",
			expected: []Token{
				{IDENT, "SELECT"},
				{WS, " "},
				{WILDCARD, "*"},
				{WS, " "},
				{IDENT, "FROM"},
				{WS, " "},
				{IDENT, "ONLY"},
				{WS, " "},
				{IDENT, "tab1"},
				{WS, " "},
				{IDENT, "where"},
				{WS, " "},
				{IDENT, "id"},
				{WS, " "},
				{OPERATOR, "="},
				{WS, " "},
				{NUMBER, "1"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lexer := New(tt.input, tt.lexerOpts...)
			tokens := lexer.ScanAll()
			assert.Equal(t, tt.expected, tokens)
		})
	}
}

func TestLexerUnicode(t *testing.T) {
	tests := []struct {
		input     string
		expected  []Token
		lexerOpts []lexerOption
	}{
		{
			input: `Descripció_CAT`,
			expected: []Token{
				{IDENT, `Descripció_CAT`},
			},
		},
		{
			input: `世界`,
			expected: []Token{
				{IDENT, `世界`},
			},
		},
		{
			input: `こんにちは`,
			expected: []Token{
				{IDENT, `こんにちは`},
			},
		},
		{
			input: `안녕하세요`,
			expected: []Token{
				{IDENT, `안녕하세요`},
			},
		},
		{
			input: `über`,
			expected: []Token{
				{IDENT, `über`},
			},
		},
		{
			input: `résumé`,
			expected: []Token{
				{IDENT, `résumé`},
			},
		},
		{
			input: `"über"`,
			expected: []Token{
				{QUOTED_IDENT, `"über"`},
			},
		},
	}

	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			lexer := New(tt.input, tt.lexerOpts...)
			tokens := lexer.ScanAll()
			assert.Equal(t, tt.expected, tokens)
		})
	}
}

func ExampleLexer() {
	query := "SELECT * FROM users WHERE id = 1"
	lexer := New(query)
	tokens := lexer.ScanAll()
	fmt.Println(tokens)
	// Output: [{6 SELECT} {2  } {9 *} {2  } {6 FROM} {2  } {6 users} {2  } {6 WHERE} {2  } {6 id} {2  } {8 =} {2  } {5 1}]
}
