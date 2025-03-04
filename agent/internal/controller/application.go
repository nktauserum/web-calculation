package controller

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/nktauserum/web-calculation/shared"
)

type Agent struct {
	Port int
}

func NewAgent(port int) *Agent {
	return &Agent{
		Port: port,
	}
}

func (app *Agent) Run() error {
	log.Println("Agent started!")
	computingPower := os.Getenv("COMPUTING_POWER")
	workers, err := strconv.Atoi(computingPower)
	if err != nil {
		workers = 1
	}

	var wg sync.WaitGroup

	processedTasks := sync.Map{}

	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				// Get task from orchestrator
				resp, err := http.Get("http://orchestrator:8080/internal/task")
				if err != nil {
					log.Printf("Error getting task: %v", err)
					time.Sleep(time.Second)
					continue
				}

				if resp.StatusCode == http.StatusNoContent {
					time.Sleep(time.Second)
					continue
				}

				body, err := io.ReadAll(resp.Body)
				if err != nil {
					log.Printf("Error reading response: %v", err)
					continue
				}

				var task shared.Task
				err = json.Unmarshal(body, &task)
				if err != nil {
					log.Printf("Error unmarshalling body: %v", err)
					continue
				}

				_, exists := processedTasks.Load(task.ID)
				if exists {
					//log.Printf("Задача %d уже в обработке", task.ID)
					continue
				}

				processedTasks.Store(task.ID, "")

				log.Printf("Получена задача %d", task.ID)

				// Calculate expression
				result, err := calculateExpression(task)
				if err != nil {
					log.Printf("Error calculating expression: %v", err)
					continue
				}

				var answer = shared.TaskResult{
					ID:     task.ID,
					Result: result,
				}

				data, err := json.Marshal(&answer)
				if err != nil {
					log.Printf("Error marshalling result: %v", err)
					continue
				}

				// Send result back
				_, err = http.Post("http://localhost:8080/internal/task", "application/json", bytes.NewBuffer(data))
				if err != nil {
					log.Printf("Error sending result: %v", err)
				}

				processedTasks.Delete(task.ID)
			}
		}()
	}

	wg.Wait()
	return nil
}

/*
Время выполнения операций задается переменными среды в миллисекундах
TIME_ADDITION_MS - время выполнения операции сложения в миллисекундах
TIME_SUBTRACTION_MS - время выполнения операции вычитания в миллисекундах
TIME_MULTIPLICATIONS_MS - время выполнения операции умножения в миллисекундах
TIME_DIVISIONS_MS - время выполнения операции деления в миллисекундах
*/

func calculateExpression(task shared.Task) (float64, error) {
	firstarg, err := strconv.ParseFloat(task.FirstArgument, 64)
	if err != nil {
		log.Printf("Error parsing first argument: %v", err)
		return 0, err
	}
	secondarg, err := strconv.ParseFloat(task.SecondArgument, 64)
	if err != nil {
		log.Printf("Error parsing second argument: %v", err)
		return 0, err
	}

	switch task.Operator {
	case '+':
		// задержка для сложения
		var addition_time, err = time.ParseDuration(os.Getenv("TIME_ADDITION_MS"))
		if err != nil {
			log.Printf("Error parsing addition_time: %v", err)
			return 0, err
		}
		time.Sleep(addition_time)
		return firstarg + secondarg, nil
	case '-':
		// задержка для вычитания
		var subtraction_time, err = time.ParseDuration(os.Getenv("TIME_SUBTRACTION_MS"))
		if err != nil {
			log.Printf("Error parsing subtraction_time: %v", err)
			return 0, err
		}
		time.Sleep(subtraction_time)
		return firstarg - secondarg, nil
	case '*':
		// задержка для умножения
		var multiplication_time, err = time.ParseDuration(os.Getenv("TIME_MULTIPLICATIONS_MS"))
		if err != nil {
			log.Printf("Error parsing multiplication_time: %v", err)
			return 0, err
		}
		time.Sleep(multiplication_time)
		return firstarg * secondarg, nil
	case '/':
		// задержка для деления
		var division_time, err = time.ParseDuration(os.Getenv("TIME_DIVISIONS_MS"))
		if err != nil {
			log.Printf("Error parsing division_time: %v", err)
			return 0, err
		}
		time.Sleep(division_time)

		if secondarg == 0 {
			return 0, fmt.Errorf("на ноль делить нельзя")
		}
		return firstarg / secondarg, nil

	}
	return 0, nil
}
