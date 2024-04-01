package host

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"
)

func probeDarwin(darwinProber prober) (*Host, error) {
	var host Host
	session, err := darwinProber.sshClient.NewSession()
	if err != nil {
		return nil, fmt.Errorf("failed to create session: %s", err)
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
		return nil, fmt.Errorf("failed to run: %s", err)
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
		host.OS = substrs[0]
		host.OSVersion = substrs[1]
	}
	if kernelVersionMatch != nil {
		substrs := strings.Split(kernelVersionMatch[1], " ")
		host.Kernel = substrs[0]
		host.KernelVersion = substrs[1]
	}
	if computerNameMatch != nil {
		host.Hostname = computerNameMatch[1]
	}
	if archMatch != nil {
		host.Arch = archMatch[1]
	}
	if memSizeMatch != nil {
		memSize := memSizeMatch[1]
		memSizeInt := 0
		fmt.Sscanf(memSize, "%d", &memSizeInt)
		host.MemoryTotal = fmt.Sprintf("%d GB", memSizeInt/1024/1024/1024)
	}
	if diskSizeMatch != nil {
		host.DiskTotal = diskSizeMatch[1]
	}

	// fmt.Println(host)
	return &host, nil
}
