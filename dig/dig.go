package dig

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/berquerant/dql/errors"
)

type Instr int

//go:generate stringer -type Instr -output dig_stringer_generated.go

const (
	// InstrCancel cancels the invocation of Dig.
	InstrCancel Instr = iota
	// InstrSkipDir skips digging the children of the directory.
	InstrSkipDir
	// InstrContinue digs the children of the directory recursively.
	InstrContinue
)

type (
	FileInfoHandler func(FileInfo) Instr

	FileInfo interface {
		Name() string
		Size() int64
		Mode() fs.FileMode
		ModTime() time.Time
		IsDir() bool
	}

	fileInfo struct {
		name string
		stat fs.FileInfo
	}
)

// Override name.
func (s *fileInfo) Name() string       { return s.name }
func (s *fileInfo) Size() int64        { return s.stat.Size() }
func (s *fileInfo) Mode() fs.FileMode  { return s.stat.Mode() }
func (s *fileInfo) ModTime() time.Time { return s.stat.ModTime() }
func (s *fileInfo) IsDir() bool        { return s.stat.IsDir() }

// Digger provides recursive file search operations.
type Digger interface {
	// Dig searches the information of a file or a directory recursively.
	Dig(name string, handler FileInfoHandler) error
}

// New returns a new Digger.
func New() Digger { return &digger{} }

type digger struct{}

func (s *digger) Dig(name string, handler FileInfoHandler) error {
	p, err := filepath.Abs(name)
	if err != nil {
		return errors.Wrap(err, "digger dig %s", name)
	}
	if err := s.dig(p, handler); err != nil && !errors.Is(err, errDone) {
		return err
	}
	return nil
}

var (
	errDone = errors.New("dig done")
)

func (s *digger) dig(name string, handler FileInfoHandler) error {
	stat, err := os.Stat(name)
	if err != nil {
		return errors.Wrap(err, "digger cannot get stat of %s", name)
	}
	info := &fileInfo{
		name: name,
		stat: stat,
	}
	instr := handler(info)
	switch instr {
	case InstrCancel:
		// cancel dig invocations.
		return errDone
	case InstrSkipDir:
		// skip digging the children of the directory.
		return nil
	case InstrContinue:
		// dig the children of the directory.
		if !info.IsDir() {
			return nil
		}
		dir, err := os.Open(name)
		if err != nil {
			return errors.Wrap(err, "digger cannot open directory %s", name)
		}
		defer dir.Close()
		children, err := dir.Readdirnames(0)
		if err != nil {
			return errors.Wrap(err, "digger cannot read children of %s", name)
		}
		sort.Strings(children)
		for _, c := range children {
			p := filepath.Join(name, c)
			if err := s.dig(p, handler); err != nil {
				if errors.Is(err, errDone) {
					return nil
				}
				return err
			}
		}
		return nil
	default:
		panic(fmt.Sprintf("dig encountered unknown Instr %s", instr))
	}
}
