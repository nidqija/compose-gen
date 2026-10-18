package features


import (
	"os/exec"
	"fmt"
	"os"
)

// function to run docker compose with given subcommands and arguments
// it will check if docker is installed and will run the docker compose command
// if not , we will output a message telling user to install docker LOLOLOLOL
func RunCompose(subcmd string , defaultArgs[] string , userArgs[] string){
	

	if _, err := exec.LookPath("docker") ; err != nil {
         fmt.Println("Docker is not installed , please install Docker via this link : https://www.docker.com/products/docker-desktop/")
		os.Exit(1)
	}

	var composeArgs []string
	
	if subcmd != "" {
		composeArgs = append(composeArgs, subcmd)
	}

	if len(userArgs) == 0 {
        composeArgs = append(composeArgs, defaultArgs...)
	} else {
		composeArgs = append(composeArgs, userArgs...)
	}

	finalArgs := append([]string{"compose"} , composeArgs...)

	cmd := exec.Command("docker" , finalArgs...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		if _, ok := err.(*exec.ExitError); ok {
			fmt.Println("Error running docker compose: " , err)
			os.Exit(1)
		} else {
			fmt.Println("Unexpected error running docker compose: " , err)
			os.Exit(1)
		}
	}

	

}


func PrintHelp() {
	fmt.Println("Usage: ./compose-gen [command]")
	fmt.Println("Commands:")
	fmt.Println("  up          Start services")
	fmt.Println("  down        Stop services")
	fmt.Println("  build       Build services")
	fmt.Println("  ps          List services")
	fmt.Println("  logs        View logs")
	fmt.Println("  help        Show this help")
}