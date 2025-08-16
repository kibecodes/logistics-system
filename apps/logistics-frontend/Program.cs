using Microsoft.Extensions.DependencyInjection;
using Microsoft.AspNetCore.Components.Web;
using Microsoft.AspNetCore.Components.WebAssembly.Hosting;
using Microsoft.AspNetCore.Components.Authorization;
using logistics_frontend;
using logistics_frontend.Services.CustomAuthStateProvider;
using System.Text.Json;
using logistics_frontend.Services.AuthHeaderHandler;

var builder = WebAssemblyHostBuilder.CreateDefault(args);
builder.RootComponents.Add<App>("#app");
builder.RootComponents.Add<HeadOutlet>("head::after");

// Fetch appsettings.json to get API base URL
var tempHttpClient = new HttpClient
{
    BaseAddress = new Uri(builder.HostEnvironment.BaseAddress)
};

using var settingsResponse = await tempHttpClient.GetAsync("appsettings.json");
if (!settingsResponse.IsSuccessStatusCode)
{
    throw new Exception("Failed to load appsettings.json");
}

var json = await settingsResponse.Content.ReadAsStringAsync();

var config = JsonSerializer.Deserialize<Dictionary<string, string>>(json)
    ?? throw new Exception("Failed to parse appsettings.json");

if (!config.TryGetValue("ApiBaseUrl", out var apiBaseUrl) || string.IsNullOrWhiteSpace(apiBaseUrl))
    throw new Exception("ApiBaseUrl not found or invalid in appsettings.json");

Console.WriteLine($"API Base URL set to: {apiBaseUrl}");


// Register Auth Token Handler
builder.Services.AddScoped<AuthHeaderHandler>();

// Register named HttpClients
// 1. Anonymous Client (no Bearer header)
builder.Services.AddHttpClient("AnonymousApi", client =>
{
    client.BaseAddress = new Uri(apiBaseUrl);
});

// 2. Authenticated Client (Bearer token auto-injected)
builder.Services.AddHttpClient("AuthenticatedApi", client =>
{
    client.BaseAddress = new Uri(apiBaseUrl);
}).AddHttpMessageHandler<AuthHeaderHandler>();


// Register app services
builder.Services.AddScoped<UserService>();
builder.Services.AddScoped<OrderService>();
builder.Services.AddScoped<OrderService>();
builder.Services.AddScoped<DropdownDataService>();
builder.Services.AddScoped<DriverService>();
builder.Services.AddScoped<PaymentService>();
builder.Services.AddScoped<DeliveryService>();
builder.Services.AddScoped<FeedbackService>();
builder.Services.AddScoped<InventoryService>();
builder.Services.AddScoped<UserSessionService>();
builder.Services.AddScoped<NotificationService>();
builder.Services.AddScoped<ToastService>();

// Authentication & Authorization
builder.Services.AddOptions();
builder.Services.AddAuthorizationCore();
builder.Services.AddScoped<AuthenticationStateProvider, CustomAuthStateProvider>();
builder.Services.AddScoped<CustomAuthStateProvider>();

await builder.Build().RunAsync();   
