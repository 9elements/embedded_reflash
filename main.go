package main

import (
	"flag"
	"log"
	"os"
	"strings"

	"github.com/facebook/openbmc/tools/flashy/lib/flash/flashcp"
	"github.com/moby/sys/mountinfo"
	"golang.org/x/sys/unix"
)

func main() {
	dev := flag.String("device", "/dev/mtd0", "Flash device to write to")
	img := flag.String("image", "/tmp/flash.img", "Flash image to write")
	skipRemount := flag.Bool("skip-remount", false, "Skip remounting rootfs read-only")
	flag.Parse()

	if os.Geteuid() != 0 {
		log.Fatal("This program must be run as root!")
	}

	if !*skipRemount {
		log.Println("Blocking file system access to flash...")
		if _, err := os.Stat("/proc/sysrq-trigger"); os.IsNotExist(err) {
			log.Println("No SysRq support, trying to remount rootfs read-only instead...")

			minfo, err := mountinfo.GetMounts(RWFilter)
			if err != nil {
				log.Fatalf("Failed to get mount info: %v", err)
			}

			for _, i := range minfo {
				log.Printf("Remounting %q read-only", i.Mountpoint)
				if err := unix.Mount("", i.Mountpoint, "", unix.MS_REMOUNT|unix.MS_RDONLY, ""); err != nil {
					log.Fatalf("Failed to remount %q: %v", i.Mountpoint, err)
				}
			}
		} else {
			if err := os.WriteFile("/proc/sysrq-trigger", []byte("u"), 00644); err != nil {
				log.Fatalf("Failed to remount rootfs: %v", err)
			}
		}
	}

	log.Println("Overwriting flash contents via mtd access...")
	if err := flashcp.FlashCp(*img, *dev, 0); err != nil {
		log.Fatalf("Failed to overwrite flash: %v", err)
	}

	log.Println("Resetting system now...")
	if err := unix.Reboot(unix.LINUX_REBOOT_CMD_RESTART); err != nil {
		log.Fatalf("Failed to reboot machine: %v", err)
	}
}

func RWFilter(info *mountinfo.Info) (bool, bool) {
	var pseudoFs bool
	switch info.FSType {
	case "devtmpfs", "proc", "sysfs", "cgroup", "cgroup2", "pstore", "bpf", "securityfs", "debugfs", "tracefs", "mqueue", "hugetlbfs", "configfs", "efivarfs", "autofs", "binfmt_misc", "fusectl", "rpc_pipefs", "devpts", "tmpfs":
		pseudoFs = true
	default:
		pseudoFs = false
	}

	if pseudoFs || strings.Contains(info.Options, "ro") || info.Mountpoint == "/" {
		return true, false
	}

	return false, false
}
