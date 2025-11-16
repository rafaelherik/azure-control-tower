package resource

import (
	"fmt"
	"sync"
)

// Registry manages resource type handlers
type Registry struct {
	fetchers map[string]Fetcher
	handlers map[string]ResourceHandler
	mu       sync.RWMutex
}

// NewRegistry creates a new resource registry
func NewRegistry() *Registry {
	return &Registry{
		fetchers: make(map[string]Fetcher),
		handlers: make(map[string]ResourceHandler),
	}
}

// Register registers a resource type fetcher
func (r *Registry) Register(resourceType string, fetcher Fetcher) {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.fetchers[resourceType] = fetcher
}

// GetFetcher returns the fetcher for a resource type
func (r *Registry) GetFetcher(resourceType string) (Fetcher, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	fetcher, ok := r.fetchers[resourceType]
	if !ok {
		return nil, fmt.Errorf("resource type %s not registered", resourceType)
	}
	return fetcher, nil
}

// ListResourceTypes returns all registered resource types
func (r *Registry) ListResourceTypes() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	types := make([]string, 0, len(r.fetchers))
	for t := range r.fetchers {
		types = append(types, t)
	}
	return types
}

// RegisterHandler registers a resource type handler
func (r *Registry) RegisterHandler(handler ResourceHandler) {
	r.mu.Lock()
	defer r.mu.Unlock()
	resourceType := handler.GetResourceType()
	r.handlers[resourceType] = handler
}

// GetHandler returns the handler for a resource type
func (r *Registry) GetHandler(resourceType string) (ResourceHandler, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	handler, ok := r.handlers[resourceType]
	if !ok {
		return nil, fmt.Errorf("resource type %s handler not registered", resourceType)
	}
	return handler, nil
}

// GetHandlerOrDefault returns the handler for a resource type, or a default handler if not found
func (r *Registry) GetHandlerOrDefault(resourceType string) ResourceHandler {
	r.mu.RLock()
	defer r.mu.RUnlock()
	handler, ok := r.handlers[resourceType]
	if !ok {
		// Return default handler if available
		if defaultHandler, hasDefault := r.handlers[""]; hasDefault {
			return defaultHandler
		}
		return nil
	}
	return handler
}

// GetSupportedResourceTypes returns all resource types that have handlers and can navigate to list
func (r *Registry) GetSupportedResourceTypes() []string {
	r.mu.RLock()
	defer r.mu.RUnlock()
	types := make([]string, 0)
	for resourceType, handler := range r.handlers {
		if handler.CanNavigateToList() {
			types = append(types, resourceType)
		}
	}
	return types
}

