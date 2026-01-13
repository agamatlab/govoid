package helper

import (
	"log"
	"os/exec"
	"strings"
)

func ReverseSlice(slice []string) {
	for i, j := 0, len(slice)-1; i < j; i, j = i+1, j-1 {
		slice[i], slice[j] = slice[j], slice[i] // Swap elements
	}
}

func KillApp(appName string) error {
    cmd := exec.Command("pkill", "-x", "-i", appName)
    return cmd.Run()
}

func GetApps() []string {
	fullCommand := `lsappinfo visibleProcessList | tr -d '(),"=' | sed 's/visibleProcessList//g' | xargs -n 1 lsappinfo info -only name | grep "LSDisplayName" | cut -d '"' -f 4`
	cmd := exec.Command("sh", "-c", fullCommand)
	out, errout := cmd.Output()

	if errout != nil {
		log.Printf("Error output")
	} else{
		stringout := string(out)

		apps := strings.Split(stringout, "\n")
		newApps := make([]string, 0)
		for _, app := range apps {
			if strings.TrimSpace(app) != "" {
				newApps = append(newApps, app)
			}
		}
		return newApps
	}
	return nil
	
}

