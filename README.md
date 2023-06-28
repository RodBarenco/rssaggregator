# My First Golang Project

This is my first project using Golang. It was a great learning experience, and I would like to acknowledge the valuable contribution of the FreeCodeCamp tutorial by Wagslane, which played a crucial role in helping me get started.

## Project Modifications

I made some modifications to the project, replacing Goose and SQLC with GORM. While working on the project, I encountered compatibility issues with SQLC on Windows 11, so I opted for GORM as an alternative.

## CLI Enhancement

One potential improvement for the project is to add a command-line interface (CLI) that provides a user-friendly way to interact with the server. This can enhance usability and make it easier for users to interact with the application.

## Docker Scripts

I have provided two ready-to-use scripts, one for Windows and another for Linux. These scripts can be adapted based on the user's Docker Compose version if they choose to run the database using Docker. Please make sure to customize the scripts according to your specific environment requirements.

## Test Suite

I have initiated the setup for a test suite but have not yet written the actual tests. Adding a comprehensive test suite is an essential step to ensure the reliability and correctness of the application. Feel free to contribute by adding appropriate tests.

## Caveat

To modify the files rss.go, ./db/post.go, and scraper.go according to the structure of an XML or any other file you wish to aggregate, follow these steps:

Open the file rss.go and locate the code responsible for parsing the XML feed. This code should be responsible for extracting relevant information from the XML and creating blog posts. Modify this code to match the structure of the XML or the file you want to aggregate.

Similarly, open the file ./db/post.go and find the code that handles storing the blog posts in the database. Adjust this code to ensure that the extracted information from the XML or the file is stored correctly in the database.

Finally, open the file scraper.go and identify the code that initiates the scraping process. Make any necessary changes to this code to fetch the XML or the desired file from the URL https://wagslane.dev/index.xml or any other URL you wish to use for aggregation.

By modifying these three files (rss.go, ./db/post.go, and scraper.go), you can customize the structure and source of the content you want to aggregate. Ensure that the modifications align with the structure of the XML or the file you intend to use.

If you have any further questions, feel free to ask!

## Getting Started

To get started with the project, follow these steps:

1. Clone the repository to your local machine.
2. Install any necessary dependencies.
3. Set up your development environment.
4. Customize the configuration files if needed.
5. Run the application and start exploring the features.

Feel free to reach out if you have any questions or need further assistance!


### Docker Compose (docker-compose.yml)

```yaml
version: '2.17.3'
services:
  dev-db:
    container_name: rssagg-dev
    image: postgres:14.7
    ports:
      - 5434:5432
    environment:
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: 123
      POSTGRES_DB: go
    networks:
      - goserver

  test-db:
    container_name: rssagg-test
    image: postgres:14.7
    ports:
      - 5436:5432
    environment:
      POSTGRES_USER: postgres
      POSTGRES_PASSWORD: 123
      POSTGRES_DB: go
    networks:
      - goserver

networks:
  goserver:
```

### Script to run the project on Linux (init.sh)

```sh
#!/bin/bash

processFilePath="$PWD/rssaggregator.exe"
pidFilePath="$PWD/process.pid"
dockerComposeFilePath="$PWD/docker-compose.yml"
dockerDevContainerName="dev-db"
dockerTestContainerName="test-db"
envDev=".env"
envTest=".test.env"

# Function to check if Docker container is running
IsDockerContainerRunning() {
    local containerName=$1
    local containerStatus=$(docker inspect --format='{{.State.Status}}' "$containerName" 2> /dev/null)
    if [[ "$containerStatus" == "running" ]]; then
        return 0
    else
        return 1
    fi
}

# Prompt to select the application version
echo "Select the application version:"
echo -e "1. Dev" "\033[36m"
echo -e "2. Test" "\033[34m"
read -p "Enter the corresponding number for the version: " versionChoice

# Check Docker database status
if [[ "$versionChoice" == "1" ]]; then
    serviceName=$dockerDevContainerName
    envFilePath=$envDev
    echo -e "Environment variable: $envFilePath" "\033[36m"
elif [[ "$versionChoice" == "2" ]]; then
    serviceName=$dockerTestContainerName
    envFilePath=$envTest
    echo -e "Environment variable: $envFilePath" "\033[34m"
else
    echo "Invalid option. Exiting the script."
    exit 1
fi

# Compile Go code
go build
if [[ $? -eq 0 ]]; then
    if ! IsDockerContainerRunning "$serviceName"; then
        echo "Starting the Docker container..."
        docker-compose -f "$dockerComposeFilePath" up -d "$serviceName"
        # Wait for a few seconds for the container to fully start
        sleep 2
    fi

    # Start the application process
    nohup "$processFilePath" "$envFilePath" >/dev/null 2>&1 &
    echo "Process started with PID $!" "\033[35m"
    echo "$!" > "$pidFilePath"

    # Wait for the process to finish
    wait $!

    # Ask for confirmation to delete the container upon pressing Ctrl + C
    echo "Do you want to delete the container before exiting? (1 - Yes, 2 - No)"
    echo -e "Note: If you choose YES, the container will be permanently deleted:" "\033[31m\c"
    read -p " " containerChoice
    if [[ "$containerChoice" == "1" ]]; then
        echo "Stopping the Docker container..."
        docker-compose -f "$dockerComposeFilePath" rm -f -s -v "$serviceName"
    fi
else
    echo "Build failed"
fi

```


