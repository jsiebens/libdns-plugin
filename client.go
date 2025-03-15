package plugin

import (
	"github.com/hashicorp/go-plugin"
	"os/exec"
)

func NewClient(cmd *exec.Cmd) (*plugin.Client, Provider, error) {
	client := plugin.NewClient(&plugin.ClientConfig{
		HandshakeConfig: Handshake,
		Plugins:         PluginMap,
		Cmd:             cmd,
		AllowedProtocols: []plugin.Protocol{
			plugin.ProtocolNetRPC,
			plugin.ProtocolGRPC,
		},
	})

	// Connect via RPC
	rpcClient, err := client.Client()
	if err != nil {
		client.Kill()
		return nil, nil, err
	}

	// Request the plugin
	raw, err := rpcClient.Dispense(ProviderPluginName)
	if err != nil {
		client.Kill()
		return nil, nil, err
	}

	provider := raw.(Provider)

	return client, provider, nil
}
