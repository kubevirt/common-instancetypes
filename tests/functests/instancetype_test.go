package tests

import (
	"bytes"
	"context"
	"crypto/ed25519"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"golang.org/x/crypto/ssh"
	extv1 "k8s.io/apiextensions-apiserver/pkg/apis/apiextensions/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8srand "k8s.io/apimachinery/pkg/util/rand"
	"k8s.io/client-go/util/jsonpath"
	"k8s.io/utils/ptr"
	v1 "kubevirt.io/api/core/v1"
	instancetypev1beta1 "kubevirt.io/api/instancetype/v1beta1"
	"kubevirt.io/client-go/kubecli"
)

const (
	sshPort = 22
)

type testFn func(kubecli.KubevirtClient, string)

var _ = Describe("Common instance types func tests", func() {
	var (
		vm  *v1.VirtualMachine
		err error
	)

	JustAfterEach(func() {
		// On failure dump the current state of the VirtualMachine into the test output
		// Useful for debugging when the namespace has already been cleaned up
		if CurrentSpecReport().Failed() && vm != nil {
			vm, err = virtClient.VirtualMachine(testNamespace).Get(context.Background(), vm.Name, metav1.GetOptions{})
			if err != nil && errors.IsNotFound(err) {
				GinkgoWriter.Printf("VM %s defined but not created", vm.Name)
				return
			}
			Expect(err).ToNot(HaveOccurred())

			var status []byte
			status, err = json.MarshalIndent(vm, "", "\t")
			Expect(err).ToNot(HaveOccurred())
			GinkgoWriter.Print(string(status))
		}
	})

	AfterEach(func() {
		if vm == nil {
			return
		}
		err = virtClient.VirtualMachine(testNamespace).Delete(context.Background(), vm.Name,
			metav1.DeleteOptions{GracePeriodSeconds: ptr.To[int64](0)})
		if err != nil && !errors.IsNotFound(err) {
			Expect(err).ToNot(HaveOccurred())
		}
	})

	It("[test_id:10735] VirtualMachine using an instancetype can be created", func() {
		for _, instancetype := range getClusterInstancetypes(virtClient) {
			vm = randomVM(&v1.InstancetypeMatcher{Name: instancetype.Name}, nil, v1.RunStrategyHalted)
			vm, err = virtClient.VirtualMachine(testNamespace).Create(context.Background(), vm, metav1.CreateOptions{})
			Expect(err).ToNot(HaveOccurred())
		}
	})

	Context("VirtualMachine using a preference with resource requirements", func() {
		var skipPreference = map[string]any{
			"legacy":    nil,
			"linux":     nil,
			"linux.efi": nil,
		}

		clusterPreferencesWithRequirements :=
			func() []instancetypev1beta1.VirtualMachineClusterPreference {
				clusterPreferences, err := virtClient.VirtualMachineClusterPreference().List(context.Background(), metav1.ListOptions{})
				Expect(err).ToNot(HaveOccurred())
				Expect(clusterPreferences.Items).ToNot(BeEmpty())

				var completePreferences []instancetypev1beta1.VirtualMachineClusterPreference
				for _, preference := range clusterPreferences.Items {
					if _, ok := skipPreference[preference.Name]; !ok {
						Expect(preference.Spec.Requirements).ToNot(BeNil())
						completePreferences = append(completePreferences, preference)
					}
				}
				return completePreferences
			}

		It("[test_id:10736] is rejected if it does not provide enough memory resources", func() {
			createInstancetype(8, "tiny-instancetype-memory", "64M")
			instanceTypeMatcher := v1.InstancetypeMatcher{
				Name: "tiny-instancetype-memory",
				Kind: "VirtualMachineInstancetype",
			}
			for _, preference := range clusterPreferencesWithRequirements() {
				vm = randomVM(&instanceTypeMatcher, &v1.PreferenceMatcher{Name: preference.Name}, v1.RunStrategyHalted)
				_, err = virtClient.VirtualMachine(testNamespace).Create(context.Background(), vm, metav1.CreateOptions{})
				Expect(err).To(MatchError(
					fmt.Sprintf(
						"admission webhook \"virtualmachine-validator.kubevirt.io\" denied the request: "+
							"failure checking preference requirements: "+
							"insufficient Memory resources of 64M provided by instance type, preference requires %s",
						preference.Spec.Requirements.Memory.Guest.String(),
					)))
			}
		})

		It("[test_id:TODO] is rejected if it does not provide enough cpu resources", func() {
			createInstancetype(1, "tiny-instancetype-cpu", "64M")
			instanceTypeMatcher := v1.InstancetypeMatcher{
				Name: "tiny-instancetype-cpu",
				Kind: "VirtualMachineInstancetype",
			}
			for _, preference := range clusterPreferencesWithRequirements() {
				if preference.Spec.Requirements.CPU == nil || preference.Spec.Requirements.CPU.Guest < 2 {
					continue
				}
				vm = randomVM(&instanceTypeMatcher, &v1.PreferenceMatcher{Name: preference.Name}, v1.RunStrategyHalted)
				_, err = virtClient.VirtualMachine(testNamespace).Create(context.Background(), vm, metav1.CreateOptions{})
				Expect(err).To(MatchError(MatchRegexp("1 vCPU(?:s)? provided by (?:the )?instance type")))
			}
		})

		It("[test_id:10737] can be created when enough resources are provided", func() {
			for _, preference := range clusterPreferencesWithRequirements() {
				instanceTypeName := "it-for-" + preference.Name
				createInstancetype(preference.Spec.Requirements.CPU.Guest, instanceTypeName, preference.Spec.Requirements.Memory.Guest.String())
				instanceTypeMatcher := v1.InstancetypeMatcher{
					Name: instanceTypeName,
					Kind: "VirtualMachineInstancetype",
				}

				vm = randomVM(&instanceTypeMatcher, &v1.PreferenceMatcher{Name: preference.Name}, v1.RunStrategyHalted)
				vm, err = virtClient.VirtualMachine(testNamespace).Create(context.Background(), vm, metav1.CreateOptions{})
				Expect(err).ToNot(HaveOccurred())
			}
		})

		It("[test_id:TODO] verifies all preferences have at least one compatible instance type", func() {
			preferenceInstancetypeMap := map[string]string{
				"centos.stream9.dpdk": "u1.2xlarge",
				"rhel.8.dpdk":         "u1.2xlarge",
				"rhel.9.dpdk":         "u1.2xlarge",
				"rhel.9.realtime":     "u1.xlarge",
				"windows.11":          "u1.2xmedium",
				"windows.11.virtio":   "u1.2xmedium",
				"oraclelinux":         "u1.2xmedium",
			}
			for _, preference := range clusterPreferencesWithRequirements() {
				instancetype := "u1.medium"
				if pInstancetype, ok := preferenceInstancetypeMap[preference.Name]; ok {
					instancetype = pInstancetype
				}
				vm = randomVM(&v1.InstancetypeMatcher{Name: instancetype}, &v1.PreferenceMatcher{Name: preference.Name}, v1.RunStrategyHalted)
				vm, err = virtClient.VirtualMachine(testNamespace).Create(context.Background(), vm, metav1.CreateOptions{})
				Expect(err).ToNot(HaveOccurred())
			}
		})
	})

	Context("VirtualMachine using a preference is able to boot", func() {
		var privKey ed25519.PrivateKey

		expectSSHToRunCommandOnLinux := func(username string) testFn {
			return func(virtClient kubecli.KubevirtClient, vmName string) {
				var signer ssh.Signer
				signer, err = ssh.NewSignerFromKey(privKey)
				Expect(err).ToNot(HaveOccurred())
				expectSSHToRunCommand(virtClient, vmName, username, ssh.PublicKeys(signer))
			}
		}

		BeforeEach(func() {
			_, privKey, err = ed25519.GenerateKey(nil)
			Expect(err).ToNot(HaveOccurred())
		})

		DescribeTable("a Linux guest with", func(containerDisk, instancetype string, preferences map[string]string, testFns []testFn) {
			preference, hasArch := preferences[preferenceArch]
			if !hasArch {
				Skip(fmt.Sprintf("skipping as no preference provided for arch %s", preferenceArch))
			}
			vm = randomVM(&v1.InstancetypeMatcher{Name: instancetype}, &v1.PreferenceMatcher{Name: preference}, v1.RunStrategyAlways)
			addContainerDisk(vm, containerDisk)
			addCloudInitWithAuthorizedKey(vm, privKey)
			vm, err = virtClient.VirtualMachine(testNamespace).Create(context.Background(), vm, metav1.CreateOptions{})
			Expect(err).ToNot(HaveOccurred())
			expectVMToBeReady(virtClient, vm.Name, defaultVMReadyTimeout)
			for _, testFn := range testFns {
				testFn(virtClient, vm.Name)
			}
		},
			Entry("[test_id:10738] Fedora", fedoraContainerDisk, "u1.small",
				map[string]string{"amd64": "fedora", "arm64": "fedora.arm64", "s390x": "fedora.s390x"},
				[]testFn{expectGuestAgentToBeConnected, expectSSHToRunCommandOnLinux("fedora")}),
			Entry("[test_id:10745] CentOS Stream 9", centosStream9ContainerDisk, "u1.small",
				map[string]string{"amd64": "centos.stream9", "arm64": "centos.stream9", "s390x": "centos.stream9"},
				[]testFn{expectGuestAgentToBeConnected, expectSSHToRunCommandOnLinux("cloud-user")}),
			Entry("[test_id:TODO] CentOS Stream 10", centosStream10ContainerDisk, "u1.small",
				map[string]string{"amd64": "centos.stream10", "arm64": "centos.stream10", "s390x": "centos.stream10"},
				[]testFn{expectGuestAgentToBeConnected, expectSSHToRunCommandOnLinux("cloud-user")}),
			Entry("[test_id:TODO] RHEL 9", rhel9ContainerDisk, "u1.small",
				map[string]string{"amd64": "rhel.9", "arm64": "rhel.9.arm64", "s390x": "rhel.9.s390x"},
				[]testFn{expectGuestAgentToBeConnected, expectSSHToRunCommandOnLinux("cloud-user")}),
			Entry("[test_id:TODO] RHEL 10", rhel10ContainerDisk, "u1.small",
				map[string]string{"amd64": "rhel.10", "arm64": "rhel.10.arm64", "s390x": "rhel.10.s390x"},
				[]testFn{expectGuestAgentToBeConnected, expectSSHToRunCommandOnLinux("cloud-user")}),
			Entry("[test_id:10741] Ubuntu 18.04", ubuntu1804ContainerDisk, "u1.small",
				map[string]string{"amd64": "ubuntu"},
				[]testFn{expectSSHToRunCommandOnLinux("ubuntu")}),
			Entry("[test_id:10742] Ubuntu 20.04", ubuntu2004ContainerDisk, "u1.small",
				map[string]string{"amd64": "ubuntu", "arm64": "ubuntu", "s390x": "ubuntu"},
				[]testFn{expectSSHToRunCommandOnLinux("ubuntu")}),
			Entry("[test_id:10743] Ubuntu 22.04", ubuntu2204ContainerDisk, "u1.small",
				map[string]string{"amd64": "ubuntu", "arm64": "ubuntu", "s390x": "ubuntu"},
				[]testFn{expectSSHToRunCommandOnLinux("ubuntu")}),
			Entry("[test_id:TODO] Ubuntu 24.04", ubuntu2404ContainerDisk, "u1.small",
				map[string]string{"amd64": "ubuntu", "arm64": "ubuntu", "s390x": "ubuntu"},
				[]testFn{expectSSHToRunCommandOnLinux("ubuntu")}),
			Entry("[test_id:TODO] OpenSUSE Tumbleweed", openSUSETumbleweedContainerDisk, "u1.small",
				map[string]string{"amd64": "opensuse.tumbleweed"},
				[]testFn{expectGuestAgentToBeConnected, expectSSHToRunCommandOnLinux("opensuse")}),
			Entry("[test_id:TODO] OpenSUSE Leap 15", openSUSELeap15ContainerDisk, "u1.small",
				map[string]string{"amd64": "opensuse.leap", "arm64": "opensuse.leap"},
				[]testFn{expectGuestAgentToBeConnected, expectSSHToRunCommandOnLinux("opensuse")}),
			Entry("[test_id:TODO] Oracle Linux 8 amd64", ol8ContainerDisk, "u1.2xmedium",
				map[string]string{"amd64": "oraclelinux"},
				[]testFn{expectGuestAgentToBeConnected, expectSSHToRunCommandOnLinux("cloud-user")}),
			Entry("[test_id:TODO] Oracle Linux 8 arm64", ol8ContainerDisk, "u1.2xmedium",
				map[string]string{"arm64": "oraclelinux"},
				[]testFn{expectGuestAgentToBeConnected, expectSSHToRunCommandOnLinux("opc")}),
			Entry("[test_id:TODO] Oracle Linux 9 amd64", ol9ContainerDisk, "u1.2xmedium",
				map[string]string{"amd64": "oraclelinux"},
				[]testFn{expectGuestAgentToBeConnected, expectSSHToRunCommandOnLinux("cloud-user")}),
			Entry("[test_id:TODO] Oracle Linux 9 arm64", ol9ContainerDisk, "u1.2xmedium",
				map[string]string{"arm64": "oraclelinux"},
				[]testFn{expectGuestAgentToBeConnected, expectSSHToRunCommandOnLinux("opc")}),
			Entry("[test_id:TODO] SLES 15", sles15ContainerDisk, "u1.small", map[string]string{"amd64": "sles"},
				[]testFn{expectGuestAgentToBeConnected, expectSSHToRunCommandOnLinux("sles")}),
			Entry("[test_id:TODO] Debian 11", debian11ContainerDisk, "u1.small",
				map[string]string{"amd64": "debian", "arm64": "debian"}, []testFn{expectSSHToRunCommandOnLinux("debian")}),
			Entry("[test_id:TODO] Debian 12", debian12ContainerDisk, "u1.small",
				map[string]string{"amd64": "debian", "arm64": "debian"}, []testFn{expectSSHToRunCommandOnLinux("debian")}),
		)

		DescribeTable("a Windows guest with", func(containerDisk, instancetype string, preferences map[string]string, testFns []testFn) {
			preference, hasArch := preferences[preferenceArch]
			if !hasArch {
				Skip(fmt.Sprintf("skipping as no preference provided for arch %s", preferenceArch))
			}
			vm = randomVM(&v1.InstancetypeMatcher{Name: instancetype}, &v1.PreferenceMatcher{Name: preference}, v1.RunStrategyAlways)
			addContainerDisk(vm, containerDisk)
			vm, err = virtClient.VirtualMachine(testNamespace).Create(context.Background(), vm, metav1.CreateOptions{})
			Expect(err).ToNot(HaveOccurred())
			expectVMToBeReady(virtClient, vm.Name, windowsReadyTimeout)
			for _, testFn := range testFns {
				testFn(virtClient, vm.Name)
			}
		},
			Entry("[test_id:10739] Validation OS", validationOsContainerDisk, "u1.2xmedium",
				map[string]string{"amd64": "windows.11"}, []testFn{expectSSHToRunCommandOnWindows}),
			Entry("[test_id:????] Windows 10", windows10ContainerDisk, "u1.2xmedium",
				map[string]string{"amd64": "windows.10.virtio"},
				[]testFn{expectGuestAgentToBeConnected, expectSSHToRunCommandOnWindows}),
			Entry("[test_id:????] Windows 11", windows11ContainerDisk, "u1.2xmedium",
				map[string]string{"amd64": "windows.11.virtio"},
				[]testFn{expectGuestAgentToBeConnected, expectSSHToRunCommandOnWindows}),
			Entry("[test_id:????] Windows Server 2016", windows2k16ContainerDisk, "u1.2xmedium",
				map[string]string{"amd64": "windows.2k16.virtio"},
				[]testFn{expectGuestAgentToBeConnected, expectSSHToRunCommandOnWindows}),
			Entry("[test_id:????] Windows Server 2019", windows2k19ContainerDisk, "u1.2xmedium",
				map[string]string{"amd64": "windows.2k19.virtio"},
				[]testFn{expectGuestAgentToBeConnected, expectSSHToRunCommandOnWindows}),
			Entry("[test_id:????] Windows Server 2022", windows2k22ContainerDisk, "u1.2xmedium",
				map[string]string{"amd64": "windows.2k22.virtio"},
				[]testFn{expectGuestAgentToBeConnected, expectSSHToRunCommandOnWindows}),
			Entry("[test_id:????] Windows Server 2025", windows2k25ContainerDisk, "u1.2xmedium",
				map[string]string{"amd64": "windows.2k25.virtio"},
				[]testFn{expectGuestAgentToBeConnected, expectSSHToRunCommandOnWindows}),
			Entry("[test_id:????] Windows XP", windowsXpContainerDisk, "u1.nano",
				map[string]string{"amd64": "windows.xp"},
				[]testFn{expectGuestAgentToBeConnected, expectSSHToRunCommandOnWindows}),
		)
	})

	Context("with deprecated fields used", func() {
		It("[test_id:11951] rejects them in instance type", func() {
			deprecatedFields := getDeprecatedFieldFromCrd(virtClient, "virtualmachineinstancetypes.instancetype.kubevirt.io")

			clusterInstancetype, err := virtClient.VirtualMachineClusterInstancetype().List(context.Background(), metav1.ListOptions{})
			Expect(err).ToNot(HaveOccurred())
			Expect(clusterInstancetype.Items).ToNot(BeEmpty())

			for _, instancetype := range clusterInstancetype.Items {
				for _, deprecatedFieldPath := range deprecatedFields {
					Expect(isFieldUsed(deprecatedFieldPath, serialize(instancetype))).To(BeFalse())
				}
			}
		})

		It("[test_id:11952] rejects them in preference", func() {
			deprecatedFields := getDeprecatedFieldFromCrd(virtClient, "virtualmachineclusterpreferences.instancetype.kubevirt.io")
			clusterPreferences, err := virtClient.VirtualMachineClusterPreference().List(context.Background(), metav1.ListOptions{})
			Expect(err).ToNot(HaveOccurred())
			Expect(clusterPreferences.Items).ToNot(BeEmpty())

			for _, preference := range clusterPreferences.Items {
				for _, deprecatedFieldPath := range deprecatedFields {
					Expect(isFieldUsed(deprecatedFieldPath, serialize(preference))).To(BeFalse())
				}
			}
		})
	})

})

