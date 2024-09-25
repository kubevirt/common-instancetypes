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
	defaultRHEL9ContainerDisk              = "registry:5000/rhel-guest-image:9"
	defaultRHEL10ContainerDisk             = "registry:5000/rhel-guest-image:10"
	defaultUbuntu1804ContainerDisk         = "quay.io/containerdisks/ubuntu:18.04"
	defaultUbuntu2004ContainerDisk         = "quay.io/containerdisks/ubuntu:20.04"
	defaultUbuntu2204ContainerDisk         = "quay.io/containerdisks/ubuntu:22.04"
	defaultUbuntu2404ContainerDisk         = "quay.io/containerdisks/ubuntu:24.04"
	defaultOpenSUSETumbleweedContainerDisk = "quay.io/containerdisks/opensuse-tumbleweed:1.0.0"
	defaultOpenSUSELeap15ContainerDisk     = "quay.io/containerdisks/opensuse-leap:15.6"
	defaultSLES15ContainerDisk             = "registry:5000/sles15-container-disk:latest"
	defaultValidationOsContainerDisk       = "registry:5000/validation-os-container-disk:latest"
	defaultWindows10ContainerDisk          = "registry:5000/windows10-container-disk:latest"
	defaultWindows11ContainerDisk          = "registry:5000/windows11-container-disk:latest"
	defaultWindows2k16ContainerDisk        = "registry:5000/windows2k16-container-disk:latest"
	defaultWindows2k19ContainerDisk        = "registry:5000/windows2k19-container-disk:latest"
	defaultWindows2k22ContainerDisk        = "registry:5000/windows2k22-container-disk:latest"
	defaultWindows2k25ContainerDisk        = "registry:5000/windows2k25-container-disk:latest"

	defaultVMReadyTimeout = 300 * time.Second
)

var (
	afterSuiteReporters []Reporter
	virtClient          kubecli.KubevirtClient

	fedoraContainerDisk             string
	centos7ContainerDisk            string
	centosStream8ContainerDisk      string
	centosStream9ContainerDisk      string
	rhel9ContainerDisk              string
	rhel10ContainerDisk             string
	ubuntu1804ContainerDisk         string
	ubuntu2004ContainerDisk         string
	ubuntu2204ContainerDisk         string
	ubuntu2404ContainerDisk         string
	validationOsContainerDisk       string
	windows10ContainerDisk          string
	windows11ContainerDisk          string
	windows2k16ContainerDisk        string
	windows2k19ContainerDisk        string
	windows2k22ContainerDisk        string
	windows2k25ContainerDisk        string
	openSUSETumbleweedContainerDisk string
	openSUSELeap15ContainerDisk     string
	sles15ContainerDisk             string

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
	flag.StringVar(&sles15ContainerDisk, "sles15-container-disk",
		defaultSLES15ContainerDisk, "SLES 15 container disk used by functional tests")
	flag.StringVar(&validationOsContainerDisk, "validation-os-container-disk",
		defaultValidationOsContainerDisk, "Validation OS container disk used by functional tests")
	flag.StringVar(&windows10ContainerDisk, "windows-10-container-disk",
		defaultWindows10ContainerDisk, "Windows 10 container disk used by functional tests")
	flag.StringVar(&windows11ContainerDisk, "windows-11-container-disk",
		defaultWindows11ContainerDisk, "Windows 11 container disk used by functional tests")
	flag.StringVar(&windows2k16ContainerDisk, "windows-2k16-container-disk",
		defaultWindows2k16ContainerDisk, "Windows Server 2016 container disk used by functional tests")
	flag.StringVar(&windows2k19ContainerDisk, "windows-2k19-container-disk",
		defaultWindows2k19ContainerDisk, "Windows Server 2019 container disk used by functional tests")
	flag.StringVar(&windows2k22ContainerDisk, "windows-2k22-container-disk",
		defaultWindows2k22ContainerDisk, "Windows Server 2022 container disk used by functional tests")
	flag.StringVar(&windows2k25ContainerDisk, "windows-2k25-container-disk",
		defaultWindows2k25ContainerDisk, "Windows Server 2025 container disk used by functional tests")
	flag.DurationVar(&windowsReadyTimeout, "windows-ready-timeout",
		defaultVMReadyTimeout, "Duration after Windows VM will timeout")
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

	namespaceObj := &core.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: testNamespace,
		},
	}
	_, err = virtClient.CoreV1().Namespaces().Create(context.TODO(), namespaceObj, metav1.CreateOptions{})
	Expect(err).ToNot(HaveOccurred())

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
