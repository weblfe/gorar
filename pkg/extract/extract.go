package extract

import (
	"fmt"
	"github.com/weblfe/gorar/pkg/detect"
	"io/fs"
	"sync"
)

func Extract(input, outputDir string, fss ...fs.FS) error {
	var opts = []Option{
		WithHandler(GetDefaultHandlers()...),
	}
	if len(fss) > 0 && fss[0] != nil {
		opts = append(opts, WithRoot(fss[0]))
	}
	archives := New(opts...)
	return archives.Extract(input, outputDir)
}

type (
	ArchiveCmd struct {
		root     fs.FS
		handlers map[string]Handler
	}

	Handler interface {
		Extension() string
		List(source string, fss ...fs.FS) []fs.FileInfo
		Extract(source, output string, fss ...fs.FS) error
	}
	Option func(cmd *ArchiveCmd)

	containers struct {
		s     *sync.Mutex
		lists map[string]Handler
	}
)

func New(opts ...Option) *ArchiveCmd {
	cmd := &ArchiveCmd{
		root:     nil,
		handlers: make(map[string]Handler),
	}
	for _, opt := range opts {
		opt(cmd)
	}
	return cmd
}

func WithRoot(root fs.FS) Option {
	return func(cmd *ArchiveCmd) {
		cmd.root = root
	}
}

func WithHandler(handlers ...Handler) Option {
	return func(cmd *ArchiveCmd) {
		for _, h := range handlers {
			if h == nil || h.Extension() == "" {
				continue
			}
			cmd.handlers[h.Extension()] = h
		}
	}
}

func (c *ArchiveCmd) Register(handlers ...Handler) *ArchiveCmd {
	for _, h := range handlers {
		if h != nil && h.Extension() != "" {
			c.handlers[h.Extension()] = h
		}
	}
	return c
}

func (c *ArchiveCmd) Extract(source, output string) error {
	var (
		err       error
		extension string
	)
	if extension, err = detect.Detect(source, c.root); err != nil {
		return err
	}
	h, ok := c.handlers[extension]
	if !ok {
		return fmt.Errorf("unsupported archive type: %s", extension)
	}
	return h.Extract(source, output, c.root)
}

var (
	register = &containers{
		s:     &sync.Mutex{},
		lists: make(map[string]Handler),
	}
)

func GetDefaultHandlers() []Handler {
	return register.ListHandler()
}

func Register(h Handler) {
	register.Register(h)
}

func (c *containers) Register(h Handler) {
	if h == nil || h.Extension() == "" {
		return
	}
	c.s.Lock()
	defer c.s.Unlock()
	c.lists[h.Extension()] = h
}

func (c *containers) GetHandler(name string) Handler {
	c.s.Lock()
	defer c.s.Unlock()
	return c.lists[name]
}

func (c *containers) ListHandler() []Handler {
	c.s.Lock()
	defer c.s.Unlock()
	var handlers []Handler
	for _, h := range c.lists {
		handlers = append(handlers, h)
	}
	return handlers
}
