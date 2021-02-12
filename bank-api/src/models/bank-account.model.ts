import {
  Entity,
  Column,
  CreateDateColumn,
  PrimaryGeneratedColumn,
  BeforeInsert
} from "typeorm";

import { v4 as uuidv4 } from "uuid";

@Entity({
  name: 'bank_accounts'
})
export class BankAccount {

  @PrimaryGeneratedColumn('uuid')
  id: string;

  @Column({ name: 'account_number' })
  accountNumber: string;

  @Column({ name: 'owner_name' })
  ownerName: string;

  @Column()
  balance: number;

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
  initBalance() {
    if (this.balance) {
      return;
    }

    this.balance = 0;
  }
}
