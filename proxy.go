package main

import (
	"log"
	"net/http"
	"net/url"
	"os"
)

// createHTTPClient creates an HTTP client with proxy support if HTTP_PROXY or HTTPS_PROXY is set.
func createHTTPClient() *http.Client {
	client := &http.Client{}

	// Check for HTTP_PROXY or HTTPS_PROXY
	proxyURL := getProxyURL()
	if proxyURL != "" {
		log.Printf("Using proxy: %s\n", proxyURL)
		parsedProxyURL, err := url.Parse(proxyURL)
		if err != nil {
			log.Println("Error parsing proxy URL:", err)
			return client
		}

		// Set the proxy transport for the client
		client.Transport = &http.Transport{
			Proxy: http.ProxyURL(parsedProxyURL),
		}
	}

	return client
}

// getProxyURL returns the value of HTTP_PROXY or HTTPS_PROXY (or lower case) if set.
func getProxyURL() string {
	if proxyURL, exists := os.LookupEnv("HTTP_PROXY"); exists {
		return proxyURL
	}
	if proxyURL, exists := os.LookupEnv("HTTPS_PROXY"); exists {
		return proxyURL
	}
	if proxyURL, exists := os.LookupEnv("http_proxy"); exists {
		return proxyURL
	}
	if proxyURL, exists := os.LookupEnv("https_proxy"); exists {
		return proxyURL
	}
	return ""
}
