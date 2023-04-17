package xerrgroup

import (
	"context"
	"fmt"
	"os"
	"testing"

	"golang.org/x/sync/errgroup"
)

var (
	Web    = fakeSearch("web")
	Image  = fakeSearch("image")
	Video1 = fakeSearch("video1")
	Video2 = fakeSearch("video2")
	Video3 = fakeSearch("video3")
	Video4 = fakeSearch("video4")
	Video5 = fakeSearch("video5")
	Video6 = fakeSearch("video6")
	Video7 = fakeSearch("video7")
	Video8 = fakeSearch("video8")
)

type (
	Result string
	Search func(ctx context.Context, query string) (Result, error)
)

func fakeSearch(kind string) Search {
	return func(_ context.Context, query string) (Result, error) {
		return Result(fmt.Sprintf("%s result for %q", kind, query)), nil
	}
}

func TestParallel(t *testing.T) {
	Google := func(ctx context.Context, query string) ([]Result, error) {
		g, ctx := errgroup.WithContext(ctx)

		searches := []Search{Web, Image, Video1, Video2, Video3, Video4, Video5, Video6, Video7, Video8}
		results := make([]Result, len(searches))
		for i, search := range searches {
			i, search := i, search // https://golang.org/doc/faq#closures_and_goroutines
			g.Go(func() error {
				result, err := search(ctx, query)
				if err == nil {
					results[i] = result
				}
				return err
			})
		}
		if err := g.Wait(); err != nil {
			return nil, err
		}
		return results, nil
	}

	results, err := Google(context.Background(), "golang")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return
	}
	for _, result := range results {
		fmt.Println(result)
	}
}
