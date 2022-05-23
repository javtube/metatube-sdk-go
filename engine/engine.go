package engine

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"gorm.io/gorm"

	"github.com/javtube/javtube-sdk-go/common/fetch"
	"github.com/javtube/javtube-sdk-go/model"
	javtube "github.com/javtube/javtube-sdk-go/provider"
)

type Engine struct {
	db             *gorm.DB
	movieProviders map[string]javtube.MovieProvider
	actorProviders map[string]javtube.ActorProvider
}

func New(db *gorm.DB, timeout time.Duration) *Engine {
	return &Engine{
		db:             db,
		actorProviders: initActorProviders(timeout),
		movieProviders: initMovieProviders(timeout),
	}
}

// initActorProviders initializes actor providers.
func initActorProviders(timeout time.Duration) (providers map[string]javtube.ActorProvider) {
	providers = make(map[string]javtube.ActorProvider)
	javtube.RangeActorFactory(func(name string, factory javtube.ActorFactory) {
		provider := factory()
		if s, ok := provider.(javtube.RequestTimeoutSetter); ok {
			s.SetRequestTimeout(timeout)
		}
		providers[strings.ToUpper(name)] = factory()
	})
	return
}

// initMovieProviders initializes movie providers.
func initMovieProviders(timeout time.Duration) (providers map[string]javtube.MovieProvider) {
	providers = make(map[string]javtube.MovieProvider)
	javtube.RangeMovieFactory(func(name string, factory javtube.MovieFactory) {
		provider := factory()
		if s, ok := provider.(javtube.RequestTimeoutSetter); ok {
			s.SetRequestTimeout(timeout)
		}
		providers[strings.ToUpper(name)] = provider
	})
	return
}

func (e *Engine) AutoMigrate(v bool) error {
	if !v {
		return nil
	}
	return e.db.AutoMigrate(
		&model.MovieInfo{},
		&model.ActorInfo{})
}

func (e *Engine) Fetch(url string, provider javtube.Provider) (*http.Response, error) {
	// Provider which implements Fetcher interface should be
	// used to fetch all its corresponding resources.
	if fetcher, ok := provider.(javtube.Fetcher); ok {
		return fetcher.Fetch(url)
	}
	return fetch.Fetch(url)
}

func (e *Engine) IsActorProvider(name string) (ok bool) {
	_, ok = e.actorProviders[strings.ToUpper(name)]
	return
}

func (e *Engine) GetActorProvider(name string) (javtube.ActorProvider, error) {
	provider, ok := e.actorProviders[strings.ToUpper(name)]
	if !ok {
		return nil, fmt.Errorf("actor provider not found: %s", name)
	}
	return provider, nil
}

func (e *Engine) MustGetActorProvider(name string) javtube.ActorProvider {
	provider, err := e.GetActorProvider(name)
	if err != nil {
		panic(err)
	}
	return provider
}

func (e *Engine) IsMovieProvider(name string) (ok bool) {
	_, ok = e.movieProviders[strings.ToUpper(name)]
	return
}

func (e *Engine) GetMovieProvider(name string) (javtube.MovieProvider, error) {
	provider, ok := e.movieProviders[strings.ToUpper(name)]
	if !ok {
		return nil, fmt.Errorf("movie provider not found: %s", name)
	}
	return provider, nil
}

func (e *Engine) MustGetMovieProvider(name string) javtube.MovieProvider {
	provider, err := e.GetMovieProvider(name)
	if err != nil {
		panic(err)
	}
	return provider
}
