package buf

type (
	Strings interface {
		Add(v string)
		Get() []string
	}

	stringsImpl struct {
		buf []string
	}
)

func NewStrings() Strings { return &stringsImpl{buf: []string{}} }

func (s *stringsImpl) Add(v string)  { s.buf = append(s.buf, v) }
func (s *stringsImpl) Get() []string { return s.buf }
