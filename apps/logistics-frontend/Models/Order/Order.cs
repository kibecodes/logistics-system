using System.ComponentModel.DataAnnotations;

namespace logistics_frontend.Models.Order
{
    public class Order 
    {
        public Guid ID { get; set; }
        public Guid CustomerID { get; set; }
        public string PickupLocation { get; set; } = string.Empty;
        public string DeliveryLocation { get; set; } = string.Empty;
        public string OrderStatus { get; set; } = "pending";
        public DateTime CreatedAt { get; set; }
        public DateTime UpdatedAt { get; set; }
    }

    public class CreateOrderRequest 
    {
        [Required]
        public Guid AdminID { get; set; }
        [Required]
        public Guid CustomerID { get; set; }
        [Required]
        public string PickupLocation { get; set; } = string.Empty;

        [Required]
        public string DeliveryLocation { get; set; } = string.Empty;
    }
}