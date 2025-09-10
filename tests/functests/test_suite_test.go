package tests

import (
	"context"
	"flag"
	"os"
	"testing"
	"time"

	. "github.com/onsi/ginkgo/v2"
	ginkgo_reporters "github.com/onsi/ginkgo/v2/reporters"
	. "github.com/onsi/gomega"
	qe_reporters "kubevirt.io/qe-tools/pkg/ginkgo-reporters"

	core "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"kubevirt.io/client-go/kubecli"
)

const (
	testNamespace = "common-instancetype-functest"

	defaultFedoraContainerDisk             = "quay.io/containerdisks/fedora:latest"
	defaultCentos7ContainerDisk            = "quay.io/containerdisks/centos:7-2009"
	defaultCentosStream8ContainerDisk      = "quay.io/containerdisks/centos-stream:8"
	defaultCentosStream9ContainerDisk      = "quay.io/containerdisks/centos-stream:9"
	defaultCentosStream10ContainerDisk     = "quay.io/containerdisks/centos-stream:10"
	defaultRHEL6ContainerDisk              = "registry:5000/rhel-guest-image:6"
	defaultRHEL7ContainerDisk              = "registry:5000/rhel-guest-image:7"
	defaultRHEL8ContainerDisk              = "registry:5000/rhel-guest-image:8"
	defaultRHEL9ContainerDisk              = "registry:5000/rhel-guest-image:9"
	defaultRHEL10ContainerDisk             = "registry:5000/rhel-guest-image:10"
	defaultUbuntu1804ContainerDisk         = "quay.io/containerdisks/ubuntu:18.04"
	defaultUbuntu2004ContainerDisk         = "quay.io/containerdisks/ubuntu:20.04"
	defaultUbuntu2204ContainerDisk         = "quay.io/containerdisks/ubuntu:22.04"
	defaultUbuntu2404ContainerDisk         = "quay.io/containerdisks/ubuntu:24.04"
	defaultOpenSUSETumbleweedContainerDisk = "quay.io/containerdisks/opensuse-tumbleweed:1.0.0"
	defaultOpenSUSELeap15ContainerDisk     = "quay.io/containerdisks/opensuse-leap:15.6"
	defaultDebian11ContainerDisk           = "quay.io/containerdisks/debian:11"
	defaultDebian12ContainerDisk           = "quay.io/containerdisks/debian:12"
	defaultDebian13ContainerDisk           = "quay.io/containerdisks/debian:13"
	defaultOL8ContainerDisk                = "registry:5000/oraclelinux:8.10"
	defaultOL9ContainerDisk                = "registry:5000/oraclelinux:9.5"
	defaultSLES15SP6ContainerDisk          = "registry:5000/sles15sp6-container-disk:latest"
	defaultSLES15SP7ContainerDisk          = "registry:5000/sles15sp7-container-disk:latest"
	defaultValidationOsContainerDisk       = "registry:5000/validation-os-container-disk:latest"
	defaultWindows7ContainerDisk           = "registry:5000/windows7-container-disk:latest"
	defaultWindows10ContainerDisk          = "registry:5000/windows10-container-disk:latest"
	defaultWindows11ContainerDisk          = "registry:5000/windows11-container-disk:latest"
	defaultWindows2k3ContainerDisk         = "registry:5000/windows2k3-container-disk:latest"
	defaultWindows2k8I386ContainerDisk     = "registry:5000/windows2k8-container-disk:i386"
	defaultWindows2k8Amd64ContainerDisk    = "registry:5000/windows2k8-container-disk:amd64"
	defaultWindows2k8R2ContainerDisk       = "registry:5000/windows2k8r2-container-disk:latest"
	defaultWindows2k12ContainerDisk        = "registry:5000/windows2k12-container-disk:latest"
	defaultWindows2k12R2ContainerDisk      = "registry:5000/windows2k12r2-container-disk:latest"
	defaultWindows2k16ContainerDisk        = "registry:5000/windows2k16-container-disk:latest"
	defaultWindows2k19ContainerDisk        = "registry:5000/windows2k19-container-disk:latest"
	defaultWindows2k22ContainerDisk        = "registry:5000/windows2k22-container-disk:latest"
	defaultWindows2k25ContainerDisk        = "registry:5000/windows2k25-container-disk:latest"
	defaultWindowsXpContainerDisk          = "registry:5000/windowsxp-container-disk:latest"

	defaultVMReadyTimeout = 300 * time.Second
)

