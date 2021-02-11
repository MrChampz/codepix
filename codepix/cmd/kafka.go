package cmd

import (
	ckafka "github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/MrChampz/codepix/application/kafka"
	"github.com/MrChampz/codepix/infrastructure/db"
	"github.com/spf13/cobra"
	"os"
)

// kafkaCmd represents the kafka command
var kafkaCmd = &cobra.Command{
	Use:   "kafka",
	Short: "Start consuming transactions using Apache Kafka",
	Run: func(cmd *cobra.Command, args []string) {
		database := db.ConnectDB(os.Getenv("env"))
		deliveryChan := make(chan ckafka.Event)

		producer := kafka.NewKafkaProducer()
		go kafka.DeliveryReport(deliveryChan)

		processor := kafka.NewKafkaProcessor(database, producer, deliveryChan)
		processor.Consume()
	},
}

func init() {
	rootCmd.AddCommand(kafkaCmd)
}
