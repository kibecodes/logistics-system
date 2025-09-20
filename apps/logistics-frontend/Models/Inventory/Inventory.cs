using System.ComponentModel.DataAnnotations;
using System.Text.Json.Serialization;

namespace logistics_frontend.Models.Inventory
{
    public class Inventory
    {
        [JsonPropertyName("id")]
        public Guid ID { get; set; }

        [JsonPropertyName("admin_id")]
        public Guid AdminID { get; set; }

        [JsonPropertyName("name")]
        public string Name { get; set; } = string.Empty;

        [JsonPropertyName("slug")]
        public string? Slug { get; set; }

        [JsonPropertyName("category")]
        public string Category { get; set; } = string.Empty;

        [JsonPropertyName("stock")]
        public int Stock { get; set; }

        [JsonPropertyName("price_amount")]
        public long PriceAmount { get; set; }

        [JsonPropertyName("price_currency")]
        public string PriceCurrency { get; set; } = string.Empty;

        [JsonPropertyName("images")]
        public string Images { get; set; } = string.Empty;

        [JsonPropertyName("unit")]
        public string Unit { get; set; } = string.Empty;

        [JsonPropertyName("packaging")]
        public string Packaging { get; set; } = string.Empty;

        [JsonPropertyName("description")]
        public string Description { get; set; } = string.Empty;

        [JsonPropertyName("location")]
        public string Location { get; set; } = string.Empty;

        [JsonPropertyName("created_at")]
        public DateTime CreatedAt { get; set; }

        [JsonPropertyName("updated_at")]
        public DateTime UpdatedAt { get; set; }

    }

    public class CreateInventoryRequest
    {
        [Required(ErrorMessage = "AdminID is required")]
        [JsonPropertyName("admin_id")]
        public Guid AdminID { get; set; }

        [Required(ErrorMessage = "Name is required")]
        [JsonPropertyName("name")]
        public string Name { get; set; } = string.Empty;

        [Required(ErrorMessage = "Category is required")]
        [JsonPropertyName("category")]
        public string Category { get; set; } = string.Empty;

        [JsonPropertyName("slug")]
        public string? Slug { get; set; }

        [Required(ErrorMessage = "Stock is required")]
        [JsonPropertyName("stock")]
        public int Stock { get; set; }

        [Required(ErrorMessage = "Price amount is required")]
        [JsonPropertyName("price_amount")]
        public long PriceAmount { get; set; }

        // Default to KES
        [JsonPropertyName("price_currency")]
        public string PriceCurrency { get; set; } = "KES";

        [Required(ErrorMessage = "Unit is required")]
        [JsonPropertyName("unit")]
        public string Unit { get; set; } = string.Empty;

        [Required(ErrorMessage = "Packaging is required")]
        [JsonPropertyName("packaging")]
        public string Packaging { get; set; } = string.Empty;

        [Required(ErrorMessage = "Description is required")]
        [JsonPropertyName("description")]
        public string Description { get; set; } = string.Empty;

        [Required(ErrorMessage = "Location is required")]
        [JsonPropertyName("location")]
        public string Location { get; set; } = string.Empty;

        [Required(ErrorMessage = "Min of 3 images is required")]
        [JsonPropertyName("images")]
        public string Images { get; set; } = string.Empty;

    }

    public class StorePublicView
    {
        public string AdminName { get; set; } = string.Empty;
        public List<Inventory> Products { get; set; } = new();
    }
}
