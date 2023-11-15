package tests

import (
	"bytes"
	"io"
	"os"
	"testing"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"k8s.io/apimachinery/pkg/util/yaml"
	instancetypev1beta1 "kubevirt.io/api/instancetype/v1beta1"
)

const (
	clusterInstanceTypesBundlePath = "../../_build/common-clusterinstancetypes-bundle.yaml"
	clusterPreferencesBundlePath   = "../../_build/common-clusterpreferences-bundle.yaml"
)

var (
	virtualMachineClusterInstanceTypes []instancetypev1beta1.VirtualMachineClusterInstancetype
	virtualMachineClusterPreferences   []instancetypev1beta1.VirtualMachineClusterPreference
)

// The following functions are taken from SSP operator
// https://github.com/kubevirt/ssp-operator/blob/c3de0957e0446f4b5ab23947d124b9ac9107d0a3/internal/operands/common-instancetypes/reconcile.go#L83
// https://github.com/kubevirt/ssp-operator/blob/c3de0957e0446f4b5ab23947d124b9ac9107d0a3/internal/operands/common-instancetypes/reconcile.go#L99
// https://github.com/kubevirt/ssp-operator/blob/c3de0957e0446f4b5ab23947d124b9ac9107d0a3/internal/operands/common-instancetypes/reconcile.go#L107
type clusterType interface {
	instancetypev1beta1.VirtualMachineClusterInstancetype | instancetypev1beta1.VirtualMachineClusterPreference
}

func decodeResources[C clusterType](b []byte) ([]C, error) {
	decoder := yaml.NewYAMLOrJSONDecoder(bytes.NewReader(b), 1024)
	var bundle []C
	for {
		bundleResource := new(C)
		err := decoder.Decode(bundleResource)
		if err == io.EOF {
			return bundle, nil
		}
		if err != nil {
			return nil, err
		}
		bundle = append(bundle, *bundleResource)
	}
}

func FetchBundleResource[C clusterType](path string) ([]C, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	return decodeResources[C](file)
}

func fetchResourcesFromBundle() ([]instancetypev1beta1.VirtualMachineClusterInstancetype, []instancetypev1beta1.VirtualMachineClusterPreference, error) {
	virtualMachineClusterInstancetypes, err := FetchBundleResource[instancetypev1beta1.VirtualMachineClusterInstancetype](clusterInstanceTypesBundlePath)
	if err != nil {
		return nil, nil, err
	}
	virtualMachineClusterPreferences, err := FetchBundleResource[instancetypev1beta1.VirtualMachineClusterPreference](clusterPreferencesBundlePath)
	if err != nil {
		return nil, nil, err
	}
	return virtualMachineClusterInstancetypes, virtualMachineClusterPreferences, err
}

// end of functions taken from SSP operator

func loadBundles() {
	var err error
	virtualMachineClusterInstanceTypes, virtualMachineClusterPreferences, err = fetchResourcesFromBundle()
	// It is expected that make export is run first
	Expect(err).ToNot(HaveOccurred())
}

var _ = BeforeSuite(func() {
	loadBundles()

	Expect(len(virtualMachineClusterInstanceTypes)).ToNot(BeZero())
	Expect(len(virtualMachineClusterPreferences)).ToNot(BeZero())
})

func TestFunctional(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Unit test suite")
}
