import { Controller, Get, HttpCode } from '@nestjs/common';
import { AppService } from './app.service';

@Controller('/')
export class AppController {
  constructor(private readonly appService: AppService) {}

  @Get()
  @HttpCode(200)
  async writeFile() {
    //call second service
    const str = await this.appService.getTokenByExternalService();

    //write in minio
    const info = this.appService.writeFileInMinio(str);
    return info;
  }
}
