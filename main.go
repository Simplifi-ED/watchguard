package main

import (
	"context"
	"log"
	"strings"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/dynamic/dynamicinformer"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
)

var preferredVersions = map[string]string{
	"flowcontrol.apiserver.k8s.io/v1beta3": "flowcontrol.apiserver.k8s.io/v1",
}

func main() {
	// Get in-cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		log.Fatalf("Failed to get cluster config: %v", err)
	}

	// Create discovery client
	discoveryClient, err := discovery.NewDiscoveryClientForConfig(config)
	if err != nil {
		log.Fatalf("Failed to create discovery client: %v", err)
	}

	// Create dynamic client
	dynamicClient, err := dynamic.NewForConfig(config)
	if err != nil {
		log.Fatalf("Failed to create dynamic client: %v", err)
	}

	// Get all API resources
	_, resourceList, err := discoveryClient.ServerGroupsAndResources()
	if err != nil {
		log.Fatalf("Failed to get server resources: %v", err)
	}

	// Create dynamic informer factory
	factory := dynamicinformer.NewDynamicSharedInformerFactory(dynamicClient, time.Minute)
	// Track all GVRs we're watching
	watchedResources := make(map[schema.GroupVersionResource]struct{})

	// Set up informers for all resources
	for _, apiResourceList := range resourceList {
		if apiResourceList == nil {
			continue
		}

		// Skip deprecated API versions if a preferred version exists
		if preferredVersion, isDeprecated := preferredVersions[apiResourceList.GroupVersion]; isDeprecated {
			log.Printf("Skipping deprecated API version %s in favor of %s",
				apiResourceList.GroupVersion, preferredVersion)
			continue
		}

		gv, err := schema.ParseGroupVersion(apiResourceList.GroupVersion)
		if err != nil {
			log.Printf("Warning: Failed to parse GroupVersion %s: %v",
				apiResourceList.GroupVersion, err)
			continue
		}

		for _, apiResource := range apiResourceList.APIResources {
			// Skip subresources and non-watchable resources
			if strings.Contains(apiResource.Name, "/") || !containsString(apiResource.Verbs, "watch") {
				continue
			}

			gvr := schema.GroupVersionResource{
				Group:    gv.Group,
				Version:  gv.Version,
				Resource: apiResource.Name,
			}

			// Skip if already watching this resource
			if _, exists := watchedResources[gvr]; exists {
				continue
			}

			log.Printf("Setting up watcher for %s/%s", gvr.GroupVersion(), gvr.Resource)
			watchedResources[gvr] = struct{}{}

			informer := factory.ForResource(gvr).Informer()
			resourceCopy := apiResource

			informer.AddEventHandler(cache.ResourceEventHandlerFuncs{
				AddFunc: func(obj interface{}) {
					metadata, ok := obj.(metav1.Object)
					if !ok {
						return
					}
					log.Printf("[ADDED] %s.%s: %s/%s",
						resourceCopy.Name,
						gv.Group,
						metadata.GetNamespace(),
						metadata.GetName())
				},
				UpdateFunc: func(old, new interface{}) {
					metadata, ok := new.(metav1.Object)
					if !ok {
						return
					}
					log.Printf("[UPDATED] %s.%s: %s/%s",
						resourceCopy.Name,
						gv.Group,
						metadata.GetNamespace(),
						metadata.GetName())
				},
				DeleteFunc: func(obj interface{}) {
					metadata, ok := obj.(metav1.Object)
					if !ok {
						return
					}
					log.Printf("[DELETED] %s.%s: %s/%s",
						resourceCopy.Name,
						gv.Group,
						metadata.GetNamespace(),
						metadata.GetName())
				},
			})
		}
	}

	// Start all informers
	ctx := context.Background()
	factory.Start(ctx.Done())
	factory.WaitForCacheSync(ctx.Done())

	log.Printf("Actively watching %d resource types", len(watchedResources))
	<-ctx.Done()
}

func containsString(slice []string, s string) bool {
	for _, item := range slice {
		if item == s {
			return true
		}
	}
	return false
}
