package tests

import (
	"fmt"
	"strconv"
	"strings"

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

	skipLabels := map[string]bool{
		"instancetype.kubevirt.io/arch":                         true,
		"instancetype.kubevirt.io/os-type":                      true,
		"instancetype.kubevirt.io/deprecated":                   true,
		"instancetype.kubevirt.io/common-instancetypes-version": true,
		"instancetype.kubevirt.io/version":                      true,
		"instancetype.kubevirt.io/class":                        true,
		"instancetype.kubevirt.io/icon-pf":                      true,
	}

	instanceTypeFunctionMap := map[string]func(string, string, instancetypev1beta1.VirtualMachineClusterInstancetype) error{
		"instancetype.kubevirt.io/cpu":                   checkCPU,
		"instancetype.kubevirt.io/memory":                checkMemory,
		"instancetype.kubevirt.io/size":                  checkSize,
		"instancetype.kubevirt.io/hugepages":             checkHugepages,
		"instancetype.kubevirt.io/gpus":                  checkGPUs,
		"instancetype.kubevirt.io/numa":                  checkNuma,
		"instancetype.kubevirt.io/dedicatedCPUPlacement": checkDedicatedCPUPlacement,
		"instancetype.kubevirt.io/isolateEmulatorThread": checkIsolateEmulatorThread,
		"instancetype.kubevirt.io/realtime":              checkRealtime,
		"instancetype.kubevirt.io/vendor":                instancetypeCheckVendor,
	}

	preferenceFunctionMap := map[string]func(string, string, instancetypev1beta1.VirtualMachineClusterPreference) error{
		"openshift.io/display-name":       checkDisplayName,
		"instancetype.kubevirt.io/vendor": preferenceCheckVendor,
	}

	Context("VirtualMachineClusterPreference", func() {
		It("check version", func() {
			for _, preference := range loadedVirtualMachineClusterPreferences {
				Expect(preference.APIVersion).To(Equal(expectedVersion))
			}
		})

		It("check if labels match resources", func() {
			for _, preference := range loadedVirtualMachineClusterPreferences {
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
			for _, instanceType := range loadedVirtualMachineClusterInstanceTypes {
				Expect(instanceType.APIVersion).To(Equal(expectedVersion))
			}
		})

		It("check if labels match resources", func() {
			for _, instanceType := range loadedVirtualMachineClusterInstanceTypes {
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

func checkSize(labelValue, labelName string, instanceType instancetypev1beta1.VirtualMachineClusterInstancetype) error {
	if strings.Split(instanceType.Name, ".")[1] != labelValue {
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

func checkRealtime(labelValue, labelName string, instanceType instancetypev1beta1.VirtualMachineClusterInstancetype) error {
	boolValue, err := strconv.ParseBool(labelValue)
	if err != nil {
		return err
	}

	if (instanceType.Spec.CPU.Realtime == nil && boolValue) ||
		(instanceType.Spec.CPU.Realtime != nil && !boolValue) {
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
