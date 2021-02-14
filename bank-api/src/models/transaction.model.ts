import {
  BeforeInsert,
  Column,
  CreateDateColumn,
  Entity,
  JoinColumn,
  ManyToOne,
  PrimaryGeneratedColumn
} from "typeorm";
import { v4 as uuidv4 } from "uuid";

import { BankAccount } from "./bank-account.model";

export enum TransactionStatus {
  pending = 'pending',
  completed = 'completed',
  error = 'error',
}

export enum TransactionOperation {
  debit = 'debit',
  credit = 'credit',
}

@Entity({ name: 'transactions' })
export class Transaction {

  @PrimaryGeneratedColumn('uuid')
  id: string;

  @Column({ name: 'external_id' })
  externalId: string;

  @Column()
  amount: number;

  @Column()
  description: string;

  @ManyToOne(() => BankAccount)
  @JoinColumn({ name: 'bank_account_id' })
  bankAccount: BankAccount;

  @Column({ name: 'bank_account_id' })
  bankAccountId: string;

  @ManyToOne(() => BankAccount)
  @JoinColumn({ name: 'bank_account_from_id' })
  bankAccountFrom: BankAccount;

  @Column({ name: 'bank_account_from_id' })
  bankAccountFromId: string;

  @Column({ name: 'pix_key_key' })
  pixKeyKey: string;

  @Column({ name: 'pix_key_kind' })
  pixKeyKind: string;

  @Column()
  status: TransactionStatus = TransactionStatus.pending;

  @Column()
  operation: TransactionOperation

  @CreateDateColumn({ name: 'created_at', type: 'timestamp' })
  createdAt: Date;

  @BeforeInsert()
  generateId() {
    if (this.id) {
      return;
    }
    this.id = uuidv4();
  }

  @BeforeInsert()
  generateExternalId() {
    if (this.externalId) {
      return;
    }
    this.externalId = uuidv4();
  }
}
