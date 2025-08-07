using logistics_frontend.Models.User;
using logistics_frontend.Services.AuthHeaderHandler;
using System.Net.Http.Json;

public class UserService
{
    private readonly HttpClient _http;
    public UserService(HttpClient http)
    {
        _http = http;

    }

    public async Task<User> GetUserByID(Guid id)
    {
        var user = await _http.GetFromJsonAsync<User>($"users/by-id/{id}");
        return user ?? new User();
    }
    public async Task<List<User>> GetAllUsers()
    {
        var users = await _http.GetFromJsonAsync<List<User>>("users/all_users");
        return users ?? new List<User>();
    }
}