func getClusterInstancetypes(virtClient kubecli.KubevirtClient) []instancetypev1beta1.VirtualMachineClusterInstancetype {
	clusterInstancetypes, err := virtClient.VirtualMachineClusterInstancetype().List(context.Background(), metav1.ListOptions{})
	Expect(err).ToNot(HaveOccurred())
	Expect(clusterInstancetypes.Items).ToNot(BeEmpty())
	return clusterInstancetypes.Items
}

func createInstancetype(cpu uint32, name, memory string) {
	instancetype := &instancetypev1beta1.VirtualMachineInstancetype{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
		Spec: instancetypev1beta1.VirtualMachineInstancetypeSpec{
			CPU: instancetypev1beta1.CPUInstancetype{
				Guest: cpu,
			},
			Memory: instancetypev1beta1.MemoryInstancetype{
				Guest: resource.MustParse(memory),
			},
		},
	}
	_, err := virtClient.VirtualMachineInstancetype(testNamespace).Create(context.Background(), instancetype, metav1.CreateOptions{})
	Expect(err).ToNot(HaveOccurred())
}

func randomVM(
	instancetype *v1.InstancetypeMatcher,
	preference *v1.PreferenceMatcher,
	runStrategy v1.VirtualMachineRunStrategy,
) *v1.VirtualMachine {
	name := "test-vm-" + k8srand.String(5)
	return &v1.VirtualMachine{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: testNamespace,
		},
		Spec: v1.VirtualMachineSpec{
			RunStrategy:  &runStrategy,
			Instancetype: instancetype,
			Preference:   preference,
			Template: &v1.VirtualMachineInstanceTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{"name": name},
				},
				Spec: v1.VirtualMachineInstanceSpec{
					TerminationGracePeriodSeconds: ptr.To[int64](0),
					Domain: v1.DomainSpec{
						Devices: v1.Devices{
							Interfaces: []v1.Interface{{
								Name: "default",
								InterfaceBindingMethod: v1.InterfaceBindingMethod{
									Masquerade: &v1.InterfaceMasquerade{},
								},
							}},
						},
					},
					Networks: []v1.Network{{
						Name: "default",
						NetworkSource: v1.NetworkSource{
							Pod: &v1.PodNetwork{},
						},
					}},
				},
			},
		},
	}
}

