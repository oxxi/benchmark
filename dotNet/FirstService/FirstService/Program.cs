using FirstService;


var builder = WebApplication.CreateBuilder(args);
var config = builder.Configuration.GetSection("MinioConfig").Get<MinioConfig>() ?? new MinioConfig();
var service = builder.Configuration.GetSection("SecondService").Get<string>() ?? "";

builder.Services.AddSingleton(new ClientService(service));
builder.Services.AddSingleton(new MinioService(config));

var app = builder.Build();


app.MapGet("/", async (ClientService clientService, MinioService minioService) =>
{
    string result = await clientService.GetIdFromService();

    var response = await  minioService.UploadFile(result);

    return response;
});

app.Run();
