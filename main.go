package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
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

func readFile(path string) ([]task, error) {
	file, err := os.OpenFile(path, os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var fileTasks []task

	data, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}

	if len(data) != 0 {
		err = json.Unmarshal(data, &fileTasks)
		if err != nil {
			return nil, err
		}
	}

	return fileTasks, nil
}

func saveFile(path string, fileTasks []task) error {
	data, err := json.MarshalIndent(fileTasks, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}

func commandAdd(cmd command) error {

	fileTasks, err := readFile("tasks.json")
	if err != nil {
		return err
	}

	newTask := task{
		Id:          len(fileTasks) + 1,
		Description: strings.Trim(strings.Join(cmd.arguments, " "), "\""),
		Status:      "todo",
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}

	fileTasks = append(fileTasks, newTask)

	return saveFile("tasks.json", fileTasks)
}

func commandUpdate(cmd command) error {

	taskID, err := strconv.Atoi(cmd.arguments[0])
	if err != nil {
		return err
	}

	newDesc := strings.Trim(strings.Join(cmd.arguments[1:], " "), "\"")

	fileTasks, err := readFile("tasks.json")
	if err != nil {
		return err
	}

	fileTasks[taskID-1].Description = newDesc

	return saveFile("tasks.json", fileTasks)
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
		cmds.register("update", commandUpdate)

		err := cmds.run(command{
			command:   input[0],
			arguments: args,
		})

		if err != nil {
			log.Fatal(err)
		}

	}
}
