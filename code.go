package main

func NewCode() *Code {
	return &Code{Token: make([][]byte, 0), C: make([]*Code, 0), AnnotationStart: make([][]byte, 0)}
}

type Code struct {
	K                []byte
	V                []byte
	AnnotationStart  [][]byte
	AnnotationLeft   []byte
	AnnotationRight  []byte
	LeftParenthesis  int
	RightParenthesis int
	Model            int
	Token            [][]byte
	C                []*Code
}

func printIndent(level int) []byte {
	buf := make([]byte, level*4)
	for i := 0; i < len(buf); i++ {
		buf[i] = ' '
	}
	return buf
}

func (c *Code) Display(level int) []byte {
	s := printIndent(level)
	for _, s2 := range c.AnnotationStart {
		s = append(s, '#')
		s = append(s, s2...)
		s = append(s, '\n')
		s = append(s, printIndent(level)...)
	}

	s = append(s, c.K...)
	s = append(s, ' ')
	switch {
	case c.V != nil:
		s = append(s, '=')
		s = append(s, ' ')
		s = append(s, c.V...)
		if c.AnnotationRight != nil {
			s = append(s, ' ')
			s = append(s, '#')
			s = append(s, []byte(c.AnnotationRight)...)
		}
	case len(c.C) > 0:
		s = append(s, '=')
		s = append(s, ' ')
		s = append(s, '{')
		if c.AnnotationLeft != nil {
			s = append(s, ' ')
			s = append(s, '#')
			s = append(s, []byte(c.AnnotationLeft)...)
		}
		s = append(s, '\n')
		for _, code := range c.C {
			s = append(s, code.Display(level+1)...)
		}
		s = append(s, printIndent(level)...)
		s = append(s, '}')
		if c.AnnotationRight != nil {
			s = append(s, ' ')
			s = append(s, '#')
			s = append(s, c.AnnotationRight...)
		}
	}
	s = append(s, '\n')

	return s
}

func (c *Code) Add(token []byte) int {
	switch  {
	//case "=":
	//	c.Model = 1
	case byteCom(token,"{"):
		c.LeftParenthesis++
		if c.Model == 2 {
			c.Token = append(c.Token, token)
		}
		if c.Model == 1 {
			c.Model = 2
		}
	case byteCom(token,"}"):
		c.RightParenthesis++
		if c.RightParenthesis == c.LeftParenthesis {
			c.parse()
		} else {
			c.Token = append(c.Token, token)
		}
	default:
		// 注释
		if token[0] == '#' {
			if token[1] == '0' && c.Model == 0 {
				c.AnnotationStart = append(c.AnnotationStart, token[2:])
				return 1
			}
			if token[1] == '1' {
				if c.Model == 0 {
					return 0
				}
				if c.Model == -1 {
					c.AnnotationRight = token[2:]
					return 1
				}
				if c.Model == 2 && len(c.Token) == 0 {
					c.AnnotationLeft = token[2:]
					return 1
				}
			}
		}
		switch c.Model {
		case 0:
			if byteCom(token , "=") {
				c.Model = 1
				return 1
			}
			if c.K != nil {
				c.Model = -1
				return 0
			}
			c.K = token

		case 1:
			c.V = token
			c.Model = -1
		case 2:
			c.Token = append(c.Token, token)
		case -1:
			PanicString("code has complete!")
		}

	}
	return 1
}
func byteCom(b []byte, s string) bool {
	if len(b) == len(s) {
		for i, b2 := range b {
			if s[i] != b2 {
				return false
			}
		}
		return true
	}
	return false
}

func (c *Code) parse() {
	if c.Model != 2 {
		PanicString("code grammatical error!")
	}
	this := NewCode()
	for _, s := range c.Token {
		m := this.Add(s)
		if this.Model == -1 {
			c.C = append(c.C, this)
			this = NewCode()
			if m == 0 {
				this.Add(s)
			}
		}
		if m == 0 && s[0] == '#' {
			c.C[len(c.C)-1].Add(s)
		}
	}
	c.Model = -1

}
