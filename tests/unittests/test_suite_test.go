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
	loadedVirtualMachineClusterInstanceTypes []instancetypev1beta1.VirtualMachineClusterInstancetype
	loadedVirtualMachineClusterPreferences   []instancetypev1beta1.VirtualMachineClusterPreference
)

// The following functions are taken from SSP operator
// https://github.com/kubevirt/ssp-operator/blob/main/internal/operands/common-instancetypes/reconcile.go
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

func fetchResourcesFromBundle() ([]instancetypev1beta1.VirtualMachineClusterInstancetype,
	[]instancetypev1beta1.VirtualMachineClusterPreference, error) {
	virtualMachineClusterInstancetypes, err :=
		FetchBundleResource[instancetypev1beta1.VirtualMachineClusterInstancetype](clusterInstanceTypesBundlePath)
	if err != nil {
		return nil, nil, err
	}
	virtualMachineClusterPreferences, err :=
		FetchBundleResource[instancetypev1beta1.VirtualMachineClusterPreference](clusterPreferencesBundlePath)
	if err != nil {
		return nil, nil, err
	}
	return virtualMachineClusterInstancetypes, virtualMachineClusterPreferences, err
}

// end of functions taken from SSP operator

func loadBundles() {
	var err error
	loadedVirtualMachineClusterInstanceTypes, loadedVirtualMachineClusterPreferences, err = fetchResourcesFromBundle()
	// It is expected that make export is run first
	Expect(err).ToNot(HaveOccurred())
}

var _ = BeforeSuite(func() {
	loadBundles()

	Expect(len(loadedVirtualMachineClusterInstanceTypes)).ToNot(BeZero())
	Expect(len(loadedVirtualMachineClusterPreferences)).ToNot(BeZero())
})

func TestUnit(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Unit test suite")
}
