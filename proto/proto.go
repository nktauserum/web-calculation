package proto

import (
	"context"
	"fmt"

	"github.com/nktauserum/web-calculation/orchestrator/pkg/service"
	tsk "github.com/nktauserum/web-calculation/orchestrator/pkg/task"
	"github.com/nktauserum/web-calculation/proto/pb"
	"github.com/nktauserum/web-calculation/shared"
)

type Server struct {
	pb.TaskServiceServer
}

func (s *Server) GetAvailableTask(context.Context, *pb.Empty) (*pb.Task, error) {
	queue := service.GetQueue()

	tasks := queue.GetTasks()
	if len(tasks) == 0 {
		return nil, nil
	}

	var finalTask *shared.Task
	var found bool

	for _, task := range tasks {
		if !task.Status && tsk.Complete(task) {
			finalTask = queue.FindTask(task.ID)
			found = true
			break
		}
	}

	if !found {
		return &pb.Task{Status: false}, nil
	}

	return &pb.Task{
		Id:            finalTask.ID,
		Arg1:          finalTask.FirstArgument,
		Arg2:          finalTask.SecondArgument,
		Operator:      string(finalTask.Operator),
		OperationTime: finalTask.OperationTime,
		Status:        true,
		Result:        finalTask.Result,
	}, nil
}

func (s *Server) CompleteTask(ctx context.Context, taskResult *pb.TaskResult) (*pb.Empty, error) {
	queue := service.GetQueue()
	queue.Done(taskResult.Id, taskResult.Result)
	fmt.Printf("Задача %d успешно выполнена!\n", taskResult.Id)
	return &pb.Empty{}, nil
}
