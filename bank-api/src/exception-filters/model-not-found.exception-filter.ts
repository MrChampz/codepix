import { ArgumentsHost, ExceptionFilter } from "@nestjs/common";
import { EntityNotFoundError } from "typeorm";
// @ts-ignore
import { Response } from "express";

export class ModelNotFoundExceptionFilter implements ExceptionFilter {
  catch(exception: EntityNotFoundError, host: ArgumentsHost) {
    const context = host.switchToHttp();
    const response = context.getResponse<Response>();
    return response.status(404).json({
      error: {
        error: "Not Found",
        message: exception.message,
      }
    });
  }
}