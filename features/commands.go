package features


import (
	"os/exec"
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

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


// function to read the generated docker-compose.yml and cloudflare log 
// and extract the cloudflare tunnel URL
// and update the .env file with the tunnel URL
func FetchTunnelURL(){

	// read docker-compose.yml file
	composeData , err := os.ReadFile("docker-compose.yml")


	// if there is no cloudflare service being found from dockerfile
	// return 
	if err != nil || !strings.Contains(string(composeData), "cloudflare"){
		return
	}


	// set target env file path
	targetEnv := "configs/.env"

	// check if the target env file path is exist , if not , return
	if _, err := os.Stat(targetEnv);os.IsNotExist(err){
		fmt.Println("Filepath is not exist")
		return
	}
    
	// regex to find cloudflare tunnel URL
	re := regexp.MustCompile(`https://[a-zA-Z0-9-]+\.trycloudflare\.com`)

	// log the process
	fmt.Println("Fetching Cloudflare Tunnel URL...")


	// init a timeout to find the generated url
	timeout := time.After(10 * time.Second)

	// ticker to check the generated url for every 500 millisecond
	ticker := time.NewTicker(500 * time.Millisecond)

	// stop the ticker when the function is done
	defer ticker.Stop()


	// init a tunnel url variable
	var tunnelURL string 

	// start the ticker
	for {

		// we use select as switches
		// it is like an if else statement
		// case 
		select {

		// case if timeout ( not found )
		case <- timeout:
			fmt.Println("\nTimed out waiting for Cloudflare Tunnel URL")
			return

		// case if ticker is triggered ( found )
		case <-ticker.C:
			fmt.Print(".")

			// read the generated tunnel log
		if data , err := os.ReadFile("logs/tunnel.log"); err == nil {
			if match := re.FindString(string(data)); match != "" {
				tunnelURL = match
				break
			}
		}
		}

		// if tunnel url is found , break the loop
		if tunnelURL != "" {
			break
		}

	
	}

	// print the found tunnel url from the logs for debugging
	fmt.Println("\nCloudflare Tunnel URL: ", tunnelURL)

	// read the target env file
	data , err := os.ReadFile(targetEnv)

	if err != nil {
		fmt.Println("Error reading env file: ", err)
		return
	}

	// init content variable by converting bytes into string
	content := string(data)

	// init a regex to find the url line in the env file
	urlLineRegex := regexp.MustCompile(`(?m)^APP_BASE_URL=.*$`)

	// init a updated content variable
	var updatedContent string
	

	// check if the url line exists in the env file
	if urlLineRegex.MatchString(content){

		// if found , replace the url line with the new tunnel url found
		updatedContent = urlLineRegex.ReplaceAllString(content , "APP_BASE_URL="+tunnelURL)
	} else {

		// if not found , append the new tunnel url in the env file
		updatedContent = strings.TrimRight(content , "\n") + "\n" + "APP_BASE_URL="+ tunnelURL
	}

	// write the updated content into the env file
	if err := os.WriteFile(targetEnv, []byte(updatedContent), 0644); err != nil {
		fmt.Println("Error writing env file: ", err)
		return
	}

	// log the updated url for debugging purposes
	fmt.Printf("Updated %s with Cloudflare Tunnel URL with APP_BASE_URL=%s\n" , targetEnv , tunnelURL)
	
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