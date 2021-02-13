import {
  BeforeInsert,
  Column,
  CreateDateColumn,
  Entity,
  JoinColumn,
  ManyToMany,
  PrimaryGeneratedColumn
} from "typeorm";
import { v4 as uuidv4 } from "uuid";

import { BankAccount } from "./bank-account.model";

export enum PixKeyKind {
  cpf = "cpf",
  email = "email",
}

@Entity({ name: 'pix_keys' })
export class PixKey {

  @PrimaryGeneratedColumn('uuid')
  id: string;

  @Column()
  kind: PixKeyKind;

  @Column()
  key: string;

  @ManyToMany(() => BankAccount)
  @JoinColumn({ name: 'bankAccountId' })
  bankAccount: BankAccount;

  @Column({ name: 'bank_account_id' })
  bankAccountId: string;

  @CreateDateColumn({ name: 'created_at', type: 'timestamp' })
  createdAt: Date;

  @BeforeInsert()
  generateId() {
    if (this.id) {
      return;
    }
    this.id = uuidv4();
  }
}
