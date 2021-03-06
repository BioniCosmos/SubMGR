package main

import (
    "flag"
    "fmt"
    "os"
)

func main() {
    // Flags for "add"
    addCmd := flag.NewFlagSet("add", flag.ExitOnError)
    addConfig := addCmd.String("config", "./config.json", "The path to the configuration file")
    addNode := addCmd.String("node", "", "A node name specified in the configuration file")
    addEmail := addCmd.String("email", "", "User's email address")
    addId := addCmd.String("uuid", "", "User's UUID")
    addLevel := addCmd.Int("level", 0, "User's level")

    // Flags for "rm"
    rmCmd := flag.NewFlagSet("rm", flag.ExitOnError)
    rmConfig := rmCmd.String("config", "./config.json", "The path to the configuration file")
    rmNode := rmCmd.String("node", "", "A node name specified in the configuration file")
    rmEmail := rmCmd.String("email", "", "User's email address")

    // Flags for "load"
    loadCmd := flag.NewFlagSet("load", flag.ExitOnError)
    loadConfig := loadCmd.String("config", "./config.json", "The path to the configuration file")
    loadUser := loadCmd.String("user", "./user.json", "The path to user's information file")

    // Flags for "empty"
    emptyCmd := flag.NewFlagSet("empty", flag.ExitOnError)
    emptyConfig := emptyCmd.String("config", "./config.json", "The path to the configuration file")
    emptyUser := emptyCmd.String("user", "./user.json", "The path to user's information file")

    // Flags for "sub"
    subCmd := flag.NewFlagSet("sub", flag.ExitOnError)
    subUser := subCmd.String("user", "./user.json", "The path to user's information file")

    // Flags for "server"
    serverCmd := flag.NewFlagSet("server", flag.ExitOnError)
    serverUser := serverCmd.String("user", "./user.json", "The path to user's information file")
    serverListen := serverCmd.String("listen", "127.0.0.1:8080", "The listening address and port")
    serverPath := serverCmd.String("path", "", "The extra path of the server")

    if len(os.Args) < 2 {
        fmt.Println("No arguments specified")
        fmt.Println("Enter submgr add/rm/load/empty/sub/server -h for detailed information.")
        os.Exit(1)
    }

    switch os.Args[1] {
    case "add":
        addCmd.Parse(os.Args[2:])
        add(*addConfig, *addNode, *addEmail, *addId, *addLevel)
    case "rm":
        rmCmd.Parse(os.Args[2:])
        rm(*rmConfig, *rmNode, *rmEmail)
    case "load":
        loadCmd.Parse(os.Args[2:])
        load(*loadConfig, *loadUser)
    case "empty":
        emptyCmd.Parse(os.Args[2:])
        empty(*emptyConfig, *emptyUser)
    case "sub":
        subCmd.Parse(os.Args[2:])
        sub(*subUser)
    case "server":
        serverCmd.Parse(os.Args[2:])
        server(*serverUser, *serverListen, *serverPath)
    default:
        fmt.Println("Wrong arguments")
        fmt.Println("Enter submgr add/rm/load/empty/sub/server -h for detailed information.")
        os.Exit(1)
    }
}