var (
	afterSuiteReporters []Reporter
	virtClient          kubecli.KubevirtClient

	fedoraContainerDisk             string
	centos7ContainerDisk            string
	centosStream8ContainerDisk      string
	centosStream9ContainerDisk      string
	centosStream10ContainerDisk     string
	rhel6ContainerDisk              string
	rhel7ContainerDisk              string
	rhel8ContainerDisk              string
	rhel9ContainerDisk              string
	rhel10ContainerDisk             string
	ubuntu1804ContainerDisk         string
	ubuntu2004ContainerDisk         string
	ubuntu2204ContainerDisk         string
	ubuntu2404ContainerDisk         string
	validationOsContainerDisk       string
	windows7ContainerDisk           string
	windows10ContainerDisk          string
	windows11ContainerDisk          string
	windows2k3ContainerDisk         string
	windows2k8I386ContainerDisk     string
	windows2k8Amd64ContainerDisk    string
	windows2k8R2ContainerDisk       string
	windows2k12ContainerDisk        string
	windows2k12R2ContainerDisk      string
	windows2k16ContainerDisk        string
	windows2k19ContainerDisk        string
	windows2k22ContainerDisk        string
	windows2k25ContainerDisk        string
	windowsXpContainerDisk          string
	openSUSETumbleweedContainerDisk string
	openSUSELeap15ContainerDisk     string
	sles15SP6ContainerDisk          string
	sles15SP7ContainerDisk          string
	debian11ContainerDisk           string
	debian12ContainerDisk           string
	debian13ContainerDisk           string
	ol8ContainerDisk                string
	ol9ContainerDisk                string

	preferenceArch      string
	windowsReadyTimeout time.Duration
)

