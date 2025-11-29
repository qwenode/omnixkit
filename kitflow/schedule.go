package kitflow

import (
    "context"

    "github.com/rs/zerolog/log"
    "go.temporal.io/api/enums/v1"
    "go.temporal.io/sdk/client"
)

// ScheduleConfig 定时任务配置
type ScheduleConfig struct {
    // Name 任务名称，同时作为 Schedule ID
    Name string
    // Cron cron 表达式列表
    Cron []string
    // Workflow 工作流名称
    Workflow string
    // TaskQueue 任务队列名称
    TaskQueue string
    // Args 工作流参数
    Args []any
    // Overlap 重叠策略，默认 SKIP
    Overlap enums.ScheduleOverlapPolicy
}

// ToScheduleOptions 转换为 Temporal ScheduleOptions
func (c *ScheduleConfig) ToScheduleOptions() client.ScheduleOptions {
    overlap := c.Overlap
    if overlap == enums.SCHEDULE_OVERLAP_POLICY_UNSPECIFIED {
        overlap = enums.SCHEDULE_OVERLAP_POLICY_SKIP
    }
    taskQueue := c.TaskQueue
    if taskQueue == "" {
        taskQueue = "default"
    }
    return client.ScheduleOptions{
        ID: c.Name,
        Spec: client.ScheduleSpec{
            CronExpressions: c.Cron,
        },
        Action: &client.ScheduleWorkflowAction{
            ID:        c.Name,
            Workflow:  c.Workflow,
            TaskQueue: taskQueue,
            Args:      c.Args,
        },
        Overlap: overlap,
    }
}

// ScheduleRegistry 定时任务注册器
type ScheduleRegistry struct {
    configs []ScheduleConfig
}

// NewScheduleRegistry 创建注册器
func NewScheduleRegistry() *ScheduleRegistry {
    return &ScheduleRegistry{
        configs: make([]ScheduleConfig, 0),
    }
}

// Add 添加定时任务配置
func (r *ScheduleRegistry) Add(config ScheduleConfig) *ScheduleRegistry {
    r.configs = append(r.configs, config)
    return r
}

// AddMany 批量添加定时任务配置
func (r *ScheduleRegistry) AddMany(configs ...ScheduleConfig) *ScheduleRegistry {
    r.configs = append(r.configs, configs...)
    return r
}

// Register 注册所有定时任务到 Temporal
// closeClient: 是否在注册完成后关闭客户端
func (r *ScheduleRegistry) Register(cli client.Client, closeClient bool) {
    if closeClient {
        defer cli.Close()
    }
    ctx := context.Background()

    // 获取已存在的 schedule
    exist := r.listExisting(ctx, cli)

    // 注册不存在的 schedule
    for _, config := range r.configs {
        if exist[config.Name] {
            continue
        }
        opts := config.ToScheduleOptions()
        handle, err := cli.ScheduleClient().Create(ctx, opts)
        if err != nil {
            log.Err(err).Str("Name", config.Name).Msg("failed to create scheduled task")
            continue
        }
        log.Info().Str("Name", config.Name).Str("ID", handle.GetID()).Msg("scheduled task created successfully")
    }
}

// listExisting 获取已存在的 schedule ID 列表
func (r *ScheduleRegistry) listExisting(ctx context.Context, cli client.Client) map[string]bool {
    exist := make(map[string]bool)
    list, err := cli.ScheduleClient().List(ctx, client.ScheduleListOptions{
        PageSize: 1000,
    })
    if err != nil {
        log.Err(err).Msg("failed to retrieve schedule list")
        return exist
    }
    for list.HasNext() {
        next, err := list.Next()
        if err != nil {
            break
        }
        exist[next.ID] = true
    }
    return exist
}
