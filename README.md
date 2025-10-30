# Hostable Digital Cookbook

### Currently only pulls from allrecipes.com

### Compile from source
This project was developed with Golang 1.25.3 on Windows. Building with Linux is not currently supported, but should not be difficult with the basic `go build` for Linux.
To build the recipe_tool.exe, run the following command in the root of the project folder.
```
make
```
The resulting executable is placed in ${PROJECT_ROOT}\bin\

### Using the cookbook
To use the REPL-like command tool, enter the numbers corresponding to the presented instructions. To use the website, you must start *both* the web server and API server.
**Before starting the servers, you must first generate the config file. This can be done by selecting option 3. The web server does not currently support any other port besides 8000**
