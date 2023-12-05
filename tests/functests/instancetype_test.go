package tests

import (
	"context"
	"fmt"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/rand"
	"k8s.io/utils/ptr"
	v1 "kubevirt.io/api/core/v1"
	instancetypev1beta1 "kubevirt.io/api/instancetype/v1beta1"
)

const vmReadyTimeout = 300 * time.Second

var _ = Describe("Common instance types func tests", func() {

	Context("VirtualMachine using a preference with resource requirements", func() {
		It("is rejected if it does not provide enough resources", func() {
			createInstancetype(2, "tiny-instancetype", "64M")
			intanceTypeMatcher := v1.InstancetypeMatcher{
				Name: "tiny-instancetype",
				Kind: "VirtualMachineInstancetype",
			}
			vm := randomVMWithInstancetype(&intanceTypeMatcher)

			clusterPreferences, err := virtClient.VirtualMachineClusterPreference().List(context.Background(), metav1.ListOptions{})
			Expect(err).ToNot(HaveOccurred())
			Expect(clusterPreferences.Items).ToNot(BeEmpty())

			for _, preference := range clusterPreferences.Items {
				vm.Name = fmt.Sprintf("test-vm-%s", preference.Name)
				vm.Spec.Preference = &v1.PreferenceMatcher{
					Name: preference.Name,
				}

				_, err := virtClient.VirtualMachine(testNamespace).Create(context.Background(), vm)
				Expect(err).To(MatchError(
					fmt.Sprintf(
						"admission webhook \"virtualmachine-validator.kubevirt.io\" denied the request: "+
							"failure checking preference requirements: "+
							"insufficient Memory resources of 64M provided by instance type, preference requires %s",
						preference.Spec.Requirements.Memory.Guest.String(),
					)))
			}
		})

		It("can be created when enough resources are provided", func() {
			intanceTypeMatcher := v1.InstancetypeMatcher{
				Name: "u1.large",
			}
			vm := randomVMWithInstancetype(&intanceTypeMatcher)

			clusterPreferences, err := virtClient.VirtualMachineClusterPreference().List(context.Background(), metav1.ListOptions{})
			Expect(err).ToNot(HaveOccurred())
			Expect(clusterPreferences.Items).ToNot(BeEmpty())

			for _, preference := range clusterPreferences.Items {
				vm.Name = fmt.Sprintf("test-vm-%s", preference.Name)
				vm.Spec.Preference = &v1.PreferenceMatcher{
					Name: preference.Name,
				}

				vm, err = virtClient.VirtualMachine(testNamespace).Create(context.Background(), vm)
				Expect(err).ToNot(HaveOccurred())
				Eventually(func(g Gomega) {
					vm, err = virtClient.VirtualMachine(testNamespace).Get(context.Background(), vm.Name, &metav1.GetOptions{})
					g.Expect(err).ToNot(HaveOccurred())
					g.Expect(vm.Status.Ready).To(BeTrue())
				}, vmReadyTimeout, time.Second).Should(Succeed())

				err = virtClient.VirtualMachine(testNamespace).Delete(context.Background(), vm.Name, &metav1.DeleteOptions{})
				Expect(err).ToNot(HaveOccurred())
			}
		})
	})
})

func randomVMWithInstancetype(instancetypeMatcher *v1.InstancetypeMatcher) *v1.VirtualMachine {
	name := "test-vm-" + rand.String(5)

	return &v1.VirtualMachine{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: testNamespace,
		},
		Spec: v1.VirtualMachineSpec{
			Running:      ptr.To(true),
			Instancetype: instancetypeMatcher,
			Template: &v1.VirtualMachineInstanceTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels:    map[string]string{"name": name},
					Namespace: testNamespace,
				},
				Spec: v1.VirtualMachineInstanceSpec{
					Volumes: []v1.Volume{
						{
							Name: "disk",
							VolumeSource: v1.VolumeSource{
								ContainerDisk: &v1.ContainerDiskSource{
									Image: "quay.io/containerdisks/fedora:latest",
								},
							},
						},
					},
				},
			},
		},
	}
}

func createInstancetype(cpu int, name, memory string) {
	instancetype := &instancetypev1beta1.VirtualMachineInstancetype{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
		Spec: instancetypev1beta1.VirtualMachineInstancetypeSpec{
			CPU: instancetypev1beta1.CPUInstancetype{
				Guest: uint32(cpu),
			},
			Memory: instancetypev1beta1.MemoryInstancetype{
				Guest: resource.MustParse(memory),
			},
		},
	}
	_, err := virtClient.VirtualMachineInstancetype(testNamespace).Create(context.Background(), instancetype, metav1.CreateOptions{})
	Expect(err).ToNot(HaveOccurred())
}
