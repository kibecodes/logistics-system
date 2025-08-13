using System.Net.Http.Json;
using System.Text.Json;
using logistics_frontend.Models.Errors;
using logistics_frontend.Models.Inventory;

public class InventoryService
{
    private readonly HttpClient _http;
    public InventoryService(IHttpClientFactory httpClientFactory)
    {
        _http = httpClientFactory.CreateClient("AuthenticatedApi");
    }

    public async Task<ServiceResult<HttpResponseMessage>> AddInventory(CreateInventoryRequest inventory)
    {
        try
        {
            var response = await _http.PostAsJsonAsync("inventories/create", inventory);
            if (response.IsSuccessStatusCode)
            {
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

    public async Task<ServiceResult<Inventory>> GetInventoryByID(Guid inventory_id)
    {
        return await GetFromJsonSafe<Inventory>($"inventories/by-id/{inventory_id}");
    }

    public async Task<ServiceResult<List<Inventory>>> GetInventoriesByName(string name)
    {
        var encodedName = Uri.EscapeDataString(name);
        return await GetFromJsonSafe<List<Inventory>>($"inventories/by-name?name={encodedName}");
    }

    public async Task<ServiceResult<List<Inventory>>> GetAllInventories()
    {
        return await GetFromJsonSafe<List<Inventory>>("inventories/all_inventories?limit=10&offset=0");
    }

    public async Task<ServiceResult<List<Inventory>>> GetInventoriesByCategory(string category)
    {
        var encodedCategory = Uri.EscapeDataString(category);
        return await GetFromJsonSafe<List<Inventory>>($"inventories/by-category?category={encodedCategory}");
    }

    public async Task<ServiceResult<List<string>>> GetCategories()
    {
        return await GetFromJsonSafe<List<string>>("inventories/categories");
    }

    public async Task<bool> DeleteInventory(Guid id)
    {
        var res = await _http.DeleteAsync($"inventories/{id}");
        return res.IsSuccessStatusCode;
    }


    // Generic method for GET + JSON
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

}