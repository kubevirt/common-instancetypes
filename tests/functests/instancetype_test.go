package tests

import (
	"context"
	"crypto/ed25519"
	"fmt"
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"

	"golang.org/x/crypto/ssh"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8srand "k8s.io/apimachinery/pkg/util/rand"
	"k8s.io/utils/ptr"
	v1 "kubevirt.io/api/core/v1"
	instancetypev1beta1 "kubevirt.io/api/instancetype/v1beta1"
	"kubevirt.io/client-go/kubecli"
)

const (
	vmReadyTimeout = 300 * time.Second
	sshPort        = 22
)

type testFn func(kubecli.KubevirtClient, string)

var _ = Describe("Common instance types func tests", func() {
	var (
		vm  *v1.VirtualMachine
		err error
	)

	AfterEach(func() {
		err = virtClient.VirtualMachine(testNamespace).Delete(context.Background(), vm.Name,
			&metav1.DeleteOptions{GracePeriodSeconds: ptr.To[int64](0)})
		if err != nil && !errors.IsNotFound(err) {
			Expect(err).ToNot(HaveOccurred())
		}
	})

	It("[test_id:10735] VirtualMachine using an instancetype can be created", func() {
		for _, instancetype := range getClusterInstancetypes(virtClient) {
			vm = randomVM(&v1.InstancetypeMatcher{Name: instancetype.Name}, nil, false)
			vm, err = virtClient.VirtualMachine(testNamespace).Create(context.Background(), vm)
			Expect(err).ToNot(HaveOccurred())
		}
	})

	Context("VirtualMachine using a preference with resource requirements", func() {
		It("[test_id:10736] is rejected if it does not provide enough resources", func() {
			createInstancetype(2, "tiny-instancetype", "64M")
			instanceTypeMatcher := v1.InstancetypeMatcher{
				Name: "tiny-instancetype",
				Kind: "VirtualMachineInstancetype",
			}
			for _, preference := range getClusterPreferences(virtClient) {
				vm = randomVM(&instanceTypeMatcher, &v1.PreferenceMatcher{Name: preference.Name}, false)
				_, err = virtClient.VirtualMachine(testNamespace).Create(context.Background(), vm)
				Expect(err).To(MatchError(
					fmt.Sprintf(
						"admission webhook \"virtualmachine-validator.kubevirt.io\" denied the request: "+
							"failure checking preference requirements: "+
							"insufficient Memory resources of 64M provided by instance type, preference requires %s",
						preference.Spec.Requirements.Memory.Guest.String(),
					)))
			}
		})

		It("[test_id:10737] can be created when enough resources are provided", func() {
			for _, preference := range getClusterPreferences(virtClient) {
				vm = randomVM(&v1.InstancetypeMatcher{Name: "u1.large"}, &v1.PreferenceMatcher{Name: preference.Name}, false)
				vm, err = virtClient.VirtualMachine(testNamespace).Create(context.Background(), vm)
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

		DescribeTable("a Linux guest with", func(containerDisk, preference string, testFns []testFn) {
			vm = randomVM(&v1.InstancetypeMatcher{Name: "u1.small"}, &v1.PreferenceMatcher{Name: preference}, true)
			addContainerDisk(vm, containerDisk)
			addCloudInitWithAuthorizedKey(vm, privKey)
			vm, err = virtClient.VirtualMachine(testNamespace).Create(context.Background(), vm)
			Expect(err).ToNot(HaveOccurred())
			expectVMToBeReady(virtClient, vm.Name)
			for _, testFn := range testFns {
				testFn(virtClient, vm.Name)
			}
		},
			Entry("[test_id:10738] Fedora", fedoraContainerDisk, "fedora",
				[]testFn{expectGuestAgentToBeConnected, expectSSHToRunCommandOnLinux("fedora")}),
			Entry("[test_id:10740] CentOS 7", centos7ContainerDisk, "centos.7",
				[]testFn{expectGuestAgentToBeConnected, expectSSHToRunCommandOnLinux("centos")}),
			Entry("[test_id:10744] CentOS Stream 8", centosStream8ContainerDisk, "centos.stream8",
				[]testFn{expectGuestAgentToBeConnected, expectSSHToRunCommandOnLinux("centos")}),
			Entry("[test_id:10745] CentOS Stream 9", centosStream9ContainerDisk, "centos.stream9",
				[]testFn{expectGuestAgentToBeConnected, expectSSHToRunCommandOnLinux("cloud-user")}),
			Entry("[test_id:10741] Ubuntu 18.04", ubuntu1804ContainerDisk, "ubuntu",
				[]testFn{expectSSHToRunCommandOnLinux("ubuntu")}),
			Entry("[test_id:10742] Ubuntu 20.04", ubuntu2004ContainerDisk, "ubuntu",
				[]testFn{expectSSHToRunCommandOnLinux("ubuntu")}),
			Entry("[test_id:10743] Ubuntu 22.04", ubuntu2204ContainerDisk, "ubuntu",
				[]testFn{expectSSHToRunCommandOnLinux("ubuntu")}),
			Entry("[test_id:TODO] OpenSUSE Tumbleweed", openSUSETumbleweedContainerDisk, "opensuse.tumbleweed",
				[]testFn{expectGuestAgentToBeConnected, expectSSHToRunCommandOnLinux("opensuse")}),
			Entry("[test_id:TODO] OpenSUSE Leap 15.6", openSUSELeapContainerDisk, "opensuse.leap",
				[]testFn{expectGuestAgentToBeConnected, expectSSHToRunCommandOnLinux("opensuse")}),
		)

		DescribeTable("a Windows guest with", func(containerDisk, preference string, testFns []testFn) {
			vm = randomVM(&v1.InstancetypeMatcher{Name: "u1.large"}, &v1.PreferenceMatcher{Name: preference}, true)
			addContainerDisk(vm, containerDisk)
			vm, err = virtClient.VirtualMachine(testNamespace).Create(context.Background(), vm)
			Expect(err).ToNot(HaveOccurred())
			expectVMToBeReady(virtClient, vm.Name)
			for _, testFn := range testFns {
				testFn(virtClient, vm.Name)
			}
		},
			Entry("[test_id:10739] Validation OS", validationOsContainerDisk, "windows.11",
				[]testFn{expectSSHToRunCommandOnWindows}),
			Entry("[test_id:????] Windows 10", windows10ContainerDisk, "windows.10.virtio",
				[]testFn{expectGuestAgentToBeConnected, expectSSHToRunCommandOnWindows}),
			Entry("[test_id:????] Windows 11", windows11ContainerDisk, "windows.11.virtio",
				[]testFn{expectGuestAgentToBeConnected, expectSSHToRunCommandOnWindows}),
			Entry("[test_id:????] Windows Server 2016", windows2k16ContainerDisk, "windows.2k16.virtio",
				[]testFn{expectGuestAgentToBeConnected}),
			Entry("[test_id:????] Windows Server 2019", windows2k19ContainerDisk, "windows.2k19.virtio",
				[]testFn{expectGuestAgentToBeConnected, expectSSHToRunCommandOnWindows}),
			Entry("[test_id:????] Windows Server 2022", windows2k22ContainerDisk, "windows.2k22.virtio",
				[]testFn{expectGuestAgentToBeConnected, expectSSHToRunCommandOnWindows}),
		)
	})
})

func getClusterInstancetypes(virtClient kubecli.KubevirtClient) []instancetypev1beta1.VirtualMachineClusterInstancetype {
	clusterInstancetypes, err := virtClient.VirtualMachineClusterInstancetype().List(context.Background(), metav1.ListOptions{})
	Expect(err).ToNot(HaveOccurred())
	Expect(clusterInstancetypes.Items).ToNot(BeEmpty())
	return clusterInstancetypes.Items
}

func getClusterPreferences(virtClient kubecli.KubevirtClient) []instancetypev1beta1.VirtualMachineClusterPreference {
	clusterPreferences, err := virtClient.VirtualMachineClusterPreference().List(context.Background(), metav1.ListOptions{})
	Expect(err).ToNot(HaveOccurred())
	Expect(clusterPreferences.Items).ToNot(BeEmpty())
	return clusterPreferences.Items
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

func randomVM(instancetype *v1.InstancetypeMatcher, preference *v1.PreferenceMatcher, running bool) *v1.VirtualMachine {
	name := "test-vm-" + k8srand.String(5)
	return &v1.VirtualMachine{
		ObjectMeta: metav1.ObjectMeta{
			Name:      name,
			Namespace: testNamespace,
		},
		Spec: v1.VirtualMachineSpec{
			Running:      ptr.To(running),
			Instancetype: instancetype,
			Preference:   preference,
			Template: &v1.VirtualMachineInstanceTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{"name": name},
				},
				Spec: v1.VirtualMachineInstanceSpec{
					TerminationGracePeriodSeconds: ptr.To[int64](0),
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

func expectVMToBeReady(virtClient kubecli.KubevirtClient, vmName string) {
	Eventually(func(g Gomega) {
		vm, err := virtClient.VirtualMachine(testNamespace).Get(context.Background(), vmName, &metav1.GetOptions{})
		g.Expect(err).ToNot(HaveOccurred())
		g.Expect(vm.Status.Ready).To(BeTrue())
	}, vmReadyTimeout, 10*time.Second).Should(Succeed())
}

func expectGuestAgentToBeConnected(virtClient kubecli.KubevirtClient, vmName string) {
	Eventually(func(g Gomega) {
		_, err := virtClient.VirtualMachineInstance(testNamespace).GuestOsInfo(context.Background(), vmName)
		g.Expect(err).ToNot(HaveOccurred())
	}, vmReadyTimeout, 10*time.Second).Should(Succeed())
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
	}, vmReadyTimeout, 10*time.Second).Should(Succeed())
}
