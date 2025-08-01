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
    public async Task<List<User>> GetAllUsers()
    {
        var users = await _http.GetFromJsonAsync<List<User>>("users/all_users");
        return users ?? new List<User>();
    }
}
