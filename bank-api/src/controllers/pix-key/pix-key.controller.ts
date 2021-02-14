import {
  Body,
  Controller,
  Get, HttpCode,
  Inject,
  InternalServerErrorException, NotFoundException,
  Param,
  ParseUUIDPipe,
  Post, Query,
  UnprocessableEntityException,
  ValidationPipe
} from '@nestjs/common';

import { Repository } from "typeorm";
import { InjectRepository } from "@nestjs/typeorm";
import { ClientGrpc } from "@nestjs/microservices";

import { PixKey } from "../../models/pix-key.model";
import { PixKeyDto } from "../../dto/pix-key.dto";
import { BankAccount } from "../../models/bank-account.model";
import { PixService } from "../../grpc-types/pix-service.grpc";
import { PixKeyExistsDto } from "../../dto/pix-key-exists.dto";

@Controller('bank-accounts/:bankAccountId/pix-keys')
export class PixKeyController {
  constructor(
    @InjectRepository(PixKey)
    private pixKeyRepo: Repository<PixKey>,

    @InjectRepository(BankAccount)
    private bankAccountRepo: Repository<BankAccount>,

    @Inject('CODEPIX_PACKAGE')
    private client: ClientGrpc
  ) {}

  @Get()
  index(
    @Param('bankAccountId', new ParseUUIDPipe({ version: '4' }))
    bankAccountId: string,
  ) {
    return this.pixKeyRepo.find({
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
    @Param('bankAccountId', new ParseUUIDPipe({ version: '4' }))
    bankAccountId: string,
    @Body(new ValidationPipe({ errorHttpStatusCode: 422 }))
    body: PixKeyDto
  ) {
    await this.bankAccountRepo.findOneOrFail(bankAccountId);

    const pixService: PixService = this.client.getService('PixService');
    const notFound = this.checkPixKeyNotFound(body);
    if (!notFound) {
      throw new UnprocessableEntityException("PixKey already exists");
    }

    const createdPixKey = await pixService.registerPixKey({
      ...body,
      accountId: bankAccountId,
    }).toPromise();

    if (createdPixKey.error) {
      throw new InternalServerErrorException(createdPixKey.error);
    }

    const pixKey = this.pixKeyRepo.create({
      id: createdPixKey.id,
      bankAccountId: bankAccountId,
      ...body
    });

    return await this.pixKeyRepo.save(pixKey);
  }

  @Get('exists')
  @HttpCode(204)
  async exists(
    @Query(new ValidationPipe({ errorHttpStatusCode: 422 }))
    params: PixKeyExistsDto
  ) {
    const pixService: PixService = this.client.getService('PixService');
    try {
      await pixService.find(params).toPromise();
    } catch (e) {
      if (e.details == "no key was found") {
        throw new NotFoundException(e.details);
      }

      throw new InternalServerErrorException("Server not available");
    }
  }

  private async checkPixKeyNotFound(params: { key: string, kind: string }) {
    const pixService: PixService = this.client.getService('PixService');
    try {
      await pixService.find(params).toPromise();
      return false;
    } catch (e) {
      if (e.details == "no key was found") {
        return true;
      }

      throw new InternalServerErrorException("Server not available");
    }
  }
}
