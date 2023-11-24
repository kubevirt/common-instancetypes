//go:build tools

package tools

import (
	_ "github.com/mikefarah/yq/v4"
	_ "github.com/yannh/kubeconform/cmd/kubeconform"
	_ "sigs.k8s.io/kustomize/kustomize/v5"
)
