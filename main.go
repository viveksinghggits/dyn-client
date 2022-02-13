package main

import (
	"context"
	"flag"
	"fmt"
	"time"

	"github.com/viveksinghggits/kluster/pkg/apis/viveksingh.dev/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/dynamic/dynamicinformer"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func main() {
	kubeconfig := flag.String("kubeconfig", "/home/vivek/.kube/config", "location to your kubeconfig file")
	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		// handle error
		fmt.Printf("erorr %s building config from flags\n", err.Error())
		config, err = rest.InClusterConfig()
		if err != nil {
			fmt.Printf("error %s, getting inclusterconfig", err.Error())
		}
	}

	dynClient, err := dynamic.NewForConfig(config)
	if err != nil {
		fmt.Printf("error %s, getting dyn client\n", err.Error())
	}

	unsObject, err := dynClient.Resource(schema.GroupVersionResource{
		Group:    "viveksingh.dev",
		Version:  "v1alpha1",
		Resource: "klusters",
	}).Namespace("default").Get(context.Background(), "kluster-0", metav1.GetOptions{})
	if err != nil {
		fmt.Printf("error %s gettign resource from dyn client\n", err.Error())
	}

	k := v1alpha1.Kluster{}
	// getting and setting fields on unsObject
	fmt.Printf("Go teh object %s\n", unsObject.GetName())

	// how do I convert unsObject into a typed object
	err = runtime.DefaultUnstructuredConverter.FromUnstructured(unsObject.UnstructuredContent(), &k)
	if err != nil {
		fmt.Printf("error %s, converting unstructured to kluster type", err.Error())
	}

	infFactory := dynamicinformer.NewDynamicSharedInformerFactory(dynClient, 10*time.Minute)

	c := newController(dynClient, infFactory)
	infFactory.Start(make(<-chan struct{}))
	c.run(make(<-chan struct{}))
	fmt.Printf("the concrete type that we got is %+v\n", k)
}
