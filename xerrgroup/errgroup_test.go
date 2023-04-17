package xerrgroup

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	"golang.org/x/sync/errgroup"
	"gotest.tools/assert"
)

func GetGoid() int64 {
	var (
		buf [64]byte
		n   = runtime.Stack(buf[:], false)
		stk = strings.TrimPrefix(string(buf[:n]), "goroutine")
	)

	idField := strings.Fields(stk)[0]
	id, err := strconv.Atoi(idField)
	if err != nil {
		panic(fmt.Errorf("can not get goroutine id: %v", err))
	}

	return int64(id)
}

var (
	datarange = []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j"}
	randIndex = rand.Int31n(10)
)

func calc(index int, val string) (string, error) {
	if randIndex == int32(index) {
		return "", errors.New("invalid index")
	}

	return val, nil
}

func TestErrGroupWithCtx(t *testing.T) {
	wg, ctx := errgroup.WithContext(context.Background())
	result := make(map[string]bool)
	var mu sync.Mutex
	for i, v := range datarange {
		index, val := i, v
		wg.Go(func() error {
			gid := GetGoid()
			select {
			case <-ctx.Done():
				fmt.Printf("goroutine: %d 未执行，msg: %s\n", gid, ctx.Err())
				return nil
			default:

			}
			data, err := calc(index, val)
			if err != nil {
				return fmt.Errorf("在g: %d报错, %s\n", gid, err)
			}

			fmt.Printf("[%s] 执行: %d\n", data, gid)
			mu.Lock()
			result[data] = true
			mu.Unlock()

			fmt.Printf("正常退出: %d\n", gid)

			return nil
		})
	}

	if err := wg.Wait(); err != nil {
		fmt.Println(err)
	}

	fmt.Println("运行结束", result)

	// first nil err
	_, ok := result[datarange[randIndex]]
	assert.Equal(t, ok, false)
}

func TestErrGroupNoCtx(t *testing.T) {
	var wg errgroup.Group

	result := make(map[string]bool)
	var mu sync.Mutex

	for i, v := range datarange {
		index, val := i, v
		wg.Go(func() error {
			gid := GetGoid()

			data, err := calc(index, val)
			if err != nil {
				return fmt.Errorf("在g: %d报错, %s\n", gid, err)
			}

			fmt.Printf("[%s] 执行: %d\n", data, gid)
			mu.Lock()
			result[data] = true
			mu.Unlock()

			fmt.Printf("正常退出: %d\n", GetGoid())

			return nil
		})
	}

	if err := wg.Wait(); err != nil {
		fmt.Println(err)
	}

	fmt.Println("运行结束", result)

	// first nil err
	_, ok := result[datarange[randIndex]]
	assert.Equal(t, ok, false)
}

// https://go.dev/play/p/hXWjtN4uj06
func TestCtxGroup(t *testing.T) {
	errs, ctx := errgroup.WithContext(context.Background())

	errs.Go(func() error {
		// This routine shouldn't print Long Wait, because by the time its execution completes the context was canceled
		select {
		case <-ctx.Done():
			// Handle Cancelation
			return ctx.Err()
		case <-time.After(5000 * time.Millisecond):
			fmt.Printf("Long Wait successful")
			return nil
		}
	})

	errs.Go(func() error {
		// This routine print Short Wait, because this routine completes before context was canceled
		select {
		case <-ctx.Done():
			// Handle Cancelation
			return ctx.Err()
		case <-time.After(10 * time.Millisecond):
			fmt.Printf("Short Wait successful\n")
			return nil
		}
	})
	errs.Go(func() error {
		select {
		case <-ctx.Done():
			// Handle Cancelation
			return ctx.Err()
		case <-time.After(100 * time.Millisecond):
			return fmt.Errorf("Stop routine called\n")
		}
	})

	if err := errs.Wait(); err != nil {
		fmt.Printf("Error received is %+v", err)
	}
}
