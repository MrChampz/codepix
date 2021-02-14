package kafka

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/MrChampz/codepix/application/factory"
	"github.com/MrChampz/codepix/application/usecase"
	"github.com/MrChampz/codepix/domain/model"
	"os"
	appmodel "github.com/MrChampz/codepix/application/model"
	ckafka "github.com/confluentinc/confluent-kafka-go/kafka"
)

type KafkaProcessor struct {
	Database *gorm.DB
	Producer *ckafka.Producer
	DeliveryChan chan ckafka.Event
}

func NewKafkaProcessor(
	database *gorm.DB,
	producer *ckafka.Producer,
	deliveryChan chan ckafka.Event,
) *KafkaProcessor {
	return &KafkaProcessor {
		Database: 			database,
		Producer: 			producer,
		DeliveryChan:		deliveryChan,
	}
}

func (processor *KafkaProcessor) Consume() {
	configMap := &ckafka.ConfigMap {
		"bootstrap.servers":	os.Getenv("kafkaBootstrapServers"),
		"group.id": 				  os.Getenv("kafkaConsumerGroupId"),
		"auto.offset.reset":	"earliest",
	}

	consumer, err := ckafka.NewConsumer(configMap)
	if err != nil {
		panic(err)
	}

	topics := []string {
		os.Getenv("kafkaTransactionTopic"),
		os.Getenv("kafkaTransactionConfirmationTopic"),
	}
	consumer.SubscribeTopics(topics, nil)

	fmt.Println("kafka consumer has been started")
	for {
		msg, err := consumer.ReadMessage(-1)
		if err == nil {
			processor.processMessage(msg)
		}
	}
}

func (processor *KafkaProcessor) processMessage(msg *ckafka.Message) {
	transactionsTopic := "transactions"
	transactionConfirmationTopic := "transaction_confirmation"
	
	switch topic := *msg.TopicPartition.Topic; topic {
	case transactionsTopic:
		processor.processTransaction(msg)
	case transactionConfirmationTopic:
		processor.processTransactionConfirmation(msg)
	default:
		fmt.Println("not a valid topic", string(msg.Value))
	}
}

func (processor *KafkaProcessor) processTransaction(msg *ckafka.Message) error {
	transaction := appmodel.NewTransaction()
	err := transaction.ParseJson(msg.Value)
	if err != nil {
		return err
	}

	transactionUseCase := factory.TransactionUseCaseFactory(processor.Database)

	createdTransaction, err := transactionUseCase.Register(
		transaction.ID,
		transaction.AccountID,
		transaction.Amount,
		transaction.PixKeyTo,
		transaction.PixKeyToKind,
		transaction.Description,
	)
	if err != nil {
		fmt.Println("error registering transaction", err)
		return err
	}

	topic := "bank" + createdTransaction.PixKeyTo.Account.Bank.Code
	transaction.ID = createdTransaction.ID
	transaction.Status = model.TransactionPending

	transactionJson, err := transaction.ToJson()
	if err != nil {
		return err
	}

	err = Publish(string(transactionJson), topic, processor.Producer, processor.DeliveryChan)
	if err != nil {
		return err
	}

	return nil
}

func (processor *KafkaProcessor) processTransactionConfirmation(msg *ckafka.Message) error {
	transaction := appmodel.NewTransaction()
	err := transaction.ParseJson(msg.Value)
	if err != nil {
		return err
	}

	transactionUseCase := factory.TransactionUseCaseFactory(processor.Database)

	if transaction.Status == model.TransactionConfirmed {
		err = processor.confirmTransaction(transaction, transactionUseCase)
		if err != nil {
			return err
		}
	} else if transaction.Status == model.TransactionCompleted {
		_, err := transactionUseCase.Complete(transaction.ID)
		if err != nil {
			return err
		}
	}
	return nil
}

func (processor *KafkaProcessor) confirmTransaction(
	transaction *appmodel.Transaction,
	transactionUseCase usecase.TransactionUseCase,
) error {
	confirmedTransaction, err := transactionUseCase.Confirm(transaction.ID)
	if err != nil {
		return err
	}

	topic := "bank" + confirmedTransaction.AccountFrom.Bank.Code
	
	transactionJson, err := transaction.ToJson()
	if err != nil {
		return err
	}

	err = Publish(string(transactionJson), topic, processor.Producer, processor.DeliveryChan)
	if err != nil {
		return err
	}

	return nil
}