package kitflow

import (
	"context"
	"time"

	"github.com/rs/zerolog/log"
	"go.temporal.io/api/workflowservice/v1"
	"go.temporal.io/sdk/client"
	"google.golang.org/protobuf/types/known/durationpb"
)

// CreateNamespace 创建 Temporal 命名空间
// 如果 namespace 为 "default" 或空字符串，则跳过创建
// retention 为工作流执行保留时间，传 0 则默认 7 天
func CreateNamespace(hostPort, namespace string, retention time.Duration) error {
	if namespace == "default" || namespace == "" {
		log.Warn().Msg("namespace is default or empty, skip creation")
		return nil
	}

	namespaceClient, err := client.NewNamespaceClient(client.Options{HostPort: hostPort})
	if err != nil {
		return err
	}

	if retention <= 0 {
		retention = time.Hour * 24 * 7
	}
	err = namespaceClient.Register(context.Background(), &workflowservice.RegisterNamespaceRequest{
		Namespace:                        namespace,
		WorkflowExecutionRetentionPeriod: durationpb.New(retention),
	})
	if err != nil {
		return err
	}
	log.Info().Msg("created successfully")
	return nil
}