//nolint:gochecknoinits
func init() {
	kubecli.Init()

	flag.StringVar(&fedoraContainerDisk, "fedora-container-disk",
		defaultFedoraContainerDisk, "Fedora container disk used by functional tests")
	flag.StringVar(&centos7ContainerDisk, "centos-7-container-disk",
		defaultCentos7ContainerDisk, "CentOS 7 container disk used by functional tests")
	flag.StringVar(&centosStream8ContainerDisk, "centos-stream-8-container-disk",
		defaultCentosStream8ContainerDisk, "CentOS Stream 8 container disk used by functional tests")
	flag.StringVar(&centosStream9ContainerDisk, "centos-stream-9-container-disk",
		defaultCentosStream9ContainerDisk, "CentOS Stream 9 container disk used by functional tests")
	flag.StringVar(&centosStream10ContainerDisk, "centos-stream-10-container-disk",
		defaultCentosStream10ContainerDisk, "CentOS Stream 10 container disk used by functional tests")
	flag.StringVar(&rhel6ContainerDisk, "rhel-6-container-disk",
		defaultRHEL6ContainerDisk, "RHEL 6 container disk used by functional tests")
	flag.StringVar(&rhel7ContainerDisk, "rhel-7-container-disk",
		defaultRHEL7ContainerDisk, "RHEL 7 container disk used by functional tests")
	flag.StringVar(&rhel8ContainerDisk, "rhel-8-container-disk",
		defaultRHEL8ContainerDisk, "RHEL 8 container disk used by functional tests")
	flag.StringVar(&rhel9ContainerDisk, "rhel-9-container-disk",
		defaultRHEL9ContainerDisk, "RHEL 9 container disk used by functional tests")
	flag.StringVar(&rhel10ContainerDisk, "rhel-10-container-disk",
		defaultRHEL10ContainerDisk, "RHEL 10 container disk used by functional tests")
	flag.StringVar(&ubuntu1804ContainerDisk, "ubuntu-1804-container-disk",
		defaultUbuntu1804ContainerDisk, "Ubuntu 18.04 container disk used by functional tests")
	flag.StringVar(&ubuntu2004ContainerDisk, "ubuntu-2004-container-disk",
		defaultUbuntu2004ContainerDisk, "Ubuntu 20.04 container disk used by functional tests")
	flag.StringVar(&ubuntu2204ContainerDisk, "ubuntu-2204-container-disk",
		defaultUbuntu2204ContainerDisk, "Ubuntu 22.04 container disk used by functional tests")
	flag.StringVar(&ubuntu2404ContainerDisk, "ubuntu-2404-container-disk",
		defaultUbuntu2404ContainerDisk, "Ubuntu 24.04 container disk used by functional tests")
	flag.StringVar(&openSUSETumbleweedContainerDisk, "opensuse-tumbleweed-container-disk",
		defaultOpenSUSETumbleweedContainerDisk, "OpenSUSE Tumbleweed container disk used by functional tests")
	flag.StringVar(&openSUSELeap15ContainerDisk, "opensuse-leap-container-disk",
		defaultOpenSUSELeap15ContainerDisk, "OpenSUSE Leap container disk used by functional tests")
	flag.StringVar(&debian11ContainerDisk, "debian-11-container-disk",
		defaultDebian11ContainerDisk, "Debian 11 container disk used by functional tests")
	flag.StringVar(&debian12ContainerDisk, "debian-12-container-disk",
		defaultDebian12ContainerDisk, "Debian 12 container disk used by functional tests")
	flag.StringVar(&debian13ContainerDisk, "debian-13-container-disk",
		defaultDebian13ContainerDisk, "Debian 13 container disk used by functional tests")
	flag.StringVar(&ol8ContainerDisk, "ol-8-container-disk",
		defaultOL8ContainerDisk, "Oracle Linux 8 container disk used by function tests")
	flag.StringVar(&ol9ContainerDisk, "ol-9-container-disk",
		defaultOL9ContainerDisk, "Oracle Linux 9 container disk used by function tests")
	flag.StringVar(&sles15SP6ContainerDisk, "sles15sp6-container-disk",
		defaultSLES15SP6ContainerDisk, "SLES 15 SP 6 container disk used by functional tests")
	flag.StringVar(&sles15SP7ContainerDisk, "sles15sp7-container-disk",
		defaultSLES15SP7ContainerDisk, "SLES 15 SP 7 container disk used by functional tests")
	flag.StringVar(&validationOsContainerDisk, "validation-os-container-disk",
		defaultValidationOsContainerDisk, "Validation OS container disk used by functional tests")
	flag.StringVar(&windows7ContainerDisk, "windows-7-container-disk",
		defaultWindows7ContainerDisk, "Windows 7 container disk used by functional tests")
	flag.StringVar(&windows10ContainerDisk, "windows-10-container-disk",
		defaultWindows10ContainerDisk, "Windows 10 container disk used by functional tests")
	flag.StringVar(&windows11ContainerDisk, "windows-11-container-disk",
		defaultWindows11ContainerDisk, "Windows 11 container disk used by functional tests")
	flag.StringVar(&windows2k3ContainerDisk, "windows-2k3-container-disk",
		defaultWindows2k3ContainerDisk, "Windows Server 2003 container disk used by functional tests")
	flag.StringVar(&windows2k8I386ContainerDisk, "windows-2k8-i386-container-disk",
		defaultWindows2k8I386ContainerDisk, "Windows Server 2008 i368 container disk used by functional tests")
	flag.StringVar(&windows2k8Amd64ContainerDisk, "windows-2k8-amd64-container-disk",
		defaultWindows2k8Amd64ContainerDisk, "Windows Server 2008 amd64 container disk used by functional tests")
	flag.StringVar(&windows2k8R2ContainerDisk, "windows-2k8r2-container-disk",
		defaultWindows2k8R2ContainerDisk, "Windows Server 2008 R2 container disk used by functional tests")
	flag.StringVar(&windows2k12ContainerDisk, "windows-2k12-container-disk",
		defaultWindows2k12ContainerDisk, "Windows Server 2012 container disk used by functional tests")
	flag.StringVar(&windows2k12R2ContainerDisk, "windows-2k12r2-container-disk",
		defaultWindows2k12R2ContainerDisk, "Windows Server 2012 R2 container disk used by functional tests")
	flag.StringVar(&windows2k16ContainerDisk, "windows-2k16-container-disk",
		defaultWindows2k16ContainerDisk, "Windows Server 2016 container disk used by functional tests")
	flag.StringVar(&windows2k19ContainerDisk, "windows-2k19-container-disk",
		defaultWindows2k19ContainerDisk, "Windows Server 2019 container disk used by functional tests")
	flag.StringVar(&windows2k22ContainerDisk, "windows-2k22-container-disk",
		defaultWindows2k22ContainerDisk, "Windows Server 2022 container disk used by functional tests")
	flag.StringVar(&windows2k25ContainerDisk, "windows-2k25-container-disk",
		defaultWindows2k25ContainerDisk, "Windows Server 2025 container disk used by functional tests")
	flag.StringVar(&windowsXpContainerDisk, "windows-xp-container-disk",
		defaultWindowsXpContainerDisk, "Windows XP container disk used by functional tests")
	flag.StringVar(&preferenceArch, "preference-arch", "", "Architecture to test preferences for")
	flag.DurationVar(&windowsReadyTimeout, "windows-ready-timeout",
		defaultVMReadyTimeout, "Duration after Windows VM will timeout")
}