func addContainerDisk(vm *v1.VirtualMachine, image string) {
	if vm.Spec.Template == nil {
		vm.Spec.Template = &v1.VirtualMachineInstanceTemplateSpec{}
	}

	vol := v1.Volume{
		Name: "containerdisk-" + k8srand.String(5),
		VolumeSource: v1.VolumeSource{
			ContainerDisk: &v1.ContainerDiskSource{
				Image: image,
			},
		},
	}
	vm.Spec.Template.Spec.Volumes = append(vm.Spec.Template.Spec.Volumes, vol)
}

func addCloudInitWithAuthorizedKey(vm *v1.VirtualMachine, privKey ed25519.PrivateKey) {
	pubKey, err := ssh.NewPublicKey(privKey.Public())
	Expect(err).ToNot(HaveOccurred())

	if vm.Spec.Template == nil {
		vm.Spec.Template = &v1.VirtualMachineInstanceTemplateSpec{}
	}

	vol := v1.Volume{
		Name: "cloudinitdisk",
		VolumeSource: v1.VolumeSource{
			CloudInitNoCloud: &v1.CloudInitNoCloudSource{
				UserData: fmt.Sprintf("#cloud-config\nssh_authorized_keys:\n  - %s", string(ssh.MarshalAuthorizedKey(pubKey))),
			},
		},
	}
	vm.Spec.Template.Spec.Volumes = append(vm.Spec.Template.Spec.Volumes, vol)
}

