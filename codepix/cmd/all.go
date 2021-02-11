package cmd

import (
	"github.com/MrChampz/codepix/application/grpc"
	"github.com/MrChampz/codepix/application/kafka"
	"github.com/MrChampz/codepix/infrastructure/db"
	"github.com/spf13/cobra"
	"os"
	ckafka "github.com/confluentinc/confluent-kafka-go/kafka"
)

var grpcPort int

// allCmd represents the all command
var allCmd = &cobra.Command{
	Use:   "all",
	Short: "Run whole codepix application",
	Run: func(cmd *cobra.Command, args []string) {
		database := db.ConnectDB(os.Getenv("env"))
		deliveryChan := make(chan ckafka.Event)
		
		go grpc.StartGrpcServer(database, grpcPort)

		producer := kafka.NewKafkaProducer()
		go kafka.DeliveryReport(deliveryChan)

		processor := kafka.NewKafkaProcessor(database, producer, deliveryChan)
		processor.Consume()
	},
}

func init() {
	rootCmd.AddCommand(allCmd)
	grpcCmd.Flags().IntVarP(&grpcPort, "grpc-port", "g", 50051, "gRPC server port")
}
