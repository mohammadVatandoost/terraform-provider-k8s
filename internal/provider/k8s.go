package provider

import (
	"context"
	"fmt"
	"os"

	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/mohammadVatandoost/terraform-provider-k8s/pkg/utils"

	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

func CreateClusterClient(ctx context.Context) (kubernetes.Interface, error) {
	tflog.Info(ctx, "***** CreateClusterClient ******")
	homeDie, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	kubeConfigPath := homeDie + "/.kube/config"
	var config *rest.Config
	if utils.FileExists(kubeConfigPath) {
		tflog.Info(ctx, fmt.Sprintf("kube config file exist in path: %v", kubeConfigPath))
		config, err = clientcmd.BuildConfigFromFlags("", kubeConfigPath)
		if err != nil {
			return nil, err
		}
	} else {
		// creates the in-cluster config
		config, err = rest.InClusterConfig()
		if err != nil {
			return nil, err
		}
	}

	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	pods, err := clientset.CoreV1().Pods("test").List(ctx, v1.ListOptions{})
	if err != nil {
		tflog.Error(ctx, err.Error())
		return nil, err
	}

	tflog.Info(ctx, fmt.Sprintf("** CreateClusterClient pods: %v", pods.Items))

	return clientset, nil
}