func expectVMToBeReady(virtClient kubecli.KubevirtClient, vmName string, vmReadyTimeout time.Duration) {
	Eventually(func(g Gomega) {
		vm, err := virtClient.VirtualMachine(testNamespace).Get(context.Background(), vmName, metav1.GetOptions{})
		g.Expect(err).ToNot(HaveOccurred())
		g.Expect(vm.Status.Ready).To(BeTrue())
	}, vmReadyTimeout, 10*time.Second).Should(Succeed())
}

func expectGuestAgentToBeConnected(virtClient kubecli.KubevirtClient, vmName string) {
	Eventually(func(g Gomega) {
		_, err := virtClient.VirtualMachineInstance(testNamespace).GuestOsInfo(context.Background(), vmName)
		g.Expect(err).ToNot(HaveOccurred())
	}, defaultVMReadyTimeout, 10*time.Second).Should(Succeed())
}

func expectSSHToRunCommandOnWindows(virtClient kubecli.KubevirtClient, vmName string) {
	expectSSHToRunCommand(virtClient, vmName, "Administrator", ssh.Password("Administrator"))
}

func expectSSHToRunCommand(virtClient kubecli.KubevirtClient, vmName, username string, authMethod ssh.AuthMethod) {
	// Test SSH while deliberately ignoring insecure host keys
	config := &ssh.ClientConfig{
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), //nolint:gosec
		User:            username,
		Auth: []ssh.AuthMethod{
			authMethod,
		},
	}

	Eventually(func(g Gomega) {
		tunnel, err := virtClient.VirtualMachineInstance(testNamespace).PortForward(vmName, sshPort, "tcp")
		g.Expect(err).ToNot(HaveOccurred())

		conn, chans, reqs, err := ssh.NewClientConn(tunnel.AsConn(), fmt.Sprintf("vm/%s.%s:%d", vmName, testNamespace, sshPort), config)
		g.Expect(err).ToNot(HaveOccurred())

		session, err := ssh.NewClient(conn, chans, reqs).NewSession()
		g.Expect(err).ToNot(HaveOccurred())

		err = session.Run("echo hello")
		g.Expect(err).ToNot(HaveOccurred())
	}, defaultVMReadyTimeout, 10*time.Second).Should(Succeed())
}

