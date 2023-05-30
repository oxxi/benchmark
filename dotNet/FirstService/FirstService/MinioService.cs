using Microsoft.Extensions.Options;
using Minio;

using System.Text;


namespace FirstService
{
    public class MinioService
    {
        private readonly MinioClient _client;
        private readonly MinioConfig _config;
        
        public MinioService(MinioConfig options)
        {
            _config = options;
            _client = new MinioClient().WithEndpoint(_config.Uri).WithCredentials(_config.MinioAccessKey, _config.MinioSecretKey).WithSSL(_config.MinioSSL).Build();
            
        }




        public async Task<object> UploadFile(string id)
        {
            //create file
            string fileName = string.Format("dotNet{0}.txt", id);

            var file =  new MemoryStream(Encoding.UTF8.GetBytes(id));
            PutObjectArgs putObjet = new PutObjectArgs().WithBucket(_config.MinioBucket)
                                              .WithObject(fileName)
                                              .WithStreamData(file)
                                              .WithObjectSize(file.Length)
                                              .WithContentType("application/octet-stream");

            //instance minio
            var minioResponse = await _client.PutObjectAsync(putObjet);
            


            return new { etag = minioResponse.Etag, varsionId = "" };

        }
        
}


    
}
