package tests

import (
	"fmt"
	"strconv"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"k8s.io/apimachinery/pkg/api/resource"
	instancetypev1beta1 "kubevirt.io/api/instancetype/v1beta1"
)

const (
	expectedVersion          = "instancetype.kubevirt.io/v1beta1"
	expectedVendor           = "kubevirt.io"
	instanceTypeErrorMessage = "%s label does not match in instancetype: %s"
	preferenceErrorMessage   = "%s label does not match in preference: %s"
)

var _ = Describe("Common instance types unit tests", func() {

	var (
		testsClusterPreferences  []instancetypev1beta1.VirtualMachineClusterPreference
		testsClusterInstanceType []instancetypev1beta1.VirtualMachineClusterInstancetype
	)

	skipLabels := map[string]bool{
		"instancetype.kubevirt.io/os-type":                      true,
		"instancetype.kubevirt.io/deprecated":                   true,
		"instancetype.kubevirt.io/common-instancetypes-version": true,
		"instancetype.kubevirt.io/version":                      true,
		"instancetype.kubevirt.io/class":                        true,
	}

	instanceTypeFunctionMap := map[string]func(string, string, instancetypev1beta1.VirtualMachineClusterInstancetype) error{
		"instancetype.kubevirt.io/cpu":                   checkCPU,
		"instancetype.kubevirt.io/memory":                checkMemory,
		"instancetype.kubevirt.io/hugepages":             checkHugepages,
		"instancetype.kubevirt.io/gpus":                  checkGPUs,
		"instancetype.kubevirt.io/numa":                  checkNuma,
		"instancetype.kubevirt.io/dedicatedCPUPlacement": checkDedicatedCPUPlacement,
		"instancetype.kubevirt.io/isolateEmulatorThread": checkIsolateEmulatorThread,
		"instancetype.kubevirt.io/vendor":                instancetypeCheckVendor,
	}

	preferenceFunctionMap := map[string]func(string, string, instancetypev1beta1.VirtualMachineClusterPreference) error{
		"openshift.io/display-name":       checkDisplayName,
		"instancetype.kubevirt.io/vendor": preferenceCheckVendor,
	}

	BeforeEach(func() {
		// Restore resources to make sure we have unchanged data
		copy(testsClusterPreferences, loadedVirtualMachineClusterPreferences)
		copy(testsClusterInstanceType, loadedVirtualMachineClusterInstanceTypes)
	})

	Context("VirtualMachineClusterPreference", func() {
		It("check version", func() {
			for _, preference := range testsClusterPreferences {
				Expect(preference.APIVersion).To(Equal(expectedVersion))
			}
		})

		It("[test_id:10799]check if labels match resources", func() {
			for _, preference := range testsClusterPreferences {
				for key, value := range preference.Labels {
					if skipLabels[key] {
						continue
					}
					Expect(preferenceFunctionMap).To(HaveKey(key))
					Expect(preferenceFunctionMap[key](value, key, preference)).To(Succeed())
				}
			}
		})

	})

	Context("VirtualMachineClusterInstanceType", func() {
		It("check version", func() {
			for _, instanceType := range testsClusterInstanceType {
				Expect(instanceType.APIVersion).To(Equal(expectedVersion))
			}
		})

		It("[test_id:10801]check if labels match resources", func() {
			for _, instanceType := range testsClusterInstanceType {
				for key, value := range instanceType.Labels {
					if skipLabels[key] {
						continue
					}
					Expect(instanceTypeFunctionMap).To(HaveKey(key))
					Expect(instanceTypeFunctionMap[key](value, key, instanceType)).To(Succeed())
				}
			}
		})
	})
})

func checkCPU(labelValue, labelName string, instanceType instancetypev1beta1.VirtualMachineClusterInstancetype) error {
	expectedCPU, err := strconv.ParseUint(labelValue, 10, 32)
	if err != nil {
		return err
	}

	if uint64(instanceType.Spec.CPU.Guest) != expectedCPU {
		return fmt.Errorf(instanceTypeErrorMessage, labelName, instanceType.Name)
	}

	return nil
}

