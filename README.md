# Online workshop: Building Concurrent Web Applications
This is the repository for the [GopherCon 2022 workshop titled "Building Concurrent Web Applications"](https://www.gophercon.com/agenda/session/970935).

**Time:** Thu Oct 06, 11:00 AM - 3:00 PM CDT 

This repository contains the "Digital Ice Cream" van web application, which is the starting point of the workshop. 

**Warning:** The current implementation is not concurrency safe, as it will be fixed during the workshop.

## Instructors
[Adelina Simion](https://twitter.com/classic_addetz)\
Technology Evangelist @ [Form3](https://twitter.com/Form3Tech)

[Joseph Woodward](https://twitter.com/_josephwoodward)\
Senior Software Engineer @ [Form3](https://twitter.com/Form3Tech)

## Setup instructions
- Install and configure an editor for Go.
- Install a Go environment with Go 1.18 or later. Please follow the [official installation](https://go.dev/dl/) for your operating system.
- Sign up for a [GitHub](https://github.com/signup) account, if you do not already have one. 
- [Clone](https://docs.github.com/en/repositories/creating-and-managing-repositories/cloning-a-repository) this repository locally. 
- Computers should be capable of modern software development, such as access to install and run binaries, install a code editor, etc. 
- *Optional:* Install [Postman](https://www.postman.com/downloads/), as we will be sending requests to the application. We will also provide commandline alternatives usin [cURL](https://curl.se/docs/install.html).

## Start the application
- Follow [Setup instructions](#setup-instructions).
- From the root of your cloned repository, run the following command in the terminal. The server will be listening at `http://localhost:3000/orders`. 
```bash
$ go run cmd/server/main.go
```

## Run Load Tests
- Follow the [Setup instructions](#setup-instructions) and [Start the application](#start-the-application) steps.
- On first run, install the load testing tool. 
```bash
$ go get -u github.com/rakyll/hey
```
- From the root of the cloned repository, start a new terminal window and run the following command.
```bash
$ hey -n 5 -c 2 -z 3s -m POST -T "application/json" -d ./body.txt http://localhost:3000/orders
```

## We look forward to seeing you all at our workshop! 

![Our Logo](https://raw.githubusercontent.com/form3tech-oss/.github/master/profile/form3-logo-gopher.png)

