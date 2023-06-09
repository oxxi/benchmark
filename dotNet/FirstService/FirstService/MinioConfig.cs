﻿namespace FirstService
{
    public class MinioConfig
    {
        public string MinioEndPoint { get; set; }
        public int MinioPort { get; set; }
        public string MinioAccessKey { get; set; }
        public string MinioSecretKey { get; set; }
        public string MinioBucket { get; set; }

        public bool MinioSSL {get; set;}
        public string Uri { get=> MinioEndPoint +":"+MinioPort; }
    }
}
