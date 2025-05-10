package controller

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/joho/godotenv"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/nktauserum/web-calculation/proto/pb"
	"github.com/nktauserum/web-calculation/shared"
)

type Agent struct {
	client pb.TaskServiceClient
	conn   *grpc.ClientConn
}

func NewAgent(address string) (*Agent, error) {
	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}

	client := pb.NewTaskServiceClient(conn)
	return &Agent{
		client: client,
		conn:   conn,
	}, nil
}

func (c *Agent) GetTask() (*shared.Task, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	task, err := c.client.GetAvailableTask(ctx, &pb.Empty{})
	if err != nil {
		return nil, err
	}

	if task == nil {
		return nil, nil
	}

	if !task.Status {
		return nil, fmt.Errorf("no tasks available")
	}

	return &shared.Task{
		ID:             task.Id,
		FirstArgument:  task.Arg1,
		SecondArgument: task.Arg2,
		Operator:       rune(task.Operator[0]),
		OperationTime:  task.OperationTime,
		Status:         task.Status,
		Result:         task.Result,
	}, nil
}

func (c *Agent) CompleteTask(id int64, result float64) error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	_, err := c.client.CompleteTask(ctx, &pb.TaskResult{
		Id:     id,
		Result: result,
	})
	return err
}

func (c *Agent) Close() {
	if c.conn != nil {
		c.conn.Close()
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

	for range workers {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				// Get task from orchestrator
				task, err := app.GetTask()
				if err != nil {
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
				result, err := calculateExpression(*task)
				if err != nil {
					log.Printf("Error calculating expression: %v", err)
					continue
				}

				err = app.CompleteTask(task.ID, result)
				if err != nil {
					log.Printf("Error completing task: %v", err)
					continue
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
	if err := godotenv.Load(".env"); err != nil {
		log.Printf("Error loading .env file: %v", err)
	}

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
		var addition_time time.Duration
		if os.Getenv("TIME_ADDITION_MS") == "" {
			addition_time = time.Second * 0
		} else {
			addition_time, err = time.ParseDuration(os.Getenv("TIME_ADDITION_MS") + "ms")
			if err != nil {
				log.Printf("Error parsing addition_time: %v", err)
				return 0, err
			}
		}

		time.Sleep(addition_time)
		return firstarg + secondarg, nil
	case '-':
		var subtraction_time time.Duration
		if os.Getenv("TIME_SUBTRACTION_MS") == "" {
			subtraction_time = time.Second * 0
		} else {
			subtraction_time, err = time.ParseDuration(os.Getenv("TIME_SUBTRACTION_MS") + "ms")
			if err != nil {
				log.Printf("Error parsing TIME_SUBTRACTION_MS: %v", err)
				return 0, err
			}
		}
		time.Sleep(subtraction_time)
		return firstarg - secondarg, nil
	case '*':
		// задержка для умножения
		var multiplication_time time.Duration
		if os.Getenv("TIME_MULTIPLICATIONS_MS") == "" {
			multiplication_time = time.Second * 0
		} else {
			multiplication_time, err = time.ParseDuration(os.Getenv("TIME_MULTIPLICATIONS_MS") + "ms")
			if err != nil {
				log.Printf("Error parsing TIME_MULTIPLICATIONS_MS: %v", err)
				return 0, err
			}
		}
		time.Sleep(multiplication_time)
		return firstarg * secondarg, nil
	case '/':
		// задержка для деления
		var division_time time.Duration
		if os.Getenv("TIME_DIVISIONS_MS") == "" {
			division_time = time.Second * 0
		} else {
			division_time, err = time.ParseDuration(os.Getenv("TIME_DIVISIONS_MS") + "ms")
			if err != nil {
				log.Printf("Error parsing TIME_DIVISIONS_MS: %v", err)
				return 0, err
			}
		}
		time.Sleep(division_time)

		if secondarg == 0 {
			return 0, fmt.Errorf("на ноль делить нельзя")
		}
		return firstarg / secondarg, nil

	}
	return 0, nil
}
