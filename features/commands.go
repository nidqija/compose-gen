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
	

	// our program will look for a docker path in our device , 
	// if the docker path is not found , it will output an error and exit
	// if the docker path is found , it will run the docker compose command
	if _, err := exec.LookPath("docker") ; err != nil {
         fmt.Println("Docker is not installed , please install Docker via this link : https://www.docker.com/products/docker-desktop/")
		os.Exit(1)
	}

	// init a compose argument (subcommand up , down , build , ps , logs , )
	var composeArgs []string
	
	// if subcommand is empty , use default arguments
	if subcmd != "" {
		composeArgs = append(composeArgs, subcmd)
	}

	// if no user argument is provided , use default arguments
	// else , use user arguments
	if len(userArgs) == 0 {
        composeArgs = append(composeArgs, defaultArgs...)
	} else {
		composeArgs = append(composeArgs, userArgs...)
	}

	// join compose subcommands with default arguments
	finalArgs := append([]string{"compose"} , composeArgs...)

	// make a new command that has our joined compose subcommands
	cmd := exec.Command("docker" , finalArgs...)

	// map our program I/O to docker's I/O
	// this will allow us to see the output of the docker command
	// and also allow us to interact with the docker command
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// run the command and check for errors
	// if error is detected , exit the program
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


// function to print help message
// this function will show the usage of the program
// and the commands that can be run
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