package cmd

import (
	"fmt"
	"os"
	"os/user"

	"github.com/Beriholic/arch-mirror-gen/file"
	"github.com/Beriholic/arch-mirror-gen/mirror"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "arch-mirror-gen",
	Short: "Arch Linux 镜像列表生成工具",
	Long:  `自动获取大陆地区最优的Archlinux镜像.`,
	Run: func(cmd *cobra.Command, args []string) {
		curUser, err := user.Current()
		if err != nil {
			fmt.Println("获取当前用户失败")
		}

		if curUser.Uid != "0" {
			fmt.Println("请使用root权限运行")
			return
		}

		list, err := mirror.GetChinaMirrorList()
		if err != nil {
			fmt.Printf("获取镜像列表失败: %v\n", err)
			return
		}

		err = file.WriteFileToPacmanMirrorlist(list)
		if err != nil {
			fmt.Printf("写入镜像列表失败: %v\n", err)
			return
		}

		fmt.Println("镜像列表已写入 /etc/pacman.d/mirrorlist")
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
