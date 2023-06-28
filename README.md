# WISdoM OSS — Microservice Template
> **Warning**
> This repository only acts as template and _should never be cloned or forked_ 
> when creating a new microservice with this template.

## Using the template
1. Download this archive as `.zip` or `.tar.gz` (whatever you prefer)

2. Extract the downloaded archive to a directory of your choice and remove the
    parent folders which may have been created during the download

3. Make sure that your folder now contains at least the following file structure:
   ```
   ├── globals
   │  ├── connections.go      (contains globally available connections)
   │  ├── variables.go        (contains globally available variables)
   ├── resources
   │  ├── authConfig.json     (contains auth config)
   │  ├── environment.json    (contains the environment setup)
   │  ├── errors.json         (contains http errors)
   │  ├── queries.sql         (contains sql queries for the service)
   ├── routes
   │  ├── templates.go        (contains three template routes)
   ├── .gitignore
   ├── init.go                (contains code used during startup)
   ├── template-service.go    (contains the bootstrapping code for the service)
   ├── go.mod                 (contains the dependencies of the service)
   ├── go.sum                 (is the lockfile for the dependencies)
   ```
   
4. **Important** Change the service name

    To change the service name, you need to edit the file `globals/variables.go`
    which should contain the following line
    
    ```go 
    const ServiceName = "template-service"
    ```
   
    This line needs to be changed to your service name. This constant is
    used in logs and error handling to identify the service.
   
5. Initialize a Git Repository with `main` as default branch

    ```shell
   git init -b main
    ```
6. Add all files to the repository

    ```shell
   git add -A
    ```
   
7. Commit the template to the repository

    ```shell
   git commit -m "loading wisdom-oss/microservice-template"
   ```
   
8. Set up a remote origin for the repository

    ```shell
   git remote add origin <your-remote-url>
    ```
   
9. Push the repository to the remote origin

    ```shell
   git push origin main
    ```
   
10. :tada: You are now able to develop your new microservice

11. Change the README to the contents you desire in here
