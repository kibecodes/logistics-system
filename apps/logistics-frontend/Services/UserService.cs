using logistics_frontend.Models.Errors;
using logistics_frontend.Models.User;
using System.Net.Http.Json;
using System.Text.Json;
public class UserService
{
    private readonly HttpClient _http;
    private readonly DropdownDataService _dropdownService;
    private readonly ToastService _toastService;
    private List<User>? _cachedUsers;
    private DateTime _lastFetchTime;
    private readonly TimeSpan _cacheDuration = TimeSpan.FromMinutes(5);
    public UserService(IHttpClientFactory httpClientFactory, DropdownDataService dropdownService, ToastService toastService)
    {
        _http = httpClientFactory.CreateClient("AuthenticatedApi");;
        _dropdownService = dropdownService;
        _toastService = toastService;
    }

    public async Task<ServiceResult<HttpResponseMessage>> AddUser(CreateUserRequest user)
    {
        try
        {
            var response = await _http.PostAsJsonAsync("public/create", user);
            if (response.IsSuccessStatusCode)
            {
                InvalidateCache();
                _dropdownService.InvalidateCache();
                return ServiceResult<HttpResponseMessage>.Ok(response);
            }

            var error = await ParseError(response);
            return ServiceResult<HttpResponseMessage>.Fail(error);
        }
        catch (HttpRequestException ex)
        {
            return ServiceResult<HttpResponseMessage>.Fail($"Network error: {ex.Message}");
        }
        catch (Exception ex)
        {
            return ServiceResult<HttpResponseMessage>.Fail($"Unexpected error: {ex.Message}");
        }
    }

    public async Task<User> GetUserByID(Guid id)
    {
        var user = await _http.GetFromJsonAsync<User>($"users/by-id/{id}");
        return user ?? new User();
    }
    public async Task<ServiceResult<List<User>>> GetAllUsers()
    {
        return await GetFromJsonSafe<List<User>>("users/all_users");
    }

    public async Task<User> UpdateUser(Guid userId, string column, object value)
    {
        var requestBody = new
        {
            column,
            value
        };

        var response = await _http.PutAsJsonAsync($"users/{userId}/update", requestBody);
        if (response.IsSuccessStatusCode)
        {
            InvalidateCache();
            return await response.Content.ReadFromJsonAsync<User>() ?? new User();
        }

        return null;

    }

    public async Task<ServiceResult<List<User>>> GetAllCachedUsers(bool forceRefresh = false)
    {
        if (!forceRefresh && _cachedUsers != null && DateTime.UtcNow - _lastFetchTime < _cacheDuration)
        {
            return ServiceResult<List<User>>.Ok(_cachedUsers, fromCache: true);
        }

        var result = await GetAllUsers();
        if (result.Success)
        {
            _cachedUsers = result.Data;
            _lastFetchTime = DateTime.UtcNow;

            _toastService.ShowToast("Users fetched successfully.", ToastService.ToastLevel.Success);
        }
        else
        {
            _toastService.ShowToast("Failed to load users.", ToastService.ToastLevel.Error);
            Console.WriteLine($"error: {result.ErrorMessage}");
        }

        return result;
    }

    public void InvalidateCache()
    {
        _cachedUsers = null;
    }

    public async Task<bool> DeleteUser(Guid id)
    {
        var res = await _http.DeleteAsync($"users/{id}");
        if (res.IsSuccessStatusCode)
        {
            InvalidateCache();
        }
        return res.IsSuccessStatusCode;
    }
    

    public async Task<string> ParseError(HttpResponseMessage response)
    {
        try
        {
            var json = await response.Content.ReadAsStringAsync();
            var error = JsonSerializer.Deserialize<ErrorResponse>(json, new JsonSerializerOptions
            {
                PropertyNameCaseInsensitive = true
            });

            return error?.Detail ?? "Unknown error occurred.";
        }
        catch
        {
            return $"HTTP {(int)response.StatusCode} - {response.ReasonPhrase}";
        }
    }
    private async Task<ServiceResult<T>> GetFromJsonSafe<T>(string url)
    {
        try
        {
            var response = await _http.GetAsync(url);

            if (response.IsSuccessStatusCode)
            {
                var result = await response.Content.ReadFromJsonAsync<T>();
                return ServiceResult<T>.Ok(result ?? Activator.CreateInstance<T>());
            }

            var error = await ParseError(response);
            return ServiceResult<T>.Fail(error);
        }
        catch (HttpRequestException ex)
        {
            return ServiceResult<T>.Fail($"Network error: {ex.Message}");
        }
        catch (Exception ex)
        {
            return ServiceResult<T>.Fail($"Unexpected error: {ex.Message}");
        }
    }
}
