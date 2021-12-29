module scheduler-framework-demo

go 1.15

require (
	k8s.io/api v0.19.16
	k8s.io/apimachinery v0.19.16
	k8s.io/component-base v0.19.16
	k8s.io/klog/v2 v2.2.0
	k8s.io/kubernetes v1.19.16
)

replace (
	k8s.io/api v0.0.0 => k8s.io/api v0.19.16
	k8s.io/apiextensions-apiserver v0.0.0 => k8s.io/apiextensions-apiserver v0.19.16
	k8s.io/apimachinery v0.0.0 => k8s.io/apimachinery v0.19.16
	k8s.io/apiserver v0.0.0 => k8s.io/apiserver v0.19.16
	k8s.io/cli-runtime v0.0.0 => k8s.io/cli-runtime v0.19.16
	k8s.io/client-go v0.0.0 => k8s.io/client-go v0.19.16
	k8s.io/cloud-provider v0.0.0 => k8s.io/cloud-provider v0.19.16
	k8s.io/cluster-bootstrap v0.0.0 => k8s.io/cluster-bootstrap v0.19.16
	k8s.io/code-generator v0.0.0 => k8s.io/code-generator v0.19.16
	k8s.io/component-base v0.0.0 => k8s.io/component-base v0.19.16
	k8s.io/component-helpers v0.0.0 => k8s.io/component-helpers v0.19.16
	k8s.io/controller-manager v0.0.0 => k8s.io/controller-manager v0.19.16
	k8s.io/cri-api v0.0.0 => k8s.io/cri-api v0.19.16
	k8s.io/csi-api v0.0.0 => k8s.io/csi-api v0.19.16
	k8s.io/csi-translation-lib v0.0.0 => k8s.io/csi-translation-lib v0.19.16
	k8s.io/klog/v2 => k8s.io/klog/v2 v2.2.0
	k8s.io/kube-aggregator v0.0.0 => k8s.io/kube-aggregator v0.19.16
	k8s.io/kube-controller-manager v0.0.0 => k8s.io/kube-controller-manager v0.19.16
	k8s.io/kube-proxy v0.0.0 => k8s.io/kube-proxy v0.19.16
	k8s.io/kube-scheduler v0.0.0 => k8s.io/kube-scheduler v0.19.16
	k8s.io/kubectl => k8s.io/kubectl v0.19.16
	k8s.io/kubelet v0.0.0 => k8s.io/kubelet v0.19.16
	k8s.io/kubernetes v0.0.0 => k8s.io/kubernetes v1.19.16
	k8s.io/legacy-cloud-providers v0.0.0 => k8s.io/legacy-cloud-providers v0.19.16
	k8s.io/metrics v0.0.0 => k8s.io/metrics v0.19.16
	k8s.io/mount-utils v0.0.0 => k8s.io/mount-utils v0.19.16
	k8s.io/node-api v0.0.0 => k8s.io/node-api v0.19.16
	k8s.io/repo-infra v0.0.0 => k8s.io/repo-infra v0.19.16
	k8s.io/sample-apiserver v0.0.0 => k8s.io/sample-apiserver v0.19.16
)
