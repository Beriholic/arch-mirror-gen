package file

import (
	"os"
	"time"
)

func WriteFileToPacmanMirrorlist(mirrorList []string) error {
	curTime := time.Now().Format("2006-01-02 15:04:05")

	header := "#######################################################\n" +
		"####Automatically generated via Arch-Mirror-gen########\n" +
		"#######################################################\n" +
		"##############" + curTime + "######################\n"

	path := "/etc/pacman.d/mirrorlist"
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(header)

	if err != nil {
		return err
	}

	for _, mirror := range mirrorList {
		_, err := file.WriteString("Server = " + mirror + "\n")
		if err != nil {
			return err
		}
	}

	return nil
}
