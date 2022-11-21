package main

import (
	"context"
	"flag"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"log"
	"path/filepath"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

func main() {
	log.Println("Starting E2E test...")

	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		log.Fatal(err)
	}

	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	ds, err := ListCRDs(clientSet)
	if err != nil {
		log.Println(err)
	}

	log.Println(ds)

	pods, err := ListPodsInNamespace(clientSet, "default")
	if err != nil {
		log.Println(err)
		return
	}

	log.Println(pods)

}

func ListCRDs(clientSet kubernetes.Interface) ([]v1.Node, error) {
	nodes, err := clientSet.CoreV1().Nodes().List(context.Background(), metav1.ListOptions{})

	if err != nil {
		return nil, err
	}

	return nodes.Items, nil
}

func ListPodsInNamespace(clientSet kubernetes.Interface, namespace string) ([]v1.Pod, error) {
	pods, err := clientSet.CoreV1().Pods(namespace).List(context.Background(), metav1.ListOptions{})

	if err != nil {
		return nil, err
	}

	return pods.Items, nil
}
