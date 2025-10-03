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

        [JsonPropertyName("inventory_id")]
        public Guid InventoryID { get; set; }

        [JsonPropertyName("pickup_address")]
        public string PickupAddress { get; set; } = string.Empty;

        [JsonPropertyName("pickup_lat")]
        public double PickupLat { get; set; }

        [JsonPropertyName("pickup_lng")]
        public double PickupLng { get; set; } 

        [JsonPropertyName("delivery_address")]
        public string DeliveryAddress { get; set; } = string.Empty;

        [JsonPropertyName("delivery_lat")]
        public double DeliveryLat { get; set; }

        [JsonPropertyName("delivery_lng")]
        public double DeliveryLng { get; set; }

        [JsonPropertyName("status")]
        [JsonConverter(typeof(JsonStringEnumConverter))]
        public OrderStatus Status { get; set; }

        [JsonPropertyName("created_at")]
        public DateTime CreatedAt { get; set; }

        [JsonPropertyName("updated_at")]
        public DateTime UpdatedAt { get; set; }

        [JsonPropertyName("category")]
        public string Category { get; set; } = string.Empty;
    }

    public class CreateOrderRequest
    {
        [Required]
        [JsonPropertyName("admin_id")]
        public Guid AdminID { get; set; }

        [Required(ErrorMessage = "Pickup location is required.")]
        [JsonPropertyName("pickup_address")]
        public string PickupAddress { get; set; } = string.Empty;

        [JsonPropertyName("pickup_lat")]
        public double PickupLat { get; set; }

        [JsonPropertyName("pickup_lng")]
        public double PickupLng { get; set; }

        [Required(ErrorMessage = "Delivery location is required.")]
        [JsonPropertyName("delivery_address")]
        public string DeliveryAddress { get; set; } = string.Empty;

        [JsonPropertyName("delivery_lat")]
        public double DeliveryLat { get; set; }

        [JsonPropertyName("delivery_lng")]
        public double DeliveryLng { get; set; }

        [Required]
        [Range(1, int.MaxValue, ErrorMessage = "Quantity must be at least 1.")]
        [JsonPropertyName("quantity")]
        public int Quantity { get; set; }

        [Required(ErrorMessage = "Customer is required.")]
        [JsonPropertyName("customer_id")]
        public Guid CustomerID { get; set; }

        [Required(ErrorMessage = "Inventory is required.")]
        [JsonPropertyName("inventory_id")]
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

        [JsonPropertyName("category")]
        public string Category { get; set; } = string.Empty;
    }
}