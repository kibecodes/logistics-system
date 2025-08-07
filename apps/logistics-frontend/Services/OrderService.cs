using System.Net.Http.Json;
using logistics_frontend.Models.Order;
public class OrderService
{
    private readonly HttpClient _http;
    public OrderService(IHttpClientFactory httpClientFactory)
    {
        _http = httpClientFactory.CreateClient("AuthenticatedApi");
    }

    public async Task AddOrder(CreateOrderRequest order)
    {
        var response = await _http.PostAsJsonAsync("orders/create", order);
        response.EnsureSuccessStatusCode();
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

    public async Task<List<Order>> GetAllOrders()
    {
        var orders = await _http.GetFromJsonAsync<List<Order>>("orders/all_orders");
        return orders ?? new List<Order>();
    }

    public async Task<bool> DeleteOrder(Guid id)
    {
        var res = await _http.DeleteAsync($"orders/{id}");
        return res.IsSuccessStatusCode;
    }
}
