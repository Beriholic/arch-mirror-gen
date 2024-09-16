package mirror

import (
	"context"
	"fmt"
	"net/http"
	"regexp"
	"sort"
	"sync"
	"time"
)

type Mirror struct {
	URL     string
	Latency time.Duration
}

func SortMirrorListByLatency(mirrorList []Mirror) []string {
	sort.Slice(mirrorList, func(i, j int) bool {
		return mirrorList[i].Latency < mirrorList[j].Latency
	})

	sortedURLs := make([]string, len(mirrorList))
	for i, mirror := range mirrorList {
		sortedURLs[i] = mirror.URL
	}
	return sortedURLs
}

func ParseMirrorList(body string) ([]string, error) {
	re := regexp.MustCompile(`http[s]?://[^/]+/archlinux/[^/]+/os/[^/]+arch`)

	mirrorList := re.FindAllString(body, -1)

	return mirrorListSpeedSort(mirrorList), nil
}

func mirrorListSpeedSort(mirrorList []string) []string {

	var mirrorSpeedList []Mirror

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	var wg sync.WaitGroup
	mirrorListChan := make(chan Mirror, len(mirrorList))

	for _, mirror := range mirrorList {
		wg.Add(1)
		go pingUrl(ctx, mirror, &wg, mirrorListChan)
	}

	go func() {
		wg.Wait()
		close(mirrorListChan)
	}()

	for mirror := range mirrorListChan {
		mirrorSpeedList = append(mirrorSpeedList, mirror)
	}

	return SortMirrorListByLatency(mirrorSpeedList)
}

func pingUrl(
	ctx context.Context,
	url string,
	wg *sync.WaitGroup,
	results chan<- Mirror,
) {
	defer wg.Done()

	start := time.Now()

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		mirror := Mirror{URL: url, Latency: time.Duration(1<<63 - 1)}
		results <- mirror
		return
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		mirror := Mirror{URL: url, Latency: time.Duration(1<<63 - 1)}
		results <- mirror
		return
	}
	defer resp.Body.Close()

	latency := time.Since(start)

	fmt.Printf("URL: %s, Latency: %s\n", url, latency)
	mirror := Mirror{URL: url, Latency: latency}
	results <- mirror
}
