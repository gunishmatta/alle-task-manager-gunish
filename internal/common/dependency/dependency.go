package dependency

import (
	"alle-task-manager-gunish/internal/api/handler"
	"alle-task-manager-gunish/internal/common/config"
	"alle-task-manager-gunish/internal/common/database"
	"alle-task-manager-gunish/internal/common/kafka"
	loggingtype "alle-task-manager-gunish/internal/common/logging"
	"alle-task-manager-gunish/internal/domain/repository"
	"alle-task-manager-gunish/internal/service"
	"context"
	"errors"
)

type Container struct {
	config         *config.Config
	database       *database.Database
	taskRepository repository.TaskRepository
	taskService    *service.TaskService
	taskEventSvc   *service.TaskEventService
	kafkaProducer  *kafka.Producer
	taskHandler    *handler.TaskHandler
	kafkaConsumer  *kafka.Consumer
}

type Option func(*Container) error

func NewContainer(options ...Option) (*Container, error) {
	c := &Container{}

	for _, option := range options {
		if err := option(c); err != nil {
			return nil, err
		}
	}

	if c.config == nil {
		return nil, errors.New("config is required")
	}
	logger := loggingtype.GetLogger()
	if logger == nil {
		return nil, errors.New("logger is required")
	}

	if err := c.initializeRepositories(); err != nil {
		return nil, err
	}

	if err := c.initializeServices(); err != nil {
		return nil, err
	}

	if err := c.initializeHandlers(); err != nil {
		return nil, err
	}

	return c, nil
}

func WithConfig(cfg *config.Config) Option {
	return func(c *Container) error {
		if cfg == nil {
			return errors.New("nil config provided")
		}
		c.config = cfg
		return nil
	}
}

func WithDatabase(db *database.Database) Option {
	return func(c *Container) error {
		if db == nil {
			return errors.New("nil database provided")
		}
		c.database = db
		return nil
	}
}

func (c *Container) initializeRepositories() error {
	if c.taskRepository == nil {
		if c.database == nil {
			return errors.New("database is required to initialize repositories")
		}
		repo, err := repository.NewGormTaskRepository(c.database.Db)
		if err != nil {
			return err
		}
		c.taskRepository = repo
	}

	return nil
}

func (c *Container) initializeKafka() error {
	if c.config == nil {
		return errors.New("config is required to initialize Kafka")
	}

	if c.kafkaProducer == nil {
		producer, err := kafka.NewProducer(c.config.Kafka.Brokers)
		if err != nil {
			return err
		}
		c.kafkaProducer = producer
	}

	consumerService := service.NewTaskEventConsumerService()
	consumer, err := kafka.NewConsumer(
		c.config.Kafka.Brokers,
		"task-management-group",
		[]string{"task-events"},
		consumerService.HandleMessage,
	)
	if err != nil {
		return err
	}
	c.kafkaConsumer = consumer

	return nil
}

func (c *Container) initializeServices() error {
	if err := c.initializeKafka(); err != nil {
		return err
	}

	if c.taskEventSvc == nil {
		c.taskEventSvc = service.NewTaskEventService(c.kafkaProducer)
	}

	if c.taskService == nil {
		c.taskService = service.NewTaskService(c.taskRepository, c.taskEventSvc)
	}

	return nil
}

func (c *Container) initializeHandlers() error {
	if c.taskHandler == nil {
		c.taskHandler = handler.NewTaskHandler(c.taskService)
	}

	return nil
}

func (c *Container) InitializeKafkaConsumer(ctx context.Context) (*kafka.Consumer, error) {
	return kafka.NewConsumer(c.config.Kafka.Brokers, "task-management-group", []string{"task-events"}, service.NewTaskEventConsumerService().HandleMessage)
}

func (c *Container) TaskHandler() *handler.TaskHandler {
	return c.taskHandler
}

func (c *Container) Config() *config.Config {
	return c.config
}

func (c *Container) Database() *database.Database {
	return c.database
}

func (c *Container) TaskService() *service.TaskService {
	return c.taskService
}

func (c *Container) TaskRepository() repository.TaskRepository {
	return c.taskRepository
}

func (c *Container) Close() {
	logger := loggingtype.GetLogger()
	if c.database != nil {
		if err := c.database.Close(); err != nil {
			logger.Error("Failed to close database", "error", err)
		}
	}

	if c.kafkaProducer != nil {
		if err := c.kafkaProducer.Close(); err != nil {
			logger.Error("Failed to close Kafka producer", "error", err)
		}
	}

	if c.kafkaConsumer != nil {
		if err := c.kafkaConsumer.Close(); err != nil {
			logger.Error("Failed to close Kafka consumer", "error", err)
		}
	}
}
