package plugins

import (
	"context"
	"encoding/json"
	"fmt"

	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/klog"
	framework "k8s.io/kubernetes/pkg/scheduler/framework/v1alpha1"
	"k8s.io/kubernetes/pkg/scheduler/nodeinfo"
)

// 插件名称
const Name = "RealResourceUsage"

type Args struct {
	FavoriteColor  string `json:"favorite_color,omitempty"`
	FavoriteNumber int    `json:"favorite_number,omitempty"`
	ThanksTo       string `json:"thanks_to,omitempty"`
}

type Sample struct {
	args   *Args
	handle framework.FrameworkHandle
}

var _ framework.FilterPlugin = &Sample{}
var _ framework.ScorePlugin = &Sample{}

func (s *Sample) Name() string {
	return Name
}

// func (s *Sample) PreFilter(ctx context.Context, pod *v1.Pod) *framework.Status {
// 	klog.V(3).Infof("prefilter pod: %v", pod.Name)
// 	return framework.NewStatus(framework.Success, "")
// }

// func (s *Sample) PreFilterExtensions() framework.PreFilterExtensions {
// 	return nil
// }

func (s *Sample) Filter(ctx context.Context, state *framework.CycleState, pod *v1.Pod, nodeInfo *nodeinfo.NodeInfo) *framework.Status {
	klog.V(3).Infof("filter pod: %v, node: %v", pod.Name, nodeInfo.Node().Name)
	return framework.NewStatus(framework.Success, "")
}

func (s *Sample) Score(ctx context.Context, state *framework.CycleState, pod *v1.Pod, nodeName string) (int64, *framework.Status) {
	klog.V(10).Infof("Score pod: %v, node: %v", pod.Name, nodeName)
	nodeInfo, err := s.handle.SnapshotSharedLister().NodeInfos().Get(nodeName)
	if err != nil {
		return 0, framework.NewStatus(framework.Error, fmt.Sprintf("getting node %q from Snapshot: %v", nodeName, err))
	}

	metricsString, ok := nodeInfo.Node().Annotations[metricsAnnotation]
	if !ok {
		return 0, framework.NewStatus(framework.Error, fmt.Sprintf("node %s annotations %s is not exists.", nodeName, metricsAnnotation))
	}

	metrics := make(map[string]float64)
	if err = json.Unmarshal([]byte(metricsString), &metrics); err != nil {
		return 0, framework.NewStatus(framework.Error, fmt.Sprintf("node %s annotations %s parse failed. %v", nodeName,
			metricsAnnotation, err))
	}
	cpuMetrics := metrics[cpuMetricsKey] * 100
	nodeScore := int64(100 - cpuMetrics)

	if klog.V(10) {
		klog.Infof("%v -> %v: %v, score %d", pod.Name, nodeName, Name, nodeScore)
	}
	return nodeScore, framework.NewStatus(framework.Success, "")
}

func (s *Sample) ScoreExtensions() framework.ScoreExtensions {
	return nil
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

func New(configuration *runtime.Unknown, f framework.FrameworkHandle) (framework.Plugin, error) {
	args := &Args{}
	if err := framework.DecodeInto(configuration, args); err != nil {
		return nil, err
	}
	klog.V(3).Infof("get plugin config args: %+v", args)

	return &Sample{
		args:   args,
		handle: f,
	}, nil
}
