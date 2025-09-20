namespace logistics_frontend.Models.Store
{
    public class Store
    {
        public Guid AdminID { get; set; }
        public string AdminName { get; set; } = string.Empty;
        public Product? Products { get; set; }
    }

    public class Product
    {
        public Guid ID { get; set; }
        public string Name { get; set; } = string.Empty;
        public string Price { get; set; } = string.Empty;
        public string Image { get; set; } = string.Empty;
        public string Stock { get; set; } = string.Empty;
        public string Description { get; set; } = string.Empty;
    }
    
}