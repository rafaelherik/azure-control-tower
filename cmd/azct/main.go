package main

import (
	"context"
	"fmt"
	"os"

	"azure-control-tower/internal/auth"
	"azure-control-tower/internal/azure"
	"azure-control-tower/internal/ui"
	"azure-control-tower/pkg/resource"
)

func main() {
	ctx := context.Background()

	// Authenticate with Azure
	cred, err := auth.NewAzureAuth()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Authentication error: %v\n", err)
		os.Exit(1)
	}

	// Create Azure client
	azureClient, err := azure.NewClient(cred)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create Azure client: %v\n", err)
		os.Exit(1)
	}

	// Initialize resource registry and register handlers
	registry := resource.NewRegistry()
	
	// Register default handler (for generic resources)
	defaultHandler := resource.NewDefaultHandler()
	registry.RegisterHandler(defaultHandler)
	
	// Register storage account handler
	storageHandler := resource.NewStorageHandler()
	registry.RegisterHandler(storageHandler)

	// Register Key Vault handler
	keyVaultHandler := resource.NewKeyVaultHandler()
	registry.RegisterHandler(keyVaultHandler)

	// Create and start UI application
	app := ui.NewApp(azureClient, registry)
	if err := app.Start(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "Application error: %v\n", err)
		os.Exit(1)
	}
}
