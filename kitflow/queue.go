package kitflow

import (
    "context"
    "encoding/json"
    "time"

    "github.com/rs/zerolog/log"
    "github.com/wagslane/go-rabbitmq"
    "go.temporal.io/api/enums/v1"
    "go.temporal.io/sdk/client"
    "go.temporal.io/sdk/temporal"
)

// QueueJob 队列工作流任务定义
type QueueJob struct {
    // Id 工作流唯一ID
    Id string `json:"id"`
    // Name 工作流名称
    Name string `json:"name"`
    // Arg 工作流参数
    Arg any `json:"arg"`
    // TaskQueue Temporal任务队列名称
    TaskQueue                                string                         `json:"task_queue"`
    StartDelay                               time.Duration                  `json:"start_delay"`
    WorkflowExecutionTimeout                 time.Duration                  `json:"execution_timeout"`
    WorkflowIDReusePolicy                    enums.WorkflowIdReusePolicy    `json:"reuse_policy"`
    WorkflowIDConflictPolicy                 enums.WorkflowIdConflictPolicy `json:"conflict_policy"`
    WorkflowExecutionErrorWhenAlreadyStarted bool                           `json:"already_started"`
    RetryPolicy                              *temporal.RetryPolicy          `json:"retry_policy"`
}

// DefaultTaskQueue 默认任务队列名称
var DefaultTaskQueue = "default"

// Validate 校验任务配置
func (j *QueueJob) Validate() error {
    if j.TaskQueue == "" {
        j.TaskQueue = DefaultTaskQueue
    }
    return nil
}

// GetRetryPolicy 获取重试策略
func (j *QueueJob) GetRetryPolicy() *temporal.RetryPolicy {
    return j.RetryPolicy
}

// QueuePublisher 队列工作流发布器
type QueuePublisher struct {
    publisher *rabbitmq.Publisher
}

// PublisherOptions 发布器配置选项
type PublisherOptions struct {
    options []func(*rabbitmq.PublisherOptions)
}

// WithPublisherOption 添加原生 rabbitmq.PublisherOptions 配置
func WithPublisherOption(opt func(*rabbitmq.PublisherOptions)) func(*PublisherOptions) {
    return func(o *PublisherOptions) {
        o.options = append(o.options, opt)
    }
}

// NewQueuePublisher 创建队列发布器
func NewQueuePublisher(conn *rabbitmq.Conn, optionFuncs ...func(*PublisherOptions)) (*QueuePublisher, error) {
    opts := &PublisherOptions{}
    for _, fn := range optionFuncs {
        fn(opts)
    }
    publisher, err := rabbitmq.NewPublisher(conn, opts.options...)
    if err != nil {
        return nil, err
    }
    return &QueuePublisher{publisher: publisher}, nil
}

// Close 关闭发布器
func (p *QueuePublisher) Close() {
    p.publisher.Close()
}

// Publish 发布工作流任务到RabbitMQ队列
func (p *QueuePublisher) Publish(job *QueueJob, queue string) error {
    data, err := json.Marshal(job)
    if err != nil {
        return err
    }
    return p.publisher.Publish(data, []string{queue})
}

// PublishOnce 发布一次后自动关闭
func (p *QueuePublisher) PublishOnce(job *QueueJob, queue string) error {
    defer p.Close()
    return p.Publish(job, queue)
}

// PublishQueueOnce 便捷函数：创建发布器、发布任务、自动关闭
func PublishQueueOnce(conn *rabbitmq.Conn, job *QueueJob, queue string) error {
    pub, err := NewQueuePublisher(conn)
    if err != nil {
        return err
    }
    return pub.PublishOnce(job, queue)
}

// ConsumerOptions 消费者配置选项
type ConsumerOptions struct {
    options []func(*rabbitmq.ConsumerOptions)
}

// WithConsumerOption 添加原生 rabbitmq.ConsumerOptions 配置
func WithConsumerOption(opt func(*rabbitmq.ConsumerOptions)) func(*ConsumerOptions) {
    return func(o *ConsumerOptions) {
        o.options = append(o.options, opt)
    }
}

// ConsumeQueue 消费RabbitMQ队列并执行Temporal工作流（纯净版，无默认配置）
func ConsumeQueue(flowClient client.Client, conn *rabbitmq.Conn, queue string, optionFuncs ...func(*ConsumerOptions)) error {
    opts := &ConsumerOptions{}
    for _, fn := range optionFuncs {
        fn(opts)
    }

    consumer, err := rabbitmq.NewConsumer(conn, queue, opts.options...)
    if err != nil {
        return err
    }
    defer consumer.Close()

    ctx := context.Background()
    return consumer.Run(func(d rabbitmq.Delivery) rabbitmq.Action {
        job := new(QueueJob)
        if err = json.Unmarshal(d.Body, job); err != nil {
            log.Err(err).Msg("failed to unmarshal message payload")
            return rabbitmq.Ack
        }
        if err = job.Validate(); err != nil {
            log.Err(err).Msg("invalid job configuration")
            return rabbitmq.Ack
        }
        
        startOpts := client.StartWorkflowOptions{
            ID:                                       job.Id,
            TaskQueue:                                job.TaskQueue,
            WorkflowExecutionTimeout:                 job.WorkflowExecutionTimeout,
            WorkflowIDReusePolicy:                    job.WorkflowIDReusePolicy,
            WorkflowIDConflictPolicy:                 job.WorkflowIDConflictPolicy,
            WorkflowExecutionErrorWhenAlreadyStarted: job.WorkflowExecutionErrorWhenAlreadyStarted,
            RetryPolicy:                              job.GetRetryPolicy(),
        }
        if job.StartDelay > 0 {
            startOpts.StartDelay = job.StartDelay
        }
        wf, err := flowClient.ExecuteWorkflow(ctx, startOpts, job.Name, job.Arg)
        if temporal.IsTimeoutError(err) || temporal.IsPanicError(err) {
            log.Debug().Str("ID", d.MessageId).Msg("workflow execution failed, requeuing message")
            return rabbitmq.NackRequeue
        }
        if err != nil {
            log.Err(err).Str("ID", d.MessageId).Msg("workflow execution failed")
            return rabbitmq.Ack
        }
        _ = wf.Get(ctx, nil)
        return rabbitmq.Ack
    })
}

// ConsumeQueueDefault 消费RabbitMQ队列并执行Temporal工作流（带默认配置）
func ConsumeQueueDefault(flowClient client.Client, conn *rabbitmq.Conn, queue string, concurrency int, optionFuncs ...func(*ConsumerOptions)) error {
    defaultOpts := []func(*ConsumerOptions){
        WithConsumerOption(rabbitmq.WithConsumerOptionsExchangeName("")),
        WithConsumerOption(rabbitmq.WithConsumerOptionsQueueDurable),
        WithConsumerOption(rabbitmq.WithConsumerOptionsConcurrency(concurrency)),
        WithConsumerOption(rabbitmq.WithConsumerOptionsQOSPrefetch(concurrency * 10)),
    }
    return ConsumeQueue(flowClient, conn, queue, append(defaultOpts, optionFuncs...)...)
}
