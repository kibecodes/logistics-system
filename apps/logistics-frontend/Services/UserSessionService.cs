using System.Text.Json;
using logistics_frontend.Models.User;
using Microsoft.JSInterop;

public class UserSessionService
{
    private readonly IJSRuntime _jSRuntime;

    public UserSessionService(IJSRuntime jSRuntime)
    {
        _jSRuntime = jSRuntime;
    }

    public async Task<string?> GetUserRoleAsync()
    {
        var json = await _jSRuntime.InvokeAsync<string>("localStorage.getItem", "auth_user");
        if (string.IsNullOrWhiteSpace(json)) return null;

        var user = JsonSerializer.Deserialize<User>(json, new JsonSerializerOptions
        {
            PropertyNameCaseInsensitive = true
        });

        return user?.Role;
    }

    public async Task<User?> GetUserAsync()
    {
        var json = await _jSRuntime.InvokeAsync<string>("localStorage.getItem", "auth_user");

        return string.IsNullOrWhiteSpace(json)
        ? null
        : JsonSerializer.Deserialize<User>(json, new JsonSerializerOptions
        {
            PropertyNameCaseInsensitive = true
        });
    }
    public async Task<Guid?> GetUserIdAsync()
    {
        var user = await GetUserAsync();
        return user?.ID;
    }

    public async Task<string?> GetTokenAsync()
    {
        var token = await _jSRuntime.InvokeAsync<string>("localStorage.getItem", "auth_token");
        return string.IsNullOrWhiteSpace(token) ? null : token;
    }
}