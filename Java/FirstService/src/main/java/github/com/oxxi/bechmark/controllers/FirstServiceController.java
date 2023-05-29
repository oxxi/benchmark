package github.com.oxxi.bechmark.controllers;


import github.com.oxxi.bechmark.utils.MinioSingleton;
import io.minio.MinioClient;
import io.minio.ObjectWriteResponse;
import io.minio.PutObjectArgs;
import io.minio.errors.*;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.http.ResponseEntity;
import org.springframework.web.bind.annotation.CrossOrigin;
import org.springframework.web.bind.annotation.GetMapping;
import org.springframework.web.bind.annotation.RequestMapping;
import org.springframework.web.bind.annotation.RestController;

import java.io.ByteArrayInputStream;
import java.io.IOException;
import java.io.InputStream;
import java.net.MalformedURLException;
import java.net.URI;
import java.net.URISyntaxException;
import java.net.http.HttpClient;
import java.net.http.HttpRequest;
import java.net.http.HttpResponse;
import java.security.InvalidKeyException;
import java.security.NoSuchAlgorithmException;
import java.util.concurrent.CompletableFuture;
import java.util.concurrent.ExecutionException;



@RequestMapping("/")
@RestController
@CrossOrigin("*")
public class FirstServiceController {

    @Value("${MINIO_PORT}")
    private int port;
    @Value("${MINIO_ACCESS_KEY}")
    private String accessKey;
    @Value("${MINIO_SECRET_KEY}")
    private String secretKey;

    @Value("${MINIO_BUCKET_NAME}")
    private String bucket;

    @Value("${SECOND_SERVICE}")
    private String endPoint;
    HttpClient client;

    public FirstServiceController()  {
        client = HttpClient.newHttpClient();

    }

    @GetMapping()
    public ResponseEntity Get() throws Exception {

        String uuid = getId();

       CompletableFuture<ResponseMinio> response =  CompletableFuture.supplyAsync(()->{
            return UploadToMinio(uuid);
        });
        ResponseMinio resp =  response.get();
        return ResponseEntity.ok(resp);
    }


    private String getId() throws Exception {
        //call second service

        try{ //avoid blocking using CompletableFuture
            CompletableFuture<HttpResponse<String>> response =  this.client.sendAsync(HttpRequest.newBuilder(new URI(endPoint)).GET().build(),
                    HttpResponse.BodyHandlers.ofString());

            String body = response.get().body();
            return body;
        }catch (URISyntaxException | ExecutionException | InterruptedException e) {
            throw new Exception("error al llamar servicio 2");
        }
    }


    private ResponseMinio UploadToMinio(String id) {
        String name = String.format("javafile%s.txt",id);

        InputStream inputStream = new ByteArrayInputStream(id.getBytes());

        PutObjectArgs args = PutObjectArgs.builder().bucket(bucket)
                                                    .object(name)
                                                    .stream(inputStream,-1, 1024 * 1024 * 5 )
                                                    .contentType("application/octet-stream")
                                                    .build();
        try{
            MinioClient minioClient = MinioSingleton.getInstance(accessKey, secretKey, port).getMinioInstance();
            ObjectWriteResponse response = minioClient.putObject(args);
            ResponseMinio resp = ResponseMinio.builder().Etag(response.etag()).versionId(response.versionId()).build();
            return resp;
        }catch (MinioException | IOException | NoSuchAlgorithmException | InvalidKeyException e) {
           throw new RuntimeException("Error al subir archivo a minio");
        }

    }
}

