import { HttpService } from '@nestjs/axios';

import { Injectable } from '@nestjs/common';
import { ConfigService } from '@nestjs/config';
import { readFileSync, unlinkSync, writeFileSync } from 'fs';
import { MinioService } from 'nestjs-minio-client';
import { lastValueFrom, map } from 'rxjs';

@Injectable()
export class AppService {
  constructor(
    private readonly httpService: HttpService,
    private readonly configService: ConfigService,
    private readonly minioService: MinioService,
  ) {}

  async getTokenByExternalService(): Promise<string> {
    const value = this.httpService
      .get(this.configService.get('SECOND_SERVICE'))
      .pipe(
        map((resp) => {
          return resp.data;
        }),
      );

    const result = await lastValueFrom(value);

    return result;
  }

  async writeFileInMinio(srt: string) {
    const name = `newFile${srt}.txt`;

    writeFileSync(name, srt);

    const doc = readFileSync(name);
    const info = await this.minioService.client.putObject(
      this.configService.get('MINIO_BUCKET_NAME'),
      name,
      doc,
    );
    await unlinkSync(name);
    return info;
  }
}
