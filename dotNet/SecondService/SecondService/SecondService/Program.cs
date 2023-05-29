var builder = WebApplication.CreateBuilder(args);
var app = builder.Build();



app.MapGet("/", () =>
{
    Guid uuid = Guid.NewGuid();

    return uuid.ToString();
});

app.Run();
