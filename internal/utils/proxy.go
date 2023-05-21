package utils

import (
	"bufio"
	"log"
	"net/url"
	"os"
	"time"

	"github.com/jellydator/ttlcache/v3"
)

var i = 0
var proxies = make([]*url.URL, 0)
var assignment = ttlcache.New(
	ttlcache.WithTTL[string, int](30 * time.Minute),
)

func LoadProxy() {
	file, err := os.Open(ProxyFile)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		u, err := url.Parse("http://" + scanner.Text())
		if err != nil {
			log.Fatal(err)
		}

		proxies = append(proxies, u)
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	go assignment.Start()
}

func GetProxy(ip string) (int, *url.URL) {
	val := assignment.Get(ip)
	if val != nil {
		index := val.Value()
		return index, proxies[index]
	}

	index := i
	assignment.Set(ip, i, ttlcache.DefaultTTL)

	i += 1
	if i >= len(proxies) {
		i = 0
	}

	return index, proxies[index]
}
