package main

import (
	"os"
	"unsafe"
)

type Module struct {
	Path string
	Buf  []byte
	C    []*Code
}

func NewModule(path string) *Module {
	buf, err := os.ReadFile(path)
	panicErr(err)
	m := new(Module)
	m.Path = path
	m.Buf = []byte(*(*string)(unsafe.Pointer(&buf)))
	return m
}

func (m *Module) Parse() {
	m.C = make([]*Code, 0)
	this := NewCode()
	token := make([]byte, 0)
	mark := false
	annotation := false
	lineIndex := 0
	for index := 0; index < len(m.Buf); index++ {
		add := -1
		t := m.Buf[index]
		if mark {
			token = append(token, t)
			if t == '"' {
				mark = false
				add = this.Add(token)
				token = make([]byte, 0)

			}

			continue
		}
		if annotation {
			if t == '\n' {
				annotation = false
				if lineIndex == 0{
					token = append([]byte("#0"),token...)
				}else {
					token = append([]byte("#1"),token...)
				}
				this.Add(token)
				token = make([]byte, 0)
			} else {
				token = append(token, t)
			}
			continue
		}
		switch t {

		case '\t', ' ', '\n', '#':
			if len(token) > 0 {
				add = this.Add(token)
				token = make([]byte, 0)
			}
			if t == '#' {
				annotation = true
			}
			if t == '\n' {
				lineIndex = 0
			}
		case '"':
			lineIndex++
			token = append(token, t)
			mark = true
		case '{', '}', '=':
			lineIndex++
			if len(token) > 0 {
				add = this.Add(token)
				token = make([]byte, 0)
			}
			this.Add([]byte{t})
		default:
			token = append(token, t)
			lineIndex++
		}
		if this.Model == -1 {
			m.C = append(m.C, this)
			this = NewCode()
		}
		if add == 0 {
			this.Add(token)
		}

	}
}
func (m *Module) display() []byte {
	buf := make([]byte, 0)
	for _, code := range m.C {
		buf = append(buf, code.Display(0)...)
		buf = append(buf, '\n')
	}

	return buf
}

func (m *Module) Display() string {
	buf := m.display()
	return *(*string)(unsafe.Pointer(&buf))
}
func (m *Module) Format(path string) {
	file, err := os.Create(path)
	panicErr(err)
	defer file.Close()
	_, err = file.Write(m.display())
	panicErr(err)
}