func getClusterArch(virtClient kubecli.KubevirtClient) string {
	nodes, err := virtClient.CoreV1().Nodes().List(context.Background(), metav1.ListOptions{})
	Expect(err).ToNot(HaveOccurred())
	Expect(nodes.Items).ToNot(BeEmpty())
	return nodes.Items[0].Status.NodeInfo.Architecture
}

func checkDeployedResources() {
	virtualMachineClusterInstancetypes, err := virtClient.VirtualMachineClusterInstancetype().List(context.Background(), metav1.ListOptions{})
	Expect(err).ToNot(HaveOccurred())
	Expect(virtualMachineClusterInstancetypes.Items).ToNot(BeEmpty())

	virtualMachineClusterPreferences, err := virtClient.VirtualMachineClusterPreference().List(context.Background(), metav1.ListOptions{})
	Expect(err).ToNot(HaveOccurred())
	Expect(virtualMachineClusterPreferences.Items).ToNot(BeEmpty())
}

var _ = BeforeSuite(func() {
	var err error
	var config *rest.Config
	kubeconfigPath := os.Getenv(clientcmd.RecommendedConfigPathEnvVar)
	if kubeconfigPath != "" {
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	} else {
		config, err = kubecli.GetKubevirtClientConfig()
	}
	Expect(err).ToNot(HaveOccurred())

	virtClient, err = kubecli.GetKubevirtClientFromRESTConfig(config)
	Expect(err).ToNot(HaveOccurred())

	if preferenceArch == "" {
		preferenceArch = getClusterArch(virtClient)
	}

	namespaceObj := &core.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: testNamespace,
		},
	}
	_, err = virtClient.CoreV1().Namespaces().Create(context.TODO(), namespaceObj, metav1.CreateOptions{})
	Expect(err).ToNot(And(HaveOccurred(), Not(MatchError(errors.IsAlreadyExists, "errors.IsAlreadyExists"))))

	checkDeployedResources()
})

var _ = AfterSuite(func() {
	// Clean up namespaced resources
	err := virtClient.CoreV1().Namespaces().Delete(context.TODO(), testNamespace, metav1.DeleteOptions{})
	if err != nil && !errors.IsNotFound(err) {
		Expect(err).ToNot(HaveOccurred())
	}
})

var _ = ReportAfterSuite("TestFunctional", func(report Report) {
	for _, reporter := range afterSuiteReporters {
		ginkgo_reporters.ReportViaDeprecatedReporter(reporter, report) //nolint:staticcheck
	}
})

func TestFunctional(t *testing.T) {
	if qe_reporters.JunitOutput != "" {
		afterSuiteReporters = append(afterSuiteReporters, ginkgo_reporters.NewJUnitReporter(qe_reporters.JunitOutput))
	}
	if qe_reporters.Polarion.Run {
		afterSuiteReporters = append(afterSuiteReporters, &qe_reporters.Polarion)
	}

	RegisterFailHandler(Fail)
	RunSpecs(t, "Functional test suite")
}
