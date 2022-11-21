package main

import (
	"context"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func ListCRDs(clientSet kubernetes.Interface) ([]v1.Node, error) {
	nodes, err := clientSet.CoreV1().Nodes().List(context.Background(), metav1.ListOptions{})

	if err != nil {
		return nil, err
	}

	return nodes.Items, nil
}
