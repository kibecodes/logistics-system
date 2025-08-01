using System.ComponentModel.DataAnnotations;

namespace logistics_frontend.Models.Inventory
{
    public class Inventory
    {
        public Guid ID { get; set; }
        public Guid AdminID { get; set; }
        public string Name { get; set; } = string.Empty;
        public string? Slug { get; set; }
        public string Category { get; set; } = string.Empty;
        public int Stock { get; set; }
        public float Price { get; set; }
        public string Images { get; set; } = string.Empty;
        public string Unit { get; set; } = string.Empty;
        public string Packaging { get; set; } = string.Empty;
        public string Description { get; set; } = string.Empty;
        public string Location { get; set; } = string.Empty;
        public DateTime CreatedAt { get; set; }
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
