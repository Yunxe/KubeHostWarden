package probe

import (
	"bytes"
	"context"
	"fmt"
	"kubehostwarden/types"
	"net"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/ssh"
)

func (ph *probeHelper) probe(ctx context.Context) error {
	portStr := fmt.Sprintf("%d", ph.sshInfo.Port)

	addr := net.JoinHostPort(ph.sshInfo.EndPoint, portStr)

	config := &ssh.ClientConfig{
		User: ph.sshInfo.User,
		Auth: []ssh.AuthMethod{
			ssh.Password(ph.sshInfo.Password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		Timeout:         5 * time.Second,
	}

	client, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		return fmt.Errorf("failed to dial: %s", err)
	}
	ph.sshClient = client
	ph.host = &types.Host{}

	switch ph.sshInfo.OSType {
	case "darwin":
		err := ph.probeDarwin(ctx)
		if err != nil {
			return err
		}
		ph.host.IPAddr = ph.sshInfo.EndPoint

	case "linux":
		err := ph.probeLinux(ctx)
		if err != nil {
			return err
		}
		ph.host.IPAddr = ph.sshInfo.EndPoint
	}
	uuid := uuid.New().String()
	ph.host.Id = uuid
	return nil
}

func (ph *probeHelper) probeDarwin(ctx context.Context) error {
	session, err := ph.sshClient.NewSession()
	if err != nil {
		return fmt.Errorf("failed to create session: %s", err)
	}
	defer session.Close()

	// run command
	var b bytes.Buffer
	session.Stdout = &b
	if err := session.Run(`
	system_profiler SPSoftwareDataType;\
	echo "arch: $(uname -m)";\
	sysctl hw.memsize;\
	echo "disksize: $(diskutil list | grep "disk0" |\
	 grep -v "disk0s" | awk 'NR==2{print$3$4}')"`); err != nil {
		return fmt.Errorf("failed to run: %s", err)
	}
	output := b.String()

	// regular expression
	reSystemVersion := regexp.MustCompile(`System Version: (.*)`)
	reKernelVersion := regexp.MustCompile(`Kernel Version: (.*)`)
	reComputerName := regexp.MustCompile(`Computer Name: (.*)`)
	reArch := regexp.MustCompile(`arch: (.*)`)
	reMemSize := regexp.MustCompile(`hw.memsize: (.*)`)
	reDiskSize := regexp.MustCompile(`disksize: (.*)`)

	// match
	systemVersionMatch := reSystemVersion.FindStringSubmatch(output)
	kernelVersionMatch := reKernelVersion.FindStringSubmatch(output)
	computerNameMatch := reComputerName.FindStringSubmatch(output)
	archMatch := reArch.FindStringSubmatch(output)
	memSizeMatch := reMemSize.FindStringSubmatch(output)
	diskSizeMatch := reDiskSize.FindStringSubmatch(output)

	if systemVersionMatch != nil {
		substrs := strings.Split(systemVersionMatch[1], " ")
		ph.host.OS = substrs[0]
		ph.host.OSVersion = substrs[1]
	}
	if kernelVersionMatch != nil {
		substrs := strings.Split(kernelVersionMatch[1], " ")
		ph.host.Kernel = substrs[0]
		ph.host.KernelVersion = substrs[1]
	}
	if computerNameMatch != nil {
		ph.host.Hostname = computerNameMatch[1]
	}
	if archMatch != nil {
		ph.host.Arch = archMatch[1]
	}
	if memSizeMatch != nil {
		memSize := memSizeMatch[1]
		memSizeInt := 0
		fmt.Sscanf(memSize, "%d", &memSizeInt)
		ph.host.MemoryTotal = fmt.Sprintf("%d GB", memSizeInt/1024/1024/1024)
	}
	if diskSizeMatch != nil {
		ph.host.DiskTotal = diskSizeMatch[1]
	}

	// fmt.Println(host)
	return nil
}

func (ph *probeHelper) probeLinux(ctx context.Context) error {
	session, err := ph.sshClient.NewSession()
	if err != nil {
		return fmt.Errorf("failed to create session: %s", err)
	}
	defer session.Close()

	// run command
	var b bytes.Buffer
	session.Stdout = &b
	if err := session.Run(`
	lsb_release -a;\
	echo "kernel: $(uname -s)";\
	echo "kernel version: $(uname -r)";\
	echo "hostname: $(hostname)";\
	echo "arch: $(uname -m)";\
	echo "Mem: $(free -h | grep Mem | awk '{print $2}')";\
	echo "disk: $(df -h | awk 'NR>1 {print $2}' | grep -E '^[0-9\.]+[GM]$' | sort -h | tail -1)"`); err != nil {
		return fmt.Errorf("failed to run: %s", err)
	}
	output := b.String()

	// regular expression
	reDistributorID := regexp.MustCompile(`Distributor ID:\s+(.*)`)
	reRelease := regexp.MustCompile(`Release:\s+(.*)`)
	reKernel := regexp.MustCompile(`kernel: (.*)`)
	reKernelVersion := regexp.MustCompile(`kernel version: (.*)`)
	reHostname := regexp.MustCompile(`hostname: (.*)`)
	reArch := regexp.MustCompile(`arch: (.*)`)
	reMemSize := regexp.MustCompile(`Mem:\s+(.*)`)
	reDiskSize := regexp.MustCompile(`disk:\s+(.*)`)

	// match
	distributorIDMatch := reDistributorID.FindStringSubmatch(output)
	releaseMatch := reRelease.FindStringSubmatch(output)
	kernelMatch := reKernel.FindStringSubmatch(output)
	kernelVersionMatch := reKernelVersion.FindStringSubmatch(output)
	hostnameMatch := reHostname.FindStringSubmatch(output)
	archMatch := reArch.FindStringSubmatch(output)
	memSizeMatch := reMemSize.FindStringSubmatch(output)
	diskSizeMatch := reDiskSize.FindStringSubmatch(output)

	if distributorIDMatch != nil {
		ph.host.OS = distributorIDMatch[1]
	}
	if releaseMatch != nil {
		ph.host.OSVersion = releaseMatch[1]
	}
	if kernelMatch != nil {
		ph.host.Kernel = kernelMatch[1]
	}
	if kernelVersionMatch != nil {
		ph.host.KernelVersion = kernelVersionMatch[1]
	}
	if hostnameMatch != nil {
		ph.host.Hostname = hostnameMatch[1]
	}
	if archMatch != nil {
		ph.host.Arch = archMatch[1]
	}
	if memSizeMatch != nil {
		ph.host.MemoryTotal = memSizeMatch[1]
	}
	if diskSizeMatch != nil {
		ph.host.DiskTotal = diskSizeMatch[1]
	}

	// fmt.Println(ph.host)
	return nil
}