func getDeprecatedFieldFromCrd(virtClient kubecli.KubevirtClient, apiResourceIdentifier string) []string {
	crd, err := virtClient.ExtensionsClient().ApiextensionsV1().CustomResourceDefinitions().Get(
		context.Background(), apiResourceIdentifier, metav1.GetOptions{})
	Expect(err).ToNot(HaveOccurred())

	var nonDeprecatedVersion extv1.CustomResourceDefinitionVersion
	for _, version := range crd.Spec.Versions {
		if !version.Deprecated {
			nonDeprecatedVersion = version
		}
	}
	Expect(nonDeprecatedVersion.Name).ToNot(BeEmpty())
	schema := nonDeprecatedVersion.Schema.OpenAPIV3Schema

	return findDeprecatedFields(schema)
}

func findDeprecatedFields(schema *extv1.JSONSchemaProps) []string {
	var deprecatedFields []string

	var collectDeprecatedFields func(schema *extv1.JSONSchemaProps, path string)
	collectDeprecatedFields = func(schema *extv1.JSONSchemaProps, path string) {
		for propName, propSchema := range schema.Properties {
			fullPath := propName
			if path != "" {
				fullPath = path + "." + propName
			}
			if strings.Contains(strings.ToLower(propSchema.Description), "deprecated") {
				deprecatedFields = append(deprecatedFields, fullPath)
			}
			collectDeprecatedFields(&propSchema, fullPath)
		}
	}

	collectDeprecatedFields(schema, "")
	for i, field := range deprecatedFields {
		deprecatedFields[i] = "." + field
	}

	return deprecatedFields
}

func isFieldUsed(path string, data any) bool {
	jp := jsonpath.New("isFieldUsed")
	Expect(jp.Parse(fmt.Sprintf("{ %s }", path))).To(Succeed())
	var buf bytes.Buffer
	err := jp.Execute(&buf, data)
	if err != nil {
		return false
	}
	return err == nil
}

func serialize(data any) map[string]any {
	m, err := json.Marshal(data)
	Expect(err).ToNot(HaveOccurred())
	var obj map[string]any
	Expect(json.Unmarshal(m, &obj)).To(Succeed())
	return obj
}
