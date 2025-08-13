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

        [JsonPropertyName("price")]
        public Money Price { get; set; } = new Money();

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
        public Guid AdminID { get; set; }

        [Required(ErrorMessage = "Name is required")]
        public string Name { get; set; } = string.Empty;

        // [Required(ErrorMessage = "Slug is required")]
        public string? Slug { get; set; }

        [Required(ErrorMessage = "Category is required")]
        public string Category { get; set; } = string.Empty;

        [Required(ErrorMessage = "Stock is required")]
        public int Stock { get; set; }

        [Required(ErrorMessage = "Price is required")]
        public float Price { get; set; }

        // [Required(ErrorMessage = "Images is required")]
        public string Images { get; set; } = string.Empty;

        [Required(ErrorMessage = "Unit is required")]
        public string Unit { get; set; } = string.Empty;

        [Required(ErrorMessage = "Packaging is required")]
        public string Packaging { get; set; } = string.Empty;

        [Required(ErrorMessage = "Description is required")]
        public string Description { get; set; } = string.Empty;

        [Required(ErrorMessage = "Location is required")]
        public string Location { get; set; } = string.Empty;

    }

    public class StorePublicView
    {
        public string AdminName { get; set; } = string.Empty;
        public List<Inventory> Products { get; set; } = new();
    }
}
