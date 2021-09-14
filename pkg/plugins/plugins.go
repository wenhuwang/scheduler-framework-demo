package plugins

import (
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/klog"
	framework "k8s.io/kubernetes/pkg/scheduler/framework/v1alpha1"
	metricsv "k8s.io/metrics/pkg/client/clientset/versioned"
)

// 插件名称
const Name = "sample-plugin"

type Args struct {
	FavoriteColor  string `json:"favorite_color,omitempty"`
	FavoriteNumber int    `json:"favorite_number,omitempty"`
	ThanksTo       string `json:"thanks_to,omitempty"`
	KubeConfig     string `json:"kubeconfig,omitempty"`
}

type Sample struct {
	kubeClient    *kubernetes.Clientset
	metricsClient *metricsv.Clientset
	args          *Args
	handle        framework.FrameworkHandle
}

func (s *Sample) Name() string {
	return Name
}

func (s *Sample) PreFilter(pc *framework.PluginContext, pod *v1.Pod) *framework.Status {
	klog.V(3).Infof("prefilter pod: %v", pod.Name)
	return framework.NewStatus(framework.Success, "")
}

func (s *Sample) Filter(pc *framework.PluginContext, pod *v1.Pod, nodeName string) *framework.Status {
	klog.V(3).Infof("filter pod: %v, node: %v", pod.Name, nodeName)
	// if nodeName == "shtl009063227" {
	// 	return framework.NewStatus(framework.Unschedulable, "")
	// }
	return framework.NewStatus(framework.Success, "")
}

func (s *Sample) Score(pc *framework.PluginContext, pod *v1.Pod, nodeName string) (int, *framework.Status) {
	klog.V(3).Infof("Score pod: %v, node: %v", pod.Name, nodeName)
	nodeInfo, ok := s.handle.NodeInfoSnapshot().NodeInfoMap[nodeName]
	if !ok {
		return framework.MinNodeScore, framework.NewStatus(framework.Error, "")
	}

	nodeMetrics, err := s.metricsClient.MetricsV1beta1().NodeMetricses().Get(nodeName, metav1.GetOptions{})
	if err != nil {
		klog.Errorf("get node %s metrics err: %+v", nodeName, err)
	}

	// Calculation formula Details:
	// (cpu((100-Usage*100)/10) + memory((100-usage*100)/10))/2
	cpuMetrics := float64(100*nodeMetrics.Usage.Cpu().MilliValue()) / float64(nodeInfo.AllocatableResource().MilliCPU)
	memMetrics := float64(100*nodeMetrics.Usage.Memory().Value()) / float64(nodeInfo.AllocatableResource().Memory)
	klog.V(3).Infof("node name: %s, cpu metrics %f, mem metrics %f", nodeInfo.Node().Name, cpuMetrics, memMetrics)

	nodeScore := ((100-cpuMetrics)/10 + (100-memMetrics)/10) / 2
	return int(nodeScore), framework.NewStatus(framework.Success, "")
}

func (s *Sample) NormalizeScore(pc *framework.PluginContext, pod *v1.Pod, nodeScoreList framework.NodeScoreList) *framework.Status {
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

func New(configuration *runtime.Unknown, f framework.FrameworkHandle) (framework.Plugin, error) {
	args := &Args{}
	if err := framework.DecodeInto(configuration, args); err != nil {
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
