import { Module } from '@nestjs/common';
import { ConfigModule } from "@nestjs/config";
import { ConsoleModule } from "nestjs-console";
import { TypeOrmModule } from "@nestjs/typeorm";
import { ClientsModule, Transport } from "@nestjs/microservices";

import { join } from 'path';

import { FixturesCommand } from "./fixtures/fixtures.command";
import { AppController } from './app.controller';
import { AppService } from './app.service';

import { BankAccount } from "./models/bank-account.model";
import { PixKey } from "./models/pix-key.model";

import { MyFirstController } from "./controllers/my-first/my-first.controller";
import { BankAccountController } from './controllers/bank-account/bank-account.controller';
import { PixKeyController } from './controllers/pix-key/pix-key.controller';

@Module({
  imports: [
    ConfigModule.forRoot(),
    ConsoleModule,
    TypeOrmModule.forRoot({
      type: process.env.TYPEORM_CONNECTION as any,
      host: process.env.TYPEORM_HOST,
      port: parseInt(process.env.TYPEORM_PORT),
      username: process.env.TYPEORM_USERNAME,
      password: process.env.TYPEORM_PASSWORD,
      database: process.env.TYPEORM_DATABASE,
      entities: [BankAccount, PixKey]
    }),
    ClientsModule.register([
      {
        name: 'CODEPIX_PACKAGE',
        transport: Transport.GRPC,
        options: {
          url: process.env.GRPC_URL,
          package: 'github.com.mrchampz.codepix',
          protoPath: [join(__dirname, 'protofiles/pixKey.proto')]
        },
      }
    ]),
    TypeOrmModule.forFeature([BankAccount, PixKey])
  ],
  controllers: [AppController, MyFirstController, BankAccountController, PixKeyController],
  providers: [AppService, FixturesCommand],
})
export class AppModule {}
