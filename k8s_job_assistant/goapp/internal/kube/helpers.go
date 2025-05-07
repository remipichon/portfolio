package kube

import (
	"flag"
	"fmt"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"path/filepath"
)

func newTrue() *bool {
	b := true
	return &b
}

func newFalse() *bool {
	b := false
	return &b
}

// InitClient instantiate a Kubernetes client based on local kubeconfig if exists
// or default to in-cluster config
func InitKubeClient() *kubernetes.Clientset {
	var kubeconfigPath string

	// Optional flag for kubeconfig
	if home := homedir.HomeDir(); home != "" {
		flag.StringVar(&kubeconfigPath, "kubeconfig", filepath.Join(home, ".kube", "config"),
			"(optional) absolute path to the kubeconfig file")
	} else {
		flag.StringVar(&kubeconfigPath, "kubeconfig", "",
			"(optional) absolute path to the kubeconfig file")
	}
	flag.Parse()

	var config *rest.Config
	var err error

	// Use kubeconfig if explicitly provided or exists
	if kubeconfigPath != "" {
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfigPath)
		if err != nil {
			fmt.Printf("Warning: Failed to load kubeconfig: %v\nFalling back to in-cluster config.\n", err)
		}
	}

	// If config is still nil, fallback to in-cluster config
	if config == nil {
		config, err = rest.InClusterConfig()
		if err != nil {
			panic(fmt.Errorf("failed to load both kubeconfig and in-cluster config: %w", err))
		}
	}

	// Create clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(fmt.Errorf("failed to create Kubernetes client: %w", err))
	}

	return clientset
}
