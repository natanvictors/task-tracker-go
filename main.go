package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"
)

type task struct {
	Id          int       `json:"id"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type command struct {
	command   string
	arguments []string
}

type commands struct {
	registeredCommands map[string]func(command) error
}

func (c *commands) register(name string, f func(command) error) {
	c.registeredCommands[name] = f
}

func (c *commands) run(cmd command) error {
	f, ok := c.registeredCommands[cmd.command]
	if !ok {
		return errors.New("command not found")
	}

	return f(cmd)
}

func commandExit(cmd command) error {
	fmt.Println("Quitting program...")
	os.Exit(0)
	return nil
}

func commandAdd(cmd command) error {
	file, err := os.OpenFile("tasks.json", os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	var fileTasks []task

	data, err := io.ReadAll(file)
	if err != nil {
		return nil
	}

	if len(data) != 0 {
		err = json.Unmarshal(data, &fileTasks)
		if err != nil {
			return err
		}
	}

	newTask := task{
		Id:          len(fileTasks) + 1,
		Description: cmd.arguments[0],
		Status:      "todo",
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}

	fileTasks = append(fileTasks, newTask)

	returnData, err := json.MarshalIndent(fileTasks, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile("tasks.json", returnData, 0644)
}

func main() {

	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("task-cli ")
		scanner.Scan()
		input := strings.Fields(scanner.Text())

		args := input[1:]

		cmds := commands{
			registeredCommands: make(map[string]func(command) error),
		}

		cmds.register("exit", commandExit)
		cmds.register("add", commandAdd)

		err := cmds.run(command{
			command:   input[0],
			arguments: args,
		})

		if err != nil {
			log.Fatal(err)
		}

	}
}
