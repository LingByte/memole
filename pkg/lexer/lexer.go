package lexer

// TokenType 用int类型表示不同的词法单元类型
type TokenType int

// 定义常量枚举（Go没有真正的枚举，用常量组实现）
const (
	// iota 从0开始自动递增
	Illegal TokenType = iota // 非法字符
	EOF                      // 文件结束

	// 标识符和字面量
	Identifier // 变量名
	Int        // 整型
	String     // 字符串字面量

	// 运算符
	Assign // =
	Plus   // +
	Minus  // -
	Multiply // *
	Slash    // /

	// 分隔符
	Comma     // ,
	Semicolon // ;
	Dot       // .

	// 括号
	LParen // (
	RParen // )

	// 布尔类型和比较运算符
	True  // true
	False // false
	GT    // >
	LT    // <
	EQ    // ==
	NotEQ // !=
	Bang  // !

	// 关键字与其他符号
	If       // if
	Else     // else
	Ty       // ty
	Stru     // stru
	Lbrace   // {
	Rbrace   // }
	Function // fn
	Return   // return
	Package  // package
	Import   // import
	While    // while
)

// Token 结构体表示词法单元（学习Go结构体定义）
type Token struct {
	Type    TokenType // 类型
	Literal string    // 字面值
	Pos     int       // 在输入中的位置（用于错误报告）
}

// Lexer 结构体保存词法分析状态（学习结构体方法）
type Lexer struct {
	input   string // 输入字符串
	pos     int    // 当前读取位置
	readPos int    // 下一个读取位置（用于预读）
	ch      byte   // 当前字符
}

// New 创建新Lexer实例（学习构造函数写法）
func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar() // 初始化读取第一个字符
	return l
}

// readChar 读取下一个字符（学习方法定义）
func (l *Lexer) readChar() {
	if l.readPos >= len(l.input) {
		l.ch = 0 // ASCII的NUL字符表示EOF
	} else {
		l.ch = l.input[l.readPos]
	}
	l.pos = l.readPos
	l.readPos++
}

// NextToken 获取下一个词法单元（核心方法）
func (l *Lexer) NextToken() Token {
	var tok Token

	// 跳过空白字符
	l.skipWhitespace()

	// 根据当前字符生成不同Token
	switch l.ch {
	case '=':
		if l.peekChar() == '=' {
			l.readChar()
			tok = Token{Type: EQ, Literal: "=="}
		} else {
			tok = newToken(Assign, l.ch)
		}
	case '+':
		tok = newToken(Plus, l.ch)
	case '-':
		tok = newToken(Minus, l.ch)
	case '!':
		if l.peekChar() == '=' {
			l.readChar()
			tok = Token{Type: NotEQ, Literal: "!="}
		} else {
			tok = newToken(Bang, l.ch)
		}
	case '*':
		tok = newToken(Multiply, l.ch)
	case '/':
		tok = newToken(Slash, l.ch)
	case '<':
		tok = newToken(LT, l.ch)
	case '>':
		tok = newToken(GT, l.ch)
	case '(':
		tok = newToken(LParen, l.ch)
	case ')':
		tok = newToken(RParen, l.ch)
	case ',':
		tok = newToken(Comma, l.ch)
	case ';':
		tok = newToken(Semicolon, l.ch)
	case '.':
		tok = newToken(Dot, l.ch)
	case '{':
		tok = newToken(Lbrace, l.ch)
	case '}':
		tok = newToken(Rbrace, l.ch)
	case '"':
		tok.Literal = l.readString()
		tok.Type = String
		l.readChar() // 跳过结束的引号
		return tok
	case 0:
		tok.Literal = ""
		tok.Type = EOF
	default:
		if isLetter(l.ch) {
			tok.Literal = l.readIdentifier()
			switch tok.Literal {
			case "true":
				tok.Type = True
			case "false":
				tok.Type = False
			case "if":
				tok.Type = If
			case "else":
				tok.Type = Else
			case "ty":
				tok.Type = Ty
			case "stru":
				tok.Type = Stru
			case "fn":
				tok.Type = Function
			case "return":
				tok.Type = Return
			case "package":
				tok.Type = Package
			case "import":
				tok.Type = Import
			case "while":
				tok.Type = While
			default:
				tok.Type = Identifier
			}
			return tok // 提前返回，已更新位置
		} else if isDigit(l.ch) {
			tok.Type = Int
			tok.Literal = l.readNumber()
			return tok
		} else {
			tok = newToken(Illegal, l.ch)
		}
	}

	l.readChar()
	return tok
}

// 辅助函数们（学习函数定义）
func newToken(tokenType TokenType, ch byte) Token {
	return Token{
		Type:    tokenType,
		Literal: string(ch),
		Pos:     -1, // 实际使用需要记录位置
	}
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

func (l *Lexer) readIdentifier() string {
	pos := l.pos
	for isLetter(l.ch) {
		l.readChar()
	}
	return l.input[pos:l.pos]
}

func isLetter(ch byte) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_'
}

func (l *Lexer) readNumber() string {
	pos := l.pos
	for isDigit(l.ch) {
		l.readChar()
	}
	return l.input[pos:l.pos]
}

func isDigit(ch byte) bool {
	return '0' <= ch && ch <= '9'
}

// readString 读取字符串字面量
func (l *Lexer) readString() string {
	var result []byte
	for {
		l.readChar()
		if l.ch == '"' || l.ch == 0 {
			break
		}
		// 处理转义字符
		if l.ch == '\\' {
			l.readChar() // 跳过转义字符
			// 将转义字符转换为实际字符
			switch l.ch {
			case '"':
				result = append(result, '"')
			case '\\':
				result = append(result, '\\')
			case 'n':
				result = append(result, '\n')
			case 't':
				result = append(result, '\t')
			case 'r':
				result = append(result, '\r')
			default:
				result = append(result, l.ch)
			}
		} else {
			result = append(result, l.ch)
		}
	}
	return string(result)
}

// 添加Token类型到字符串的映射
var tokenStrings = map[TokenType]string{
	Illegal:    "Illegal",
	EOF:        "EOF",
	Identifier: "Identifier",
	Int:        "Int",
	String:     "String",
	Assign:     "Assign",
	Plus:       "Plus",
	Minus:      "Minus",
	Multiply:   "Multiply",
	Slash:      "Slash",
	Comma:      "Comma",
	Semicolon:  "Semicolon",
	Dot:        "Dot",
	LParen:     "LParen",
	RParen:     "RParen",
	True:       "True",
	False:      "False",
	GT:         "GT",
	LT:         "LT",
	EQ:         "EQ",
	NotEQ:      "NotEQ",
	Bang:       "Bang",
	If:         "If",
	Else:       "Else",
	Ty:         "Ty",
	Stru:       "Stru",
	Lbrace:     "Lbrace",
	Rbrace:     "Rbrace",
	Function:   "Function",
	Return:     "Return",
	Package:    "Package",
	Import:     "Import",
	While:      "While",
}

// 为TokenType实现String()方法
func (tt TokenType) String() string {
	if s, ok := tokenStrings[tt]; ok {
		return s
	}
	return "Unknown"
}

// 新增peekChar方法
func (l *Lexer) peekChar() byte {
	if l.readPos >= len(l.input) {
		return 0
	}
	return l.input[l.readPos]
}
