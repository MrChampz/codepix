import {
  Body,
  Controller,
  Get,
  Inject,
  OnModuleDestroy,
  OnModuleInit,
  Param,
  ParseUUIDPipe,
  Post,
  ValidationPipe
} from '@nestjs/common';
import { InjectRepository } from "@nestjs/typeorm";
import { Repository } from "typeorm";
import { ClientKafka, MessagePattern, Payload } from "@nestjs/microservices";
import { Producer } from "kafkajs";

import { PixKey } from "../../models/pix-key.model";
import { BankAccount } from "../../models/bank-account.model";

import {
  Transaction,
  TransactionOperation,
  TransactionStatus
} from "../../models/transaction.model";
import { TransactionDto } from "../../dto/transaction.dto";

@Controller('bank-accounts/:bankAccountId/transactions')
export class TransactionController implements OnModuleInit, OnModuleDestroy {

  private kafkaProducer: Producer;

  constructor(
    @InjectRepository(BankAccount)
    private bankAccountRepo: Repository<BankAccount>,

    @InjectRepository(PixKey)
    private pixKeyRepo: Repository<PixKey>,

    @InjectRepository(Transaction)
    private transactionRepo: Repository<Transaction>,

    @Inject('TRANSACTION_SERVICE')
    private kafkaClient: ClientKafka
  ) {}

  async onModuleInit() {
    this.kafkaProducer = await this.kafkaClient.connect();
  }

  async onModuleDestroy() {
    await this.kafkaProducer.disconnect();
  }

  @Get()
  async index(
    @Param(
      'bankAccountId',
      new ParseUUIDPipe({ version: '4', errorHttpStatusCode: 422 })
    )
    bankAccountId: string,
  ) {
    await this.bankAccountRepo.findOneOrFail(bankAccountId);

    return await this.transactionRepo.find({
      where: {
        bankAccountId: bankAccountId,
      },
      order: {
        createdAt: 'DESC',
      },
    });
  }

  @Post()
  async store(
    @Param(
      'bankAccountId',
      new ParseUUIDPipe({ version: '4', errorHttpStatusCode: 422 })
    )
    bankAccountId: string,

    @Body(new ValidationPipe({ errorHttpStatusCode: 422 }))
    body: TransactionDto,
  ) {
    await this.bankAccountRepo.findOneOrFail(bankAccountId);

    let transaction = this.transactionRepo.create({
      ...body,
      amount: body.amount * -1,
      bankAccountId: bankAccountId,
      operation: TransactionOperation.debit,
    });

    transaction = await this.transactionRepo.save(transaction);

    const sendData = {
      id: transaction.externalId,
      accountId: bankAccountId,
      amount: body.amount,
      pixKeyTo: body.pixKeyKey,
      pixKeyToKind: body.pixKeyKind,
      description: body.description,
    }

    await this.kafkaProducer.send({
      topic: 'transactions',
      messages: [
        {
          key: 'transactions',
          value: JSON.stringify(sendData)
        },
      ],
    });

    return transaction;
  }

  @MessagePattern(`bank${ process.env.BANK_CODE }`)
  async processTransaction(@Payload() message) {
    if (message.value.status === TransactionStatus.pending) {
      await this.receivedTransaction(message.value);
    }

    if (message.value.status === 'confirmed') {
      await this.confirmedTransaction(message.value);
    }
  }

  async receivedTransaction(data) {
    const pixKey = await this.pixKeyRepo.findOneOrFail({
      where: {
        key: data.pixKeyTo,
        kind: data.pixKeyToKind,
      }
    });

    const transaction = this.transactionRepo.create({
      externalId: data.id,
      amount: data.amount,
      description: data.description,
      bankAccountId: pixKey.bankAccountId,
      bankAccountFromId: data.accountId,
      operation: TransactionOperation.credit,
      status: TransactionStatus.completed,
    });

    await this.transactionRepo.save(transaction)

    const sendData = {
      ...data,
      status: 'confirmed',
    }

    await this.kafkaProducer.send({
      topic: 'transaction_confirmation',
      messages: [
        {
          key: 'transaction_confirmation',
          value: JSON.stringify(sendData),
        }
      ],
    });
  }

  async confirmedTransaction(data) {
    const transaction = await this.transactionRepo.findOneOrFail({
      where: {
        externalId: data.id,
      }
    });

    await this.transactionRepo.update(
      { externalId: data.id },
      {
        status: TransactionStatus.completed,
      },
    );

    const sendData = {
      id: data.id,
      accountId: transaction.bankAccountId,
      amount: Math.abs(transaction.amount),
      pixKeyTo: transaction.pixKeyKey,
      pixKeyToKind: transaction.pixKeyKind,
      description: transaction.description,
      status: TransactionStatus.completed,
    }

    await this.kafkaProducer.send({
      topic: 'transaction_confirmation',
      messages: [
        {
          key: 'transaction_confirmation',
          value: JSON.stringify(sendData),
        },
      ],
    });
  }
}