### Script to run the project on Windows (init.ps1)

```ps1
$processFilePath = "$PSScriptRoot\rssaggregator.exe"
$pidFilePath = "$PSScriptRoot\process.pid"
$dockerComposeFilePath = "$PSScriptRoot\docker-compose.yml"
$dockerDevContainerName = "dev-db"
$dockerTestContainerName = "test-db"
$envDev = ".env"
$envTest = ".test.env"

# Function to check if the Docker container is running
function IsDockerContainerRunning($containerName) {
    $containerStatus = docker inspect --format='{{.State.Status}}' $containerName 2> $null
    return ($containerStatus -eq "running")
}

# Prompt to select the application version
Write-Host "Select the application version:"
Write-Host "1." -NoNewline
Write-Host " Dev" -ForegroundColor Cyan
Write-Host "2." -NoNewline
Write-Host " Test" -ForegroundColor Blue
$versionChoice = Read-Host "Enter the corresponding number for the version"

# Check the status of the Docker database
if ($versionChoice -eq "1") {
    $serviceName = $dockerDevContainerName
    $envFilePath = $envDev
    Write-Host "Environment variable: " $envFilePath -ForegroundColor Cyan
} elseif ($versionChoice -eq "2") {
    $serviceName = $dockerTestContainerName
    $envFilePath = $envTest
    Write-Host "Environment variable: " $envFilePath -ForegroundColor Blue
} else {
    Write-Host "Invalid option. Exiting the script."
    exit
}

# Command to build the Go code
go build
if ($?) {
    if (-not (IsDockerContainerRunning $serviceName)) {
        Write-Host "Starting the Docker container..."
        docker-compose -f $dockerComposeFilePath up -d $serviceName
        # Wait for a few seconds for the container to fully start
        Start-Sleep -Seconds 2
    }

    # Start the application process
    $process = Start-Process -FilePath $processFilePath -NoNewWindow -PassThru -ArgumentList $envFilePath
    $process.Id | Out-File $pidFilePath
    Write-Host "Process started with PID $($process.Id)" -ForegroundColor Magenta

    $job = Start-Job -ScriptBlock {
        param($processId)
        Wait-Process -Id $processId
    } -ArgumentList $process.Id
    Write-Host "Job started with ID $($job.Id)" -ForegroundColor DarkGray

    try {
        Wait-Job -Job $job
    } finally {
        # Prompt when pressing Ctrl + C to delete the container
        Write-Host "Do you want to delete the container before exiting? (1 - Yes, 2 - No)"
        Write-Host "Note: If you choose YES, the container will be permanently deleted: " -ForegroundColor Red -NoNewline
        $containerChoice = Read-Host
        if ($containerChoice -eq "1") {
            Write-Host "Stopping the Docker container..."
            docker-compose -f $dockerComposeFilePath rm -f -s -v $serviceName
        }

        Stop-Job -Job $job | Remove-Job
    }
} else {
    Write-Host "Build failed"
}
```
### Exemple for .env and .test.env - make sure to use diferent ports

```.env
PORT=8000
DATABASE_URL="postgresql://postgres:123@localhost:5434/go"
```