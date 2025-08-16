using logistics_frontend.Models.Errors;
using logistics_frontend.Models.Order;

public class DropdownDataService
{
    private readonly OrderService _orderService;
    private readonly ToastService _toastService;

    public List<Customer> Customers { get; private set; } = new();
    public List<AllInventory> Inventories { get; private set; } = new();

    private DateTime _lastFetchTime;
    private readonly TimeSpan _cacheDuration = TimeSpan.FromMinutes(10);
    public DropdownDataService(OrderService orderService, ToastService toastService)
    {
        _orderService = orderService;
        _toastService = toastService;
    }

    public async Task<ServiceResult<bool>> LoadCachedDropdownData(bool forceRefresh = false)
    {
        if (!forceRefresh && Customers.Count > 0 && Inventories.Count > 0
            && DateTime.UtcNow - _lastFetchTime < _cacheDuration)
        {
            return ServiceResult<bool>.Ok(true);
        }

        var result = await _orderService.GetDropdownMenuData();
        if (result.Success && result.Data != null)
        {
            Customers = result.Data.Customers ?? new();
            Inventories = result.Data.Inventories ?? new();
            _lastFetchTime = DateTime.UtcNow;            
            return ServiceResult<bool>.Ok(true);
        }

        _toastService.ShowToast("Failed to load dropdown data.", ToastService.ToastLevel.Warning);
        return ServiceResult<bool>.Fail(result.ErrorMessage ?? "Failed to load dropdown data.");
    }

    public void InvalidateCache()
    {
        Customers.Clear();
        Inventories.Clear();
        _lastFetchTime = DateTime.MinValue;
    }
}