func checkMemory(labelValue, labelName string, instanceType instancetypev1beta1.VirtualMachineClusterInstancetype) error {
	memory, err := resource.ParseQuantity(labelValue)
	if err != nil {
		return err
	}

	if instanceType.Spec.Memory.Guest != memory {
		return fmt.Errorf(instanceTypeErrorMessage, labelName, instanceType.Name)
	}

	return nil
}

func checkHugepages(labelValue, labelName string, instanceType instancetypev1beta1.VirtualMachineClusterInstancetype) error {
	boolValue, err := strconv.ParseBool(labelValue)
	if err != nil {
		return err
	}

	if (instanceType.Spec.Memory.Hugepages == nil && boolValue) || (instanceType.Spec.Memory.Hugepages != nil && !boolValue) {
		return fmt.Errorf(instanceTypeErrorMessage, labelName, instanceType.Name)
	}

	return nil
}

func checkGPUs(labelValue, labelName string, instanceType instancetypev1beta1.VirtualMachineClusterInstancetype) error {
	boolValue, err := strconv.ParseBool(labelValue)
	if err != nil {
		return err
	}

	if (instanceType.Spec.GPUs == nil && boolValue) || (instanceType.Spec.GPUs != nil && !boolValue) {
		return fmt.Errorf(instanceTypeErrorMessage, labelName, instanceType.Name)
	}

	return nil
}

func checkNuma(labelValue, labelName string, instanceType instancetypev1beta1.VirtualMachineClusterInstancetype) error {
	boolValue, err := strconv.ParseBool(labelValue)
	if err != nil {
		return err
	}

	if (instanceType.Spec.CPU.NUMA == nil && boolValue) || (instanceType.Spec.CPU.NUMA != nil && !boolValue) {
		return fmt.Errorf(instanceTypeErrorMessage, labelName, instanceType.Name)
	}

	return nil
}

func checkDedicatedCPUPlacement(labelValue, labelName string, instanceType instancetypev1beta1.VirtualMachineClusterInstancetype) error {
	boolValue, err := strconv.ParseBool(labelValue)
	if err != nil {
		return err
	}

	if (instanceType.Spec.CPU.DedicatedCPUPlacement == nil && boolValue) ||
		(instanceType.Spec.CPU.DedicatedCPUPlacement != nil && !boolValue) {
		return fmt.Errorf(instanceTypeErrorMessage, labelName, instanceType.Name)
	}

	return nil
}

func checkIsolateEmulatorThread(labelValue, labelName string, instanceType instancetypev1beta1.VirtualMachineClusterInstancetype) error {
	boolValue, err := strconv.ParseBool(labelValue)
	if err != nil {
		return err
	}

	if (instanceType.Spec.CPU.IsolateEmulatorThread == nil && boolValue) ||
		(instanceType.Spec.CPU.IsolateEmulatorThread != nil && !boolValue) {
		return fmt.Errorf(instanceTypeErrorMessage, labelName, instanceType.Name)
	}

	return nil
}

func checkDisplayName(labelValue, labelName string, preference instancetypev1beta1.VirtualMachineClusterPreference) error {
	if preference.Name != labelValue {
		return fmt.Errorf(instanceTypeErrorMessage, labelName, preference.Name)
	}

	return nil
}

func preferenceCheckVendor(labelValue, labelName string, instanceType instancetypev1beta1.VirtualMachineClusterPreference) error {
	if labelValue != expectedVendor {
		return fmt.Errorf(preferenceErrorMessage, labelName, instanceType.Name)
	}

	return nil
}

func instancetypeCheckVendor(labelValue, labelName string, instanceType instancetypev1beta1.VirtualMachineClusterInstancetype) error {
	if labelValue != expectedVendor {
		return fmt.Errorf(instanceTypeErrorMessage, labelName, instanceType.Name)
	}

	return nil
}
