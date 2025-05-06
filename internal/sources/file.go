package sources

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/fsnotify/fsnotify"
	"github.com/go-playground/validator/v10"
	"github.com/goccy/go-yaml"
)

type FileSource struct {
	filePath    string
	cfg         fileSourceConfig
	fileWatcher *fsnotify.Watcher
}

var _ Source = (*FileSource)(nil)

func NewFileSource(filePath string) (*FileSource, error) {
	src := FileSource{
		filePath: filePath,
	}

	if err := src.readFile(); err != nil {
		return nil, fmt.Errorf("unable to read file source: %w", err)
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, fmt.Errorf("unable to start file watcher for source: %w", err)
	}

	if err := watcher.Add(src.filePath); err != nil {
		return nil, fmt.Errorf("unable to watch source file: %w", err)
	}

	src.fileWatcher = watcher

	go func() {
		//nolint:staticcheck // This is the recommended style
		for {
			// TODO: can we report errors in some way?
			select {
			case event, ok := <-src.fileWatcher.Events:
				if !ok {
					return
				}

				if !event.Has(fsnotify.Write) {
					continue
				}

				if event.Name != src.filePath {
					continue
				}

				if err := src.readFile(); err != nil {
					continue
				}
			}
		}
	}()

	return &src, nil
}

func (src *FileSource) GetRedirectForKey(
	ctx context.Context,
	key string,
) (Redirect, error) {
	redirect, ok := src.cfg.Redirects[key]
	if !ok {
		return Redirect{}, errNoSuchRedirectKey
	}

	return Redirect(redirect), nil
}

func (src *FileSource) GetAllRedirects(ctx context.Context) (map[string]Redirect, error) {
	redirects := make(map[string]Redirect)
	for key, redirect := range src.cfg.Redirects {
		redirects[key] = Redirect(redirect)
	}

	return redirects, nil
}

func (src *FileSource) Close() {
	_ = src.fileWatcher.Close()
}

func (src *FileSource) readFile() error {
	file, err := os.Open(src.filePath)
	if err != nil {
		return fmt.Errorf("unable to open source file: %w", err)
	}

	var cfg fileSourceConfig
	if err := yaml.NewDecoder(file).Decode(&cfg); err != nil {
		return fmt.Errorf("unable to parse source file: %w", err)
	}

	if err := structValidator.Struct(cfg); err != nil {
		return fmt.Errorf("invalid source configuration: %w", err)
	}

	src.cfg = cfg

	return nil
}

var (
	errNoSuchRedirectKey = errors.New("no such redirect matching key")
)

type fileSourceConfig struct {
	Redirects map[string]fileSourceConfigRedirect `yaml:"redirects" validate:"required"`
}

type fileSourceConfigRedirect struct {
	URLPattern string `yaml:"urlPattern" validate:"required,url"`
}

var structValidator = validator.New(
	validator.WithRequiredStructEnabled(),
)
