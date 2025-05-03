package forwarder

import (
	"bytes"
	"log"
	"net/http"
	"sync"
)

type Forwarder struct {
	URLs []string
}

func NewForwarder(urls []string) *Forwarder {
	return &Forwarder{URLs: urls}
}

func (f *Forwarder) Forward(data []byte) {
	var wg sync.WaitGroup

	for _, url := range f.URLs {
		wg.Add(1)
		go func(url string) {
			defer wg.Done()
			resp, err := http.Post(url, "application/json", bytes.NewReader(data))
			if err != nil {
				log.Printf("Failed to forward to %s: %v", url, err)
				return
			}
			defer resp.Body.Close()
			log.Printf("Forwarded to %s, status: %s", url, resp.Status)
		}(url)
	}
	wg.Wait()
}
