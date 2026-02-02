package job

import (
	"github.com/anuragShingare30/go-boilerplate/internal/config"
	"github.com/hibiken/asynq"
	"github.com/rs/zerolog"
)

// Client is exposed; server and logger is internal use only
type JobService struct {
	Client *asynq.Client
	server *asynq.Server
	logger *zerolog.Logger
}

// this created new job service with concurrent workers to resolve the jobs
func NewJobService(cfg *config.Config, logger *zerolog.Logger) (*JobService){
	redisAddr := cfg.Redis.Address

	client := asynq.NewClient(asynq.RedisClientOpt{
		Addr: redisAddr,
	})

	server := asynq.NewServer(asynq.RedisClientOpt{Addr: redisAddr}, asynq.Config{
		Concurrency: 10, // represents number of concurrent workers
		Queues: map[string]int{
			"critical": 6,
			"default": 3,
			"low": 1,
		},
	})

	return &JobService{
		Client: client,
		server: server,
		logger: logger,
	}
}

// start the job server
func (job *JobService) StartServer() error{
	// Register task handlers
	mux := asynq.NewServeMux()
	mux.HandleFunc(TaskWelcome, job.handleWelcomeEmailTask)

	job.logger.Info().Msg("Starting background Job Server")
	if err := job.server.Start(mux); err != nil{
		return err
	}

	return nil
}

// stop the job server
func (job *JobService) StopServer(){
	job.logger.Info().Msg("Stopping the Job Server")
	job.server.Shutdown()
	job.Client.Close()
}