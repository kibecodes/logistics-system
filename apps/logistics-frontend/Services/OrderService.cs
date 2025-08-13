using System.Net.Http.Json;
using System.Text.Json;
using logistics_frontend.Models.Order;
using logistics_frontend.Models.Errors;
public class OrderService
{
    private readonly HttpClient _http;
    public OrderService(IHttpClientFactory httpClientFactory)
    {
        _http = httpClientFactory.CreateClient("AuthenticatedApi");
    }

    public async Task<ServiceResult<HttpResponseMessage>> AddOrder(CreateOrderRequest order)
    {
        try
        {
            var response = await _http.PostAsJsonAsync("orders/create", order);
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

    public async Task<List<Order>> GetOrderByID(Guid id)
    {
        var order = await _http.GetFromJsonAsync<List<Order>>($"orders/{id}");
        return order ?? new List<Order>();
    }


    public async Task<List<Order>> GetOrdersByCustomer(Guid customerId)
    {
        var orders = await _http.GetFromJsonAsync<List<Order>>($"orders/{customerId}");
        return orders ?? new List<Order>();
    }

    public async Task<ServiceResult<DropdownData>> GetDropdownMenuData()
    {
        return await GetFromJsonSafe<DropdownData>("orders/form-data");
    }
    
    public async Task<Order> UpdateOrder(Guid orderId, string column, object value)
    {
        var requestBody = new
        {
            column,
            value
        };

        var response = await _http.PutAsJsonAsync($"orders/{orderId}/update", requestBody);
        if (response.IsSuccessStatusCode)
        {
            return await response.Content.ReadFromJsonAsync<Order>() ?? new Order();
        }

        return null;

    }

    public async Task<ServiceResult<List<Order>>> GetAllOrders()
    {
        return await GetFromJsonSafe<List<Order>>("orders/all_orders");
    }

    public async Task<bool> DeleteOrder(Guid id)
    {
        var res = await _http.DeleteAsync($"orders/{id}");
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

    // Cache service
    public class DropdownDataService
    {
        public List<Customer> Customers { get; private set; } = new();
        public List<AllInventory> Inventories { get; private set; } = new();

        private readonly OrderService _orderService;
        public DropdownDataService(OrderService orderService)
        {
            _orderService = orderService;
        }

        public async Task LoadIfEmptyAsync()
        {
            if (Customers.Count == 0 || Inventories.Count == 0)
            {
                var result = await _orderService.GetDropdownMenuData();
                if (result.Success && result.Data != null)
                {
                    Customers = result.Data.Customers ?? new();
                    Inventories = result.Data.Inventories ?? new();
                }
            }
        }
    }
}

