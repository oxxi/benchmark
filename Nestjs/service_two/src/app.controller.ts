import { Controller, Get, HttpCode } from '@nestjs/common';
import { v4 as uuidv4 } from 'uuid';
import { AppService } from './app.service';

@Controller('/')
export class AppController {
  constructor(private readonly appService: AppService) {}

  @Get()
  @HttpCode(200)
  generate(): string {
    const uuid = uuidv4();

    return uuid;
  }
}
