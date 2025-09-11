using System.ComponentModel.DataAnnotations;
using System.Text.Json.Serialization;

namespace logistics_frontend.Models.Order
{
    public class Order
    {
        [JsonPropertyName("id")]
        public Guid ID { get; set; }

        [JsonPropertyName("quantity")]
        public int Quantity { get; set; }

        [JsonPropertyName("customer_id")]
        public Guid CustomerID { get; set; }

        [JsonPropertyName("pickup_location")]
        public string PickupLocation { get; set; } = string.Empty;

        [JsonPropertyName("delivery_location")]
        public string DeliveryLocation { get; set; } = string.Empty;

        [JsonPropertyName("order_status")]
        [JsonConverter(typeof(JsonStringEnumConverter))]
        public OrderStatus OrderStatus { get; set; }

        [JsonPropertyName("created_at")]
        public DateTime CreatedAt { get; set; }

        [JsonPropertyName("updated_at")]
        public DateTime UpdatedAt { get; set; }
    }

    public class CreateOrderRequest
    {
        [Required]
        public Guid AdminID { get; set; }

        [Required(ErrorMessage = "Pickup location is required.")]
        public string PickupLocation { get; set; } = string.Empty;

        [Required(ErrorMessage = "Delivery location is required.")]
        public string DeliveryLocation { get; set; } = string.Empty;

        [Required]
        [Range(1, int.MaxValue, ErrorMessage = "Quantity must be at least 1.")]
        public int Quantity { get; set; }

        [Required(ErrorMessage = "Customer is required.")]
        public Guid CustomerID { get; set; }

        [Required(ErrorMessage = "Inventory is required.")]
        public Guid InventoryID { get; set; }
    }

    public enum OrderStatus
    {
        Pending,
        Assigned,
        InTransit,
        Delivered,
        Cancelled
    }

    public class DropdownData
    {
        public List<Customer> Customers { get; set; } = new();
        public List<AllInventory> Inventories { get; set; } = new();
        
    }

    public class Customer
    {
        [JsonPropertyName("id")]
        public Guid ID { get; set; }

        [JsonPropertyName("name")]
        public string Name { get; set; } = string.Empty;
    }

    public class AllInventory
    {
        [JsonPropertyName("id")]
        public Guid ID { get; set; }

        [JsonPropertyName("name")]
        public string Name { get; set; } = string.Empty;

        [JsonPropertyName("admin_id")]
        public Guid AdminID { get; set; }
    }
}