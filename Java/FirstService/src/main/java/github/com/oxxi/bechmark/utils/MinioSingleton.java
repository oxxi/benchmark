package github.com.oxxi.bechmark.utils;


import io.minio.MinioClient;
import org.springframework.stereotype.Service;

import java.net.URL;

@Service
public final class MinioSingleton {

    private static MinioSingleton instance;


    private MinioClient minioClient;
    private  MinioSingleton(){}

    public synchronized static MinioSingleton getInstance(String host,String accessKey,String secretKey, int port, boolean secure ) {
        if (instance == null) {
            instance = new MinioSingleton();
            instance.minioClient = MinioClient.builder().endpoint(host,port,secure)
                    .credentials(accessKey,secretKey).build();
        }
        return  instance;
    }

    public MinioClient getMinioInstance() {
        return  instance.minioClient;
    }

}
