package plugins

import (
	"context"
	"fmt"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/klog"
	framework "k8s.io/kubernetes/pkg/scheduler/framework"
	metricsv "k8s.io/metrics/pkg/client/clientset/versioned"
)

// 插件名称
const Name = "sample-plugin"

type Sample struct {
	kubeClient    *kubernetes.Clientset
	metricsClient *metricsv.Clientset
	args          *NodeResourcesUsageArgs
	handle        framework.Handle
}

func (s *Sample) Name() string {
	return Name
}

func (s *Sample) PreFilter(ctx context.Context, pod *v1.Pod) *framework.Status {
	klog.V(3).Infof("prefilter pod: %v", pod.Name)
	return framework.NewStatus(framework.Success, "")
}

func (s *Sample) Filter(ctx context.Context, pod *v1.Pod, nodeName string) *framework.Status {
	klog.V(3).Infof("filter pod: %v, node: %v", pod.Name, nodeName)
	// if nodeName == "shtl009063227" {
	// 	return framework.NewStatus(framework.Unschedulable, "")
	// }
	return framework.NewStatus(framework.Success, "")
}

func (s *Sample) Score(ctx context.Context, pod *v1.Pod, nodeName string) (int64, *framework.Status) {
	klog.V(3).Infof("Score pod: %v, node: %v", pod.Name, nodeName)
	nodeInfo, err := s.handle.SnapshotSharedLister().NodeInfos().Get(nodeName)
	if err != nil {
		return framework.MinNodeScore, framework.NewStatus(framework.Error, "")
	}

	nodeMetrics, err := s.metricsClient.MetricsV1beta1().NodeMetricses().Get(ctx, nodeName, metav1.GetOptions{})
	if err != nil {
		klog.Errorf("get node %s metrics err: %+v", nodeName, err)
	}

	// Calculation formula Details:
	// (cpu((100-Usage*100)/10) + memory((100-usage*100)/10))/2
	cpuMetrics := float64(100*nodeMetrics.Usage.Cpu().MilliValue()) / float64(nodeInfo.Allocatable.MilliCPU)
	memMetrics := float64(100*nodeMetrics.Usage.Memory().Value()) / float64(nodeInfo.Allocatable.Memory)
	klog.V(3).Infof("node name: %s, cpu metrics %f, mem metrics %f", nodeInfo.Node().Name, cpuMetrics, memMetrics)

	nodeScore := ((100-cpuMetrics)/10 + (100-memMetrics)/10) / 2
	return int64(nodeScore), framework.NewStatus(framework.Success, "")
}

func (s *Sample) NormalizeScore(ctx context.Context, pod *v1.Pod, nodeScoreList framework.NodeScoreList) *framework.Status {
	klog.V(3).Infof("Normalize Score pod: %v, node Score List: %v", pod.Name, nodeScoreList)
	return nil
}

// func (s *Sample) PreBind(pc *framework.PluginContext, pod *v1.Pod, nodeName string) *framework.Status {
// 	if nodeInfo, ok := s.handle.NodeInfoSnapshot().NodeInfoMap[nodeName]; !ok {
// 		return framework.NewStatus(framework.Error, fmt.Sprintf("prebind get node info error: %+v", nodeName))
// 	} else {
// 		klog.V(3).Infof("prebind node info: %+v", nodeInfo.Node())
// 		return framework.NewStatus(framework.Success, "")
// 	}
// }

func getargs(obj runtime.Object) (*NodeResourcesUsageArgs, error) {
	args, ok := obj.(*NodeResourcesUsageArgs)
	if !ok {
		return &NodeResourcesUsageArgs{}, fmt.Errorf("want args to be of type Args, got %T", obj)
	}
	return args, nil
}

func New(configuration runtime.Object, f framework.Handle) (framework.Plugin, error) {
	args, err := getargs(configuration)
	if err != nil {
		return nil, err
	}
	klog.V(3).Infof("get plugin config args: %+v", args)

	config, err := clientcmd.BuildConfigFromFlags("", args.KubeConfig)
	if err != nil {
		klog.Errorf("get config err: %+v", err)
	}

	kubeClient, err := kubernetes.NewForConfig(config)
	if err != nil {
		klog.Errorf("get kubeClient err: %+v", err)
	}

	metricsClient, err := metricsv.NewForConfig(config)
	if err != nil {
		klog.Errorf("get metricsClient err: %+v", err)
	}

	return &Sample{
		kubeClient:    kubeClient,
		metricsClient: metricsClient,
		args:          args,
		handle:        f,
	}, nil
}
