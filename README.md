# libdns-plugin

A Go-based plugin system for [libdns](https://github.com/libdns/libdns) DNS providers.

## Overview

libdns-plugin enables DNS providers to be implemented as standalone plugins that can be loaded dynamically by applications using the libdns library. Built on HashiCorp's go-plugin library with GRPC communication, it provides a flexible architecture where DNS provider implementations can be distributed as separate binaries.

The system supports all standard libdns operations (get, set, append, delete records) and uses JSON for provider configuration.

## Usage

### Implementing a Plugin

```go
package main

import (
    "github.com/libdns/libdns"
    "github.com/jsiebens/libdns-plugin"
)

type MyDNSProvider struct {
    // Your provider implementation
}

// Implement the necessary libdns interfaces
func (p *MyDNSProvider) GetRecords(ctx context.Context, zone string) ([]libdns.Record, error) {
    // Implementation
}

// Other required methods...

func main() {
    provider := &MyDNSProvider{}
    plugin.Serve(provider)
}
```

### Using a Plugin

```go
package main

import (
	"context"
	"github.com/jsiebens/libdns-plugin"
)

func main() {
	client, provider, err := plugin.NewClient("path/to/plugin")
	if err != nil {
		// Handle error
	}
	defer client.Close()

	// Configure the provider
	err = provider.Configure(context.Background(), []byte(`{"api_token": "your-token"}`))
	if err != nil {
		// Handle error
	}

	// Use the provider
	records, err := provider.GetRecords(context.Background(), "example.com")
	// ...
}
```