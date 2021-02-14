import { IsNotEmpty, IsNumber, IsOptional, IsString, Min } from "class-validator";

export class TransactionDto {
  @IsString()
  @IsNotEmpty()
  pixKeyKey: string;

  @IsString()
  @IsNotEmpty()
  pixKeyKind: string;

  @IsString()
  @IsOptional()
  description: string = null;

  @IsNumber({ maxDecimalPlaces: 2 })
  @Min(0.01)
  @IsNotEmpty()
  readonly amount: number;
}